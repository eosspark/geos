package producer_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	Chain "github.com/eosspark/eos-go/chain"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
)

func (impl *ProducerPluginImpl) CalculateNextBlockTime(producerName *common.AccountName, currentBlockTime types.BlockTimeStamp) *common.TimePoint {
	var result common.TimePoint

	chain := impl.Self.chain()

	hbs := chain.HeadBlockState()
	activeSchedule := hbs.ActiveSchedule.Producers

	// determine if this producer is in the active schedule and if so, where
	var itr *types.ProducerKey
	var producerIndex uint32
	for index, asp := range activeSchedule {
		if asp.ProducerName == *producerName {
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
	currentWatermark, hasCurrentWatermark := impl.ProducerWatermarks[*producerName]
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
		result = types.BlockTimeStamp(minSlot).ToTimePoint()
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
		result = types.BlockTimeStamp(nextBlockSlot).ToTimePoint()

	}
	return &result
}

func (impl *ProducerPluginImpl) CalculatePendingBlockTime() common.TimePoint {
	chain := impl.Self.chain()
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

func (impl *ProducerPluginImpl) StartBlock() (EnumStartBlockRusult, bool) {
	chain := impl.Self.chain()

	if chain.GetReadMode() == Chain.READONLY {
		return EnumStartBlockRusult(waiting), false
	}

	hbs := chain.HeadBlockState()

	//Schedule for the next second's tick regardless of chain state
	// If we would wait less than 50ms (1/10 of block_interval), wait for the whole block interval.
	now := common.Now()
	blockTime := impl.CalculatePendingBlockTime()

	impl.PendingBlockMode = EnumPendingBlockMode(producing)

	// Not our turn
	lastBlock := uint32(types.NewBlockTimeStamp(blockTime))%uint32(common.DefaultConfig.ProducerRepetitions) == uint32(common.DefaultConfig.ProducerRepetitions)-1
	scheduleProducer := hbs.GetScheduledProducer(types.NewBlockTimeStamp(blockTime))
	currentWatermark, hasCurrentWatermark := impl.ProducerWatermarks[scheduleProducer.ProducerName]
	_, hasSignatureProvider := impl.SignatureProviders[scheduleProducer.BlockSigningKey]
	irreversibleBlockAge := impl.GetIrreversibleBlockAge()

	// If the next block production opportunity is in the present or future, we're synced.
	if !impl.ProductionEnabled {
		impl.PendingBlockMode = EnumPendingBlockMode(speculating)

	} else if !impl.Producers.Contains(scheduleProducer.ProducerName) {
		impl.PendingBlockMode = EnumPendingBlockMode(speculating)

	} else if !hasSignatureProvider {
		log.Error("Not producing block because I don't have the private key for %s", scheduleProducer.BlockSigningKey)
		impl.PendingBlockMode = EnumPendingBlockMode(speculating)

	} else if impl.ProductionPaused {
		log.Error("Not producing block because production is explicitly paused")
		impl.PendingBlockMode = EnumPendingBlockMode(speculating)

	} else if impl.MaxIrreversibleBlockAgeUs >= 0 && irreversibleBlockAge >= impl.MaxIrreversibleBlockAgeUs {
		log.Error("Not producing block because the irreversible block is too old [age:%ds, max:%ds]", irreversibleBlockAge.Count() / 1e6, impl.MaxIrreversibleBlockAgeUs.Count() / 1e6 )
		impl.PendingBlockMode = EnumPendingBlockMode(speculating)
	}

	if impl.PendingBlockMode == EnumPendingBlockMode(producing) {
		if hasCurrentWatermark {
			if currentWatermark >= hbs.BlockNum+1 {
				log.Error("Not producing block because \"%s\" signed a BFT confirmation OR block at a higher block number (%d) than the current fork's head (%d)",
					scheduleProducer.ProducerName,
					currentWatermark,
					hbs.BlockNum)
				impl.PendingBlockMode = EnumPendingBlockMode(speculating)
			}

		}
	}

	if impl.PendingBlockMode == EnumPendingBlockMode(speculating) {
		headBlockAge := now.Sub(chain.HeadBlockTime())
		if headBlockAge > common.Seconds(5) {
			return EnumStartBlockRusult(waiting), lastBlock
		}
	}

	Try(func() {
		blocksToConfirm := uint16(0)

		if impl.PendingBlockMode == EnumPendingBlockMode(producing) {
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
		chain.StartBlock(types.NewBlockTimeStamp(blockTime), blocksToConfirm)

	}).FcLogAndDrop().End()

	pbs := chain.PendingBlockState()

	if pbs != nil {

		if impl.PendingBlockMode == EnumPendingBlockMode(producing) && pbs.BlockSigningKey != scheduleProducer.BlockSigningKey {
			log.Error("Block Signing Key is not expected value, reverting to speculative mode! [expected: \"%s\", actual: \"%s\"", scheduleProducer.BlockSigningKey, pbs.BlockSigningKey)
			impl.PendingBlockMode = EnumPendingBlockMode(speculating)
		}

		// attempt to play persisted transactions first
		isExhausted := false

		// remove all persisted transactions that have now expired
		for byTrxId, byExpire := range impl.PersistentTransactions {
			if byExpire <= pbs.Header.Timestamp.ToTimePoint() {
				delete(impl.PersistentTransactions, byTrxId)
			}
		}

		origPendingTxnSize := len(impl.PendingIncomingTransactions)

		if len(impl.PersistentTransactions) > 0 || impl.PendingBlockMode == EnumPendingBlockMode(producing) {
			unappliedTrxs := chain.GetUnappliedTransactions()

			if len(impl.PersistentTransactions) > 0 {
				for i, trx := range unappliedTrxs {
					if _, has := impl.PersistentTransactions[trx.ID]; has {
						// this is a persisted transaction, push it into the block (even if we are speculating) with
						// no deadline as it has already passed the subjective deadlines once and we want to represent
						// the state of the chain including this transaction
						trace := chain.PushTransaction(trx, common.MaxTimePoint(), 0)
						if trace != nil {
							return EnumStartBlockRusult(failed), lastBlock
						}
					}

					// remove it from further consideration as it is applied
					unappliedTrxs[i] = nil
				}
			}

			if impl.PendingBlockMode == EnumPendingBlockMode(producing) {
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

					deadline := common.Now().AddUs(common.Microseconds(impl.MaxTransactionTimeMs))
					deadlineIsSubjective := false
					if impl.MaxTransactionTimeMs < 0 || impl.PendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
						deadlineIsSubjective = true
						deadline = blockTime
					}

					trace := chain.PushTransaction(trx, deadline, 0)
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

		if impl.PendingBlockMode == EnumPendingBlockMode(producing) {
			for byTrxId, byExpire := range impl.BlacklistedTransactions {
				if byExpire <= common.Now() {
					delete(impl.BlacklistedTransactions, byTrxId)
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
				for impl.IncomingTrxWeight >= 1.0 && origPendingTxnSize > 0 && len(impl.PendingIncomingTransactions) > 0 {
					e := impl.PendingIncomingTransactions[0]
					impl.PendingIncomingTransactions = impl.PendingIncomingTransactions[1:]
					origPendingTxnSize--
					impl.IncomingTrxWeight -= 1.0
					impl.OnIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
				}

				if blockTime <= common.Now() {
					isExhausted = true
					break
				}

				if _, has := impl.BlacklistedTransactions[trx]; has {
					continue
				}

				deadline := common.Now().AddUs(common.Microseconds(impl.MaxTransactionTimeMs))
				deadlineIsSubjective := false
				if impl.MaxTransactionTimeMs < 0 || impl.PendingBlockMode == EnumPendingBlockMode(producing) && blockTime < deadline {
					deadlineIsSubjective = true
					deadline = blockTime
				}

				trace := chain.PushScheduledTransaction(&trx, deadline, 0)
				if trace.Except != nil {
					if failureIsSubjective(trace.Except, deadlineIsSubjective) {
						isExhausted = true
					} else {
						expiration := common.Now().AddUs(common.Seconds(0) /*TODO chain.get_global_properties().configuration.deferred_trx_expiration_window*/)
						// this failed our configured maximum transaction time, we don't want to replay it add it to a blacklist
						impl.BlacklistedTransactions[trx] = expiration
					}
				}

				//TODO catch exception

				impl.IncomingTrxWeight += impl.IncomingDeferRadio
				if origPendingTxnSize <= 0 {
					impl.IncomingTrxWeight = 0.0
				}
			}
		} ///scheduled transactions

		if isExhausted || blockTime <= common.Now() {
			return EnumStartBlockRusult(exhausted), lastBlock
		} else {
			// attempt to apply any pending incoming transactions
			impl.IncomingTrxWeight = 0.0
			if origPendingTxnSize > 0 && len(impl.PendingIncomingTransactions) > 0 {
				e := impl.PendingIncomingTransactions[0]
				impl.PendingIncomingTransactions = impl.PendingIncomingTransactions[1:]
				origPendingTxnSize--
				impl.OnIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
				if blockTime <= common.Now() {
					return EnumStartBlockRusult(exhausted), lastBlock
				}
			}
			return EnumStartBlockRusult(succeeded), lastBlock
		}
	}
	return EnumStartBlockRusult(failed), lastBlock
}

func (impl *ProducerPluginImpl) ScheduleProductionLoop() {
	chain := impl.Self.chain()
	impl.Timer.Cancel()

	result, lastBlock := impl.StartBlock()

	if result == EnumStartBlockRusult(failed) {
		log.Error("Failed to start a pending block, will try again later")
		impl.Timer.ExpiresFromNow(common.Microseconds(common.DefaultConfig.BlockIntervalUs / 10))

		// we failed to start a block, so try again later?
		impl.timerCorelationId++
		cid := impl.timerCorelationId
		impl.Timer.AsyncWait(func(err error) {
			if impl != nil && err == nil && cid == impl.timerCorelationId {
				impl.ScheduleProductionLoop()
			}
		})

	} else if result == EnumStartBlockRusult(waiting) {
		if impl.Producers.Size() > 0 && !impl.ProductionDisabledByPolicy() {
			log.Debug("Waiting till another block is received and scheduling Speculative/Production Change")
			impl.ScheduleDelayedProductionLoop(types.NewBlockTimeStamp(impl.CalculatePendingBlockTime()))
		} else {
			log.Debug("Waiting till another block is received")
			// nothing to do until more blocks arrive
		}

	} else if impl.PendingBlockMode == EnumPendingBlockMode(producing) {
		// we succeeded but block may be exhausted
		if result == EnumStartBlockRusult(succeeded) {
			// ship this block off no later than its deadline
			EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "producing without pending_block_state, start_block succeeded")
			deadline := chain.PendingBlockTime().TimeSinceEpoch()
			if lastBlock {
				deadline += common.Microseconds(impl.LastBlockTimeOffsetUs)
			} else {
				deadline += common.Microseconds(impl.ProduceTimeOffsetUs)
			}
			impl.Timer.ExpiresAt(deadline)
			log.Debug("Scheduling Block Production on Normal Block #%d for %s", chain.PendingBlockState().BlockNum, deadline)
		} else {
			EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "producing without pending_block_state")
			expectTime := chain.PendingBlockTime().SubUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
			// ship this block off up to 1 block time earlier or immediately
			if common.Now() >= expectTime {
				impl.Timer.ExpiresFromNow(0)
				log.Debug("Scheduling Block Production on Exhausted Block #%d immediately", chain.PendingBlockState().BlockNum)
			} else {
				impl.Timer.ExpiresAt(expectTime.TimeSinceEpoch())
				log.Debug("Scheduling Block Production on Exhausted Block #%d at %s", chain.PendingBlockState().BlockNum, expectTime)
			}
		}

		impl.timerCorelationId++
		cid := impl.timerCorelationId
		impl.Timer.AsyncWait(func(err error) {
			if impl != nil && err == nil && cid == impl.timerCorelationId {
				res := impl.MaybeProduceBlock()
				log.Debug("Producing Block #%d returned: %v", chain.PendingBlockState().BlockNum, res)
			}
		})

	} else if impl.PendingBlockMode == EnumPendingBlockMode(speculating) && impl.Producers.Size() > 0 && !impl.ProductionDisabledByPolicy() {
		log.Debug("Speculative Block Created; Scheduling Speculative/Production Change")
		EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "speculating without pending_block_state")
		pbs := chain.PendingBlockState()
		impl.ScheduleDelayedProductionLoop(pbs.Header.Timestamp)

	} else {
		log.Debug( "Speculative Block Created")
	}
}

func (impl *ProducerPluginImpl) ScheduleDelayedProductionLoop(currentBlockTime types.BlockTimeStamp) {
	// if we have any producers then we should at least set a timer for our next available slot
	var wakeUpTime *common.TimePoint
	impl.Producers.Each(func(index int, value interface{}) {
		p := value.(common.AccountName)
		nextProducerBlockTime := impl.CalculateNextBlockTime(&p, currentBlockTime)
		if nextProducerBlockTime != nil {
			producerWakeupTime := nextProducerBlockTime.SubUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
			if wakeUpTime != nil {
				// wake up with a full block interval to the deadline
				if *wakeUpTime > producerWakeupTime {
					*wakeUpTime = producerWakeupTime
				}
			} else {
				wakeUpTime = &producerWakeupTime
			}
		}
	})

	if wakeUpTime != nil {
		log.Debug("Scheduling Speculative/Production Change at %s", wakeUpTime)
		impl.Timer.ExpiresAt(wakeUpTime.TimeSinceEpoch())

		impl.timerCorelationId++
		cid := impl.timerCorelationId
		impl.Timer.AsyncWait(func(err error) {
			if impl != nil && err == nil && cid == impl.timerCorelationId {
				fmt.Println("===re loop")
				impl.ScheduleProductionLoop()
			}
		})
	} else {
		log.Debug("Speculative Block Created; Not Scheduling Speculative/Production, no local producers had valid wake up times")
	}

}

func (impl *ProducerPluginImpl) MaybeProduceBlock() bool {
	defer func() {
		impl.ScheduleProductionLoop()
	}()

	returning, r := false, false
	Try(func() {
		impl.ProduceBlock()
		returning, r = true, true
	}).Catch(func(e GuardExceptions) {
		//TODO: app().get_plugin<chain_plugin>().handle_guard_exception(e);
		returning, r = true, false
	}).FcLogAndDrop().End()

	if returning {
		return r
	}

	log.Debug("Aborting block due to produce_block error")
	chain := impl.Self.chain()
	chain.AbortBlock()
	return false
}

func (impl *ProducerPluginImpl) ProduceBlock() {
	EosAssert(impl.PendingBlockMode == EnumPendingBlockMode(producing), &ProducerException{}, "called produce_block while not actually producing")
	chain := impl.Self.chain()
	pbs := chain.PendingBlockState()
	EosAssert(pbs != nil, &MissingPendingBlockState{}, "pending_block_state does not exist but it should, another plugin may have corrupted it")

	signatureProvider := impl.SignatureProviders[pbs.BlockSigningKey]
	EosAssert(signatureProvider != nil, &ProducerPrivKeyNotFound{}, "Attempting to produce a block for which we don't have the private key")

	chain.FinalizeBlock()
	chain.SignBlock(func(d crypto.Sha256) ecc.Signature {
		defer makeDebugTimeLogger()
		return signatureProvider(d)
	})

	chain.CommitBlock(true)

	newBs := chain.HeadBlockState()
	impl.ProducerWatermarks[newBs.Header.Producer] = chain.HeadBlockNum()

	log.Info("Produced block %s...#%d @ %s signed by %s [trxs: %d, lib: %d, confirmed: %d]\n",
		newBs.BlockId.String()[0:16], newBs.BlockNum, newBs.Header.Timestamp, common.S(uint64(newBs.Header.Producer)),
		len(newBs.SignedBlock.Transactions), chain.LastIrreversibleBlockNum(), newBs.Header.Confirmed)
}
