package types

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/config"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
)

type ActionReceipt struct {
	Receiver       common.AccountName            `json:"receiver"`
	ActDigest      common.SHA256Bytes            `json:"act_digest"`
	GlobalSequence uint64                        `json:"global_sequence"`
	RecvSequence   uint64                        `json:"recv_sequence"`
	AuthSequence   map[common.AccountName]uint64 `json:"auth_sequence"`
	CodeSequence   uint32                        `json:"code_sequence"` //TODO
	ABISequence    uint32                        `json:"abi_sequence"`
}
type PendingState struct {
	DBSeesion         *eosiodb.Session `json:"db_session"`
	PendingBlockState BlockState       `json:"pending_block_state"`
	Actions           []ActionReceipt  `json:"actions"`
	BlockStatus       BlockStatus      `json:"block_status"`
	Valid             bool             `json:"valid"`
}

func NewPendingState(db *eosiodb.DataBase) *PendingState {
	pending := PendingState{}
	/*db, err := eosiodb.NewDatabase(config.DefaultConfig.BlockDir, "eos.db", true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:",err)
	}
	defer db.Close()*/
	session := db.StartSession()

	pending.DBSeesion = session
	pending.Valid = true
	return &pending
}

func GetInstance() *PendingState {
	pending := PendingState{}
	db, err := eosiodb.NewDataBase(config.DefaultBlocksDirName, config.DBFileName, true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:", err)
	}
	defer db.Close()
	session := db.StartSession()
	if err != nil {
		fmt.Println(err.Error())
	}
	pending.DBSeesion = session
	pending.Valid = false
	return &pending
}

func Reset(pending *PendingState) {
	pending = nil
	log.Info("destory pending")
}
