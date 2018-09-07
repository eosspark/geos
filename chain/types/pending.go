package types

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/chain/config"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/rlp"
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

type ProducerScheduleType struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type SharedProducerScheduleType struct {
	Version   uint32
	Producers []ProducerKey
}

func NewPendingState(db eosiodb.Database) *PendingState {
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
	pending.DBSeesion = session
	pending.Valid = true
	return &pending
}

func GetInstance() *PendingState {
	pending := PendingState{}
	db, err := eosiodb.NewDatabase(config.DefaultBlocksDirName, config.DBFileName, true)
	if err != nil {
		log.Error("pending NewPendingState is error detail:", err)
	}
	defer db.Close()
	session, err := db.Start_Session()
	if err != nil {
		fmt.Println(err.Error())
	}
	pending.DBSeesion = session
	pending.Valid = false
	return &pending
}

func (bhs *BlockState) SetNewProducers(pending ProducerScheduleType) {
	if pending.Version == bhs.ActiveSchedule.Version+1 {
		log.Error("wrong producer schedule version specified")
		return
	}
	bhs.Header.NewProducers = pending
	tmp, _ := rlp.EncodeToBytes(bhs.Header.NewProducers)
	bhs.PendingScheduleHash = sigDigest(tmp)
	bhs.PendingSchedule = bhs.Header.NewProducers
	bhs.PendingScheduleLibNum = bhs.BlockNum
}

func sigDigest(p []byte) (id [4]uint64) {
	h := sha256.New()
	_, _ = h.Write(p)
	tmp := h.Sum(nil)
	for i := range id {
		id[i] = binary.LittleEndian.Uint64(tmp[i*8 : (i+1)*8])
	}
	log.Info("hash:", id)
	return
}

func (spst *SharedProducerScheduleType) Clear() {
	spst.Version = 0
	spst.Producers = []ProducerKey{}
}

func (spst *SharedProducerScheduleType) SharedroducerScheduleType(a ProducerScheduleType) *ProducerScheduleType {
	var result ProducerScheduleType = ProducerScheduleType{}
	spst.Version = a.Version
	spst.Producers = nil
	spst.Producers = a.Producers
	for i := 0; i < len(a.Producers); i++ {
		spst.Producers[i] = a.Producers[i]
	}
	return &result
}

func (spst *SharedProducerScheduleType) ProducerScheduleType() *ProducerScheduleType {
	var result ProducerScheduleType = ProducerScheduleType{}
	result.Version = spst.Version
	if len(result.Producers) == 0 {
		result.Producers = spst.Producers
	} else {
		var step int = len(result.Producers)
		for _, p := range spst.Producers {
			result.Producers[step] = p
			step++
		}
	}
	return &result
}

func Reset(pending *PendingState) {
	pending = nil
	log.Info("destory pending")
}
