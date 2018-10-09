package types

import (
	"github.com/eosspark/eos-go/common"
)

type BaseActionTrace struct {
	Receipt       ActionReceipt
	Act           Action
	ContextFree	  bool //default false
	Elapsed       common.Microseconds
	CpuUsage      uint64
	Console       string
	TotalCpuUsage uint64                   /// total of inline_traces[x].cpu_usage + cpu_usage
	TrxId         common.TransactionIdType ///< the transaction that generated this action
	BlockNum      uint32
	BlockTime     common.BlockTimeStamp
	ProducerBlockId common.BlockIdType
	AccountRamDeltas []AccountDelta
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
	Scheduled       bool				//false
	ActionTraces    []ActionTrace
	FailedDtrxTrace TransactionTrace
	//TODO exception
	/*fc::optional<fc::exception>                except;
	std::exception_ptr                         except_ptr;*/
}

type AccountDelta struct {
	Account common.AccountName
	Delta   int64
}

func NewBaseActionTrace(ar *ActionReceipt) *BaseActionTrace {
	bat := BaseActionTrace{}
	bat.Receipt = *ar
	bat.BlockNum=0
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

/*func Test(){
	bat:=BaseActionTrace{ActionReceipt{},Action{},common.Microseconds(60),
	10,"10",10,common.TransactionIdType{}}
	bat.Receipt=ActionReceipt{}
}*/
