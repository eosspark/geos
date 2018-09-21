package types

import "github.com/eosspark/eos-go/common"

type GlobalPropertyObject struct {
	ID                       common.BlockIDType         `storm:"unique" json:"id"`
	ProposedScheduleBlockNum uint32                     `json:"proposed_schedule_block_num"`
	ProposedSchedule         SharedProducerScheduleType `json:"proposed_schedule"`
	Configuration            common.Config              //TODO
}

type DynamicGlobalPropertyObject struct {
	ID                   int64  `storm:"unique" json:"id"` //c++ chainbase.hpp id_type
	GlobalActionSequence uint64 `json:"global_action_sequence"`
}

type GlobalPropertyMultiIndex struct {
	GlobalPropertyObject
	ID int64		`storm:"unique" json:"id"`
}

type DynamicGlobalPropertyMultiIndex struct {
	DynamicGlobalPropertyObject
	ID int64	`storm:"unique" json:"id"`
}