package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"reflect"
)

type BlockState struct {
	BlockHeaderState `multiIndex:"inline"`
	SignedBlock      *SignedBlock `multiIndex:"inline"`
	Validated        bool         `json:"validated"`
	InCurrentChain   bool         `json:"in_current_chain"`
	Trxs             []*TransactionMetadata
}

func NewBlockState(cur *BlockHeaderState) *BlockState {
	return &BlockState{*cur, &SignedBlock{},
		false, false, make([]*TransactionMetadata, 0)}
}

func NewBlockState2(prev *BlockHeaderState, when BlockTimeStamp) *BlockState {
	bs := &BlockState{
		BlockHeaderState: *prev.GenerateNext(when),
		SignedBlock:      &SignedBlock{},
	}
	bs.SignedBlock.SignedBlockHeader = bs.Header
	return bs
}

func NewBlockState3(prev *BlockHeaderState, b *SignedBlock, trust bool) *BlockState {
	return &BlockState{
		BlockHeaderState: *prev.Next(b.SignedBlockHeader, trust),
		SignedBlock:      b,
	}
}

//for treeset
var BlockIdTypes = reflect.TypeOf(common.BlockIdType(*crypto.NewSha256Nil()))

//for treeset
var BlockNumType = reflect.TypeOf(uint32(0))
