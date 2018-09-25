package histtory_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
)

type AccountHistoryObject struct {
	ID                 types.IdType       `storm:"id,increment" json:"ById"`
	Account            common.AccountName `storm:"index"`
	ActionSequenceNum  uint64
	AccountSequenceNum int32
	ByAccountActionSeq ByAccountActionSeq `storm:"index"`
	/*ByAccountActionSeq struct {
		Account            common.AccountName
		AccountSequenceNum int32
	}`storm:"index"`*/
}
/*search condition*/
type ByAccountActionSeq struct {
	Account            common.AccountName
	AccountSequenceNum int32
}

type ActionHistoryObject struct {
	ID                 types.IdType `storm:"id,increment"`
	ActionSequenceNum  uint64       `storm:"by_action_sequence_num unique"`
	PackedActionTrace  *string
	BlockNum           uint32
	BlockTime          common.BlockTimeStamp
	TrxId              common.TransactionIDType
	//c++ no param
	ByTrxId struct {
		TrxId             common.TransactionIDType
		ActionSequenceNum uint64
	} `storm:"unique"`
}

func AddAccountHistoryObject(db *eosiodb.DataBase, aho *AccountHistoryObject) {

	fmt.Println("params:", aho)
	err := db.Insert(&aho)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func GetAccountHistoryObjectById(db *eosiodb.DataBase, id types.IdType) {
	result := AccountHistoryObject{}
	err := db.Find("ById", id, &result)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GetAccountHistoryObjectById result :", result)
}

func GetAccountByAccountActionSeq(db *eosiodb.DataBase, asNum int32, account common.AccountName) {

	aho := AccountHistoryObject{}
	aho.ByAccountActionSeq.AccountSequenceNum = asNum
	aho.ByAccountActionSeq.Account = account
	result := AccountHistoryObject{}
	err := db.Find("condition", aho.ByAccountActionSeq, &result)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GetAccountByAccActSeq result:", result)
}

func GetActionHistoryObjectByActSeqNum(db *eosiodb.DataBase, asNum uint64) *ActionHistoryObject{
	result := ActionHistoryObject{}
	err := db.Find("ActionSequenceNum",asNum,&result)
	defer db.Close()
	if err !=nil {
		fmt.Println("GetActionHistoryObjectByActSeqNum is error detail :",err.Error())
	}
	return &result
}

func GetActionHistoryByAccActSeq(db *eosiodb.DataBase, asNum uint64, trxId common.TransactionIDType) {
	result := ActionHistoryObject{}
	result.ByTrxId.TrxId = trxId
	result.ByTrxId.ActionSequenceNum = asNum
	err := db.Find("ByAccountActionSeq", result.ByTrxId, &result)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("result:", result)

}

func GetAccountHistoryObjectByAccount(db *eosiodb.DataBase, account common.AccountName) *[]AccountHistoryObject {
	aho := AccountHistoryObject{}
	aho.Account = account
	//aho.ByAccountActionSeq.AccountSequenceNum = 4
	fmt.Println("+++++++++++++++++:", aho.ByAccountActionSeq.Account)
	var result []AccountHistoryObject
	err := db.Get("Account", aho.Account, &result)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//fmt.Print("result :", result)
	return &result
}
