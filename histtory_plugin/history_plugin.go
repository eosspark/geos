package histtory_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
)

type AccountHistoryObject struct {
	ID                 types.IdType `storm:"id,increment"`
	Account            common.AccountName
	ActionSequenceNum  uint64
	AccountSequenceNum int32
}

type AccountHistoryIndex struct {
	Id                 types.IdType `storm:"id,increment"`
	Aho                *AccountHistoryObject
	ById               types.IdType       `storm:"unique"`
	ByAccount          common.AccountName `storm:"index"`
	ByAccountActionSeq struct {
		account            common.AccountName
		accountSequenceNum int32
	} `strom:"unique"`
}

type ActionHistoryObject struct {
	ID                types.IdType `storm:"id,increment"`
	ActionSequenceNum uint64
	PackedActionTrace *string
	BlockNum          uint32
	BlockTime         common.BlockTimeStamp
	TrxId             common.TransactionIDType
}

type ActionHistoryIndex struct {
	Id                  types.IdType `storm:"id ,increment"`
	Aho                 *ActionHistoryObject
	ByActionSequenceNum uint64       `storm:"index"`
	ById                types.IdType `storm:"unique"`
	condition           struct {
		trxId             common.TransactionIDType
		actionSequenceNum uint64
	} `storm:"unique"`
}

func AddAccountHistoryIndex(db *eosiodb.DataBase, aho *AccountHistoryObject) {
	index := AccountHistoryIndex{}
	index.Aho = aho
	//aho.ID=1
	index.ById = aho.ID
	index.ByAccount = aho.Account
	index.ByAccountActionSeq.account = aho.Account
	index.ByAccountActionSeq.accountSequenceNum = aho.AccountSequenceNum
	fmt.Println(index)
	err := db.Insert(&index)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var indexs []AccountHistoryIndex
	err = db.All(&indexs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(indexs)
	defer db.Close()
}

func GetAccountHistoryIndexById(db *eosiodb.DataBase, id types.IdType) {
	index := AccountHistoryIndex{}
	err := db.Find("ById", id, &index)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(index)
	defer db.Close()
}

func GetAccountHistoryIndexByNum(db *eosiodb.DataBase, num uint64) {
	index := AccountHistoryIndex{}
	err := db.Find("ByActionSequenceNum", num, &index)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(index)
	defer db.Close()
}

func GetAccountIndexByAccountActionSeq(db *eosiodb.DataBase, asNum int32, account common.AccountName) {

	index := AccountHistoryIndex{}
	index.ByAccountActionSeq.accountSequenceNum = asNum
	index.ByAccountActionSeq.account = account
	result := AccountHistoryIndex{}
	err := db.Find("condition", index.ByAccountActionSeq, &result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
	defer db.Close()
}

func GetAccountIndexByAccount(db *eosiodb.DataBase, account common.AccountName) *AccountHistoryIndex {
	index := AccountHistoryIndex{}
	index.ByAccount = account
	err := db.Find("ByAccount", account, index)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &index
}
