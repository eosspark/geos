package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type ProducerKey struct {
	ProducerName    common.AccountName `json:"producer_name"`
	BlockSigningKey ecc.PublicKey      `json:"block_signing_key"`
}

type ProducerScheduleType struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

func (ps *ProducerScheduleType) GetProducerKey(p common.AccountName) (ecc.PublicKey, bool) {
	for _, i := range ps.Producers {
		if i.ProducerName == p {
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

func (s *SharedProducerScheduleType) SharedProducerScheduleType(a ProducerScheduleType) *SharedProducerScheduleType {
	s.Version = a.Version
	s.Producers = nil
	for i := 0; i < len(a.Producers); i++ {
		s.Producers = append(s.Producers, a.Producers[i])
	}
	return s
}

func (s *SharedProducerScheduleType) ProducerScheduleType() *ProducerScheduleType {
	var result = ProducerScheduleType{}
	result.Version = s.Version
	if len(result.Producers) == 0 {
		result.Producers = s.Producers
	} else {
		var step = len(result.Producers)
		for _, p := range s.Producers {
			result.Producers = append(result.Producers, p)
			step++
		}
	}
	return &result
}

func (p SharedProducerScheduleType) IsEmpty() bool {
	return p.Version == 0 && len(p.Producers) == 0
}

func (p ProducerScheduleType) IsEmpty() bool {
	return p.Version == 0 && len(p.Producers) == 0
}
