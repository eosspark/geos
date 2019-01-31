package multi_index

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type NodeTransactionState struct {
	ID            common.TransactionIdType
	Expires       common.TimePointSec
	PackedTxn     types.PackedTransaction
	SerializedTxn []byte
	BlockNum      uint32
	TrueBlock     uint32
	Requests      uint16
}

type TransactionState struct {
	ID              common.TransactionIdType
	IsKnownByPeer   bool
	IsNoticedToPeer bool
	BlockNum        uint32
	Expires         common.TimePointSec
	RequestedTime   common.TimePoint
}

type PeerBlockState struct {
	ID            common.BlockIdType
	BlockNum      uint32
	IsKnown       bool
	IsNoticed     bool
	RequestedTime common.TimePoint
}
