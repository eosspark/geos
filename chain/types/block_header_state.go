package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/container"
	math "github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"sort"
)

func init() {
	Assert(math.MaxUint8 >= common.DefaultConfig.MaxProducers*2/3+1, "8bit confirmations may not be able to hold all of the needed confirmations")
}

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treemap" AccountNameUint32Map(common.AccountName,uint32,common.CompareName,false)
type BlockHeaderState struct {
	ID                               common.IdType        `multiIndex:"id,increment" json:"id"`
	BlockId                          common.BlockIdType   `multiIndex:"byId,orderedUnique" json:"block_id"`
	BlockNum                         uint32               `multiIndex:"block_num,orderedUnique:byLibBlockNum,orderedUnique"`
	Header                           SignedBlockHeader    `multiIndex:"inline" json:"header"`
	DposProposedIrreversibleBlocknum uint32               `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         uint32               `multiIndex:"byLibBlockNum,orderedUnique" json:"dpos_irreversible_blocknum"`
	BftIrreversibleBlocknum          uint32               `multiIndex:"byLibBlockNum,orderedUnique" json:"bft_irreversible_blocknum"`
	PendingScheduleLibNum            uint32               `json:"pending_schedule_lib_num"`
	PendingScheduleHash              crypto.Sha256        `json:"pending_schedule_hash"`
	PendingSchedule                  ProducerScheduleType `json:"pending_schedule"`
	ActiveSchedule                   ProducerScheduleType `json:"active_schedule"`
	BlockrootMerkle                  IncrementalMerkle    `json:"blockroot_merkle"`
	ProducerToLastProduced           AccountNameUint32Map `json:"producer_to_last_produced"`
	ProducerToLastImpliedIrb         AccountNameUint32Map `json:"producer_to_last_implied_irb"`
	BlockSigningKey                  ecc.PublicKey        `json:"block_signing_key"`
	ConfirmCount                     []uint8              `json:"confirm_count"`
	Confirmations                    []HeaderConfirmation `json:"confirmations"`
}

func (b *BlockHeaderState) GetScheduledProducer(t BlockTimeStamp) ProducerKey {
	index := uint32(t) % uint32(len(b.ActiveSchedule.Producers)*common.DefaultConfig.ProducerRepetitions)
	index /= uint32(common.DefaultConfig.ProducerRepetitions)
	return b.ActiveSchedule.Producers[index]
}

func (b *BlockHeaderState) CalcDposLastIrreversible() uint32 {
	blockNums := b.ProducerToLastImpliedIrb.Values()
	/// 2/3 must be greater, so if I go 1/3 into the list sorted from low to high, then 2/3 are greater

	if len(blockNums) == 0 {
		return 0
	}
	/// TODO: update to nth_element
	sort.Sort(container.UInt32Slice(blockNums))
	return uint32(blockNums[(len(blockNums)-1)/3])
}

func (b *BlockHeaderState) GenerateNext(when BlockTimeStamp) *BlockHeaderState {
	result := &BlockHeaderState{}

	if when > 0 {
		EosAssert(when > b.Header.Timestamp, &BlockValidateException{}, "next block must be in the future")
	} else {
		when = b.Header.Timestamp + 1
	}

	result.Header.Timestamp = when
	result.Header.Previous = b.BlockId
	result.Header.ScheduleVersion = b.ActiveSchedule.Version

	proKey := b.GetScheduledProducer(when)
	result.BlockSigningKey = proKey.BlockSigningKey
	result.Header.Producer = proKey.ProducerName

	result.PendingScheduleLibNum = b.PendingScheduleLibNum
	result.PendingScheduleHash = b.PendingScheduleHash
	result.BlockNum = b.BlockNum + 1
	result.ProducerToLastProduced = *CopyFromAccountNameUint32Map(&b.ProducerToLastProduced)
	result.ProducerToLastProduced.Put(proKey.ProducerName, result.BlockNum)
	result.BlockrootMerkle = b.BlockrootMerkle
	result.BlockrootMerkle.Append(b.BlockId)

	result.ActiveSchedule = b.ActiveSchedule
	result.PendingSchedule = b.PendingSchedule
	result.DposProposedIrreversibleBlocknum = b.DposProposedIrreversibleBlocknum
	result.BftIrreversibleBlocknum = b.BftIrreversibleBlocknum

	result.ProducerToLastImpliedIrb = *CopyFromAccountNameUint32Map(&b.ProducerToLastImpliedIrb)
	result.ProducerToLastImpliedIrb.Put(proKey.ProducerName, result.DposProposedIrreversibleBlocknum)
	result.DposIrreversibleBlocknum = result.CalcDposLastIrreversible()

	/// grow the confirmed count
	// static assert math.MaxUint8 >= common.DefaultConfig.MaxProducers*2/3+1

	// This uses the previous block active_schedule because thats the "schedule" that signs and therefore confirms _this_ block
	numActiveProducers := len(b.ActiveSchedule.Producers)
	requiredConfs := uint32(numActiveProducers*2/3) + 1

	if len(b.ConfirmCount) < common.DefaultConfig.MaxTrackedDposConfirmations {
		result.ConfirmCount = make([]uint8, len(b.ConfirmCount)+1)
		copy(result.ConfirmCount, b.ConfirmCount)
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	} else {
		result.ConfirmCount = make([]uint8, len(b.ConfirmCount))
		copy(result.ConfirmCount, b.ConfirmCount[1:])
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	}

	return result
} /// generate_next

func (b *BlockHeaderState) MaybePromotePending() bool {
	if len(b.PendingSchedule.Producers) > 0 && b.DposIrreversibleBlocknum >= b.PendingScheduleLibNum {
		b.ActiveSchedule = b.PendingSchedule
		b.PendingSchedule = ProducerScheduleType{} //active_schedule = move( pending_schedule )
		//var newProducerToLastProduced = make(map[common.AccountName]uint32)
		//var newProducerToLastImpliedIrb = make(map[common.AccountName]uint32)

		newProducerToLastProduced := NewAccountNameUint32Map()

		for _, pro := range b.ActiveSchedule.Producers {
			if blockNum := b.ProducerToLastProduced.Get(pro.ProducerName); !blockNum.IsEnd() {
				newProducerToLastProduced.Put(pro.ProducerName, blockNum.Value())
			} else {
				newProducerToLastProduced.Put(pro.ProducerName, b.DposIrreversibleBlocknum)
			}
		}

		newProducerToLastImpliedIrb := NewAccountNameUint32Map()

		for _, pro := range b.ActiveSchedule.Producers {
			if blockNum := b.ProducerToLastImpliedIrb.Get(pro.ProducerName); !blockNum.IsEnd() {
				newProducerToLastImpliedIrb.Put(pro.ProducerName, blockNum.Value())
			} else {
				newProducerToLastImpliedIrb.Put(pro.ProducerName, b.DposIrreversibleBlocknum)
			}
		}

		b.ProducerToLastProduced = *newProducerToLastProduced
		b.ProducerToLastImpliedIrb = *newProducerToLastImpliedIrb
		b.ProducerToLastProduced.Put(b.Header.Producer, b.BlockNum)

		return true
	}
	return false
}

func (b *BlockHeaderState) SetNewProducers(pending SharedProducerScheduleType) {
	EosAssert(pending.Version == b.ActiveSchedule.Version+1, &ProducerScheduleException{}, "wrong producer schedule version specified")
	EosAssert(len(b.PendingSchedule.Producers) == 0, &ProducerScheduleException{},
		"cannot set new pending producers until last pending is confirmed")
	b.Header.NewProducers = &pending
	b.PendingScheduleHash = *crypto.Hash256(*b.Header.NewProducers)
	b.PendingSchedule = *b.Header.NewProducers.ProducerScheduleType()
	b.PendingScheduleLibNum = b.BlockNum
}

/**
 *  Transitions the current header state into the next header state given the supplied signed block header.
 *
 *  Given a signed block header, generate the expected template based upon the header time,
 *  then validate that the provided header matches the template.
 *
 *  If the header specifies new_producers then apply them accordingly.
 */
func (b *BlockHeaderState) Next(h SignedBlockHeader, trust bool) *BlockHeaderState {
	EosAssert(h.Timestamp != BlockTimeStamp(0), &BlockValidateException{}, "%s", h)
	EosAssert(len(h.HeaderExtensions) == 0, &BlockValidateException{}, "no supported extensions")

	EosAssert(h.Timestamp > b.Header.Timestamp, &BlockValidateException{}, "block must be later in time")
	EosAssert(h.Previous == b.BlockId, &UnlinkableBlockException{}, "block must link to current state")

	result := b.GenerateNext(h.Timestamp)

	EosAssert(result.Header.Producer == h.Producer, &WrongProducer{}, "wrong producer specified")
	EosAssert(result.Header.ScheduleVersion == h.ScheduleVersion, &ProducerScheduleException{}, "schedule_version in signed block is corrupted")

	if blockNum := b.ProducerToLastProduced.Get(h.Producer); !blockNum.IsEnd() {
		EosAssert(blockNum.Value() < result.BlockNum-uint32(h.Confirmed), &ProducerDoubleConfirm{}, "producer %s double-confirming known range", h.Producer)
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
		EosAssert(result.BlockSigningKey.Compare(result.Signee()), &WrongSigningKey{}, "block not signed by expected key, "+
			"result.block_signing_key: %s, signee: %s", result.BlockSigningKey, result.Signee())
	}

	return result
} ///next

func (b *BlockHeaderState) SetConfirmed(numPrevBlocks uint16) {
	b.Header.Confirmed = numPrevBlocks

	i := len(b.ConfirmCount) - 1
	blocksToConfirm := numPrevBlocks + 1 /// confirm the head block too
	for i >= 0 && blocksToConfirm > 0 {
		b.ConfirmCount[i]--
		if b.ConfirmCount[i] == 0 {
			blockNumFori := b.BlockNum - uint32(len(b.ConfirmCount)-1-i)
			b.DposProposedIrreversibleBlocknum = blockNumFori

			if i == len(b.ConfirmCount)-1 {
				b.ConfirmCount = make([]uint8, 0)
			} else {
				b.ConfirmCount = b.ConfirmCount[i+1:]
			}

			return
		}
		i--
		blocksToConfirm--
	}

}

func (b *BlockHeaderState) SigDigest() crypto.Sha256 {
	headerBmroot := crypto.Hash256(common.MakePair(b.Header.Digest(), b.BlockrootMerkle.GetRoot()))
	digest := crypto.Hash256(common.MakePair(headerBmroot, b.PendingScheduleHash))
	return *digest
}

func (b *BlockHeaderState) Sign(signer func(sha256 crypto.Sha256) ecc.Signature) {
	d := b.SigDigest()
	b.Header.ProducerSignature = signer(d)
	signKey, err := b.Header.ProducerSignature.PublicKey(d.Bytes())
	if err != nil {
		Throw(err)
	}
	EosAssert(b.BlockSigningKey.Compare(signKey), &WrongSigningKey{}, "block is signed with unexpected key")

}

func (b *BlockHeaderState) Signee() ecc.PublicKey {
	pk, err := b.Header.ProducerSignature.PublicKey(b.SigDigest().Bytes())
	if err != nil {
		Throw(err)
	}
	return pk
}

func (b *BlockHeaderState) AddConfirmation(conf *HeaderConfirmation) {
	for _, c := range b.Confirmations {
		EosAssert(c.Producer != conf.Producer, &ProducerDoubleConfirm{}, "block already confirmed by this producer")
	}

	key, hasKey := b.ActiveSchedule.GetProducerKey(conf.Producer)
	EosAssert(hasKey, &ProducerNotInSchedule{}, "producer not in current schedule")
	signer, err := conf.ProducerSignature.PublicKey(b.SigDigest().Bytes())
	if err != nil {
		Throw(err)
	}
	EosAssert(signer.Compare(key), &WrongSigningKey{}, "confirmation not signed by expected key")

	b.Confirmations = append(b.Confirmations, *conf)
}
