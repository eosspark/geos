package history_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/chain_plugin"
)

type AccountHistoryObject struct {
	ID                 types.IdType       `storm:"id,increment" json:"id"`
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
	ID                types.IdType `storm:"id,increment"`
	ActionSequenceNum uint64       `storm:"by_action_sequence_num unique"`
	PackedActionTrace common.HexBytes
	BlockNum          uint32
	BlockTime         common.BlockTimeStamp
	TrxId             common.TransactionIDType
	//c++ no param
	ByTrxId struct {
		TrxId             common.TransactionIDType
		ActionSequenceNum uint64
	} `storm:"unique"`
}

func AddAccountHistoryObject(db *eosiodb.DataBase, aho *AccountHistoryObject) {
	//fmt.Println("params:", aho)
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
	//fmt.Println("GetAccountHistoryObjectById result :", result)
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
	//fmt.Println("GetAccountByAccActSeq result:", result)
}

func GetActionHistoryObjectByActSeqNum(db *eosiodb.DataBase, asNum uint64) *ActionHistoryObject {
	result := ActionHistoryObject{}
	err := db.Find("ActionSequenceNum", asNum, &result)
	defer db.Close()
	if err != nil {
		fmt.Println("GetActionHistoryObjectByActSeqNum is error detail :", err.Error())
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
	//fmt.Println("result:", result)
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

type GetActionParam struct {
	AccountName common.AccountName
	Pos 	int32
	Offset  int32
}

type OrderedActionResult struct {
	GlobalActionSeq		uint64
	AccountActionSeq 	int32
	BlockNum 			uint32
	BlockTime 			common.BlockTimeStamp
	ActionTrace 		types.ActionTrace
}
type GetActionResult struct {
	Actions []OrderedActionResult
	LastIrreversibleBlock	uint32
	TimeLimitExceededError	bool
}
func GetActions(params *GetActionParam) *[]GetActionResult{
	//control := chain.GetControlInstance()
	//db := control.DB()
	//chain.get_abi_serializer_max_time
	abiSerializerMaxTime := chain_plugin.GetInstance().GetAbiSerializerMaxTime()

	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryObject is error detail:", err.Error())
	}
	defer db.Close()

	ahoList := *GetAccountHistoryObjectByAccount(db,params.AccountName)
	var pageNum, start, end int32
	start = params.Pos
	if start == 0{
		end = params.Offset		//data count per page
	}
	//end = params.Offset+5
	if len(ahoList)>0{
		pageNum = int32(len(ahoList))/end
	}
	//db.Find("ByAccountActionSeq",AccountHistoryObject,[]AccountHistoryObject)

	array := ahoList[0:3]
	fmt.Println(array)

	fmt.Println(abiSerializerMaxTime,pageNum,start,end)

	//return &array
	return nil
}
