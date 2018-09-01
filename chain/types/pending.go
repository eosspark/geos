package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"fmt"
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
	DBSeesion   eosiodb.Session  `json:"db_session"`
	PendingBlockState  BlockState `json:"pending_block_state"`
	Actions     []ActionReceipt  `json:"actions"`
	BlockStatus BlockStatus      `json:"block_status"`
	//IsExist		bool
}

//exec when start block
func NewPendingState(db eosiodb.Database) *PendingState{
	pending := PendingState{}
	/*db, err := eosiodb.NewDatabase(config.DefaultConfig.BlockDir, "eos.db", true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:",err)
	}
	defer db.Close()*/
	session, err := db.Start_Session()
	if err != nil {
		fmt.Println(err.Error())
	}
	pending.DBSeesion = *session
	return &pending
}

func Reset(){

}
