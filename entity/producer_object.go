package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type ProducerObject struct {
	ID                    common.IdType      `multiIndex:"id,increment,byKey"`
	Owner                 common.AccountName `multiIndex:"byOwner,orderedUnique"`
	LastAslot             uint64             //c++ default value 0
	SigningKey            ecc.PublicKey      `multiIndex:"byKey,orderedUnique"`
	TotalMissed           int64              //c++ default value 0
	LastConfirmedBlockNum uint32

	/// The blockchain configuration values this producer recommends
	//chain_config       configuration //TODO
}
