package producer_plugin

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/rlp"
)

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

func (pp *ProducerPlugin) onBlock(bsp *types.BlockState) {
	if bsp.Header.Timestamp.ToTimePoint() <= pp.lastSignedBlockTime {
		return
	}
	if bsp.Header.Timestamp.ToTimePoint() <= pp.startTime {
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
	newBs := bsp.GenerateNext(&newBlockHeader.Timestamp)

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

func (pp *ProducerPlugin) onIrreversibleBlock(lib *types.SignedBlock) {
	pp.irreversibleBlockTime = lib.Timestamp.ToTimePoint()
}

func (pp *ProducerPlugin) onIncomingBlock(block *types.SignedBlock) {
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
	defer pp.scheduleProductionLoop()

	// push the new block
	except := false

	chain.PushBlock(block)

	if except {
		//C++ app().get_channel<channels::rejected_block>().publish( block );
		return
	}

	if chain.HeadBlockState().Header.Timestamp.Next().ToTimePoint() >= common.Now() {
		pp.productionEnabled = true
	}

	if common.Now().Sub(block.Timestamp.ToTimePoint()) < common.Minutes(5) || block.BlockNumber()%1000 == 0 {
		//	ilog("Received block ${id}... #${n} @ ${t} signed by ${p} [trxs: ${count}, lib: ${lib}, conf: ${confs}, latency: ${latency} ms]",
		//		("p",block->producer)("id",fc::variant(block->id()).as_string().substr(8,16))
		//	("n",block_header::num_from_id(block->id()))("t",block->timestamp)
		//	("count",block->transactions.size())("lib",chain.last_irreversible_block_num())("confs", block->confirmed)("latency", (fc::time_point::now() - block->timestamp).count()/1000 ) );
	}
}

type ErrorORTrace struct {
	error error
	trace *types.TransactionTrace
}

type pendingIncomingTransaction struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(ErrorORTrace)
}

func (pp *ProducerPlugin) onIncomingTransactionAsync(trx *types.PackedTransaction, persistUntilExpired bool, next func(ErrorORTrace)) {
	if chain.PendingBlockState() == nil {
		pp.pendingIncomingTransactions = append(pp.pendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
		return
	}

	blockTime := chain.PendingBlockState().Header.Timestamp.ToTimePoint()

	sendResponse := func(response ErrorORTrace) {
		next(response)
		if response.error != nil {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(response.get<fc::exception_ptr>(), trx));
		} else {
			//C++ _transaction_ack_channel.publish(std::pair<fc::exception_ptr, packed_transaction_ptr>(nullptr, trx));
		}
	}

	id := trx.ID()
	if trx.Expiration().ToTimePoint() < blockTime {
		sendResponse(ErrorORTrace{errors.New(fmt.Sprintf("expired transaction %s", id)), nil})
		return
	}

	if chain.IsKnownUnexpiredTransaction(id) {
		sendResponse(ErrorORTrace{errors.New(fmt.Sprintf("duplicate transaction %s", id)), nil})
		return
	}

	deadline := common.Now().AddUs(common.Milliseconds(int64(pp.maxTransactionTimeMs)))
	deadlineIsSubjective := false

	if pp.maxTransactionTimeMs < 0 || pp.pendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
		deadlineIsSubjective = true
		deadline = blockTime
	}

	trace := chain.PushTransaction(types.NewTransactionMetadata(*trx), deadline)
	if trace.Except != nil {
		if failureIsSubjective(trace.Except, deadlineIsSubjective) {
			pp.pendingIncomingTransactions = append(pp.pendingIncomingTransactions, pendingIncomingTransaction{trx, persistUntilExpired, next})
		} else {
			sendResponse(ErrorORTrace{trace.Except, nil})
		}
	} else {
		if persistUntilExpired {
			// if this trx didnt fail/soft-fail and the persist flag is set, store its ID so that we can
			// ensure its applied to all future speculative blocks as well.
			pp.persistentTransactions[trx.ID()] = trx.Expiration().ToTimePoint()
		}
		sendResponse(ErrorORTrace{nil, trace})
	}
}

func (pp *ProducerPlugin) getIrreversibleBlockAge() common.Microseconds /*Microsecond*/ {
	now := common.Now()
	if now < pp.irreversibleBlockTime {
		return common.Microseconds(0)
	} else {
		return now.Sub(pp.irreversibleBlockTime)
	}
}

