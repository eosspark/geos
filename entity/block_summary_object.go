package entity

import ("github.com/eosspark/eos-go/common"
)

type BlockSummaryObject struct {
	Id      common.IdType `storm:"id,increment"`
	BlockId common.BlockIdType
}
