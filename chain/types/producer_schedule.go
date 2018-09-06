package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

type ProducerKey struct {
	AccountName     common.AccountName `json:"account_name"`
	BlockSigningKey ecc.PublicKey      `json:"block_signing_key"`
}

type ProducerScheduleType struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

func (ps *ProducerScheduleType) GetProducerKey(p common.AccountName) (ecc.PublicKey, bool) {
	for _, i := range ps.Producers {
		if i.AccountName == p {
			return i.BlockSigningKey, true
		}
	}
	return ecc.PublicKey{}, false
}

type SharedProducerScheduleType struct {
	Version   uint32
	Producers []ProducerKey
}

func (spst *SharedProducerScheduleType) Clear() {
	spst.Version = 0
	spst.Producers = []ProducerKey{}
}

func (spst *SharedProducerScheduleType) SharedProducerScheduleType(a ProducerScheduleType) *ProducerScheduleType {
	var result = ProducerScheduleType{}
	spst.Version = a.Version
	spst.Producers = nil
	//spst.Producers = a.Producers
	for i := 0; i < len(a.Producers); i++ {
		spst.Producers[i] = a.Producers[i]
	}
	return &result
}

func (spst *SharedProducerScheduleType) ProducerScheduleType() *ProducerScheduleType {
	var result = ProducerScheduleType{}
	result.Version = spst.Version
	if len(result.Producers) == 0 {
		result.Producers = spst.Producers
	} else {
		var step = len(result.Producers)
		for _, p := range spst.Producers {
			result.Producers[step] = p
			step++
		}
	}
	return &result
}
