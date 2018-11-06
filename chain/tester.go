package chain

import "github.com/eosspark/eos-go/chain/types"

type tester struct{
	control *Controller
}

func (t tester)pushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.control.AbortBlock()
	//t.control.PushBlock(b)
	return &types.SignedBlock{}
}