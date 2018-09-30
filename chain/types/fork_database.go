package types

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/config"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
)

var isActive bool = false

type ForkDatabase struct {
	db      *eosiodb.DataBase
	Index   *ForkMultiIndexType `json:"index"`
	Head    *BlockState         `json:"head"`
	DataDir string

}
var IrreversibleBlock chan BlockState = make(chan BlockState)

type ForkMultiIndexType struct {
	ByBlockID common.BlockIdType `storm:"unique" json:"id"`
	ByPrev    common.BlockIdType `storm:"index"  json:"prev"`
	//Pair<block_num,in_current_chain>
	ByBlockNum common.Tuple `storm:"index"  json:"block_num"`
	//tuple<dpos_irreversible_blocknum,bft_irreversible_blocknum,block_num>
	ByLibBlockNum common.Tuple `storm:"index"  json:"lib_block_num"`
	BlockState    BlockState   `storm:"inline"`
}

func (f *ForkDatabase) setHead(head *BlockState) *ForkDatabase {
	if f.Head == nil {
		f.Head = head
	} else if f.Head.BlockNum < head.BlockNum {
		f.Head = head
	}
	return f
}

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := ForkDatabase{}
	if !isActive {
		forkd, err := newForkDatabase(stateDir, config.ForkDBName, true)
		if err != nil {
			log.Error("GetForkDbInstance is error ,detail:", err)
		}
		forkDB = *forkd
	}
	return &forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {
	//forkdb := &ForkDatabase{}

	db, err := eosiodb.NewDataBase(path, fileName, rw)
	if err != nil {
		log.Error("newForkDatabase is error:",err)
		return nil, err
	}

	/*var indexObj []ForkMultiIndexType
	err = db.ByIndex("ID", &indexObj)
	if err != nil {
		log.Error("new forkDatabase is error detail:", err)
	}
	var size int = len(indexObj)
	if size > 0 {
		for index, value := range indexObj {
			var indexType = value
			log.Debug("init fork database :", index)

			forkdb = forkdb.setHead(&indexType.BlockState)
		}
	}
	log.Debug("indexObj:", len(indexObj))
	isActive = true //set active is true
	if len(indexObj) > 0 {
		// TODO indexObj[0]
		return &ForkDatabase{db: db, Index: &indexObj[0], Head: forkdb.Head}, err
	} else {
		return &ForkDatabase{db: db}, err
	}*/
	fmt.Println(db)
	return &ForkDatabase{db: db}, err
}

func (f *ForkDatabase) set(s BlockState){

}
func (fdb *ForkDatabase) GetBlock(id common.BlockIdType) BlockState {
	//blockId   = fdb.Index.ID
	var blockState BlockState
	err := fdb.db.Find("ID", id, blockState)
	if err != nil {
		return blockState
	}
	return blockState
}

func (fdb *ForkDatabase) GetBlockByID(blockId common.BlockIdType) (*BlockState, error) {
	var indexObj ForkMultiIndexType
	err := fdb.db.Find("ID", blockId, &indexObj)
	if err != nil {
		return nil, err
	}
	return &indexObj.BlockState, nil
}

func (fdb *ForkDatabase) GetBlockByNum(blockNum uint32) (*BlockState, error) {
	var indexObj ForkMultiIndexType
	err := fdb.db.Get("BlockNum", blockNum, &indexObj)
	if err != nil {
		return nil, err
	}

	return &indexObj.BlockState, nil
}

func (fdb *ForkDatabase) AddBlockState(blockState BlockState) *BlockState {

	var index ForkMultiIndexType = ForkMultiIndexType{ByBlockID: blockState.ID,
		ByPrev:     blockState.SignedBlock.Previous,
		ByBlockNum: common.MakeTuple(blockState.BlockNum, true),
		BlockState: blockState}

	err := fdb.db.Insert(index)
	if err != nil {
		log.Error("AddBlockState is error for detail:", err)
	}

	var libHeaderObj []BlockState
	err = fdb.db.ByIndex("blockLibNum", &libHeaderObj)
	if err != nil {
		log.Error("AddBlockState find ByIndex is error for detail:", err)
	}
	if libHeaderObj != nil && len(libHeaderObj) > 0 {
		fdb.Head = &libHeaderObj[0]
	}
	var libNum = fdb.Head.DposIrreversibleBlocknum

	var headerObj []BlockState
	err = fdb.db.ByIndex("blockNum", &headerObj)
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
	return &blockState
}
func (fdb *ForkDatabase) AddSignedBlockState(signedBlcok *SignedBlock) *BlockState {
	blockId := signedBlcok.BlockID()
	var blockState BlockState
	err := fdb.db.Get("ID", blockId, &blockState)
	if err != nil {
		log.Error("AddSignedBlockState is error,detail:", err)
	}
	if &blockState != nil {
		err := fdb.db.Get("ID", signedBlcok.Previous, blockState)
		if err != nil {
			log.Error("AddSignedBlockState is error,detail:", err)
		}
	}
	block := fdb.AddBlockState(blockState)
	return block
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

func (fdb *ForkDatabase) FetchBranchFrom(first common.BlockIdType, second common.BlockIdType) error {
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
