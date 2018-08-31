package producer_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"gopkg.in/urfave/cli.v1"
	"time"
)

type ProducerPlugin struct {
	timer              *scheduleTimer
	producers          map[common.AccountName]struct{}
	pendingBlockMode   EnumPendingBlockMode
	productionEnabled  bool
	productionPaused   bool
	signatureProviders map[ecc.PublicKey]signatureProviderType
	producerWatermarks map[common.AccountName]uint32

	persistentTransactions  transactionIdWithExpireIndex
	blacklistedTransactions transactionIdWithExpireIndex

	maxTransactionTimeMs      int32
	maxIrreversibleBlockAgeUs int32
	produceTImeOffsetUs       int32
	lastBlockTimeOffsetUs     int32
	irreversibleBlockTime     time.Time
	keosdProviderTimeoutUs    int32

	lastSignedBlockTime time.Time
	startTime           time.Time
	lastSignedBlockNum  uint32

	confirmedBlock func(signature ecc.Signature)

	pendingIncomingTransactions []tuple

	// keep a expected ratio between defer txn and incoming txn
	incomingTrxWeight  float64
	incomingDeferRadio float64 // 1:1
}

func (pp *ProducerPlugin) init() {
	pp.producers = make(map[common.AccountName]struct{})
	pp.signatureProviders = make(map[ecc.PublicKey]signatureProviderType)
	pp.producerWatermarks = make(map[common.AccountName]uint32)
	pp.persistentTransactions = make(transactionIdWithExpireIndex)
	pp.blacklistedTransactions = make(transactionIdWithExpireIndex)
	pp.timer = new(scheduleTimer)

	pp.maxIrreversibleBlockAgeUs = -1

	pp.incomingTrxWeight = 0.0
	pp.incomingDeferRadio = 1.0
}

func (pp *ProducerPlugin) IsProducerKey(key ecc.PublicKey) bool {
	privateKey := pp.signatureProviders[key]
	if privateKey != nil {
		return true
	}
	return false
}

func (pp *ProducerPlugin) SignCompact(key *ecc.PublicKey, digest common.SHA256Bytes) ecc.Signature {
	if key != nil {
		privateKeyFunc := pp.signatureProviders[*key]
		if privateKeyFunc == nil {
			panic(ErrProducerPriKeyNotFound)
		}

		return privateKeyFunc(digest)
	}
	return ecc.Signature{}
}

func (pp *ProducerPlugin) Initialize(app *cli.App) {
	pp.init()

	pp.signatureProviders[initPubKey] = func(hash []byte) ecc.Signature {
		sig, _ := initPriKey.Sign(hash)
		return sig
	}

	var producers cli.StringSlice

	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:        "enable-stale-production, e",
			Usage:       "Enable block production, even if the chain is stale.",
			Destination: &pp.productionEnabled,
		},
		cli.StringSliceFlag{
			Name:  "producer-name, p",
			Usage: "ID of producer controlled by this node(e.g. inita; may specify multiple times)",
			Value: &producers,
		},
	}

	app.Action = func(c *cli.Context) {
		if len(producers) > 0 {
			for _, p := range producers {
				pp.producers[common.AccountName(common.StringToName(p))] = struct{}{}
			}
		}
		fmt.Println(pp.producers)
	}
}

func (pp *ProducerPlugin) Startup() {
	pp.scheduleProductionLoop()
}

