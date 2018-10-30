package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
)

var PendingValid bool //singleton

type PendingState struct {
	//MaybeSession      *database.Session `json:"db_session"`
	DbSession         *MaybeSession         `json:"db_session"`
	PendingBlockState *types.BlockState     `json:"pending_block_state"`
	Actions           []types.ActionReceipt `json:"actions"`
	BlockStatus       types.BlockStatus     `json:"block_status"`
	ProducerBlockId   common.BlockIdType
}

//TODO wait modify Singleton
func NewPendingState(db database.DataBase) *PendingState {
	PendingValid = true
	pending := PendingState{}
	pending.DbSession = NewMaybeSession(db)
	return &pending
}

func (p *PendingState) Reset() {
	if PendingValid {

		p = nil
	}
}

func (p *PendingState) Push() {
	p.DbSession.Push()
}
