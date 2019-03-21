package node_transaction

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/libraries/container"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
)

//go:generate go install "github.com/eosspark/eos-go/libraries/multiindex/"
//go:generate go install "github.com/eosspark/eos-go/libraries/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/libraries/multiindex/ordered_index/..."

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/libraries/multiindex/multi_index_container" NodeTransactionIndex(ById,ByIdNode,multi_index.NodeTransactionState)
func (m *NodeTransactionIndex) GetById() *ById             { return m.super }
func (m *NodeTransactionIndex) GetByExpiry() *ByExpiry     { return m.super.super }
func (m *NodeTransactionIndex) GetByBlockNum() *ByBlockNum { return m.super.super.super }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/libraries/multiindex/ordered_index" ById(NodeTransactionIndex,NodeTransactionIndexNode,ByExpiry,ByExpiryNode,multi_index.NodeTransactionState,common.TransactionIdType,ByIdFunc,ByIdCompare,false)
var ByIdFunc = func(n multi_index.NodeTransactionState) common.TransactionIdType { return n.ID }
var ByIdCompare = crypto.Sha256Compare

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/libraries/multiindex/ordered_index" ByExpiry(NodeTransactionIndex,NodeTransactionIndexNode,ByBlockNum,ByBlockNumNode,multi_index.NodeTransactionState,common.TimePointSec,ByExpiryFunc,ByExpiryCompare,true)
var ByExpiryFunc = func(n multi_index.NodeTransactionState) common.TimePointSec { return n.Expires }
var ByExpiryCompare = func(a, b common.TimePointSec) int { return container.UInt32Comparator(uint32(a), uint32(b)) }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/libraries/multiindex/ordered_index" ByBlockNum(NodeTransactionIndex,NodeTransactionIndexNode,NodeTransactionIndexBase,NodeTransactionIndexBaseNode,multi_index.NodeTransactionState,uint32,ByBlockNumFunc,ByBlockNumCompare,true)
var ByBlockNumFunc = func(n multi_index.NodeTransactionState) uint32 { return n.BlockNum }
var ByBlockNumCompare = container.UInt32Comparator

//go:generate go build
