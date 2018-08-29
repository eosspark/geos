package producer_plugin

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"time"
)

var chain = new(mockChain)

type mockChain struct {
}

func (c mockChain) HeadBlockState() *types.BlockState {
	result := new(types.BlockState)
	header := types.SignedBlockHeader{}
	header.Timestamp = common.NewBlockTimeStamp(time.Now())
	result.Header = header
	return result

}

func (c mockChain) HeadBlockTime() time.Time {
	return time.Now()
}

func (c mockChain) PendingBlockTime() time.Time {
	return time.Now()
}

func (c mockChain) HeadBlockNum() uint32 {
	return 1
}

func (c mockChain) PendingBlockState() *types.BlockState {
	return new(types.BlockState)
}

func (c mockChain) GetUnappliedTransactions() []*types.TransactionMetadata {
	return make([]*types.TransactionMetadata, 10)
}

func (c mockChain) GetScheduledTransactions() []common.TransactionIDType {
	return make([]common.TransactionIDType, 10)
}

func (c mockChain) AbortBlock()                                                     {}
func (c mockChain) StartBlock(when common.BlockTimeStamp, confirmBlockCount uint16) {}
func (c mockChain) FinalizeBlock()                                                  {}
func (c mockChain) SignBlock(func([]byte) ecc.Signature)                            {}
func (c mockChain) CommitBlock()                                                    {}

func (c mockChain) PushTransaction(trx *types.TransactionMetadata, deadline time.Time) error {
	return nil
}
func (c mockChain) PushScheduledTransaction(trx common.TransactionIDType, deadline time.Time) error {
	return nil
}

func (c mockChain) PushBlock(b *types.SignedBlock) error {
	return nil
}
