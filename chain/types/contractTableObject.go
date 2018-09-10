package types

import "github.com/eosspark/eos-go/common"

type TableIDObject struct {
	id    uint32             `storm "unique" json:"id"`
	code  common.AccountName `storm "" json:"code"`
	scope common.ScopeName   `json:"scope"`
	table common.TableName   `json:"table"`
	payer common.AccountName `json:"payer"`
	count uint32             `json:"count"`
}

type TableIdIndex struct {
}
