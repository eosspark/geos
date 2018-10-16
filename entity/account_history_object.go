package entity

import (
	"github.com/eosspark/eos-go/common"
)

type AccountHistoryObject struct {
	ID                 common.IdType `storm:"id,increment,byId"`
	Account            common.AccountName	`storm:"unique,ByAccountActionSeq"`
	ActionSequenceNum  uint64
	AccountSequenceNum int32				`storm:"unique,ByAccountActionSeq"`
}

type ActionHistoryObject struct {
	ID                common.IdType `storm:"id,increment,byId"`
	ActionSequenceNum uint64        `storm:"orderedUnique,ByActionSequenceNum,ByTrxId"`
	PackedActionTrace common.HexBytes
	BlockNum          uint32
	BlockTime         common.BlockTimeStamp
	TrxId             common.TransactionIdType	`storm:"unique,ByTrxId"`
}

type FilterEntry struct {
	Receiver common.Name
	Action common.Name
	Actor common.Name
}

func (fe *FilterEntry) Key() common.Tuple{
	return common.MakeTuple(fe.Receiver,fe.Action,fe.Actor)
}
