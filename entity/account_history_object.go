package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

type AccountHistoryObject struct {
	ID                 common.IdType 		`multiIndex:"id,increment"`
	Account            common.AccountName	`multiIndex:"byAccountActionSeq,orderedUnique"`
	ActionSequenceNum  uint64
	AccountSequenceNum int32				`multiIndex:"byAccountActionSeq,orderedUnique"`
}

type ActionHistoryObject struct {
	ID                common.IdType			   `multiIndex:"id,increment"`
	ActionSequenceNum uint64        		   `multiIndex:"byActionSequenceNum,orderedUnique:byTrxId,orderedUnique"`
	PackedActionTrace common.HexBytes
	BlockNum          uint32
	BlockTime         types.BlockTimeStamp
	TrxId             common.TransactionIdType `multiIndex:"byTrxId,orderedUnique"`
}

//type FilterEntry struct {
//	Receiver common.Name
//	Action common.Name
//	Actor common.Name
//}
//
//func (fe *FilterEntry) Key() common.Tuple{
//	return common.MakeTuple(fe.Receiver,fe.Action,fe.Actor)
//}
