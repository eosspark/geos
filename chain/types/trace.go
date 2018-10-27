package types

import (
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
)

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
	BlockTime        common.BlockTimeStamp
	ProducerBlockId  common.BlockIdType
	AccountRamDeltas common.FlatSet
}

type ActionTrace struct {
	BaseActionTrace
	InlineTraces []ActionTrace
}

type TransactionTrace struct {
	ID              common.TransactionIdType
	BlockNum        uint32
	BlockTime       common.BlockTimeStamp
	ProducerBlockId common.BlockIdType
	Receipt         TransactionReceiptHeader
	Elapsed         common.Microseconds
	NetUsage        uint64
	Scheduled       bool //false
	ActionTraces    []ActionTrace
	FailedDtrxTrace *TransactionTrace
	//TODO exception
	Except Exception
	/*fc::optional<fc::exception>                except;
	std::exception_ptr                         except_ptr;*/
}

type AccountDelta struct {
	Account common.AccountName
	Delta   int64
}

func (a AccountDelta) Compare(first common.FlatSet, second common.FlatSet) bool {
	return first.(AccountDelta).Account <= second.(AccountDelta).Account
}

/*func search(n int, f func(int) bool) int {
	result, i, j := 0, 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if !f(h) {
			i = h + 1
		} else {
			j = h
		}
		result = h
	}
	return result
}*/

/*func (f *FlatSet) Append(account common.AccountName, delta int64) {
	ap := AccountDelta{}
	ap.Account = account
	ap.Delta = delta
	param := AccountDelta{}
	param.Account = account
	param.Delta = delta
	if len(f.data) == 0 {
		f.data = append(f.data, ap)
	} else {
		target := search(len(f.data), func(i int) bool { return f.data[i].Account < param.Account && f.data[i+1].Account > param.Account }) + 1

		first := f.data[:target+1]
		second := make([]AccountDelta, len(f.data[target+1:len(f.data)]))
		copy(second, f.data[target+1:len(f.data)])
		first = append(first, ap)
		first = append(first, second...)
		f.data = first
	}
}*/

func NewBaseActionTrace(ar *ActionReceipt) *BaseActionTrace {
	bat := BaseActionTrace{}
	bat.Receipt = *ar
	bat.BlockNum = 0
	bat.ContextFree = false
	bat.TotalCpuUsage = 0
	bat.CpuUsage = 0
	return &bat
}

func NewAccountDelta(name *common.AccountName, d int64) *AccountDelta {
	ad := AccountDelta{}
	ad.Account = *name
	ad.Delta = d
	return &ad
}
