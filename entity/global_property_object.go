package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type GlobalPropertyObject struct {
	ID                       common.IdType              `multiIndex:"id,increment"`
	ProposedScheduleBlockNum uint32
	ProposedSchedule         types.SharedProducerScheduleType
	Configuration            common.Config                    //TODO
}

type DynamicGlobalPropertyObject struct {
	ID                   common.IdType  `multiIndex:"id,increment"` //c++ chainbase.hpp id_type
	GlobalActionSequence uint64
}