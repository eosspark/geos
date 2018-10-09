package entity

import "github.com/eosspark/eos-go/common"

type TransactionObject struct {
	ID         common.IdType            `storm:"id,increment,byExpiration"` //byID,byExpiration
	Expiration common.TimePointSec      `storm:"byExpiration,unique"`       //byExpiration
	TrxID      common.TransactionIdType `storm:"unique"`                    //byTrxID
}
