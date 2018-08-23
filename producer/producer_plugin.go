package producer_plugin

import (
	"fmt"
	"github.com/eoscanada/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"time"
)

type PendingBlockMode int

const (
	producing = PendingBlockMode(iota)
	speculating
)

type signatureProviderType func([]byte) ecc.Signature

type respVariant struct {
	err error
	trx int //TODO:transaction_trace_ptr

}

type tuple struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(respVariant)
}

// producer_plugin
type ProducerPlugin struct {
	producers          map[common.AccountName]struct{}
	pendingBlockMode   PendingBlockMode
	productionEnabled  bool
	productionPaused   bool
	signatureProviders map[ecc.PublicKey]signatureProviderType
	producerWatermarks map[common.AccountName]uint32

	maxTransactionTimeMs      int32
	maxIrreversibleBlockAgeUs time.Duration
	produceTImeOffsetUs       int32
	lastBlockTimeOffsetUs     int32
	irreversibleBlockTime     time.Time
	keosdProviderTimeoutUs    time.Duration

	lastSignedBlockTime time.Time
	startTime           time.Time
	lastSignedBlockNum  uint32

	confirmedBlock func(signature ecc.Signature)

	pendingIncomingTransaction []tuple
}

func (pp *ProducerPlugin) init() {
	pp.producers = make(map[common.AccountName]struct{})
	pp.signatureProviders = make(map[ecc.PublicKey]signatureProviderType)
	pp.producerWatermarks = make(map[common.AccountName]uint32)
}

func (pp *ProducerPlugin) IsProducerKey(key ecc.PublicKey) bool {
	privateKey := pp.signatureProviders[key]
	if privateKey != nil {
		return true
	}
	return false
}

func (pp *ProducerPlugin) SignCompact(key *ecc.PublicKey, digest common.SHA256Bytes) (ecc.Signature, error) {
	if key != nil {
		privateKeyFunc := pp.signatureProviders[*key]
		if privateKeyFunc == nil {
			//EOS_ASSERT(private_key_itr != my->_signature_providers.end(), producer_priv_key_not_found, "Local producer has no private key in config.ini corresponding to public key ${key}", ("key", key));
			return ecc.Signature{}, nil
		}

		privateKeyFunc(digest)
	}
	return ecc.Signature{}, nil
}

func (pp *ProducerPlugin) Startup() {
	pp.scheduleProductionLoop()
}

func (pp *ProducerPlugin) Initialize() {
	pp.signatureProviders[ecc.PublicKey{}] = func(hash []byte) ecc.Signature {
		priKey, _ := ecc.NewPrivateKey("privateKey")
		sig, _ := priKey.Sign(hash)
		return sig
	}
}

func (pp *ProducerPlugin) Pause() {
	pp.productionPaused = true
}

func (pp *ProducerPlugin) Resume() {
	pp.productionPaused = false
	// it is possible that we are only speculating because of this policy which we have now changed
	// re-evaluate that now
	//
	if pp.pendingBlockMode == speculating {
		chain.AbortBlock()
		pp.scheduleProductionLoop()
	}
}

// producer_plugin_impl

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
var timerCorelationId uint32 = 0

// keep a expected ratio between defer txn and incoming txn
var incomingTrxWeight = 0.0
var incomingDeferRadio = 1.0 // 1:1

