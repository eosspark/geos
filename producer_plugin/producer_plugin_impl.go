package producer_plugin

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

type ProducerPluginImpl struct {
	ProductionEnabled   bool
	ProductionPaused    bool
	ProductionSkipFlags uint32

	SignatureProviders map[ecc.PublicKey]signatureProviderType
	Producers          map[common.AccountName]struct{}
	Timer              *common.Timer
	ProducerWatermarks map[common.AccountName]uint32
	PendingBlockMode   EnumPendingBlockMode

	PersistentTransactions  transactionIdWithExpireIndex
	BlacklistedTransactions transactionIdWithExpireIndex

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

	ConfirmedBlock func(signature ecc.Signature) //TODO

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

	activeProducers := make(map[common.AccountName]struct{}, len(bsp.ActiveSchedule.Producers))
	for _, p := range bsp.ActiveSchedule.Producers {
		activeProducers[p.AccountName] = struct{}{}
	}

	for producer := range impl.Producers {
		_, has := activeProducers[producer]
		if !has {
			continue
		}

		if producer != bsp.Header.Producer {
			var itr *types.ProducerKey
			for _, k := range activeProducerToSigningKey {
				if k.AccountName == producer {
					itr = &k
					break
				}
			}

			if itr != nil {
				privateKeyItr := impl.SignatureProviders[itr.BlockSigningKey]
				if privateKeyItr != nil {
					d := bsp.SigDigest()
					sig := privateKeyItr(d)
					impl.LastSignedBlockTime = bsp.Header.Timestamp.ToTimePoint()
					impl.LastSignedBlockNum = bsp.BlockNum

					impl.ConfirmedBlock(sig)
				}
			}
		}
	} //set_intersection

	// since the watermark has to be set before a block is created, we are looking into the future to
	// determine the new schedule to identify producers that have become active
	hbn := bsp.BlockNum
	newBlockHeader := bsp.Header
	newBlockHeader.Timestamp = newBlockHeader.Timestamp.Next()
	newBlockHeader.Previous = bsp.ID
	newBs := bsp.GenerateNext(&newBlockHeader.Timestamp)

	// for newly installed producers we can set their watermarks to the block they became active
	if newBs.MaybePromotePending() && bsp.ActiveSchedule.Version != newBs.ActiveSchedule.Version {
		newProducers := make(map[common.AccountName]struct{}, len(newBs.ActiveSchedule.Producers))
		for _, p := range newBs.ActiveSchedule.Producers {
			if _, has := impl.Producers[p.AccountName]; has {
				newProducers[p.AccountName] = struct{}{}
			}
		}

		for _, p := range bsp.ActiveSchedule.Producers {
			delete(newProducers, p.AccountName)
		}

		for newProducer := range newProducers {
			impl.ProducerWatermarks[newProducer] = hbn
		}
	}
}

func (impl *ProducerPluginImpl) OnIrreversibleBlock(lib *types.SignedBlock) {
	impl.IrreversibleBlockTime = lib.Timestamp.ToTimePoint()
}

func (impl *ProducerPluginImpl) OnIncomingBlock(block *types.SignedBlock) {
	//fc_dlog(_log, "received incoming block ${id}", ("id", block->id()));

	if block.Timestamp.ToTimePoint() >= (common.Now().AddUs(common.Seconds(7))) {
		panic(ErrBlockFromTheFuture)
	}
	id := block.BlockID()
	existing := chain.FetchBlockById(id)
	if existing != nil {
		return
	}

	// abort the pending block
	chain.AbortBlock()

	// make sure we restart our loop
	defer impl.ScheduleProductionLoop()

	// push the new block
	except := false

	chain.PushBlock(block)

	if except {
		//C++ app().get_channel<channels::rejected_block>().publish( block );
		return
	}

	if chain.HeadBlockState().Header.Timestamp.Next().ToTimePoint() >= common.Now() {
		impl.ProductionEnabled = true
	}

	if common.Now().Sub(block.Timestamp.ToTimePoint()) < common.Minutes(5) || block.BlockNumber()%1000 == 0 {
		//	ilog("Received block ${id}... #${n} @ ${t} signed by ${p} [trxs: ${count}, lib: ${lib}, conf: ${confs}, latency: ${latency} ms]",
		//		("p",block->producer)("id",fc::variant(block->id()).as_string().substr(8,16))
		//	("n",block_header::num_from_id(block->id()))("t",block->timestamp)
		//	("count",block->transactions.size())("lib",chain.last_irreversible_block_num())("confs", block->confirmed)("latency", (fc::time_point::now() - block->timestamp).count()/1000 ) );
	}
}

type pendingIncomingTransaction struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(interface{})
}

func (impl *ProducerPluginImpl) OnIncomingTransactionAsync(trx *types.PackedTransaction, persistUntilExpired bool, next func(interface{})) {
	if chain.PendingBlockState() == nil {
		impl.PendingIncomingTransactions = append(impl.PendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
		return
	}

	blockTime := chain.PendingBlockState().Header.Timestamp.ToTimePoint()

	sendResponse := func(response interface{}) {
		next(response)
		if _, ok := response.(error); ok {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(response.get<fc::exception_ptr>(), trx));
		} else {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(nullptr, trx));
		}
	}

	id := trx.ID()
	if trx.Expiration().ToTimePoint() < blockTime {
		sendResponse(errors.New(fmt.Sprintf("expired transaction %s", id)))
		return
	}

	if chain.IsKnownUnexpiredTransaction(id) {
		sendResponse(errors.New(fmt.Sprintf("duplicate transaction %s", id)))
		return
	}

	deadline := common.Now().AddUs(common.Milliseconds(int64(impl.MaxTransactionTimeMs)))
	deadlineIsSubjective := false

	if impl.MaxTransactionTimeMs < 0 || impl.PendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
		deadlineIsSubjective = true
		deadline = blockTime
	}

	trace := chain.PushTransaction(types.NewTransactionMetadata(*trx), deadline)
	if trace.Except != nil {
		if failureIsSubjective(trace.Except, deadlineIsSubjective) {
			impl.PendingIncomingTransactions = append(impl.PendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
		} else {
			sendResponse(trace.Except)
		}
	} else {
		if persistUntilExpired {
			// if this trx didnt fail/soft-fail and the persist flag is set, store its ID so that we can
			// ensure its applied to all future speculative blocks as well.
			impl.PersistentTransactions[trx.ID()] = trx.Expiration().ToTimePoint()
		}
		sendResponse(trace)
	}
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
