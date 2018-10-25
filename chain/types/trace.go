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
	AccountRamDeltas FlatSet
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

type FlatSet struct {
	data []AccountDelta
	less func(first, secend *AccountDelta) bool
}

func (f *FlatSet) Len() int {
	return len(f.data)
}

func (f *FlatSet) Swap(i, j int) {
	f.data[i], f.data[j] = f.data[j], f.data[i]
}

func (f *FlatSet) Less(i, j int) bool {
	return f.less(&f.data[i], &f.data[j])
}

func less(f, s *AccountDelta) bool {
	return f.Account < s.Account
}

func compare(f AccountDelta, s AccountDelta) bool {
	return f.Account < s.Account
}

func (f *FlatSet) Append(account common.AccountName, delta int64) {
	ap := AccountDelta{}
	ap.Account = account
	ap.Delta = delta
	if len(f.data) == 0 {
		f.data = append(f.data, ap)
	} else {
		for i := 0; i < len(f.data); i++ {
			if len(f.data) == 1 {
				f.data = append(f.data, ap)
				if !compare(f.data[i], ap) {
					f.Swap(i, i+1)
				}
			}
			if compare(f.data[i], ap) && !compare(f.data[i+1], ap) {
				first := f.data[:i+1]
				secend := make([]AccountDelta, len(f.data[i+1:len(f.data)]))
				copy(secend, f.data[i+1:len(f.data)])
				first = append(first, ap)
				first = append(first, secend...)
				//fmt.Println(first)
				f.data = first
			}
		}
	}
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

func NewAccountDelta(name *common.AccountName, d int64) *AccountDelta {
	ad := AccountDelta{}
	ad.Account = *name
	ad.Delta = d
	return &ad
}
