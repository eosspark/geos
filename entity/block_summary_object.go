package entity

import ("github.com/eosspark/eos-go/common"
)

type BlockSummaryObject struct {
	Id      common.IdType `multiIndex:"id,increment"`
	BlockId common.BlockIdType
}