func (pp *ProducerPlugin) productionDisabledByPolicy() bool {
	return !pp.productionEnabled || pp.productionPaused || (pp.maxIrreversibleBlockAgeUs >= 0 && pp.getIrreversibleBlockAge() >= pp.maxIrreversibleBlockAgeUs)
}

func (pp *ProducerPlugin) calculateNextBlockTime(producerName common.AccountName, currentBlockTime common.BlockTimeStamp) *common.TimePoint {
	var result common.TimePoint

	hbs := chain.HeadBlockState()
	activeSchedule := hbs.ActiveSchedule.Producers

	// determine if this producer is in the active schedule and if so, where
	var itr *types.ProducerKey
	var producerIndex uint32
	for index, asp := range activeSchedule {
		if asp.AccountName == producerName {
			itr = &asp
			producerIndex = uint32(index)
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
		blockNum := chain.PendingBlockState().BlockNum
		if chain.PendingBlockState() != nil {
			blockNum++
		}
		if currentWatermark > blockNum {
			minOffset = currentWatermark - blockNum + 1
		}
	}

	// this producers next opportuity to produce is the next time its slot arrives after or at the calculated minimum
	minSlot := uint32(currentBlockTime) + minOffset
	minSlotProducerIndex := (minSlot % (uint32(len(activeSchedule)) * uint32(common.DefaultConfig.ProducerRepetitions))) / uint32(common.DefaultConfig.ProducerRepetitions)
	if producerIndex == minSlotProducerIndex {
		// this is the producer for the minimum slot, go with that
		result = common.BlockTimeStamp(minSlot).ToTimePoint()
	} else {
		// calculate how many rounds are between the minimum producer and the producer in question
		producerDistance := producerIndex - minSlotProducerIndex
		// check for unsigned underflow
		if producerDistance > producerIndex {
			producerDistance += uint32(len(activeSchedule))
		}

		// align the minimum slot to the first of its set of reps
		firstMinProducerSlot := minSlot - (minSlot % uint32(common.DefaultConfig.ProducerRepetitions))

		// offset the aligned minimum to the *earliest* next set of slots for this producer
		nextBlockSlot := firstMinProducerSlot + (producerDistance * uint32(common.DefaultConfig.ProducerRepetitions))
		result = common.BlockTimeStamp(nextBlockSlot).ToTimePoint()

	}
	return &result
}

func (pp *ProducerPlugin) calculatePendingBlockTime() common.TimePoint {
	now := common.Now()
	var base common.TimePoint
	if now > chain.HeadBlockTime() {
		base = now
	} else {
		base = chain.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}

	return blockTime
}

func (pp *ProducerPlugin) startBlock() (EnumStartBlockRusult, bool) {
	hbs := chain.HeadBlockState()

	//Schedule for the next second's tick regardless of chain state
	// If we would wait less than 50ms (1/10 of block_interval), wait for the whole block interval.
	now := common.Now()
	blockTime := pp.calculatePendingBlockTime()

	pp.pendingBlockMode = EnumPendingBlockMode(producing)

	// Not our turn
	lastBlock := uint32(common.NewBlockTimeStamp(blockTime))%uint32(common.DefaultConfig.ProducerRepetitions) == uint32(common.DefaultConfig.ProducerRepetitions)-1
	scheduleProducer := hbs.GetScheduledProducer(common.NewBlockTimeStamp(blockTime))
	currentWatermark, hasCurrentWatermark := pp.producerWatermarks[scheduleProducer.AccountName]
	_, hasSignatureProvider := pp.signatureProviders[scheduleProducer.BlockSigningKey]
	irreversibleBlockAge := pp.getIrreversibleBlockAge()

	// If the next block production opportunity is in the present or future, we're synced.
	if !pp.productionEnabled {
		pp.pendingBlockMode = EnumPendingBlockMode(speculating)
	} else if _, has := pp.producers[scheduleProducer.AccountName]; !has {
		pp.pendingBlockMode = EnumPendingBlockMode(speculating)
	} else if !hasSignatureProvider {
		pp.pendingBlockMode = EnumPendingBlockMode(speculating)
		//elog("Not producing block because I don't have the private key for ${scheduled_key}", ("scheduled_key", scheduled_producer.block_signing_key));
	} else if pp.productionPaused {
		//elog("Not producing block because production is explicitly paused");
		pp.pendingBlockMode = EnumPendingBlockMode(speculating)
	} else if pp.maxIrreversibleBlockAgeUs >= 0 && irreversibleBlockAge >= pp.maxIrreversibleBlockAgeUs {
		//elog("Not producing block because the irreversible block is too old [age:${age}s, max:${max}s]", ("age", irreversible_block_age.count() / 1'000'000)( "max", _max_irreversible_block_age_us.count() / 1'000'000 ));
		pp.pendingBlockMode = EnumPendingBlockMode(speculating)
	}

	if pp.pendingBlockMode == EnumPendingBlockMode(producing) {
		if hasCurrentWatermark {
			if currentWatermark >= hbs.BlockNum+1 {
				/*
									elog("Not producing block because \"${producer}\" signed a BFT confirmation OR block at a higher block number (${watermark}) than the current fork's head (${head_block_num})",
					                ("producer", scheduled_producer.producer_name)
					                ("watermark", currrent_watermark_itr->second)
					                ("head_block_num", hbs->block_num));
				*/
				pp.pendingBlockMode = EnumPendingBlockMode(speculating)
			}

		}
	}

	if pp.pendingBlockMode == EnumPendingBlockMode(speculating) {
		headBlockAge := now.Sub(chain.HeadBlockTime())
		if headBlockAge > common.Seconds(5) {
			return EnumStartBlockRusult(waiting), lastBlock
		}
	}

	var blocksToConfirm uint16 = 0

	if pp.pendingBlockMode == EnumPendingBlockMode(producing) {
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

	chain.AbortBlock()
	chain.StartBlock(common.NewBlockTimeStamp(blockTime), blocksToConfirm)

	pbs := chain.PendingBlockState()

	if pbs != nil {

		if pp.pendingBlockMode == EnumPendingBlockMode(producing) && pbs.BlockSigningKey != scheduleProducer.BlockSigningKey {
			//C++ elog("Block Signing Key is not expected value, reverting to speculative mode! [expected: \"${expected}\", actual: \"${actual\"", ("expected", scheduled_producer.block_signing_key)("actual", pbs->block_signing_key));
			pp.pendingBlockMode = EnumPendingBlockMode(speculating)
		}

		// attempt to play persisted transactions first
		isExhausted := false

		// remove all persisted transactions that have now expired
		for byTrxId, byExpire := range pp.persistentTransactions {
			if byExpire <= pbs.Header.Timestamp.ToTimePoint() {
				delete(pp.persistentTransactions, byTrxId)
			}
		}

		origPendingTxnSize := len(pp.pendingIncomingTransactions)

		if len(pp.persistentTransactions) > 0 || pp.pendingBlockMode == EnumPendingBlockMode(producing) {
			unappliedTrxs := chain.GetUnappliedTransactions()

			if len(pp.persistentTransactions) > 0 {
				for i, trx := range unappliedTrxs {
					if _, has := pp.persistentTransactions[trx.ID]; has {
						// this is a persisted transaction, push it into the block (even if we are speculating) with
						// no deadline as it has already passed the subjective deadlines once and we want to represent
						// the state of the chain including this transaction
						err := chain.PushTransaction(trx, common.MaxTimePoint())
						if err != nil {
							return EnumStartBlockRusult(failed), lastBlock
						}
					}

					// remove it from further consideration as it is applied
					unappliedTrxs[i] = nil
				}
			}

			if pp.pendingBlockMode == EnumPendingBlockMode(producing) {
				for _, trx := range unappliedTrxs {
					if blockTime <= common.Now() {
						isExhausted = true
					}
					if isExhausted {
						break
					}

					if trx == nil {
						// nulled in the loop above, skip it
						continue
					}

					if trx.PackedTrx.Expiration().ToTimePoint() < pbs.Header.Timestamp.ToTimePoint() {
						// expired, drop it
						chain.DropUnappliedTransaction(trx)
						continue
					}

					deadline := common.Now().AddUs(common.Microseconds(pp.maxTransactionTimeMs))
					deadlineIsSubjective := false
					if pp.maxTransactionTimeMs < 0 || pp.pendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
						deadlineIsSubjective = true
						deadline = blockTime
					}

					trace := chain.PushTransaction(trx, deadline)
					if trace.Except != nil {
						if failureIsSubjective(trace.Except, deadlineIsSubjective) {
							isExhausted = true
						} else {
							// this failed our configured maximum transaction time, we don't want to replay it
							chain.DropUnappliedTransaction(trx)
						}
					}

					//TODO catch exception
				}
			}
		} ///unapplied transactions

		if pp.pendingBlockMode == EnumPendingBlockMode(producing) {
			for byTrxId, byExpire := range pp.blacklistedTransactions {
				if byExpire <= common.Now() {
					delete(pp.blacklistedTransactions, byTrxId)
				}
			}

			scheduledTrxs := chain.GetScheduledTransactions()

			for _, trx := range scheduledTrxs {
				if blockTime <= common.Now() {
					isExhausted = true
				}
				if isExhausted {
					break
				}

				// configurable ratio of incoming txns vs deferred txns
				for pp.incomingTrxWeight >= 1.0 && origPendingTxnSize > 0 && len(pp.pendingIncomingTransactions) > 0 {
					e := pp.pendingIncomingTransactions[0]
					pp.pendingIncomingTransactions = pp.pendingIncomingTransactions[1:]
					origPendingTxnSize--
					pp.incomingTrxWeight -= 1.0
					pp.onIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
				}

				if blockTime <= common.Now() {
					isExhausted = true
					break
				}

				if _, has := pp.blacklistedTransactions[trx]; has {
					continue
				}

				deadline := common.Now().AddUs(common.Microseconds(pp.maxTransactionTimeMs))
				deadlineIsSubjective := false
				if pp.maxTransactionTimeMs < 0 || pp.pendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
					deadlineIsSubjective = true
					deadline = blockTime
				}

				trace := chain.PushScheduledTransaction(trx, deadline)
				if trace.Except != nil {
					if failureIsSubjective(trace.Except, deadlineIsSubjective) {
						isExhausted = true
					} else {
						expiration := common.Now().AddUs(common.Seconds(0) /*TODO chain.get_global_properties().configuration.deferred_trx_expiration_window*/)
						// this failed our configured maximum transaction time, we don't want to replay it add it to a blacklist
						pp.blacklistedTransactions[trx] = expiration
					}
				}

				//TODO catch exception

				pp.incomingTrxWeight += pp.incomingDeferRadio
				if origPendingTxnSize <= 0 {
					pp.incomingTrxWeight = 0.0
				}
			}
		} ///scheduled transactions

		if isExhausted || blockTime <= common.Now() {
			return EnumStartBlockRusult(exhausted), lastBlock
		} else {
			// attempt to apply any pending incoming transactions
			pp.incomingTrxWeight = 0.0
			if origPendingTxnSize > 0 && len(pp.pendingIncomingTransactions) > 0 {
				e := pp.pendingIncomingTransactions[0]
				pp.pendingIncomingTransactions = pp.pendingIncomingTransactions[1:]
				origPendingTxnSize--
				pp.onIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
				if blockTime <= common.Now() {
					return EnumStartBlockRusult(exhausted), lastBlock
				}
			}
			return EnumStartBlockRusult(succeeded), lastBlock
		}
	}
	return EnumStartBlockRusult(failed), lastBlock
}

