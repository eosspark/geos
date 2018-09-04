package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewBlockHeaderState(t *testing.T) *BlockHeaderState {
	var initPriKey, _ = ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	var initPubKey = initPriKey.PublicKey()
	var eosio = common.AccountName(common.StringToName("eosio"))
	var yuanc = common.AccountName(common.StringToName("yuanc"))
	var tester = common.AccountName(common.StringToName("tester"))

	initSchedule := ProducerScheduleType{0, []ProducerKey{
		{eosio, initPubKey},
		{yuanc, initPubKey},
		{tester, initPubKey},
	}}

	genHeader := new(BlockHeaderState)
	genHeader.ActiveSchedule = initSchedule
	genHeader.PendingSchedule = initSchedule
	genHeader.Header.Timestamp = common.BlockTimeStamp(1162425600) //slot of 2018-6-2 00:00:00:000
	genHeader.ID, _ = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()

	assert.Equal(t, uint32(1), genHeader.BlockNum)

	return genHeader
}

func Test_BlockHeaderState_GetScheduledProducer(t *testing.T) {
	bs := NewBlockHeaderState(t)
	assert.Equal(t, "tester", common.NameToString(uint64(bs.GetScheduledProducer(100).AccountName)))
	assert.Equal(t, "eosio", common.NameToString(uint64(bs.GetScheduledProducer(110).AccountName)))
	assert.Equal(t, "yuanc", common.NameToString(uint64(bs.GetScheduledProducer(120).AccountName)))
}

func Test_BlockHeaderState_GenerateNext(t *testing.T) {
	bs := NewBlockHeaderState(t)

	t100 := common.BlockTimeStamp(1162425600 + 100)
	t2 := common.BlockTimeStamp(1162425602)

	bsNil := bs.GenerateNext(nil)
	bs100 := bs.GenerateNext(&t100)
	bs2 := bs.GenerateNext(&t2)

	assert.Equal(t, common.BlockTimeStamp(1162425601), bsNil.Header.Timestamp)
	assert.Equal(t, common.BlockTimeStamp(1162425700), bs100.Header.Timestamp)
	assert.Equal(t, common.BlockTimeStamp(1162425602), bs2.Header.Timestamp)

	bsNil.SetConfirmed(10)

	assert.Equal(t, []uint8{2}, bsNil.ConfirmCount)

	bss := bsNil.GenerateNext(nil)

	bss.SetConfirmed(2)

	assert.Equal(t, []uint8{1, 2}, bss.ConfirmCount)

}
