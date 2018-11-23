package include

import (
	. "github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type Function interface {
	call(data interface{})
}

type PreAcceptedBlockFunc struct {
	Func  func(s *SignedBlock)
}

func (p *PreAcceptedBlockFunc) call(data interface{}) {
	p.Func(data.(*SignedBlock))
}

type RejectedBlockFunc struct {
	Func func(s *SignedBlock)
}

func (r *RejectedBlockFunc) call(data interface{}) {
	r.Func(data.(*SignedBlock))
}

type AcceptedBlockHeaderFunc struct {
	Func func(b *BlockState)
}

func (a *AcceptedBlockHeaderFunc) call(data interface{}) {
	a.Func(data.(*BlockState))
}

type AcceptedBlockFunc struct {
	Func func(b *BlockState)
}

func (a *AcceptedBlockFunc) call(data interface{}) {
	a.Func(data.(*BlockState))
}

type IrreversibleBlockFunc struct {
	Func func(b *BlockState)
}

func (i *IrreversibleBlockFunc) call(data interface{}) {
	i.Func(data.(*BlockState))
}

type AcceptedTransactionFunc struct {
	Func func(t *TransactionMetadata)
}

func (a *AcceptedTransactionFunc) call(data interface{}) {
	a.Func(data.(*TransactionMetadata))
}

type AppliedTransactionFunc struct {
	Func func(t *TransactionTrace)
}

func (a *AppliedTransactionFunc) call(data interface{}) {
	a.Func(data.(*TransactionTrace))
}

type AcceptedConfirmationFunc struct {
	Func func(h *HeaderConfirmation)
}

func (a *AcceptedConfirmationFunc) call(data interface{}) {
	a.Func(data.(*HeaderConfirmation))
}



type BlockFunc struct {
	Func func(s *SignedBlock)
}

func (b *BlockFunc) call(data interface{}) {
	b.Func(data.(*SignedBlock))
}

type TransactionFunc struct {
	Func func(p *PackedTransaction)
}

func (t *TransactionFunc) call(data interface{}) {
	t.Func(data.(*PackedTransaction))
}



type AcceptedTransactionFuncExec struct {
	Func func(p *common.Pair)
}

func (a *AcceptedTransactionFuncExec) call(data interface{}) {
	a.Func(data.(*common.Pair))
}