func (pp *ProducerPlugin) scheduleProductionLoop() {
	pp.timer.Cancel()

	result, lastBlock := pp.startBlock()

	if result == EnumStartBlockRusult(failed) {
		//elog("Failed to start a pending block, will try again later");
		pp.timer.ExpiresFromNow(common.Microseconds(common.DefaultConfig.BlockIntervalUs / 10))

		// we failed to start a block, so try again later?
		timerCorelationId++
		cid := timerCorelationId
		pp.timer.AsyncWait(func() {
			if cid == timerCorelationId {
				pp.scheduleProductionLoop()
			}
		})

	} else if result == EnumStartBlockRusult(waiting) {
		if len(pp.producers) > 0 && !pp.productionDisabledByPolicy() {
			log.Debug("Waiting till another block is received and scheduling Speculative/Production Change")
			pp.scheduleDelayedProductionLoop(common.NewBlockTimeStamp(pp.calculatePendingBlockTime()))
		} else {
			log.Debug("Waiting till another block is received")
			// nothing to do until more blocks arrive
		}

	} else if pp.pendingBlockMode == EnumPendingBlockMode(producing) {
		// we succeeded but block may be exhausted
		if result == EnumStartBlockRusult(succeeded) {
			// ship this block off no later than its deadline
			if chain.PendingBlockState() == nil {
				panic("producing without pending_block_state, start_block succeeded")
			}
			epoch := chain.PendingBlockTime().TimeSinceEpoch()
			if lastBlock {
				epoch += common.Microseconds(pp.lastBlockTimeOffsetUs)
			} else {
				epoch += common.Microseconds(pp.produceTImeOffsetUs)
			}
			pp.timer.ExpiresAt(epoch)
			log.Debug(fmt.Sprintf("Scheduling Block Production on Normal Block #%dfor %s", chain.PendingBlockState().BlockNum, chain.PendingBlockTime()))
		} else {
			expectTime := chain.PendingBlockTime().SubUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
			// ship this block off up to 1 block time earlier or immediately
			if common.Now() >= expectTime {
				pp.timer.ExpiresFromNow(0)
			} else {
				pp.timer.ExpiresAt(expectTime.TimeSinceEpoch())
			}
			log.Debug("Scheduling Block Production on Exhausted Block #%d immediately", chain.PendingBlockState().BlockNum)
		}

		timerCorelationId++
		cid := timerCorelationId
		pp.timer.AsyncWait(func() {
			if cid == timerCorelationId {
				pp.maybeProduceBlock()
				log.Debug("Producing Block #${num} returned: ${res}")
				//fc_dlog(_log, "Producing Block #${num} returned: ${res}", ("num", chain.pending_block_state()->block_num)("res", res) );
			}
		})

	} else if pp.pendingBlockMode == EnumPendingBlockMode(speculating) && len(pp.producers) > 0 && !pp.productionDisabledByPolicy() {
		//fc_dlog(_log, "Specualtive Block Created; Scheduling Speculative/Production Change");
		pbs := chain.PendingBlockState()
		pp.scheduleDelayedProductionLoop(pbs.Header.Timestamp)

	} else {
		//fc_dlog(_log, "Speculative Block Created");
	}
}

