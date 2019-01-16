package chain_interface

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type NextFunction = func(common.StaticVariant)

type ChannelsType int

const (
	PreAcceptedBlock = ChannelsType(iota)
	RejectedBlock
	AcceptedBlockHeader
	AcceptedBlock
	IrreversibleBlock
	AcceptedTransaction
	AppliedTransaction
	AcceptedConfirmation

	//incoming
	Block
	Transaction

	//compat
	TransactionAck
)

type PreAcceptedBlockCaller struct {
	Caller func(s *types.SignedBlock)
}

func (p *PreAcceptedBlockCaller) Call(data ...interface{}) {
	p.Caller(data[0].(*types.SignedBlock))
}

type RejectedBlockCaller = PreAcceptedBlockCaller

type AcceptedBlockHeaderCaller struct {
	Caller func(b *types.BlockState)
}

func (a *AcceptedBlockHeaderCaller) Call(data ...interface{}) {
	a.Caller(data[0].(*types.BlockState))
}

type AcceptedBlockCaller = AcceptedBlockHeaderCaller
type IrreversibleBlockCaller = AcceptedBlockHeaderCaller

type AcceptedTransactionCaller struct {
	Caller func(t *types.TransactionMetadata)
}

func (a *AcceptedTransactionCaller) Call(data ...interface{}) {
	a.Caller(data[0].(*types.TransactionMetadata))
}

type AppliedTransactionCaller struct {
	Caller func(t *types.TransactionTrace)
}

func (a *AppliedTransactionCaller) Call(data ...interface{}) {
	a.Caller(data[0].(*types.TransactionTrace))
}

type AcceptedConfirmationCaller struct {
	Caller func(h *types.HeaderConfirmation)
}

func (a *AcceptedConfirmationCaller) Call(data ...interface{}) {
	a.Caller(data[0].(*types.HeaderConfirmation))
}

type BlockCaller struct {
	Caller func(s *types.SignedBlock)
}

func (b *BlockCaller) Call(data ...interface{}) {
	b.Caller(data[0].(*types.SignedBlock))
}

type TransactionCaller struct {
	Caller func(p *types.PackedTransaction)
}

func (t *TransactionCaller) Call(data ...interface{}) {
	t.Caller(data[0].(*types.PackedTransaction))
}

type TransactionAckCaller struct {
	Caller func(p common.Pair) //<Exception, *PackedTransaction>
}

func (a *TransactionAckCaller) Call(data ...interface{}) {
	a.Caller(data[0].(common.Pair))
}

type MethodsType int

const (
	GetBlockByNumber = MethodsType(iota)
	GetBlockById
	GetHeadBlockId
	GetLibBlockId

	GetLastIrreversibleBlockNumber

	//incoming
	BlockSync
	TransactionAsync
)

type BlockSyncCaller = BlockCaller

type TransactionAsyncCaller struct {
	Caller func(*types.PackedTransaction, bool, NextFunction /*TransactionTrace*/)
}

func (t *TransactionAsyncCaller) Call(data ...interface{}) {
	t.Caller(data[0].(*types.PackedTransaction), data[1].(bool), data[2].(NextFunction))
}
