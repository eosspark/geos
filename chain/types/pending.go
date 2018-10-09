package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
)

type PendingState struct {
	DBSeesion         *eosiodb.Session `json:"db_session"`
	PendingBlockState BlockState       `json:"pending_block_state"`
	Actions           []ActionReceipt  `json:"actions"`
	BlockStatus       BlockStatus      `json:"block_status"`
	ProducerBlockId   common.BlockIdType
}

//TODO wait modify Singleton
func NewPendingState(db *eosiodb.DataBase) *PendingState {
	pending := PendingState{}
	/*db, err := eosiodb.NewDatabase(config.DefaultConfig.BlockDir, "eos.db", true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:",err)
	}
	defer db.Close()*/
	session := db.StartSession()

	pending.DBSeesion = session
	return &pending
}
//TODO wait modify Singleton
func GetInstance() *PendingState {
	pending := PendingState{}
	db, err := eosiodb.NewDataBase(common.DefaultConfig.DefaultBlocksDirName, common.DefaultConfig.DBFileName, true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:", err)
	}
	defer db.Close()
	session := db.StartSession()
	if err != nil {
		fmt.Println(err.Error())
	}
	pending.DBSeesion = session
	return &pending
}

func Reset(pending *PendingState) {
	pending = nil
	log.Info("destory pending")
}