func (pp *ProducerPlugin) scheduleDelayedProductionLoop(currentBlockTime common.BlockTimeStamp) {
	var wakeUpTime *common.TimePoint
	for p := range pp.producers {
		nextProducerBlockTime := pp.calculateNextBlockTime(p, currentBlockTime)
		if nextProducerBlockTime != nil {
			producerWakeupTime := nextProducerBlockTime.SubUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
			if wakeUpTime != nil {
				if *wakeUpTime > producerWakeupTime {
					*wakeUpTime = producerWakeupTime
				}
			} else {
				wakeUpTime = &producerWakeupTime
			}
		}
	}

	if wakeUpTime != nil {
		//fc_dlog(_log, "Scheduling Speculative/Production Change at ${time}", ("time", wake_up_time));
		pp.timer.ExpiresAt(wakeUpTime.TimeSinceEpoch())

		timerCorelationId++
		cid := timerCorelationId
		pp.timer.AsyncWait(func() {
			if cid == timerCorelationId {
				pp.scheduleProductionLoop()
			}
		})
	} else {
		//fc_dlog(_log, "Speculative Block Created; Not Scheduling Speculative/Production, no local producers had valid wake up times");
	}

}

func (pp *ProducerPlugin) maybeProduceBlock() (res bool) {
	defer func() {
		if err := recover(); err != nil {
			chain.AbortBlock()
			res = false
		}

		pp.scheduleProductionLoop()
	}()

	pp.produceBlock()

	return true
}

func (pp *ProducerPlugin) produceBlock() {
	if pp.pendingBlockMode != EnumPendingBlockMode(producing) {
		panic(ErrProducerFail)
	}
	pbs := chain.PendingBlockState()
	if pbs == nil {
		panic(ErrMissingPendingBlockState)
	}

	signatureProvider := pp.signatureProviders[pbs.BlockSigningKey]
	if signatureProvider == nil {
		panic(ErrProducerPriKeyNotFound)
	}

	chain.FinalizeBlock()
	chain.SignBlock(func(d rlp.Sha256) ecc.Signature {
		defer makeDebugTimeLogger()
		return signatureProvider(d)
	})

	chain.CommitBlock()

	newBs := chain.HeadBlockState()
	pp.producerWatermarks[newBs.Header.Producer] = chain.HeadBlockNum()

	fmt.Printf("Produced block #%d @ %s signed by %s [trxs: %d, lib: %d, confirmed: %d]\n",
		newBs.BlockNum, newBs.Header.Timestamp.ToTimePoint(),
		common.NameToString(uint64(newBs.Header.Producer)),
		len(newBs.SignedBlock.Transactions), chain.LastIrreversibleBlockNum(), newBs.Header.Confirmed)
}
