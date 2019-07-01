package producer_plugin

import (
	Chain "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	. "github.com/eosspark/eos-go/plugins/producer_plugin/multi_index"
)

func (impl *ProducerPluginImpl) CalculateNextBlockTime(producerName common.AccountName, currentBlockTime types.BlockTimeStamp) *common.TimePoint {
	var result common.TimePoint

	chain := impl.Chain

	hbs := chain.HeadBlockState()
	activeSchedule := hbs.ActiveSchedule.Producers

	// determine if this producer is in the active schedule and if so, where
	var itr *types.ProducerKey
	var producerIndex uint32
	for index, asp := range activeSchedule {
		if asp.ProducerName == producerName {
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
	currentWatermark, hasCurrentWatermark := impl.ProducerWatermarks[producerName]
	if hasCurrentWatermark {
		blockNum := chain.HeadBlockNum()
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
	chain := impl.Chain
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

func (impl *ProducerPluginImpl) StartBlock() (StartBlockResult, bool) {
	chain := impl.Chain

	if chain.GetReadMode() == Chain.READONLY {
		return StartBlockResult(waiting), false
	}

	hbs := chain.HeadBlockState()

	//Schedule for the next second's tick regardless of chain state
	// If we would wait less than 50ms (1/10 of block_interval), wait for the whole block interval.
	now := common.Now()
	blockTime := impl.CalculatePendingBlockTime()

	impl.PendingBlockMode = PendingBlockMode(producing)

	// Not our turn
	lastBlock := uint32(types.NewBlockTimeStamp(blockTime))%uint32(common.DefaultConfig.ProducerRepetitions) == uint32(common.DefaultConfig.ProducerRepetitions)-1
	scheduleProducer := hbs.GetScheduledProducer(types.NewBlockTimeStamp(blockTime))
	currentWatermark, hasCurrentWatermark := impl.ProducerWatermarks[scheduleProducer.ProducerName]
	_, hasSignatureProvider := impl.SignatureProviders[scheduleProducer.BlockSigningKey]
	irreversibleBlockAge := impl.GetIrreversibleBlockAge()

	// If the next block production opportunity is in the present or future, we're synced.
	if !impl.ProductionEnabled {
		impl.PendingBlockMode = PendingBlockMode(speculating)

	} else if !impl.Producers.Contains(scheduleProducer.ProducerName) {
		impl.PendingBlockMode = PendingBlockMode(speculating)

	} else if !hasSignatureProvider {
		log.Error("Not producing block because I don't have the private key for %s", scheduleProducer.BlockSigningKey)
		impl.PendingBlockMode = PendingBlockMode(speculating)

	} else if impl.ProductionPaused {
		log.Error("Not producing block because production is explicitly paused")
		impl.PendingBlockMode = PendingBlockMode(speculating)

	} else if impl.MaxIrreversibleBlockAgeUs >= 0 && irreversibleBlockAge >= impl.MaxIrreversibleBlockAgeUs {
		log.Error("Not producing block because the irreversible block is too old [age:%ds, max:%ds]", irreversibleBlockAge.Count()/1e6, impl.MaxIrreversibleBlockAgeUs.Count()/1e6)
		impl.PendingBlockMode = PendingBlockMode(speculating)
	}

	if impl.PendingBlockMode == PendingBlockMode(producing) {
		if hasCurrentWatermark {
			if currentWatermark >= hbs.BlockNum+1 {
				log.Error("Not producing block because \"%s\" signed a BFT confirmation OR block at a higher block number (%d) than the current fork's head (%d)",
					scheduleProducer.ProducerName,
					currentWatermark,
					hbs.BlockNum)
				impl.PendingBlockMode = PendingBlockMode(speculating)
			}

		}
	}

	if impl.PendingBlockMode == PendingBlockMode(speculating) {
		headBlockAge := now.Sub(chain.HeadBlockTime())
		if headBlockAge > common.Seconds(5) {
			return StartBlockResult(waiting), lastBlock
		}
	}

	Try(func() {
		blocksToConfirm := uint16(0)

		if impl.PendingBlockMode == PendingBlockMode(producing) {
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

		if impl.PendingBlockMode == PendingBlockMode(producing) && pbs.BlockSigningKey != scheduleProducer.BlockSigningKey {
			log.Error("Block Signing Key is not expected value, reverting to speculative mode! [expected: \"%s\", actual: \"%s\"", scheduleProducer.BlockSigningKey, pbs.BlockSigningKey)
			impl.PendingBlockMode = PendingBlockMode(speculating)
		}

		// attempt to play persisted transactions first
		isExhausted := false

		// remove all persisted transactions that have now expired
		persistedById := impl.PersistentTransactions.GetById()
		persistedByExpire := impl.PersistentTransactions.GetByExpiry()
		if !persistedByExpire.Empty() {
			numExpiredPersistent := 0
			origCount := impl.PersistentTransactions.Size()

			for !persistedByExpire.Empty() && persistedByExpire.Begin().Value().Expiry <= pbs.Header.Timestamp.ToTimePoint() {
				txid := persistedByExpire.Begin().Value().TrxId
				if impl.PendingBlockMode == producing {
					trxTraceLog.Debug("[TRX_TRACE] Block #%d for producer %s is EXPIRING PERSISTED tx: %s",
						chain.HeadBlockNum()+1, chain.PendingBlockState().Header.Producer, txid)

				} else {
					trxTraceLog.Debug("[TRX_TRACE] Speculative execution is EXPIRING PERSISTED tx: %s", txid)
				}

				persistedByExpire.Erase(persistedByExpire.Begin())
				numExpiredPersistent++
			}

			ppLog.Debug("Processed %d persisted transactions, Expired %d", origCount, numExpiredPersistent)
		}

		origPendingTxnSize := len(impl.PendingIncomingTransactions)

		// Processing unapplied transactions...
		//
		if impl.Producers.Empty() && persistedById.Empty() {
			// if this node can never produce and has no persisted transactions,
			// there is no need for unapplied transactions they can be dropped
			chain.DropAllUnAppliedTransactions()

		} else {
			var applyTrxs []*types.TransactionMetadata
			{ // derive appliable transactions from unapplied_transactions and drop droppable transactions
				unappliedTrxs := chain.GetUnappliedTransactions()
				applyTrxs = make([]*types.TransactionMetadata, 0, len(unappliedTrxs))

				calculateTransactionCategory := func(trx *types.TransactionMetadata) EnumTxCategory {
					if trx.PackedTrx.Expiration().ToTimePoint() < pbs.Header.Timestamp.ToTimePoint() {
						return EXPIRED
					} else if _, ok := persistedById.Find(trx.ID); ok {
						return PERSISTED
					} else {
						return UNEXPIRED_UNPERSISTED
					}
				}

				for _, trx := range unappliedTrxs {
					category := calculateTransactionCategory(trx)
					if category == EXPIRED || (category == UNEXPIRED_UNPERSISTED && impl.Producers.Empty()) {
						if !impl.Producers.Empty() {
							trxTraceLog.Debug("[TRX_TRACE] Node with producers configured is dropping an EXPIRED transaction that was PREVIOUSLY ACCEPTED : %s", trx.ID)
						}
						chain.DropUnappliedTransaction(trx)
					} else if category == PERSISTED || (category == UNEXPIRED_UNPERSISTED && impl.PendingBlockMode == producing) {
						applyTrxs = append(applyTrxs, trx)
					}
				}
			}

			if len(applyTrxs) > 0 {
				numApplied := 0
				numFailed := 0
				numProcessed := 0

				for _, trx := range applyTrxs {
					if blockTime <= common.Now() {
						isExhausted = true
						break
					}

					numProcessed++

					returning := false
					Try(func() {
						deadline := common.Now().AddUs(common.Milliseconds(int64(impl.MaxTransactionTimeMs)))
						deadlineIsSubjective := false
						if impl.MaxTransactionTimeMs < 0 || (impl.PendingBlockMode == producing && blockTime < deadline) {
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
								numFailed++
							}
						} else {
							numApplied++
						}
					}).Catch(func(e GuardExceptions) {
						app.App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).HandleGuardException(e)
						returning = true

					}).FcLogAndDrop()

					if returning {
						return failed, lastBlock
					}
				}

				ppLog.Debug("Processed %d of %d previously applied transactions, Applied %d, Failed/Dropped %d",
					numProcessed, len(applyTrxs), numApplied, numFailed)

			}
		}

		if impl.PendingBlockMode == PendingBlockMode(producing) {
			blacklistById := impl.BlacklistedTransactions.GetById()
			blacklistByExpiry := impl.BlacklistedTransactions.GetByExpiry()
			now := common.Now()
			if !blacklistByExpiry.Empty() {
				numExpired := 0
				origCount := impl.BlacklistedTransactions.Size()

				for !blacklistByExpiry.Empty() && blacklistByExpiry.Begin().Value().Expiry <= now {
					blacklistByExpiry.Erase(blacklistByExpiry.Begin())
					numExpired++
				}

				ppLog.Debug("Processed %d blacklisted transactions, Expired %d",
					origCount, numExpired)
			}

			scheduleTrx := chain.GetScheduledTransactions()
			if len(scheduleTrx) > 0 {
				numApplied := 0
				numFailed := 0
				numProcessed := 0

				for _, trx := range scheduleTrx {
					if blockTime <= common.Now() {
						isExhausted = true
					}
					if isExhausted {
						break
					}

					numProcessed++

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

					if _, ok := blacklistById.Find(trx); ok {
						continue
					}

					returning := false
					Try(func() {
						deadline := common.Now().AddUs(common.Milliseconds(int64(impl.MaxTransactionTimeMs)))
						deadlineIsSubjective := false
						if impl.MaxTransactionTimeMs < 0 || (impl.PendingBlockMode == producing && blockTime < deadline) {
							deadlineIsSubjective = true
							deadline = blockTime
						}

						trace := chain.PushScheduledTransaction(&trx, deadline, 0)
						if trace.Except != nil {
							if failureIsSubjective(trace.Except, deadlineIsSubjective) {
								isExhausted = true

							} else {
								exporation := common.Now().AddUs(common.Seconds(int64(chain.GetGlobalProperties().Configuration.DeferredTrxExpirationWindow)))
								// this failed our configured maximum transaction time, we don't want to replay it add it to a blacklist
								impl.BlacklistedTransactions.Insert(TransactionIdWithExpiry{TrxId: trx, Expiry: exporation})
								numFailed++
							}
						} else {
							numApplied++
						}
					}).Catch(func(e GuardExceptions) {
						app.App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).HandleGuardException(e)
						returning = true
					}).FcLogAndDrop()

					if returning {
						return failed, lastBlock
					}

					impl.IncomingTrxWeight += impl.IncomingDeferRadio
					if origPendingTxnSize <= 0 {
						impl.IncomingTrxWeight = 0.0
					}
				}

				ppLog.Debug("Processed %d of %d scheduled transactions, Applied %d, Failed/Dropped %d",
					numProcessed, len(scheduleTrx), numApplied, numFailed)
			}

		} ///scheduled transactions

		if isExhausted || blockTime <= common.Now() {
			return StartBlockResult(exhausted), lastBlock
		} else {
			// attempt to apply any pending incoming transactions
			impl.IncomingTrxWeight = 0.0

			if len(impl.PendingIncomingTransactions) > 0 {
				ppLog.Debug("Processing ${n} pending transactions")
				for origPendingTxnSize > 0 && len(impl.PendingIncomingTransactions) > 0 {
					e := impl.PendingIncomingTransactions[0]
					impl.PendingIncomingTransactions = impl.PendingIncomingTransactions[1:]
					origPendingTxnSize--
					impl.OnIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
					if blockTime <= common.Now() {
						return StartBlockResult(exhausted), lastBlock
					}

				}
			}
			return StartBlockResult(succeeded), lastBlock
		}
	}
	return StartBlockResult(failed), lastBlock
}

func (impl *ProducerPluginImpl) ScheduleProductionLoop() {
	chain := impl.Chain
	impl.Timer.Cancel()

	result, lastBlock := impl.StartBlock()

	if result == StartBlockResult(failed) {
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

	} else if result == StartBlockResult(waiting) {
		if impl.Producers.Size() > 0 && !impl.ProductionDisabledByPolicy() {
			ppLog.Debug("Waiting till another block is received and scheduling Speculative/Production Change")
			impl.ScheduleDelayedProductionLoop(types.NewBlockTimeStamp(impl.CalculatePendingBlockTime()))
		} else {
			ppLog.Debug("Waiting till another block is received")
			// nothing to do until more blocks arrive
		}

	} else if impl.PendingBlockMode == PendingBlockMode(producing) {
		// we succeeded but block may be exhausted
		if result == StartBlockResult(succeeded) {
			// ship this block off no later than its deadline
			EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "producing without pending_block_state, start_block succeeded")
			deadline := chain.PendingBlockTime().TimeSinceEpoch()
			if lastBlock {
				deadline += common.Microseconds(impl.LastBlockTimeOffsetUs)
			} else {
				deadline += common.Microseconds(impl.ProduceTimeOffsetUs)
			}
			impl.Timer.ExpiresAt(deadline)
			ppLog.Debug("Scheduling Block Production on Normal Block %d for %v", chain.PendingBlockState().BlockNum, deadline)
		} else {
			EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "producing without pending_block_state")
			expectTime := chain.PendingBlockTime().SubUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
			// ship this block off up to 1 block time earlier or immediately
			if common.Now() >= expectTime {
				impl.Timer.ExpiresFromNow(0)
				ppLog.Debug("Scheduling Block Production on Exhausted Block #%d immediately", chain.PendingBlockState().BlockNum)
			} else {
				impl.Timer.ExpiresAt(expectTime.TimeSinceEpoch())
				ppLog.Debug("Scheduling Block Production on Exhausted Block #%d at %s", chain.PendingBlockState().BlockNum, expectTime)
			}
		}

		impl.timerCorelationId++
		cid := impl.timerCorelationId
		impl.Timer.AsyncWait(func(err error) {
			if impl != nil && err == nil && cid == impl.timerCorelationId {
				blockNum := uint32(0)
				if pending := chain.PendingBlockState(); pending != nil {
					blockNum = pending.BlockNum
				}
				res := impl.MaybeProduceBlock()
				ppLog.Debug("Producing Block #%d returned: %v", blockNum, res)
			}
		})

	} else if impl.PendingBlockMode == PendingBlockMode(speculating) && impl.Producers.Size() > 0 && !impl.ProductionDisabledByPolicy() {
		ppLog.Debug("Speculative Block Created; Scheduling Speculative/Production Change")
		EosAssert(chain.PendingBlockState() != nil, &MissingPendingBlockState{}, "speculating without pending_block_state")
		pbs := chain.PendingBlockState()
		impl.ScheduleDelayedProductionLoop(pbs.Header.Timestamp)

	} else {
		ppLog.Debug("Speculative Block Created")
	}
}

func (impl *ProducerPluginImpl) ScheduleDelayedProductionLoop(currentBlockTime types.BlockTimeStamp) {
	// if we have any producers then we should at least set a timer for our next available slot
	var wakeUpTime *common.TimePoint
	for itr := impl.Producers.Begin(); itr.HasNext(); itr.Next() {
		nextProducerBlockTime := impl.CalculateNextBlockTime(itr.Value(), currentBlockTime)
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
	}

	if wakeUpTime != nil {
		ppLog.Debug("Scheduling Speculative/Production Change at %s", wakeUpTime)
		impl.Timer.ExpiresAt(wakeUpTime.TimeSinceEpoch())

		impl.timerCorelationId++
		cid := impl.timerCorelationId

		impl.Timer.AsyncWait(func(err error) {
			if impl != nil && err == nil && cid == impl.timerCorelationId {

				impl.ScheduleProductionLoop()
			}
		})
	} else {
		ppLog.Debug("Speculative Block Created; Not Scheduling Speculative/Production, no local producers had valid wake up times")
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

	ppLog.Debug("Aborting block due to produce_block error")
	chain := impl.Chain
	chain.AbortBlock()
	return false
}

func (impl *ProducerPluginImpl) ProduceBlock() {
	EosAssert(impl.PendingBlockMode == PendingBlockMode(producing), &ProducerException{}, "called produce_block while not actually producing")
	chain := impl.Chain
	pbs := chain.PendingBlockState()
	EosAssert(pbs != nil, &MissingPendingBlockState{}, "pending_block_state does not exist but it should, another plugin may have corrupted it")

	signatureProvider := impl.SignatureProviders[pbs.BlockSigningKey]
	EosAssert(signatureProvider != nil, &ProducerPrivKeyNotFound{}, "Attempting to produce a block for which we don't have the private key")

	chain.FinalizeBlock()
	chain.SignBlock(func(d crypto.Sha256) ecc.Signature {
		defer makeDebugTimeLogger()
		return *signatureProvider(d)
	})

	chain.CommitBlock(true)

	newBs := chain.HeadBlockState()
	impl.ProducerWatermarks[newBs.Header.Producer] = chain.HeadBlockNum()

	log.Info("Produced block %s...#%d @ %s signed by %s [trxs: %d, lib: %d, confirmed: %d]",
		newBs.BlockId.String()[0:16], newBs.BlockNum, newBs.Header.Timestamp, common.S(uint64(newBs.Header.Producer)),
		len(newBs.SignedBlock.Transactions), chain.LastIrreversibleBlockNum(), newBs.Header.Confirmed)
}
