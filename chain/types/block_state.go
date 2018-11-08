package types

import "github.com/eosspark/eos-go/common"

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

func NewBlockState2(prev *BlockHeaderState, when common.BlockTimeStamp) (bs *BlockState) {
	bs = new(BlockState)
	bs.BlockHeaderState = *prev.GenerateNext(when)
	bs.SignedBlock = new(SignedBlock)
	bs.SignedBlock.SignedBlockHeader = bs.Header
	return
}

func NewBlockState3(prev *BlockHeaderState, b *SignedBlock, trust bool) (bs *BlockState) {
	bs = new(BlockState)
	bs.BlockHeaderState = *prev.Next(b.SignedBlockHeader, trust)
	bs.SignedBlock = b
	return
}