func (pp *ProducerPlugin) onBlock(bsp *types.BlockState) {
	if !bsp.Header.Timestamp.ToTimePoint().After(pp.lastSignedBlockTime) {
		return
	}
	if !bsp.Header.Timestamp.ToTimePoint().After(pp.startTime) {
		return
	}
	if bsp.BlockNum <= pp.lastSignedBlockNum {
		return
	}

	activeProducerToSigningKey := bsp.ActiveSchedule.Producers
	activeProducers := make(map[common.AccountName]struct{}, len(bsp.ActiveSchedule.Producers))
	for _, p := range bsp.ActiveSchedule.Producers {
		activeProducers[p.AccountName] = struct{}{}
	}

	for producer := range pp.producers {
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
				privateKeyItr := pp.signatureProviders[itr.BlockSigningKey]
				if privateKeyItr != nil {
					d := bsp.SigDigest()
					sig := privateKeyItr(d)
					pp.lastSignedBlockTime = bsp.Header.Timestamp.ToTimePoint()
					pp.lastSignedBlockNum = bsp.BlockNum

					pp.confirmedBlock(sig)
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
	newBs := bsp.GenerateNext(newBlockHeader.Timestamp)

	// for newly installed producers we can set their watermarks to the block they became active
	if newBs.MaybePromotePending() && bsp.ActiveSchedule.Version != newBs.ActiveSchedule.Version {
		newProducers := make(map[common.AccountName]struct{}, len(newBs.ActiveSchedule.Producers))
		for _, p := range newBs.ActiveSchedule.Producers {
			if _, has := pp.producers[p.AccountName]; has {
				newProducers[p.AccountName] = struct{}{}
			}
		}

		for _, p := range bsp.ActiveSchedule.Producers {
			delete(newProducers, p.AccountName)
		}

		for newProducer := range newProducers {
			pp.producerWatermarks[newProducer] = hbn
		}
	}

}

func (pp *ProducerPlugin) onIncomingBlock(block *types.SignedBlock) {

	//C++ EOS_ASSERT( block->timestamp < (fc::time_point::now() + fc::seconds(7))
	id, err := block.BlockID()
	//C++ auto existing = chain.fetch_block_by_id( id );
	fmt.Println(id, err)

	existing := false
	if existing {
		return
	}

	// abort the pending block
	chain.AbortBlock()

	// make sure we restart our loop
	defer pp.scheduleProductionLoop()

	// push the new block
	except := false

	chain.PushBlock(block)

	if except {
		//C++ app().get_channel<channels::rejected_block>().publish( block );
		return
	}

	//C++ if( chain.head_block_state()->header.timestamp.next().to_time_point() >= fc::time_point::now() ) {
	if !chain.HeadBlockState().Header.Timestamp.Next().ToTimePoint().Before(time.Now()) {
		pp.productionEnabled = true
	}

	//C++ log per 1000 blocks

}

func (pp *ProducerPlugin) onIncomingTransactionAsync(trx *types.PackedTransaction, persistUntilExpired bool, next func(respVariant)) {
	if chain.PendingBlockState() == nil {
		pp.pendingIncomingTransaction = append(pp.pendingIncomingTransaction, tuple{trx, persistUntilExpired, next})
		return
	}

	blockTime := chain.PendingBlockState().Header.Timestamp

	sendResponse := func(response respVariant) {
		//TODO
		next(response)
		if response.err != nil {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(response.get<fc::exception_ptr>(), trx));
		} else {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(nullptr, trx));
		}
	}

	id := trx.ID()

	fmt.Println(blockTime, sendResponse, id)

}

func (pp *ProducerPlugin) getIrreversibleBlockAge() time.Duration {
	now := time.Now()
	if now.Before(pp.irreversibleBlockTime) {
		return 0
	} else {
		return time.Duration((now.UnixNano() - pp.irreversibleBlockTime.UnixNano()) / 1e3)
	}
}

func (pp *ProducerPlugin) productionDisabledByPolicy() bool {
	return !pp.productionEnabled || pp.productionPaused || (pp.maxIrreversibleBlockAgeUs >= 0 && pp.getIrreversibleBlockAge() >= pp.maxIrreversibleBlockAgeUs)
}

type StartBlockRusult int

const (
	succeeded = StartBlockRusult(iota)
	failed
	waiting
	exhausted
)

func (pp *ProducerPlugin) calculateNextBlockTime(producerName common.AccountName, currentBlockTime common.BlockTimeStamp) *time.Time {
	hbs := chain.HeadBlockState()
	activeSchedule := hbs.ActiveSchedule.Producers

	pbs := chain.PendingBlockState()
	pbt := pbs.Header.Timestamp

	// determine if this producer is in the active schedule and if so, where
	var itr *types.ProducerKey
	var producerIndex int
	for index, asp := range activeSchedule {
		if asp.AccountName == producerName {
			itr = &asp
			producerIndex = index
			break
		}
	}

	if itr == nil {
		// this producer is not in the active producer set
		return nil
	}

	var minOffset uint32 = 1 // must at least be the "next" block

	// account for a watermark in the future which is disqualifying this producer for now
	// this is conservative assuming no blocks are dropped.  If blocks are dropped the watermark will
	// disqualify this producer for longer but it is assumed they will wake up, determine that they
	// are disqualified for longer due to skipped blocks and re-caculate their next block with better
	// information then
	currentWatermark, hasCurrentWatermark := pp.producerWatermarks[producerName]
	if hasCurrentWatermark {
		if currentWatermark > pbs.BlockNum {
			minOffset = currentWatermark - pbs.BlockNum + 1
		}
	}
	fmt.Println(minOffset, producerIndex, pbt)

	// this producers next opportuity to produce is the next time its slot arrives after or at the calculated minimum
	//TODO

	now := time.Now()
	return &now
}

