package types

import (
	. "github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
)

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treeset" AccountDeltaSet(AccountDelta,CompareAccountDelta,false)
type BaseActionTrace struct {
	Receipt          ActionReceipt
	Act              Action
	ContextFree      bool //default false
	Elapsed          common.Microseconds
	CpuUsage         uint64
	Console          string
	TotalCpuUsage    uint64                   /// total of inline_traces[x].cpu_usage + cpu_usage
	TrxId            common.TransactionIdType ///< the transaction that generated this action
	BlockNum         uint32
	BlockTime        BlockTimeStamp
	ProducerBlockId  common.BlockIdType
	AccountRamDeltas AccountDeltaSet

	Except Exception
}

type ActionTrace struct {
	BaseActionTrace
	InlineTraces []ActionTrace
}

type TransactionTrace struct {
	ID              common.TransactionIdType
	BlockNum        uint32
	BlockTime       BlockTimeStamp
	ProducerBlockId common.BlockIdType
	Receipt         TransactionReceiptHeader
	Elapsed         common.Microseconds
	NetUsage        uint64
	Scheduled       bool //false
	ActionTraces    []ActionTrace
	FailedDtrxTrace *TransactionTrace

	Except    Exception
	ExceptPtr Exception
}

func NewBaseActionTrace(ar *ActionReceipt) *BaseActionTrace {
	bat := BaseActionTrace{}
	bat.Receipt = *ar
	bat.BlockNum = 0
	bat.ContextFree = false
	bat.TotalCpuUsage = 0
	bat.CpuUsage = 0
	return &bat
}

func (t TransactionTrace) IsEmpty() bool {
	return t.ID.IsEmpty() && t.BlockNum == 0 && t.BlockTime == 0 && t.ProducerBlockId.IsEmpty() &&
		t.Receipt.IsEmpty() && t.Elapsed == 0 && t.NetUsage == 0 && t.Scheduled && len(t.ActionTraces) == 0 &&
		t.FailedDtrxTrace == nil && t.Except == nil && t.ExceptPtr == nil
}
