package producer_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	. "github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	. "github.com/eosspark/eos-go/plugins/producer_plugin/multi_index"
)

type ProducerPluginImpl struct {
	Chain *chain.Controller

	ProductionEnabled   bool
	ProductionPaused    bool
	ProductionSkipFlags uint32

	SignatureProviders map[ecc.PublicKey]signatureProviderType
	Producers          AccountNameSet
	Timer              *common.Timer
	ProducerWatermarks map[common.AccountName]uint32
	PendingBlockMode   PendingBlockMode

	PersistentTransactions  *TransactionIdWithExpiryIndex
	BlacklistedTransactions *TransactionIdWithExpiryIndex

	MaxTransactionTimeMs      int32
	MaxIrreversibleBlockAgeUs common.Microseconds
	ProduceTimeOffsetUs       int32
	LastBlockTimeOffsetUs     int32
	IrreversibleBlockTime     common.TimePoint
	KeosdProviderTimeoutUs    common.Microseconds

	LastSignedBlockTime common.TimePoint
	StartTime           common.TimePoint
	LastSignedBlockNum  uint32

	Self *ProducerPlugin

	PendingIncomingTransactions []pendingIncomingTransaction

	/*
	 * HACK ALERT
	 * Boost timers can be in a state where a handler has not yet executed but is not abortable.
	 * As this method needs to mutate state handlers depend on for proper functioning to maintain
	 * invariants for other code (namely accepting incoming transactions in a nearly full block)
	 * the handlers capture a corelation ID at the time they are set.  When they are executed
	 * they must check that correlation_id against the global ordinal.  If it does not match that
	 * implies that this method has been called with the handler in the state where it should be
	 * cancelled but wasn't able to be.
	 */
	timerCorelationId uint32

	// keep a expected ratio between defer txn and incoming txn
	IncomingTrxWeight  float64
	IncomingDeferRadio float64
}

type StartBlockResult int

const (
	succeeded = StartBlockResult(iota)
	failed
	waiting
	exhausted
)

type PendingBlockMode int

const (
	producing = PendingBlockMode(iota)
	speculating
)

type EnumTxCategory int

const (
	PERSISTED = EnumTxCategory(iota)
	UNEXPIRED_UNPERSISTED
	EXPIRED
)

type signatureProviderType = func(sha256 crypto.Sha256) *ecc.Signature

func NewProducerPluginImpl(io *asio.IoContext) *ProducerPluginImpl {
	return &ProducerPluginImpl{
		Timer:                   common.NewTimer(io),
		SignatureProviders:      make(map[ecc.PublicKey]signatureProviderType),
		Producers:               *NewAccountNameSet(),
		ProducerWatermarks:      make(map[common.AccountName]uint32),
		PersistentTransactions:  NewTransactionIdWithExpiryIndex(),
		BlacklistedTransactions: NewTransactionIdWithExpiryIndex(),
		IncomingTrxWeight:       0.0,
		IncomingDeferRadio:      1.0, // 1:1
	}
}

