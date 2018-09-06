package producer_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

var chain *mockChain

var initPriKey, _ = ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
var initPubKey = initPriKey.PublicKey()
var eosio = common.AccountName(common.StringToName("eosio"))
var yuanc = common.AccountName(common.StringToName("yuanc"))

type mockChain struct {
	head    *types.BlockState
	pending *types.BlockState
}

func init() {
	chain = new(mockChain)

	initSchedule := types.ProducerScheduleType{0, []types.ProducerKey{
		{eosio, initPubKey},
		{yuanc, initPubKey},
	}}

	genHeader := types.BlockHeaderState{}
	genHeader.ActiveSchedule = initSchedule
	genHeader.PendingSchedule = initSchedule
	genHeader.Header.Timestamp = common.NewBlockTimeStamp(common.Now())
	genHeader.ID = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()

	genHeader.ProducerToLastProduced = make(map[common.AccountName]uint32)
	genHeader.ProducerToLastImpliedIrb = make(map[common.AccountName]uint32)

	chain.head = new(types.BlockState)
	chain.head.BlockHeaderState = genHeader

	fmt.Println("now", common.Now())
	fmt.Println("init", genHeader.Header.Timestamp.ToTimePoint())
}

func (c mockChain) LastIrreversibleBlockNum() uint32 {
	return c.head.DposIrreversibleBlocknum
}

func (c mockChain) HeadBlockState() *types.BlockState {
	return c.head

}
func (c mockChain) HeadBlockTime() common.TimePoint {
	return c.head.Header.Timestamp.ToTimePoint()
}

func (c mockChain) PendingBlockTime() common.TimePoint {
	return c.pending.Header.Timestamp.ToTimePoint()
}

func (c mockChain) HeadBlockNum() uint32 {
	return c.head.BlockNum
}

func (c mockChain) PendingBlockState() *types.BlockState {
	return c.pending
}

func (c mockChain) GetUnappliedTransactions() []*types.TransactionMetadata {
	return make([]*types.TransactionMetadata, 0)
}

func (c mockChain) GetScheduledTransactions() []common.TransactionIDType {
	return make([]common.TransactionIDType, 0)
}

func (c *mockChain) AbortBlock() {
	fmt.Println("abort block...")
	if c.pending != nil {
		c.pending = nil
	}
}

func (c *mockChain) StartBlock(when common.BlockTimeStamp, confirmBlockCount uint16) {
	fmt.Println("start block...")
	chain.pending = new(types.BlockState)
	chain.pending.BlockHeaderState = *chain.head.GenerateNext(&when)
	chain.pending.SetConfirmed(confirmBlockCount)

}
func (c *mockChain) FinalizeBlock()                       { fmt.Println("finalize block...") }
func (c *mockChain) SignBlock(func([]byte) ecc.Signature) { fmt.Println("sign block...") }
func (c *mockChain) CommitBlock() {
	fmt.Println("commit block...")
	c.head = c.pending
	c.pending = nil
}

func (c mockChain) PushTransaction(trx *types.TransactionMetadata, deadline common.TimePoint) error {
	return nil
}
func (c mockChain) PushScheduledTransaction(trx common.TransactionIDType, deadline common.TimePoint) error {
	return nil
}

func (c *mockChain) PushBlock(b *types.SignedBlock) error {
	return nil
}
