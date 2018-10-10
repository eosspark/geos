package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type ProducerObject struct {
	ID                    common.IdType      `storm:"id,increment,byKey"`
	Owner                 common.AccountName `storm:"unique"` //byOwner
	LastAslot             uint64             //c++ default value 0
	SigningKey            ecc.PublicKey      `storm:"unique,byKey"`
	TotalMissed           int64              //c++ default value 0
	LastConfirmedBlockNum uint32

	/// The blockchain configuration values this producer recommends
	//chain_config       configuration //TODO
}

func (s *ProducerObject) makeTuple(signingKey ecc.PublicKey, id common.IdType) *common.Tuple {
	result := common.MakeTuple(signingKey, id)
	return &result
}
