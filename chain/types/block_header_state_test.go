package types

import (
	"testing"

	. "github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
)

func NewBlockHeaderState(t *testing.T) *BlockHeaderState {
	initPriKey, err := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	assert.NoError(t, err)

	initPubKey := initPriKey.PublicKey()
	assert.Equal(t, "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", initPubKey.String())

	eosio := common.N("eosio")
	yuanc := common.N("yuanc")
	tester := common.N("tester")

	initSchedule := ProducerScheduleType{0, []ProducerKey{
		{eosio, initPubKey},
		{yuanc, initPubKey},
		{tester, initPubKey},
	}}

	genHeader := new(BlockHeaderState)
	genHeader.ActiveSchedule = initSchedule
	genHeader.PendingSchedule = initSchedule
	genHeader.PendingScheduleHash = *crypto.Hash256(initSchedule)
	genHeader.Header.Timestamp = BlockTimeStamp(1162339200) //1162339200 slot of 2018-6-1T12:00:00 UTC
	genHeader.Header.Confirmed = 1
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()
	genHeader.ProducerToLastProduced = *NewAccountNameUint32Map()
	genHeader.ProducerToLastImpliedIrb = *NewAccountNameUint32Map()

	genHeader.BlockSigningKey = initPubKey

	assert.Equal(t, uint32(1), genHeader.BlockNum)

	return genHeader
}

func Test_BlockHeaderState_GetScheduledProducer(t *testing.T) {
	bs := NewBlockHeaderState(t)
	assert.Equal(t, "tester", common.S(uint64(bs.GetScheduledProducer(100).ProducerName)))
	assert.Equal(t, "eosio", common.S(uint64(bs.GetScheduledProducer(110).ProducerName)))
	assert.Equal(t, "yuanc", common.S(uint64(bs.GetScheduledProducer(120).ProducerName)))
}

func Test_BlockHeaderState_GenerateNext(t *testing.T) {
	bs := NewBlockHeaderState(t)

	t100 := BlockTimeStamp(1162339200 + 100)
	t2 := BlockTimeStamp(1162339200 + 2)

	bsNil := bs.GenerateNext(0)
	bs100 := bs.GenerateNext(t100)
	bs2 := bs.GenerateNext(t2)

	assert.Equal(t, BlockTimeStamp(1162339201), bsNil.Header.Timestamp)
	assert.Equal(t, BlockTimeStamp(1162339300), bs100.Header.Timestamp)
	assert.Equal(t, BlockTimeStamp(1162339202), bs2.Header.Timestamp)

	bsNil.SetConfirmed(10)

	assert.Equal(t, []uint8{2}, bsNil.ConfirmCount)

	bss := bsNil.GenerateNext(0)

	bss.SetConfirmed(2)

	assert.Equal(t, []uint8{1, 2}, bss.ConfirmCount)

}

func TestBlockHeader_BlockID(t *testing.T) {
	bs := NewBlockHeaderState(t)
	bid := bs.Header.BlockID()

	assert.EqualValues(t, 1, NumFromID(&bid))
	assert.EqualValues(t, 1, bs.Header.BlockNumber())

	bs1 := bs.GenerateNext(0)
	bid = bs1.Header.BlockID()

	assert.EqualValues(t, 2, NumFromID(&bid))
	assert.EqualValues(t, 2, bs1.Header.BlockNumber())

}

func TestBlockHeader_Digest(t *testing.T) {
	bs := NewBlockHeaderState(t)
	assert.Equal(t,
		"29bec0f6df7f0ba248b72e755d723381f8b4f829221a4f29e49fd5396a10ed0f", //calculate by eosioc++
		bs.SigDigest().String())
}

func TestBlockHeaderState_Sign(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	bs := NewBlockHeaderState(t)

	bs.Sign(func(sha256 crypto.Sha256) ecc.Signature {
		sk, _ := initPriKey.Sign(sha256.Bytes())
		return sk
	})

	assert.Equal(t, initPriKey.PublicKey(), bs.Signee())
}
