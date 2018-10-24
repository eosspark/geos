package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	"sort"
)

type BlockHeaderState struct {
	ID                               common.IdType      `multiIndex:"id,increment"`
	BlockId                          common.BlockIdType `multiIndex:"byId,orderedUnique"`
	BlockNum                         uint32             `multiIndex:"block_num,orderedUnique:byLibBlockNum,orderedNonUnique"`
	Header                           SignedBlockHeader  `multiIndex:"inline"`
	DposProposedIrreversibleBlocknum uint32             `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         uint32             `multiIndex:"byLibBlockNum,orderedNonUnique" json:"dpos_irreversible_blocknum"`
	BftIrreversibleBlocknum          uint32             `multiIndex:"byLibBlockNum,orderedNonUnique" json:"bft_irreversible_blocknum"`
	PendingScheduleLibNum            uint32             `json:"pending_schedule_lib_num"`
	PendingScheduleHash              crypto.Sha256      `json:"pending_schedule_hash"`
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

func (bs *BlockHeaderState) GenerateNext(when common.BlockTimeStamp) *BlockHeaderState {
	result := new(BlockHeaderState)

	if when != common.BlockTimeStamp(0) {
		EosAssert(when > bs.Header.Timestamp, &BlockValidateException{}, "next block must be in the future")
	} else {
		when = bs.Header.Timestamp //TODO: check
		when++
	}

	result.Header.Timestamp = when
	result.Header.Previous = bs.BlockId
	result.Header.ScheduleVersion = bs.ActiveSchedule.Version

	proKey := bs.GetScheduledProducer(when)
	result.BlockSigningKey = proKey.BlockSigningKey
	result.Header.Producer = proKey.AccountName

	result.PendingScheduleLibNum = bs.PendingScheduleLibNum
	result.PendingScheduleHash = bs.PendingScheduleHash
	result.BlockNum = bs.BlockNum + 1

	result.ProducerToLastProduced = make(map[common.AccountName]uint32, len(bs.ProducerToLastProduced))
	for k, v := range bs.ProducerToLastProduced {
		result.ProducerToLastProduced[k] = v
	}

	result.ProducerToLastImpliedIrb = make(map[common.AccountName]uint32, len(bs.ProducerToLastImpliedIrb))
	for k, v := range bs.ProducerToLastImpliedIrb {
		result.ProducerToLastImpliedIrb[k] = v
	}

	result.ProducerToLastProduced[proKey.AccountName] = result.BlockNum
	result.BlockrootMerkle = bs.BlockrootMerkle
	result.BlockrootMerkle.Append(crypto.Sha256(bs.BlockId))

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
} /// generate_next

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

func (bs *BlockHeaderState) SetNewProducers(pending SharedProducerScheduleType) {
	EosAssert(pending.Version == bs.ActiveSchedule.Version+1, &ProducerScheduleException{}, "wrong producer schedule version specified")
	EosAssert(len(bs.PendingSchedule.Producers) == 0, &ProducerScheduleException{},
		"cannot set new pending producers until last pending is confirmed")
	bs.Header.NewProducers = &pending
	bs.PendingScheduleHash = crypto.Hash256(*bs.Header.NewProducers)
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
	EosAssert(h.Timestamp != common.BlockTimeStamp(0), &BlockValidateException{}, "%s", h)
	EosAssert(len(h.HeaderExtensions) == 0, &BlockValidateException{}, "no supported extensions")

	EosAssert(h.Timestamp > bs.Header.Timestamp, &BlockValidateException{}, "block must be later in time")
	EosAssert(h.Previous == bs.BlockId, &UnlinkableBlockException{}, "block must link to current state")

	result := bs.GenerateNext(h.Timestamp)

	EosAssert(result.Header.Producer == h.Producer, &WrongProducer{}, "wrong producer specified")
	EosAssert(result.Header.ScheduleVersion == h.ScheduleVersion, &ProducerScheduleException{}, "schedule_version in signed block is corrupted")

	itr, has := bs.ProducerToLastProduced[h.Producer]
	if has && itr >= result.BlockNum-uint32(h.Confirmed) {
		EosAssert(itr < result.BlockNum-uint32(h.Confirmed), &ProducerDoubleConfirm{}, "producer %s double-confirming known range", h.Producer)
	}

	/// below this point is state changes that cannot be validated with headers alone, but never-the-less,
	/// must result in header state changes

	result.SetConfirmed(h.Confirmed)

	wasPendingPromoted := result.MaybePromotePending()

	if h.NewProducers != nil {
		EosAssert(!wasPendingPromoted, &ProducerScheduleException{}, "cannot set pending producer schedule in the same block in which pending was promoted to active")
		result.SetNewProducers(*h.NewProducers)
	}

	result.Header.ActionMRoot = h.ActionMRoot
	result.Header.TransactionMRoot = h.TransactionMRoot
	result.Header.ProducerSignature = h.ProducerSignature
	result.BlockId = result.Header.BlockID()

	if !trust {
		EosAssert(result.BlockSigningKey == result.Signee(), &WrongSigningKey{}, "block not signed by expected key, "+
			"result.block_signing_key: %s, signee: %s", result.BlockSigningKey, result.Signee())
	}

	return result
} ///next

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

func (bs *BlockHeaderState) SigDigest() crypto.Sha256 {
	headerBmroot := crypto.Hash256(common.MakePair(bs.Header.Digest(), bs.BlockrootMerkle.GetRoot()))
	digest := crypto.Hash256(common.MakePair(headerBmroot, bs.PendingScheduleHash))
	return digest
}

func (bs *BlockHeaderState) Sign(signer func(sha256 crypto.Sha256) ecc.Signature) {
	d := bs.SigDigest()
	bs.Header.ProducerSignature = signer(d)
	signKey, err := bs.Header.ProducerSignature.PublicKey(d.Bytes())
	if err != nil {
		panic(err)
	}
	EosAssert(bs.BlockSigningKey == signKey, &WrongSigningKey{}, "block is signed with unexpected key")
	if bs.BlockSigningKey != signKey {
		panic("block is signed with unexpected key")
	}
}

func (bs *BlockHeaderState) Signee() ecc.PublicKey {
	pk, err := bs.Header.ProducerSignature.PublicKey(bs.SigDigest().Bytes())
	if err != nil {
		panic(err)
	}
	return pk
}

func (bs *BlockHeaderState) AddConfirmation(conf *HeaderConfirmation) {
	for _, c := range bs.Confirmations {
		EosAssert(c.Producer != conf.Producer, &ProducerDoubleConfirm{}, "block already confirmed by this producer")
	}

	key, hasKey := bs.ActiveSchedule.GetProducerKey(conf.Producer)
	EosAssert(hasKey, &ProducerNotInSchedule{}, "producer not in current schedule")
	signer, err := conf.ProducerSignature.PublicKey(bs.SigDigest().Bytes())
	if err != nil {
		panic(err)
	}
	EosAssert(signer == key, &WrongSigningKey{}, "confirmation not signed by expected key")

	bs.Confirmations = append(bs.Confirmations, *conf)
}
