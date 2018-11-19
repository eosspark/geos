package types

import (
		"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewBlockHeaderState(t *testing.T) *BlockHeaderState {
	initPriKey, err := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	assert.NoError(t, err)

	initPubKey := initPriKey.PublicKey()
	assert.Equal(t, "EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", initPubKey.String())

	eosio := common.AccountName(common.N("eosio"))
	yuanc := common.AccountName(common.N("yuanc"))
	tester := common.AccountName(common.N("tester"))

	initSchedule := ProducerScheduleType{0, []ProducerKey{
		{eosio, initPubKey},
		{yuanc, initPubKey},
		{tester, initPubKey},
	}}

	genHeader := new(BlockHeaderState)
	genHeader.ActiveSchedule = initSchedule
	genHeader.PendingSchedule = initSchedule
	genHeader.Header.Timestamp = BlockTimeStamp(1162425600) //slot of 2018-6-2 00:00:00:000 UTC
	genHeader.Header.Confirmed = 1
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()

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

	t100 := BlockTimeStamp(1162425600 + 100)
	t2 := BlockTimeStamp(1162425602)

	bsNil := bs.GenerateNext(0)
	bs100 := bs.GenerateNext(t100)
	bs2 := bs.GenerateNext(t2)

	assert.Equal(t, BlockTimeStamp(1162425601), bsNil.Header.Timestamp)
	assert.Equal(t, BlockTimeStamp(1162425700), bs100.Header.Timestamp)
	assert.Equal(t, BlockTimeStamp(1162425602), bs2.Header.Timestamp)

	bsNil.SetConfirmed(10)

	assert.Equal(t, []uint8{2}, bsNil.ConfirmCount)

	bss := bsNil.GenerateNext(0)

	bss.SetConfirmed(2)

	assert.Equal(t, []uint8{1, 2}, bss.ConfirmCount)

}

func TestBlockHeader_BlockID(t *testing.T) {
	bs  := NewBlockHeaderState(t)
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
		"be5bfe468786e957b46741d9d7ba012d65c533a589934f6d653402e688ad66a0",
		bs.SigDigest().String())
}

func TestBlockHeaderState_Sign(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	bs := NewBlockHeaderState(t)

	//fmt.Println("===>", bs.SigDigest())
	bs.Sign(func(sha256 crypto.Sha256) ecc.Signature {
		sk, _ := initPriKey.Sign(sha256.Bytes())
		return sk
	})
	//fmt.Println("===>", bs.SigDigest())

	assert.Equal(t, initPriKey.PublicKey(), bs.Signee())

	//data := ""
	//sk,_ := initPriKey.Sign(crypto.Hash256(data).Bytes())
	//pk,_ := sk.PublicKey(crypto.Hash256(data).Bytes())
	//
	//fmt.Println("pk", pk)
	//
	//assert.Equal(t, initPriKey.PublicKey(), pk)

}
