package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
	"sort"
)

type BlockState struct {
	BlockHeaderState
	SignedBlock    SignedBlock
	Validated      bool `json:"validated"`
	InCurrentChain bool `json:"in_current_chain"`
	Trxs           []TransactionMetadata
}

type BlockHeaderState struct {
	ID                               common.BlockIdType `storm:"id,unique"`
	BlockNum                         uint32             `storm:"block_num,unique"`
	Header                           SignedBlockHeader
	DposProposedIrreversibleBlocknum uint32     `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         uint32     `json:"dpos_irreversible_blocknum"`
	BftIrreversibleBlocknum          uint32     `json:"bft_irreversible_blocknum"`
	PendingScheduleLibNum            uint32     `json:"pending_schedule_lib_num"`
	PendingScheduleHash              rlp.Sha256 `json:"pending_schedule_hash"`
	PendingSchedule                  ProducerScheduleType
	ActiveSchedule                   ProducerScheduleType
	BlockrootMerkle                  IncrementalMerkle
	ProducerToLastProduced           map[common.AccountName]uint32
	ProducerToLastImpliedIrb         map[common.AccountName]uint32
	BlockSigningKey                  ecc.PublicKey
	ConfirmCount                     []uint8              `json:"confirm_count"`
	Confirmations                    []HeaderConfirmation `json:"confirmations"`
}

func (bs *BlockHeaderState) GetScheduledProducer(t common.BlockTimeStamp) ProducerKey {
	index := uint32(t) % uint32(len(bs.ActiveSchedule.Producers)*12)
	index /= 12
	return bs.ActiveSchedule.Producers[index]
}

func (bs *BlockHeaderState) CalcDposLastIrreversible() uint32 {
	blockNums := make([]int, 0, len(bs.ProducerToLastImpliedIrb))
	for _, value := range bs.ProducerToLastImpliedIrb {
		blockNums = append(blockNums, int(value))
	}
	/// 2/3 must be greater, so if I go 1/3 into the list sorted from low to high, then 2/3 are greater

	if len(blockNums) == 0 {
		return 0
	}
	/// TODO: update to nth_element
	sort.Ints(blockNums)
	return uint32(blockNums[(len(blockNums)-1)/3])
}

func (bs *BlockHeaderState) GenerateNext(when *common.BlockTimeStamp) *BlockHeaderState {
	result := new(BlockHeaderState)

	if when != nil {
		if *when <= bs.Header.Timestamp {
			panic("next block must be in the future") //block_validate_exception
		}
	} else {
		when = &bs.Header.Timestamp
		*when++
	}

	result.Header.Timestamp = *when
	result.Header.Previous = bs.ID
	result.Header.ScheduleVersion = bs.ActiveSchedule.Version

	proKey := bs.GetScheduledProducer(*when)
	result.BlockSigningKey = proKey.BlockSigningKey
	result.Header.Producer = proKey.AccountName

	result.PendingScheduleLibNum = bs.PendingScheduleLibNum
	result.PendingScheduleHash = bs.PendingScheduleHash
	result.BlockNum = bs.BlockNum + 1
	result.ProducerToLastProduced = bs.ProducerToLastProduced
	result.ProducerToLastImpliedIrb = bs.ProducerToLastImpliedIrb
	result.ProducerToLastProduced[proKey.AccountName] = result.BlockNum
	result.BlockrootMerkle = bs.BlockrootMerkle
	result.BlockrootMerkle.Append(rlp.Sha256(bs.ID))

	result.ActiveSchedule = bs.ActiveSchedule
	result.PendingSchedule = bs.PendingSchedule
	result.DposProposedIrreversibleBlocknum = bs.DposProposedIrreversibleBlocknum
	result.BftIrreversibleBlocknum = bs.BftIrreversibleBlocknum

	result.ProducerToLastImpliedIrb[proKey.AccountName] = result.DposProposedIrreversibleBlocknum
	result.DposIrreversibleBlocknum = result.CalcDposLastIrreversible()

	/// grow the confirmed count
	if common.DefaultConfig.MaxProducers*2/3+1 > 0xff {
		panic("8bit confirmations may not be able to hold all of the needed confirmations")
	}

	// This uses the previous block active_schedule because thats the "schedule" that signs and therefore confirms _this_ block
	numActiveProducers := len(bs.ActiveSchedule.Producers)
	requiredConfs := uint32(numActiveProducers*2/3) + 1

	if len(bs.ConfirmCount) < common.DefaultConfig.MaxTrackedDposConfirmations {
		result.ConfirmCount = make([]uint8, len(bs.ConfirmCount)+1)
		copy(result.ConfirmCount, bs.ConfirmCount)
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	} else {
		result.ConfirmCount = make([]uint8, len(bs.ConfirmCount))
		copy(result.ConfirmCount, bs.ConfirmCount[1:])
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	}

	return result
}

func (bs *BlockHeaderState) MaybePromotePending() bool {
	if len(bs.PendingSchedule.Producers) > 0 && bs.DposIrreversibleBlocknum >= bs.PendingScheduleLibNum {
		bs.ActiveSchedule = bs.PendingSchedule

		var newProducerToLastProduced = make(map[common.AccountName]uint32)
		var newProducerToLastImpliedIrb = make(map[common.AccountName]uint32)
		for _, pro := range bs.ActiveSchedule.Producers {
			existing, hasExisting := bs.ProducerToLastProduced[pro.AccountName]
			if hasExisting {
				newProducerToLastProduced[pro.AccountName] = existing
			} else {
				newProducerToLastProduced[pro.AccountName] = bs.DposIrreversibleBlocknum
			}

			existingIrb, hasExistingIrb := bs.ProducerToLastImpliedIrb[pro.AccountName]
			if hasExistingIrb {
				newProducerToLastImpliedIrb[pro.AccountName] = existingIrb
			} else {
				newProducerToLastImpliedIrb[pro.AccountName] = bs.DposIrreversibleBlocknum
			}
		}

		bs.ProducerToLastProduced = newProducerToLastProduced
		bs.ProducerToLastImpliedIrb = newProducerToLastImpliedIrb
		bs.ProducerToLastProduced[bs.Header.Producer] = bs.BlockNum

		return true
	}
	return false
}

func (bs *BlockHeaderState) SetConfirmed(numPrevBlocks uint16) {
	bs.Header.Confirmed = numPrevBlocks

	i := len(bs.ConfirmCount) - 1
	blocksToConfirm := numPrevBlocks + 1 /// confirm the head block too
	for i >= 0 && blocksToConfirm > 0 {
		bs.ConfirmCount[i]--
		if bs.ConfirmCount[i] == 0 {
			blockNumFori := bs.BlockNum - uint32(len(bs.ConfirmCount)-1-i)
			bs.DposProposedIrreversibleBlocknum = blockNumFori

			if i == len(bs.ConfirmCount)-1 {
				bs.ConfirmCount = make([]uint8, 0)
			} else {
				bs.ConfirmCount = bs.ConfirmCount[i+1:]
			}

			return
		}
		i--
		blocksToConfirm--
	}

}

func (bs *BlockHeaderState) SigDigest() rlp.Sha256 {
	headerBmroot := rlp.Hash256(common.Pair{bs.Header, bs.BlockrootMerkle.GetRoot()})
	digest := rlp.Hash256(common.Pair{headerBmroot, bs.PendingScheduleHash})
	return digest
}

func (bs *BlockHeaderState) SetNewProducers(pending SharedProducerScheduleType) {
	if pending.Version != bs.ActiveSchedule.Version+1 {
		panic("wrong producer schedule version specified")
	}
	if len(pending.Producers) != 0 {
		panic("cannot set new pending producers until last pending is confirmed")
	}
	bs.Header.NewProducers = &pending
	bs.PendingScheduleHash = rlp.Hash256(*bs.Header.NewProducers)
	bs.PendingSchedule = *bs.Header.NewProducers.ProducerScheduleType()
	bs.PendingScheduleLibNum = bs.BlockNum
}

/**
 *  Transitions the current header state into the next header state given the supplied signed block header.
 *
 *  Given a signed block header, generate the expected template based upon the header time,
 *  then validate that the provided header matches the template.
 *
 *  If the header specifies new_producers then apply them accordingly.
 */
func (bs *BlockHeaderState) Next(h SignedBlockHeader, trust bool) *BlockHeaderState {
	if h.Timestamp == common.BlockTimeStamp(0) {
		panic(fmt.Sprintf("h:%s", h))
	}
	if len(h.HeaderExtensions) != 0 {
		panic("no supported extensions")
	}
	if h.Timestamp <= bs.Header.Timestamp {
		panic("block must be later in time")
	}
	if h.Previous != bs.ID {
		panic("block must link to current state")
	}

	result := bs.GenerateNext(&h.Timestamp)

	if result.Header.Producer != h.Producer {
		panic("wrong producer specified")
	}
	if result.Header.ScheduleVersion != h.ScheduleVersion {
		panic("wrong producer specified")
	}

	itr, has := bs.ProducerToLastProduced[h.Producer]
	if has {
		if itr >= result.BlockNum-uint32(h.Confirmed) {
			panic(fmt.Sprintf("producer %s double-confirming known range", h.Producer))
		}
	}

	/// below this point is state changes that cannot be validated with headers alone, but never-the-less,
	/// must result in header state changes

	result.SetConfirmed(h.Confirmed)

	wasPendingPromoted := result.MaybePromotePending()

	if h.NewProducers != nil {
		if wasPendingPromoted {
			panic("cannot set pending producer schedule in the same block in which pending was promoted to active")
		}
		result.SetNewProducers(*h.NewProducers)
	}

	result.Header.ActionMRoot = h.ActionMRoot
	result.Header.TransactionMRoot = h.TransactionMRoot
	result.Header.ProducerSignature = h.ProducerSignature
	result.ID = result.Header.BlockID()

	if !trust {
		signKey, err := result.Signee()
		if err != nil {
			panic(err)
		}
		if result.BlockSigningKey != signKey {
			panic(fmt.Sprintf("block not signed by expected key, block_signing_key:%s, signee:%s",
				result.BlockSigningKey, signKey))
		}
	}
	return result
}

func (bs *BlockHeaderState) Sign(signer func(sha256 rlp.Sha256) ecc.Signature) {
	d := bs.SigDigest()
	bs.Header.ProducerSignature = signer(d)
	signKey, err := bs.Header.ProducerSignature.PublicKey(d.Bytes())
	if err != nil {
		panic(err)
	}
	if bs.BlockSigningKey != signKey {
		panic("block is signed with unexpected key")
	}
}

func (bs *BlockHeaderState) Signee() (ecc.PublicKey, error) {
	return bs.Header.ProducerSignature.PublicKey(bs.SigDigest().Bytes())
}

func (bs *BlockHeaderState) AddConfirmation(conf HeaderConfirmation) {
	for _, c := range bs.Confirmations {
		if c.Producer == conf.Producer {
			panic("block already confirmed by this producer")
		}
	}

	key, hasKey := bs.ActiveSchedule.GetProducerKey(conf.Producer)
	if !hasKey {
		panic("producer not in current schedule")
	}

	signer, err := conf.ProducerSignature.PublicKey(bs.SigDigest().Bytes())
	if err != nil {
		panic(err)
	}
	if signer != key {
		panic("confirmation not signed by expected key")
	}

	bs.Confirmations = append(bs.Confirmations, conf)
}
