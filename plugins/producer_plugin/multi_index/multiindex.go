package multi_index

import "github.com/eosspark/eos-go/common"

type TransactionIdWithExpiry struct {
	TrxId  common.TransactionIdType
	Expiry common.TimePoint
}

//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/hashed_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TransactionIdWithExpiryIndex(ById,ByIdNode,TransactionIdWithExpiry)
func (m *TransactionIdWithExpiryIndex) GetById() *ById         { return m.super }
func (m *TransactionIdWithExpiryIndex) GetByExpiry() *ByExpiry { return m.super.super }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ById(TransactionIdWithExpiryIndex,TransactionIdWithExpiryIndexNode,ByExpiry,ByExpiryNode,TransactionIdWithExpiry,common.TransactionIdType,ByIdFunc)
var ByIdFunc = func(n TransactionIdWithExpiry) common.TransactionIdType { return n.TrxId }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByExpiry(TransactionIdWithExpiryIndex,TransactionIdWithExpiryIndexNode,TransactionIdWithExpiryIndexBase,TransactionIdWithExpiryIndexBaseNode,TransactionIdWithExpiry,common.TimePoint,ByExpiryFunc,ByExpireCompare,true)
//go:generate go build
var ByExpiryFunc = func(n TransactionIdWithExpiry) common.TimePoint { return n.Expiry }
var ByExpireCompare = func(a, b common.TimePoint) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}