func (impl *ProducerPluginImpl) OnBlock(bsp *types.BlockState) {
	if bsp.Header.Timestamp.ToTimePoint() <= impl.LastSignedBlockTime {
		return
	}
	if bsp.Header.Timestamp.ToTimePoint() <= impl.StartTime {
		return
	}
	if bsp.BlockNum <= impl.LastSignedBlockNum {
		return
	}

	activeProducerToSigningKey := bsp.ActiveSchedule.Producers

	activeProducers := NewAccountNameSet()

	for _, p := range bsp.ActiveSchedule.Producers {
		activeProducers.Add(p.ProducerName)
	}

	AccountNameSetIntersection(&impl.Producers, activeProducers, func(producer common.AccountName) {
		if producer != bsp.Header.Producer {
			var itr *types.ProducerKey
			for _, k := range activeProducerToSigningKey {
				if k.ProducerName == producer {
					itr = &k
				}
			}

			if itr != nil {
				privateKeyItr := impl.SignatureProviders[itr.BlockSigningKey]
				if privateKeyItr != nil {
					//TODO signal ConfirmedBlock
					//d := bsp.SigDigest()
					//sig := privateKeyItr(d)
					impl.LastSignedBlockTime = bsp.Header.Timestamp.ToTimePoint()
					impl.LastSignedBlockNum = bsp.BlockNum

					//impl.Self.ConfirmedBlock
				}
			}
		}
	})

	// since the watermark has to be set before a block is created, we are looking into the future to
	// determine the new schedule to identify producers that have become active
	hbn := bsp.BlockNum
	newBlockHeader := bsp.Header
	newBlockHeader.Timestamp = newBlockHeader.Timestamp.Next()
	newBlockHeader.Previous = bsp.BlockId
	newBs := bsp.GenerateNext(newBlockHeader.Timestamp)

	// for newly installed producers we can set their watermarks to the block they became active
	if newBs.MaybePromotePending() && bsp.ActiveSchedule.Version != newBs.ActiveSchedule.Version {
		newProducers := NewAccountNameSet()
		for _, p := range newBs.ActiveSchedule.Producers {
			if impl.Producers.Contains(p.ProducerName) {
				newProducers.Add(p.ProducerName)
			}
		}

		for _, p := range bsp.ActiveSchedule.Producers {
			newProducers.Remove(p.ProducerName)
		}

		newProducers.Each(func(newProducer common.AccountName) {
			impl.ProducerWatermarks[newProducer] = hbn
		})
	}
}

func (impl *ProducerPluginImpl) OnIrreversibleBlock(lib *types.SignedBlock) {
	impl.IrreversibleBlockTime = lib.Timestamp.ToTimePoint()
}

func (impl *ProducerPluginImpl) OnIncomingBlock(block *types.SignedBlock) {
	log.Debug("received incoming block %s", block.BlockID())

	EosAssert(block.Timestamp.ToTimePoint() < common.Now().AddUs(common.Seconds(7)), &BlockFromTheFuture{}, "received a block from the future, ignoring it")

	chain := impl.Chain

	/* de-dupe here... no point in aborting block if we already know the block */
	id := block.BlockID()
	existing := chain.FetchBlockById(id)
	if existing != nil {
		return
	}

	// abort the pending block
	chain.AbortBlock()

	// exceptions throw out, make sure we restart our loop
	defer func() {
		impl.ScheduleProductionLoop()
	}()

	// push the new block
	except := false

	returning := false
	Try(func() {
		chain.PushBlock(block, types.BlockStatus(types.Complete))
	}).Catch(func(e GuardExceptions) {
		//TODO: handle_guard_exception
		returning = true
		return
	}).Catch(func(e Exception) {
		log.Error(e.DetailMessage())
		except = true
	}).End()

	if returning {
		return
	}

	if except {
		//TODO:C++ app().get_channel<channels::rejected_block>().publish( block );
		return
	}

	if chain.HeadBlockState().Header.Timestamp.Next().ToTimePoint() >= common.Now() {
		impl.ProductionEnabled = true
	}

	if common.Now().Sub(block.Timestamp.ToTimePoint()) < common.Minutes(5) || block.BlockNumber()%1000 == 0 {
		log.Info("Received block %s... #%d @ %s signed by %s [trxs: %d, lib: %d, conf: %d, lantency: %d ms]\n",
			block.BlockID().String()[8:16], block.BlockNumber(), block.Timestamp, block.Producer,
			len(block.Transactions), chain.LastIrreversibleBlockNum(), block.Confirmed, (common.Now().Sub(block.Timestamp.ToTimePoint())).Count()/1000)
	}
}

type pendingIncomingTransaction struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(interface{})
}

