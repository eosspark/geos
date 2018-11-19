package include

import (
	. "github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type Function interface {
	call(data interface{})
}

type PreAcceptedBlockFunc struct {
	function  func(s *SignedBlock)
}

func (p *PreAcceptedBlockFunc) call(data interface{}) {
	p.function(data.(*SignedBlock))
}

type RejectedBlockTagFunc struct {
	function func(s *SignedBlock)
}

func (r *RejectedBlockTagFunc) call(data interface{}) {
	r.function(data.(*SignedBlock))
}

type AcceptedBlockHeaderTag struct {
	function func(b *BlockState)
}

func (a *AcceptedBlockHeaderTag) call(data interface{}) {
	a.function(data.(*BlockState))
}

type AcceptedBlockTag struct {
	function func(b *BlockState)
}

func (a *AcceptedBlockTag) call(data interface{}) {
	a.function(data.(*BlockState))
}

type IrreversibleBlockTag struct {
	function func(b *BlockState)
}

func (i *IrreversibleBlockTag) call(data interface{}) {
	i.function(data.(*BlockState))
}

type AcceptedTransactionTag struct {
	function func(t *TransactionMetadata)
}

func (a *AcceptedTransactionTag) call(data interface{}) {
	a.function(data.(*TransactionMetadata))
}

type AppliedTransactionTag struct {
	function func(t *TransactionTrace)
}

func (a *AppliedTransactionTag) call(data interface{}) {
	a.function(data.(*TransactionTrace))
}

type AcceptedConfirmationTag struct {
	function func(h *HeaderConfirmation)
}

func (a *AcceptedConfirmationTag) call(data interface{}) {
	a.function(data.(*HeaderConfirmation))
}



type BlockTag struct {
	function func(s *SignedBlock)
}

func (b *BlockTag) call(data interface{}) {
	b.function(data.(*SignedBlock))
}

type TransactionTag struct {
	function func(p *PackedTransaction)
}

func (t *TransactionTag) call(data interface{}) {
	t.function(data.(*PackedTransaction))
}



type AcceptedTransactionTagExec struct {
	function func(p *common.Pair)
}

func (a *AcceptedTransactionTagExec) call(data interface{}) {
	a.function(data.(*common.Pair))
}



