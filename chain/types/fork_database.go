package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/exception"
)

var isFdActive bool = false

type ForkDatabase struct {
	DB      database.DataBase
	//Index   *ForkMultiIndexType `json:"index"`
	Head    *BlockState         `json:"head"`
	//DataDir string
}

//var IrreversibleBlock chan BlockState = make(chan BlockState)

/*type ForkMultiIndexType struct {
	ByBlockID  common.BlockIdType `storm:"unique" json:"id"`
	ByPrev     common.BlockIdType `storm:"index"  json:"prev"`
	ByBlockNum common.Tuple       `storm:"index"  json:"block_num"`
	BlockState *BlockState         `storm:"inline"`
}*/

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := ForkDatabase{}
	if !isFdActive {
		forkd, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
		if err != nil {
			log.Error("GetForkDbInstance is error ,detail:", err)
		}
		forkd.DB.Insert("test")
		forkDB = *forkd
		//
		isFdActive = true
	}
	return &forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {

	db, err := database.NewDataBase(path + "/" + fileName)
	if err != nil {
		log.Error("newForkDatabase is error:", err)
		return nil, err
	}

	return &ForkDatabase{DB: db}, err
}

func (f *ForkDatabase) SetHead(s *BlockState) {
	if f.Head == nil {
		f.Head = s
	} else if f.Head.BlockNum < s.BlockNum {
		f.Head = s
	}
	exception.EosAssert( s.ID == s.Header.BlockID(), &exception.ForkDatabaseException{},
		"block state id:s%, is different from block state header id:s%", s.ID,s.Header.BlockID() );

	err:=f.DB.Insert(s)
	if err != nil{
		fmt.Println("ForkDB SetHead is error:",err)
	}

}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *BlockState {
	//blockId   = fdb.Index.ID
	blockState := BlockState{}
	blockState.ID = *id
	err := f.DB.Find("ID",blockState, blockState)
	if err != nil {
		fmt.Println("ForkDB GetBlock is error:",err)
	}
	return &blockState
}

func (f *ForkDatabase) AddBlockState(blockState *BlockState) *BlockState {

	/*lib := f.Head.DposIrreversibleBlocknum
	itr,err:=f.db.Get("",&BlockState{})
	if err != nil{
		//exception.EosAssert()
		fmt.Println("ForkDB AddBlockState Get Is Error:",err)
	}
	oldest:=  itr.First()
	oldest.*/
	/*err := f.db.Insert(blockState)
	if err != nil {
		log.Error("AddBlockState is error for detail:", err)
	}
	param := BlockState{}

	f.Head = f.db.GetObjects("ByLibBlockNum",)*/

	return blockState
}
func (f *ForkDatabase) AddSignedBlockState(signedBlcok *SignedBlock,trust bool) *BlockState {
	//blockId := signedBlcok.BlockID()
	blockState := BlockState{}
	/*err := fdb.db.Get("ID", blockId, &blockState)
	if err != nil {
		log.Error("AddSignedBlockState is error,detail:", err)
	}
	if &blockState != nil {
		err := fdb.db.Get("ID", signedBlcok.Previous, blockState)
		if err != nil {
			log.Error("AddSignedBlockState is error,detail:", err)
		}
	}*/
	block := f.AddBlockState(&blockState)
	return block
}
func (f *ForkDatabase) Add(c *HeaderConfirmation) {
	header:= f.GetBlock(&c.BlockId)

	exception.EosAssert( header!=nil, &exception.ForkDbBlockNotFound{}, "unable to find block id ",c.BlockId)
	header.AddConfirmation(c)

	if header.BftIrreversibleBlocknum<header.BlockNum && len(header.Confirmations)>= ((len(header.ActiveSchedule.Producers)*2)/3+1){
		f.SetBftIrreversible(c.BlockId)
	}
}
func (f *ForkDatabase) Header() *BlockState { return f.Head }

type BranchType struct {
	branch []BlockState
}

func (f *ForkDatabase) FetchBranchFrom(first common.BlockIdType, second common.BlockIdType)  {
	//result := make(map[BranchType]BranchType)
	//var firstBlock, secondBlock *BlockState
	firstBlock:= f.GetBlock(&first)

	secondBlock := f.GetBlock(&second)

	for firstBlock.BlockNum > secondBlock.BlockNum {
	}

	for secondBlock.BlockNum > firstBlock.BlockNum {
	}

	for firstBlock.Header.Previous != secondBlock.Header.Previous {
	}

	fmt.Println(firstBlock,secondBlock)
	//return err
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *BlockState {
	b := BlockState{}
	b.BlockNum = n
	//TODO wait append
	//numIdx := fdb.db.Find("ByBlockNum",b)
	return &b
}

func (f *ForkDatabase) Remove(id *common.BlockIdType) {}

func (f *ForkDatabase) SetValidity(h *BlockState, valid bool) {
	if !valid {
		f.Remove(&h.ID)
	} else {
		h.Validated = true
	}
}
func (f *ForkDatabase) MarkInCurrentChain(b *BlockState, inCurrentChain bool) {}

func (f *ForkDatabase) Prune(b *BlockState) {}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {

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
