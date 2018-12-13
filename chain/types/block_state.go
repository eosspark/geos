package types

import (
	"bytes"
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

//for treeset
var BlockIdTypes = reflect.TypeOf(common.BlockIdType(*crypto.NewSha256Nil()))

func CompareBlockId(first interface{}, second interface{}) int {
	return bytes.Compare(first.(*BlockState).BlockId.Bytes(), second.(*BlockState).BlockId.Bytes())
}

func ComparePrev(first interface{}, second interface{}) int {
	return bytes.Compare(first.(*BlockState).BlockId.Bytes(), second.(*BlockState).BlockId.Bytes())
}

//for treeset
var BlockNumType = reflect.TypeOf(uint32(0))

func CompareBlockNum(first interface{}, second interface{}) int {
	if first.(*BlockState).InCurrentChain {
		if first.(*BlockState).BlockNum == second.(*BlockState).BlockNum {
			return 0
		} else if first.(*BlockState).BlockNum < second.(*BlockState).BlockNum {
			return -1
		} else {
			return 1
		}
	} else {
		if ^first.(*BlockState).BlockNum+1 == second.(*BlockState).BlockNum {
			return 0
		} else if ^first.(*BlockState).BlockNum+1 < second.(*BlockState).BlockNum {
			return -1
		} else {
			return 1
		}
	}
}

func CompareLibNum(first interface{}, second interface{}) int {
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
