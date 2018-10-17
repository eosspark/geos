package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"testing"
		"github.com/eosspark/eos-go/crypto"
	"fmt"
)

func NewBlockHeaderState(t *testing.T) *BlockHeaderState {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()
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
	genHeader.Header.Timestamp = common.BlockTimeStamp(1162425600) //slot of 2018-6-2 00:00:00:000
	genHeader.ID = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()

	genHeader.ProducerToLastProduced = make(map[common.AccountName]uint32)
	genHeader.ProducerToLastImpliedIrb = make(map[common.AccountName]uint32)

	genHeader.BlockSigningKey = initPubKey

	assert.Equal(t, uint32(1), genHeader.BlockNum)

	return genHeader
}

func Test_BlockHeaderState_GetScheduledProducer(t *testing.T) {
	bs := NewBlockHeaderState(t)
	assert.Equal(t, "tester", common.S(uint64(bs.GetScheduledProducer(100).AccountName)))
	assert.Equal(t, "eosio", common.S(uint64(bs.GetScheduledProducer(110).AccountName)))
	assert.Equal(t, "yuanc", common.S(uint64(bs.GetScheduledProducer(120).AccountName)))
}

func Test_BlockHeaderState_GenerateNext(t *testing.T) {
	bs := NewBlockHeaderState(t)

	t100 := common.BlockTimeStamp(1162425600 + 100)
	t2 := common.BlockTimeStamp(1162425602)

	bsNil := bs.GenerateNext(0)
	bs100 := bs.GenerateNext(t100)
	bs2 := bs.GenerateNext(t2)

	assert.Equal(t, common.BlockTimeStamp(1162425601), bsNil.Header.Timestamp)
	assert.Equal(t, common.BlockTimeStamp(1162425700), bs100.Header.Timestamp)
	assert.Equal(t, common.BlockTimeStamp(1162425602), bs2.Header.Timestamp)

	bsNil.SetConfirmed(10)

	assert.Equal(t, []uint8{2}, bsNil.ConfirmCount)

	bss := bsNil.GenerateNext(0)

	bss.SetConfirmed(2)

	assert.Equal(t, []uint8{1, 2}, bss.ConfirmCount)

}

func TestBlockHeader_Digest(t *testing.T) {
	bs := NewBlockHeaderState(t)
	fmt.Println(bs.SigDigest())
	fmt.Println(bs.SigDigest())
}

func TestBlockHeaderState_Sign(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	bs := NewBlockHeaderState(t)

	fmt.Println("===>", bs.SigDigest())
	bs.Sign(func(sha256 crypto.Sha256) ecc.Signature {
		sk, _ := initPriKey.Sign(sha256.Bytes())
		return sk
	})
	fmt.Println("===>", bs.SigDigest())

	assert.Equal(t, initPriKey.PublicKey(), bs.Signee())

	//data := ""
	//sk,_ := initPriKey.Sign(crypto.Hash256(data).Bytes())
	//pk,_ := sk.PublicKey(crypto.Hash256(data).Bytes())
	//
	//fmt.Println("pk", pk)
	//
	//assert.Equal(t, initPriKey.PublicKey(), pk)

}