func (impl *ProducerPluginImpl) OnIncomingTransactionAsync(trx *types.PackedTransaction, persistUntilExpired bool, next func(interface{})) {
	chain := impl.Chain
	if chain.PendingBlockState() == nil {
		impl.PendingIncomingTransactions = append(impl.PendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
		return
	}

	blockTime := chain.PendingBlockState().Header.Timestamp.ToTimePoint()

	sendResponse := func(response interface{}) {
		next(response)
		if re, ok := response.(Exception); ok {
			//TODO C: _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(response.get<fc::exception_ptr>(), trx));
			if impl.PendingBlockMode == PendingBlockMode(producing) {
				log.Debug("[TRX_TRACE] Block %d for producer %s is REJECTING tx: %s : %s ",
					chain.HeadBlockNum()+1, chain.PendingBlockState().Header.Producer, trx.ID(), re.What())
			} else {
				log.Debug("[TRX_TRACE] Speculative execution is REJECTING tx: %s : %s ",
					trx.ID(), re.What())
			}

		} else {
			//TODO C: _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(nullptr, trx));
			if impl.PendingBlockMode == PendingBlockMode(producing) {
				log.Debug("[TRX_TRACE] Block %d for producer %s is ACCEPTING tx: %s",
					chain.HeadBlockNum()+1, chain.PendingBlockState().Header.Producer, trx.ID())

			} else {
				log.Debug("[TRX_TRACE] Speculative execution is ACCEPTING tx: %s", trx.ID())
			}
		}
	}

	id := trx.ID()
	if trx.Expiration().ToTimePoint() < blockTime {
		sendResponse(&ExpiredTxException{Elog: log.Messages{log.FcLogMessage(log.LvlError, "expired transaction %s", id)}})
		return
	}

	if chain.IsKnownUnexpiredTransaction(&id) {
		sendResponse(&TxDuplicate{Elog: log.Messages{log.FcLogMessage(log.LvlError, "duplicate transaction %s", id)}})
		return
	}

	deadline := common.Now().AddUs(common.Milliseconds(int64(impl.MaxTransactionTimeMs)))
	deadlineIsSubjective := false

	if impl.MaxTransactionTimeMs < 0 || impl.PendingBlockMode == PendingBlockMode(producing) && blockTime < deadline {
		deadlineIsSubjective = true
		deadline = blockTime
	}

	Try(func() {
		trace := chain.PushTransaction(types.NewTransactionMetadata(trx), deadline, 0)
		if trace.Except != nil {
			if failureIsSubjective(trace.Except, deadlineIsSubjective) {
				impl.PendingIncomingTransactions = append(impl.PendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
				if impl.PendingBlockMode == PendingBlockMode(producing) {
					log.Debug("[TRX_TRACE] Block %d for producer %s COULD NOT FIT, tx: %s RETRYING ",
						chain.HeadBlockNum()+1, chain.PendingBlockState().Header.Producer, trx.ID())
				} else {
					log.Debug("[TRX_TRACE] Speculative execution COULD NOT FIT tx: %s} RETRYING", trx.ID())
				}

			} else {
				sendResponse(trace.Except)
			}
		} else {
			if persistUntilExpired {
				// if this trx didnt fail/soft-fail and the persist flag is set, store its ID so that we can
				// ensure its applied to all future speculative blocks as well.
				impl.PersistentTransactions.Insert(TransactionIdWithExpiry{TrxId: trx.ID(), Expiry: trx.Expiration().ToTimePoint()})
			}
			sendResponse(trace)
		}

	}).Catch(func(e GuardExceptions) {
		//TODO: app().get_plugin<chain_plugin>().handle_guard_exception(e);

	}).CatchAndCall(sendResponse).End()
}

func (impl *ProducerPluginImpl) GetIrreversibleBlockAge() common.Microseconds {
	now := common.Now()
	if now < impl.IrreversibleBlockTime {
		return common.Microseconds(0)
	} else {
		return now.Sub(impl.IrreversibleBlockTime)
	}
}

func (impl *ProducerPluginImpl) ProductionDisabledByPolicy() bool {
	return !impl.ProductionEnabled || impl.ProductionPaused || (impl.MaxIrreversibleBlockAgeUs >= 0 && impl.GetIrreversibleBlockAge() >= impl.MaxIrreversibleBlockAgeUs)
}
