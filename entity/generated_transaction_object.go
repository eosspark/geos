package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"math/big"
)

type GeneratedTransaction struct {
	TrxId      common.TransactionIdType
	Sender     common.AccountName
	SenderId   big.Int //c++ uint128_t
	Payer      common.AccountName
	DelayUntil common.TimePoint
	Expiration common.TimePoint
	Published  common.TimePoint
	PackedTrx  []byte
}

type GeneratedTransactionObject struct {
	Id         common.IdType             `storm:"id,increment"`
	TrxId      common.TransactionIdType `storm:"unique"`
	Sender     common.AccountName
	SenderId   big.Int //c++ uint128_t
	Payer      common.AccountName
	DelayUntil common.TimePoint
	Expiration common.TimePoint
	Published  common.TimePoint
	PackedTrx  common.HexBytes //c++ shared_string
	/*expiration、Id*/
	ByExpiration common.Tuple
	/*DelayUntil、Id*/
	ByDelay common.Tuple
	/*Sender、SenderId*/
	BySenderId common.Tuple
}

/* c++
type billable_size struct {
	const uint64_t overhead = overhead_per_row_per_index_ram_bytes * 5  ///< overhead for 5x indices internal-key, txid, expiration, delay, sender_id
	const uint64_t value = 96 + 4 + overhead ///< 96 bytes for our constant size fields, 4 bytes for a varint for packed_trx size and 96 bytes of implementation overhead
}*/

//const overhead_per_row_per_index_ram_bytes uint32 = 32

func (g *GeneratedTransactionObject) GetBillableSize() uint64 {
	overhead := overhead_per_row_per_index_ram_bytes * 5
	value := 96 + 4 + overhead
	return uint64(value)
}

func GetGTOByTrxId(db *database.DataBase, trxId common.TransactionIdType) *GeneratedTransactionObject {
	gto := GeneratedTransactionObject{}
	//err := db.Find("TrxId", trxId, gto)
	//if err != nil {
	//	fmt.Println(GetGTOByTrxId)
	//}
	return &gto
}

func GetGeneratedTransactionObjectByExpiration(db *database.DataBase, be common.Tuple) *GeneratedTransactionObject {
	gto := GeneratedTransactionObject{}
	//err := db.Find("ByExpiration", be, &gto)
	//if err != nil {
	//	fmt.Println("GetGeneratedTransactionObjectByExpiration is error :", err.Error())
	//}
	return &gto
}

func GetGeneratedTransactionObjectByDelay(db *database.DataBase, be common.Tuple) *GeneratedTransactionObject {
	gto := GeneratedTransactionObject{}
	//err := db.Find("ByDelay", be, &gto)
	//if err != nil {
	//	fmt.Println("GetGeneratedTransactionObjectByDelay is error :", err.Error())
	//}
	return &gto
}

func GetGeneratedTransactionObjectBySenderId(db *database.DataBase, be common.Tuple) *GeneratedTransactionObject {
	gto := GeneratedTransactionObject{}
	//err := db.Find("BySenderId", be, &gto)
	//if err != nil {
	//	fmt.Println("GetGeneratedTransactionObjectBySenderId is error :", err.Error())
	//}
	return &gto
}

func GeneratedTransactions(gto *GeneratedTransactionObject) *GeneratedTransaction {
	gt := GeneratedTransaction{}
	gt.TrxId = gto.TrxId
	gt.Sender = gto.Sender
	gt.SenderId = gto.SenderId
	gt.Payer = gto.Payer
	gt.DelayUntil = gto.DelayUntil
	gt.Expiration = gto.Expiration
	gt.Published = gto.Published
	gt.PackedTrx[0] = gto.PackedTrx[0]
	gt.PackedTrx[1] = gto.PackedTrx[len(gto.PackedTrx)]
	return &gt
}
