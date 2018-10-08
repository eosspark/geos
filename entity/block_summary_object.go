package entity

import "github.com/eosspark/eos-go/chain/types"

type BlockSummaryObject struct {
	Id      types.IdType `storm:"id , increment"`
	BlockId types.IdType
}
