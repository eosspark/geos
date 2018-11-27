package types

import (
	"bytes"
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

func NewBlockState2(prev *BlockHeaderState, when BlockTimeStamp) (bs *BlockState) {
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

func CompareBlockId(first *BlockState, second *BlockState) int {
	/*fir := first.(*BlockState)
	sec := second.(*BlockState)*/
	result := bytes.Compare(first.BlockId.Bytes(), second.BlockId.Bytes())
	return result
}

func ComparePrev(first *BlockState, second *BlockState) int {
	/*fir := first.(*BlockState)
	sec := second.(*BlockState)*/

	if first.BlockNum == second.BlockNum {
		return 0
	} else if first.BlockNum < second.BlockNum {
		return -1
	} else {
		return 1
	}

}

func CompareBlockNum(first *BlockState, second *BlockState) int {
	/*fir := first.(*BlockState)
	sec := second.(*BlockState)*/
	if first.InCurrentChain /* && sec.InCurrentChain*/ {
		if first.BlockNum == second.BlockNum {
			return 0
		} else if first.BlockNum < second.BlockNum {
			return -1
		} else {
			return 1
		}
	} else {
		if ^first.BlockNum+1 == second.BlockNum {
			return 0
		} else if ^first.BlockNum+1 < second.BlockNum {
			return -1
		} else {
			return 1
		}
	}
}

func CompareLibNum(first *BlockState, second *BlockState) int {
	//by_lib_block_num
	if first.DposIrreversibleBlocknum == second.DposIrreversibleBlocknum &&
		first.BftIrreversibleBlocknum == second.BftIrreversibleBlocknum &&
		first.BlockNum == second.BlockNum {
		return 0
	} else if first.DposIrreversibleBlocknum > second.DposIrreversibleBlocknum ||
		first.BftIrreversibleBlocknum > second.BftIrreversibleBlocknum ||
		first.BlockNum > second.BlockNum {
		return -1
	} else {
		return 1
	}
}
