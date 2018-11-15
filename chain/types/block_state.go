package types

import (
	"bytes"
	"github.com/eosspark/eos-go/common"
)

type BlockState struct {
	BlockHeaderState `multiIndex:"inline"`
	SignedBlock      *SignedBlock `multiIndex:"inline"`
	Validated        bool         `json:"validated"`
	InCurrentChain   bool         `json:"in_current_chain"`
	Trxs             []*TransactionMetadata
}

func NewBlockState(cur *BlockHeaderState) *BlockState {
	//a := new(TransactionMetadata)
	//b :=make([]*TransactionMetadata,1)
	//b[0]=a

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

func (b BlockState) GetKey() []byte {
	return b.BlockId.Bytes()
}

func (b *BlockState) ElementObject() {

}

func CompareBlockId(first common.ElementObject, second common.ElementObject) int {
	fir := first.(*BlockState)
	sec := second.(*BlockState)
	result := bytes.Compare(fir.BlockId.Bytes(), sec.BlockId.Bytes())
	return result
}

func ComparePrev(first common.ElementObject, second common.ElementObject) int {
	fir := first.(*BlockState)
	sec := second.(*BlockState)

	if fir.BlockNum == sec.BlockNum {
		return 0
	} else if fir.BlockNum < sec.BlockNum {
		return -1
	} else {
		return 1
	}

}

func CompareBlockNum(first common.ElementObject, second common.ElementObject) int {
	fir := first.(*BlockState)
	sec := second.(*BlockState)
	if fir.InCurrentChain /* && sec.InCurrentChain*/ {
		if fir.BlockNum == sec.BlockNum {
			return 0
		} else if fir.BlockNum < sec.BlockNum {
			return -1
		} else {
			return 1
		}
	} else {
		if ^fir.BlockNum+1 == sec.BlockNum {
			return 0
		} else if ^fir.BlockNum+1 < sec.BlockNum {
			return -1
		} else {
			return 1
		}
	}
}

func CompareLibNum(first common.ElementObject, second common.ElementObject) int {
	//by_lib_block_num
	if first.(*BlockState).DposIrreversibleBlocknum == second.(*BlockState).DposIrreversibleBlocknum &&
		first.(*BlockState).BftIrreversibleBlocknum == second.(*BlockState).BftIrreversibleBlocknum &&
		first.(*BlockState).BlockNum == second.(*BlockState).BlockNum {
		return 0
	} else if first.(*BlockState).DposIrreversibleBlocknum > second.(*BlockState).DposIrreversibleBlocknum ||
		first.(*BlockState).BftIrreversibleBlocknum > second.(*BlockState).BftIrreversibleBlocknum ||
		first.(*BlockState).BlockNum > second.(*BlockState).BlockNum {
		return -1
	} else {
		return 1
	}
}
