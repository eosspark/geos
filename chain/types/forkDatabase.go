package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
)

type ForkDatabase struct {
	database eosiodb.DataBase   `json:"database"`
	Index    ForkMultiIndexType `json:"index"`
	Head     BlockState   `json:"head"`
}

type ForkMultiIndexType struct {
	ID          common.BlockIDType `storm:"unique" json:"id"`
	Prev        common.BlockIDType `storm:"index"  json:"prev"`
	BlockNum    uint32             `storm:"index"  json:"block_num"`
	LibBlockNum uint32             `storm:"index"  json:"lib_block_num"`
	BlockState  BlockState   `storm:"inline"`
}

func setHead(forkdb ForkDatabase, head BlockState) *ForkDatabase {
	if &forkdb.Head == nil {
		forkdb.Head = head
	} else if forkdb.Head.BlockNum < head.BlockNum {
		forkdb.Head = head
	}
	return &forkdb
}

func NewForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {
	forkdb := new(ForkDatabase)
	/*
		_,err := os.Stat(path+fileName)
		if err != nil{

		}
	*/
	db, err := eosiodb.NewDataBase(path, fileName, rw)
	if err != nil {
		return nil, err
	}
	var indexObj []ForkMultiIndexType
	err = db.ByIndex("ID", &indexObj)
	if err != nil {
		log.Error("new forkDatabase is error detail:", err)
	}
	var size int = len(indexObj)
	if size > 0 {
		for index, value := range indexObj {
			var indexType = value
			log.Debug("init fork database :", index)

			forkdb = setHead(*forkdb, indexType.BlockState)
		}
	}
	log.Debug("indexObj:", len(indexObj))
	if len(indexObj) > 0 {
		// TODO indexObj[0]
		return &ForkDatabase{database: *db, Index: indexObj[0], Head: forkdb.Head}, err
	} else {
		return &ForkDatabase{database: *db}, err
	}
}

func (fdb *ForkDatabase) GetBlock(id common.BlockIDType) (BlockState, error) {
	//blockId   = fdb.Index.ID
	var blockState BlockState
	err := fdb.database.Find("ID", id, blockState)
	if err != nil {
		return blockState, err
	}
	return blockState, nil
}

func (fdb *ForkDatabase) GetBlockByID(blockId common.BlockIDType) (*BlockState, error) {
	var indexObj ForkMultiIndexType
	err := fdb.database.Find("ID", blockId, &indexObj)
	if err != nil {
		return nil, err
	}
	return &indexObj.BlockState, nil
}

func (fdb *ForkDatabase) GetBlockByNum(blockNum uint32) (*BlockState, error) {
	var indexObj ForkMultiIndexType
	err := fdb.database.Get("BlockNum", blockNum, &indexObj)
	if err != nil {
		return nil, err
	}

	return &indexObj.BlockState, nil
}

func (fdb *ForkDatabase) AddBlockState(blockState BlockState) (*BlockState, error) {

	var index ForkMultiIndexType = ForkMultiIndexType{ID: blockState.ID,
		Prev:       blockState.SignedBlock.Previous,
		BlockNum:   blockState.BlockNum,
		BlockState: blockState}

	err := fdb.database.Insert(index)
	if err != nil {
		log.Error("AddBlockState is error for detail:", err)
	}

	var libHeaderObj []BlockState
	err = fdb.database.ByIndex("blockLibNum", &libHeaderObj)
	if err != nil {
		log.Error("AddBlockState find ByIndex is error for detail:", err)
	}
	if libHeaderObj != nil && len(libHeaderObj) > 0 {
		fdb.Head = libHeaderObj[0]
	}
	var libNum = fdb.Head.DposIrreversibleBlocknum

	var headerObj []BlockState
	err = fdb.database.ByIndex("blockNum", &headerObj)
	if err != nil {
		log.Error("AddBlockState find ByIndex is error for detail:", err)
	}
	var oldBlock BlockState
	if headerObj != nil && len(headerObj) > 0 {
		oldBlock = headerObj[0]
	}
	var num = oldBlock.BlockNum

	if num < libNum {
		//TODO delete
	}
	//if fdb.BlockState.DposIrreversibleBlocknum <
	return &blockState, err
}
func (fdb *ForkDatabase) AddSignedBlockState(signedBlcok *SignedBlock) (*BlockState, error) {
	blockId:= signedBlcok.BlockID()
	var blockState BlockState
	err := fdb.database.Get("ID", blockId, &blockState)
	if err != nil {
		log.Error("AddSignedBlockState is error,detail:", err)
	}
	if &blockState != nil {
		err := fdb.database.Get("ID", signedBlcok.Previous, blockState)
		if err != nil {
			log.Error("AddSignedBlockState is error,detail:", err)
		}
	}
	block, er := fdb.AddBlockState(blockState)
	return block, er
}
func (fdb *ForkDatabase) Add(c HeaderConfirmation) {
	header, err := fdb.GetBlockByID(c.BlockId)
	if err != nil {
		log.Error("forkDatabase add header confirmation is error ,detail:", err)
	}
	fmt.Println(header)
	header.AddConfirmation(c) //TODO
}

type BranchType struct {
	branch []BlockState
}

func (fdb *ForkDatabase) FetchBranchFrom(first common.BlockIDType, second common.BlockIDType) error {
	//result := make(map[BranchType]BranchType)
	var firstBlock, secondBlock *BlockState
	firstBlock, er := fdb.GetBlockByID(first)
	if er != nil {
		log.Error("FetchBranchFrom is error for detail:", er)
	}
	secondBlock, err := fdb.GetBlockByID(second)
	if err != nil {
		log.Error("FetchBranchFrom is error for detail:", err)
	}
	for firstBlock.BlockNum > secondBlock.BlockNum {
	}

	for secondBlock.BlockNum > firstBlock.BlockNum {
	}

	for firstBlock.Header.Previous != secondBlock.Header.Previous {
	}

	return err
}

/*func main(){

	db,err := eosiodb.NewDatabase("./","test.mat",true)
	if err != nil{
		fmt.Println("test")
		return
	}
	defer db.Close()
	fmt.Print("test")
}*/