func (pp *ProducerPlugin) calculatePendingBlockTime() time.Time {
	now := time.Now()
	var base time.Time
	if now.After(chain.HeadBlockTime()) {
		base = now
	} else {
		base = chain.HeadBlockTime()
	}
	minTimeToNextBlock := int64(chain.BlockIntervalUs) - base.UnixNano()/1e3%int64(chain.BlockIntervalUs)
	blockTime := base.Add(time.Microsecond * time.Duration(minTimeToNextBlock))

	if blockTime.Sub(now) < time.Microsecond*time.Duration(chain.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime.Add(time.Microsecond * time.Duration(chain.BlockIntervalUs))
	}

	return blockTime
}

func (pp *ProducerPlugin) startBlock() StartBlockRusult {
	fmt.Println("start_block")

	hbs := chain.HeadBlockState()

	now := time.Now()
	blockTime := pp.calculatePendingBlockTime()

	pp.pendingBlockMode = producing

	// Not our turn
	scheduleProducer := hbs.GetScheduledProducer(common.NewBlockTimeStamp(blockTime))
	currentWatermark, hasCurrentWatermark := pp.producerWatermarks[scheduleProducer.AccountName]
	_, hasSignatureProvider := pp.signatureProviders[scheduleProducer.BlockSigningKey]
	irreversibleBlockAge := pp.getIrreversibleBlockAge()

	// If the next block production opportunity is in the present or future, we're synced.
	if !pp.productionEnabled {
		pp.pendingBlockMode = speculating
	} else if _, has := pp.producers[scheduleProducer.AccountName]; !has {
		pp.pendingBlockMode = speculating
	} else if !hasSignatureProvider {
		pp.pendingBlockMode = speculating
		//elog("Not producing block because I don't have the private key for ${scheduled_key}", ("scheduled_key", scheduled_producer.block_signing_key));
	} else if pp.productionPaused {
		//elog("Not producing block because production is explicitly paused");
		pp.pendingBlockMode = speculating
	} else if pp.maxIrreversibleBlockAgeUs >= 0 && irreversibleBlockAge >= pp.maxIrreversibleBlockAgeUs {
		//elog("Not producing block because the irreversible block is too old [age:${age}s, max:${max}s]", ("age", irreversible_block_age.count() / 1'000'000)( "max", _max_irreversible_block_age_us.count() / 1'000'000 ));
		pp.pendingBlockMode = speculating
	}

	if pp.pendingBlockMode == producing {
		if hasCurrentWatermark {
			if currentWatermark >= hbs.BlockNum+1 {
				/*
								elog("Not producing block because \"${producer}\" signed a BFT confirmation OR block at a higher block number (${watermark}) than the current fork's head (${head_block_num})",
				                ("producer", scheduled_producer.producer_name)
				                ("watermark", currrent_watermark_itr->second)
				                ("head_block_num", hbs->block_num));
				*/
				pp.pendingBlockMode = speculating
			}

		}
	}

	if pp.pendingBlockMode == speculating {
		headBlockAge := now //- chain.head_block_time();
		if headBlockAge.Unix() > 5 {
			return waiting
		}
	}

	var blocksToConfirm uint16 = 0

	if pp.pendingBlockMode == producing {
		// determine how many blocks this producer can confirm
		// 1) if it is not a producer from this node, assume no confirmations (we will discard this block anyway)
		// 2) if it is a producer on this node that has never produced, the conservative approach is to assume no
		//    confirmations to make sure we don't double sign after a crash TODO: make these watermarks durable?
		// 3) if it is a producer on this node where this node knows the last block it produced, safely set it -UNLESS-
		// 4) the producer on this node's last watermark is higher (meaning on a different fork)
		if hasCurrentWatermark {
			if currentWatermark < hbs.BlockNum {
				if hbs.BlockNum-currentWatermark >= 0xffff {
					blocksToConfirm = 0xffff
				} else {
					blocksToConfirm = uint16(hbs.BlockNum - currentWatermark)
				}
			}
		}
	}
	fmt.Println(blocksToConfirm)

	chain.AbortBlock()
	chain.StartBlock(common.NewBlockTimeStamp(blockTime), blocksToConfirm)

	pbs := chain.PendingBlockState()

	if pbs != nil {

		if pp.pendingBlockMode == producing && pbs.BlockSigningKey == scheduleProducer.BlockSigningKey {
			//C++ elog("Block Signing Key is not expected value, reverting to speculative mode! [expected: \"${expected}\", actual: \"${actual\"", ("expected", scheduled_producer.block_signing_key)("actual", pbs->block_signing_key));
			pp.pendingBlockMode = speculating
		}

		// attempt to play persisted transactions first
		//exhausted := false
	}

	return failed
}