func (pp *ProducerPlugin) Shutdown() {
	pp.timer.cancel()
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

func (pp *ProducerPlugin) Paused() bool {
	return pp.productionPaused
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

func (pp *ProducerPlugin) onIncomingBlock(block *types.SignedBlock) {

	if !block.Timestamp.ToTimePoint().Before(time.Now().Add(7 * time.Second)) {
		panic(ErrBlockFromTheFuture)
	}
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

	if !chain.HeadBlockState().Header.Timestamp.Next().ToTimePoint().Before(time.Now()) {
		pp.productionEnabled = true
	}

	//C++ log per 1000 blocks

}

func (pp *ProducerPlugin) onIncomingTransactionAsync(trx *types.PackedTransaction, persistUntilExpired bool, next func(respVariant)) {
	if chain.PendingBlockState() == nil {
		pp.pendingIncomingTransactions = append(pp.pendingIncomingTransactions, tuple{trx, persistUntilExpired, next})
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

	fmt.Println(blockTime, &sendResponse, id)

}

func (pp *ProducerPlugin) getIrreversibleBlockAge() int32 /*Microsecond*/ {
	now := time.Now()
	if now.Before(pp.irreversibleBlockTime) {
		return 0
	} else {
		return int32((now.UnixNano() - pp.irreversibleBlockTime.UnixNano()) / 1e3)
	}
}

func (pp *ProducerPlugin) productionDisabledByPolicy() bool {
	return !pp.productionEnabled || pp.productionPaused || (pp.maxIrreversibleBlockAgeUs >= 0 && pp.getIrreversibleBlockAge() >= pp.maxIrreversibleBlockAgeUs)
}

func (pp *ProducerPlugin) calculateNextBlockTime(producerName common.AccountName, currentBlockTime common.BlockTimeStamp) *time.Time {
	var result time.Time

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

func (pp *ProducerPlugin) calculatePendingBlockTime() time.Time {
	now := time.Now()
	var base time.Time
	if now.After(chain.HeadBlockTime()) {
		base = now
	} else {
		base = chain.HeadBlockTime()
	}
	minTimeToNextBlock := int64(common.DefaultConfig.BlockIntervalUs) - base.UnixNano()/1e3%int64(common.DefaultConfig.BlockIntervalUs)
	blockTime := base.Add(time.Microsecond * time.Duration(minTimeToNextBlock))

	if blockTime.Sub(now) < time.Microsecond*time.Duration(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime.Add(time.Microsecond * time.Duration(common.DefaultConfig.BlockIntervalUs))
	}

	return blockTime
}

func (pp *ProducerPlugin) startBlock() (EnumStartBlockRusult, bool) {
	hbs := chain.HeadBlockState()

	now := time.Now()
	blockTime := pp.calculatePendingBlockTime()

	pp.pendingBlockMode = producing

	// Not our turn
	lastBlock := uint32(common.NewBlockTimeStamp(blockTime))%uint32(common.DefaultConfig.ProducerRepetitions) == uint32(common.DefaultConfig.ProducerRepetitions)-1
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
			return waiting, lastBlock
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

	chain.AbortBlock()
	chain.StartBlock(common.NewBlockTimeStamp(blockTime), blocksToConfirm)

	pbs := chain.PendingBlockState()

	if pbs != nil {

		if pp.pendingBlockMode == producing && pbs.BlockSigningKey != scheduleProducer.BlockSigningKey {
			//C++ elog("Block Signing Key is not expected value, reverting to speculative mode! [expected: \"${expected}\", actual: \"${actual\"", ("expected", scheduled_producer.block_signing_key)("actual", pbs->block_signing_key));
			pp.pendingBlockMode = speculating
		}

		// attempt to play persisted transactions first
		isExhausted := false

		// remove all persisted transactions that have now expired
		for byTrxId, byExpire := range pp.persistentTransactions {
			if !byExpire.After(pbs.Header.Timestamp.ToTimePoint()) {
				delete(pp.persistentTransactions, byTrxId)
			}
		}

		origPendingTxnSize := len(pp.pendingIncomingTransactions)

		if len(pp.persistentTransactions) > 0 || pp.pendingBlockMode == producing {
			unappliedTrxs := chain.GetUnappliedTransactions()

			if len(pp.persistentTransactions) > 0 {
				for i, trx := range unappliedTrxs {
					if _, has := pp.persistentTransactions[trx.ID]; has {
						// this is a persisted transaction, push it into the block (even if we are speculating) with
						// no deadline as it has already passed the subjective deadlines once and we want to represent
						// the state of the chain including this transaction
						err := chain.PushTransaction(trx, time.Unix(1e15, 0))
						if err != nil {
							return failed, lastBlock
						}
					}

					// remove it from further consideration as it is applied
					unappliedTrxs[i] = nil
				}
			}

			if pp.pendingBlockMode == producing {
				for _, trx := range unappliedTrxs {
					if !blockTime.After(time.Now()) {
						isExhausted = true
					}
					if isExhausted {
						break
					}

					if trx == nil {
						// nulled in the loop above, skip it
						continue
					}

					//TODO
					/*C++
					if (trx->packed_trx.expiration() < pbs->header.timestamp.to_time_point()) {
						// expired, drop it
						chain.drop_unapplied_transaction(trx);
						continue;
					}*/

					deadline := time.Now().Add(time.Millisecond * time.Duration(pp.maxTransactionTimeMs))
					deadlineIsSubjective := false
					fmt.Println(deadlineIsSubjective)
					if pp.maxTransactionTimeMs < 0 || pp.pendingBlockMode == producing && blockTime.Before(deadline) {
						deadlineIsSubjective = true
						deadline = blockTime
					}

					chain.PushTransaction(trx, deadline)
					//TODO
					//C++
					/*if (trace->except) {
						if (failure_is_subjective(*trace->except, deadline_is_subjective)) {
							isExhausted = true;
						} else {
							// this failed our configured maximum transaction time, we don't want to replay it
							chain.drop_unapplied_transaction(trx);
						}
					}*/
				}
			}
		} ///unapplied transactions

		if pp.pendingBlockMode == producing {
			for byTrxId, byExpire := range pp.blacklistedTransactions {
				if !byExpire.After(time.Now()) {
					delete(pp.blacklistedTransactions, byTrxId)
				}
			}

			scheduledTrxs := chain.GetScheduledTransactions()

			for _, trx := range scheduledTrxs {
				if !blockTime.After(time.Now()) {
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

				if !blockTime.After(time.Now()) {
					isExhausted = true
					break
				}

				if _, has := pp.blacklistedTransactions[trx]; has {
					continue
				}

				deadline := time.Now().Add(time.Millisecond * time.Duration(pp.maxTransactionTimeMs))
				deadlineIsSubjective := false
				fmt.Println(deadlineIsSubjective)
				if pp.maxTransactionTimeMs < 0 || pp.pendingBlockMode == producing && blockTime.Before(deadline) {
					deadlineIsSubjective = true
					deadline = blockTime
				}

				chain.PushScheduledTransaction(trx, deadline)
				//TODO
				//C++
				/*if (trace->except) {
					if (failure_is_subjective(*trace->except, deadline_is_subjective)) {
						isExhausted = true;
					} else {
						auto expiration = fc::time_point::now() + fc::seconds(chain.get_global_properties().configuration.deferred_trx_expiration_window);
						// this failed our configured maximum transaction time, we don't want to replay it add it to a blacklist
						_blacklisted_transactions.insert(transaction_id_with_expiry{trx, expiration});
					}
				}*/
				pp.incomingTrxWeight += pp.incomingDeferRadio
				if origPendingTxnSize <= 0 {
					pp.incomingTrxWeight = 0.0
				}
			}
		} ///scheduled transactions

		if isExhausted || !blockTime.After(time.Now()) {
			return exhausted, lastBlock
		} else {
			// attempt to apply any pending incoming transactions
			pp.incomingTrxWeight = 0.0
			if origPendingTxnSize > 0 && len(pp.pendingIncomingTransactions) > 0 {
				e := pp.pendingIncomingTransactions[0]
				pp.pendingIncomingTransactions = pp.pendingIncomingTransactions[1:]
				origPendingTxnSize--
				pp.onIncomingTransactionAsync(e.packedTransaction, e.persistUntilExpired, e.next)
				if !blockTime.After(time.Now()) {
					return exhausted, lastBlock
				}
			}
			return succeeded, lastBlock
		}
	}
	return failed, lastBlock
}

func (pp *ProducerPlugin) scheduleProductionLoop() {
	pp.timer.cancel()

	result, lastBlock := pp.startBlock()

	if result == failed {
		//elog("Failed to start a pending block, will try again later");
		pp.timer.expiresFromNow(time.Microsecond * (time.Duration(common.DefaultConfig.BlockIntervalUs / 10)))

		// we failed to start a block, so try again later?
		timerCorelationId++
		cid := timerCorelationId
		pp.timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
			pp.scheduleProductionLoop()
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
		if result == succeeded {
			// ship this block off no later than its deadline
			epoch := chain.PendingBlockTime().UnixNano() / 1e3
			if lastBlock {
				epoch += int64(pp.lastBlockTimeOffsetUs)
			} else {
				epoch += int64(pp.produceTImeOffsetUs)
			}
			pp.timer.expiresAt(epoch)
			//fc_dlog(_log, "Scheduling Block Production on Normal Block #${num} for ${time}", ("num", chain.pending_block_state()->block_num)("time",chain.pending_block_time()));
		} else {
			expectTime := chain.PendingBlockTime().UnixNano()/1e3 - int64(common.DefaultConfig.BlockIntervalUs)
			// ship this block off up to 1 block time earlier or immediately
			if time.Now().UnixNano()/1e3 >= expectTime {
				pp.timer.expiresFromNow(0)
			} else {
				pp.timer.expiresAt(expectTime)
			}
			//fc_dlog(_log, "Scheduling Block Production on Exhausted Block #${num} immediately", ("num", chain.pending_block_state()->block_num));
		}

		timerCorelationId++
		cid := timerCorelationId
		pp.timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
			pp.maybeProduceBlock()
			//fc_dlog(_log, "Producing Block #${num} returned: ${res}", ("num", chain.pending_block_state()->block_num)("res", res) );
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
			producerWakeupTime := time.Unix(0, nextProducerBlockTime.UnixNano()-int64(time.Microsecond*time.Duration(common.DefaultConfig.BlockIntervalUs)))
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
		pp.timer.expiresAt(wakeUpTime.UnixNano() / 1e3)

		timerCorelationId++
		cid := timerCorelationId
		pp.timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
			pp.scheduleProductionLoop()
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

func (pp *ProducerPlugin) produceBlock() error {
	if pp.pendingBlockMode != producing {
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
	chain.SignBlock(func(d []byte) ecc.Signature {
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
	//	newBs.Id, newBs.BlockNum, newBs.Header.Timestamp, newBs.Header.Producer, newBs)

	/*ilog("Produced block ${id}... #${n} @ ${t} signed by ${p} [trxs: ${count}, lib: ${lib}, confirmed: ${confs}]",
	  ("p",new_bs->header.producer)("id",fc::variant(new_bs->id).as_string().substr(0,16))
	  ("n",new_bs->block_num)("t",new_bs->header.timestamp)
	  ("count",new_bs->block->transactions.size())("lib",chain.last_irreversible_block_num())("confs", new_bs->header.confirmed));*/
	return nil
}