func (pp *ProducerPlugin) scheduleProductionLoop() {
	result := pp.startBlock()

	if result == failed {
		// we failed to start a block, so try again later?
		timerCorelationId++
		cid := timerCorelationId
		time.AfterFunc(time.Microsecond*(time.Duration(chain.BlockIntervalUs/10)), func() {
			if cid == timerCorelationId {
				pp.scheduleProductionLoop()
			}
		})

	} else if result == waiting {
		if len(pp.producers) > 0 && !pp.productionDisabledByPolicy() {
			//fc_dlog(_log, "Waiting till another block is received and scheduling Speculative/Production Change");
			pp.scheduleDelayedProductionLoop(common.NewBlockTimeStamp(pp.calculatePendingBlockTime()))
		} else {
			//fc_dlog(_log, "Waiting till another block is received");
			// nothing to do until more blocks arrive
		}

	} else if pp.pendingBlockMode == producing {
		// we succeeded but block may be exhausted
		var expires time.Duration
		if result == succeeded {
			expires = time.Millisecond * time.Duration(chain.BlockIntervalMs)
		} else {
			expires = 0
		}

		timerCorelationId++
		cid := timerCorelationId
		time.AfterFunc(expires, func() {
			if cid == timerCorelationId {
				pp.maybeProduceBlock()
			}
		})

	} else if pp.pendingBlockMode == speculating && len(pp.producers) > 0 && !pp.productionDisabledByPolicy() {
		//fc_dlog(_log, "Specualtive Block Created; Scheduling Speculative/Production Change");
		pbs := chain.PendingBlockState()
		pp.scheduleDelayedProductionLoop(pbs.Header.Timestamp)

	} else {
		//fc_dlog(_log, "Speculative Block Created");
	}
}

func (pp *ProducerPlugin) scheduleDelayedProductionLoop(currentBlockTime common.BlockTimeStamp) {
	var wakeUpTime *time.Time
	for p := range pp.producers {
		nextProducerBlockTime := pp.calculateNextBlockTime(p, currentBlockTime)
		if nextProducerBlockTime != nil {
			producerWakeupTime := time.Unix(0, nextProducerBlockTime.UnixNano()-int64(time.Microsecond*time.Duration(chain.BlockIntervalUs)))
			if wakeUpTime != nil {
				if wakeUpTime.After(producerWakeupTime) {
					wakeUpTime = &producerWakeupTime
				}
			} else {
				wakeUpTime = &producerWakeupTime
			}
		}
	}

	if wakeUpTime != nil {
		//fc_dlog(_log, "Scheduling Speculative/Production Change at ${time}", ("time", wake_up_time));
		timerCorelationId++
		cid := timerCorelationId
		time.AfterFunc(wakeUpTime.Sub(time.Now()), func() {
			if cid == timerCorelationId {
				pp.scheduleProductionLoop()
			}
		})
	} else {
		//fc_dlog(_log, "Speculative Block Created; Not Scheduling Speculative/Production, no local producers had valid wake up times");
	}

}

func (pp *ProducerPlugin) maybeProduceBlock() bool {
	defer pp.scheduleProductionLoop()

	err := pp.produceBlock()

	if err == nil {
		return true
	}

	//C++ chain::controller& chain = app().get_plugin<chain_plugin>().chain();
	chain.AbortBlock()
	return false
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func (pp *ProducerPlugin) produceBlock() error {
	//fmt.Println("produced block 00000002c0b9e9e6... #", block_num)
	if pp.pendingBlockMode != producing {
	}
	pbs := chain.PendingBlockState()
	//C++ hbs := chain.HeadBlockState()
	if pbs == nil {
	}

	signatureProvider := pp.signatureProviders[pbs.BlockSigningKey]
	if signatureProvider == nil {
	}

	chain.FinalizeBlock()
	chain.SignBlock(func(d []byte) ecc.Signature {
		defer makeDebugTimeLogger()
		return signatureProvider(d)
	})

	chain.CommitBlock()

	newBs := chain.HeadBlockState()
	pp.producerWatermarks[newBs.Header.Producer] = chain.HeadBlockNum()

	//fmt.Printf("Produced blcok %v... #%d @ %v signed by %v ",
	//	newBs.Id, newBs.BlockNum, newBs.Header.Timestamp, newBs.Header.Producer, newBs)

	/*ilog("Produced block ${id}... #${n} @ ${t} signed by ${p} [trxs: ${count}, lib: ${lib}, confirmed: ${confs}]",
	  ("p",new_bs->header.producer)("id",fc::variant(new_bs->id).as_string().substr(0,16))
	  ("n",new_bs->block_num)("t",new_bs->header.timestamp)
	  ("count",new_bs->block->transactions.size())("lib",chain.last_irreversible_block_num())("confs", new_bs->header.confirmed));*/

	return nil
}
