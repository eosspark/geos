package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/log"
)

var isFdActive bool = false

type ForkDatabase struct {
	db      *database.LDataBase
	Index   *ForkMultiIndexType `json:"index"`
	Head    *BlockState         `json:"head"`
	DataDir string
}

var IrreversibleBlock chan BlockState = make(chan BlockState)

type ForkMultiIndexType struct {
	ByBlockID  common.BlockIdType `storm:"unique" json:"id"`
	ByPrev     common.BlockIdType `storm:"index"  json:"prev"`
	ByBlockNum common.Tuple       `storm:"index"  json:"block_num"`
	BlockState BlockState         `storm:"inline"`
}

func (self *ForkDatabase) setHead(head *BlockState) *ForkDatabase {
	if self.Head == nil {
		self.Head = head
	} else if self.Head.BlockNum < head.BlockNum {
		self.Head = head
	}
	return self
}

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := ForkDatabase{}
	if !isFdActive {
		forkd, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
		if err != nil {
			log.Error("GetForkDbInstance is error ,detail:", err)
		}
		forkDB = *forkd
		isFdActive = true
	}
	return &forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {
	//forkdb := &ForkDatabase{}

	db, err := database.NewDataBase(path + "/" + fileName)
	if err != nil {
		log.Error("newForkDatabase is error:", err)
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

func (self *ForkDatabase) GetBlock(id *common.BlockIdType) *BlockState {
	//blockId   = fdb.Index.ID
	blockState := BlockState{}
	/*blockState.ID = id
	err := fdb.db.Find("ID", blockState)
	if err != nil {
		return &blockState
	}*/
	return &blockState
}

func (self *ForkDatabase) GetBlockByID(blockId common.BlockIdType) (*BlockState, error) {
	indexObj := ForkMultiIndexType{}
	/*err := fdb.db.Find("ID", blockId, &indexObj)
	if err != nil {
		return nil, err
	}*/
	return &indexObj.BlockState, nil
}

func (self *ForkDatabase) GetBlockByNum(blockNum uint32) (*BlockState, error) {
	indexObj := ForkMultiIndexType{}
	/*err := fdb.db.Get("BlockNum", blockNum, &indexObj)
	if err != nil {
		return nil, err
	}
	*/
	return &indexObj.BlockState, nil
}

func (self *ForkDatabase) AddBlockState(blockState BlockState) *BlockState {

	index := ForkMultiIndexType{ByBlockID: blockState.ID,
		ByPrev:     blockState.SignedBlock.Previous,
		ByBlockNum: common.MakeTuple(blockState.BlockNum, true),
		BlockState: blockState}

	err := self.db.Insert(index)
	if err != nil {
		log.Error("AddBlockState is error for detail:", err)
	}

	/*libHeaderObj :=[]BlockState{}
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
	}*/
	//if fdb.BlockState.DposIrreversibleBlocknum <
	return &blockState
}
func (self *ForkDatabase) AddSignedBlockState(signedBlcok *SignedBlock) *BlockState {
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
	block := self.AddBlockState(blockState)
	return block
}
func (self *ForkDatabase) Add(c HeaderConfirmation) {
	header, err := self.GetBlockByID(c.BlockId)
	if err != nil {
		log.Error("forkDatabase add header confirmation is error ,detail:", err)
	}
	fmt.Println(header)
	header.AddConfirmation(c) //TODO
}
func (self *ForkDatabase) Header() *BlockState { return self.Head }

type BranchType struct {
	branch []BlockState
}

func (self *ForkDatabase) FetchBranchFrom(first common.BlockIdType, second common.BlockIdType) error {
	//result := make(map[BranchType]BranchType)
	var firstBlock, secondBlock *BlockState
	firstBlock, er := self.GetBlockByID(first)
	if er != nil {
		log.Error("FetchBranchFrom is error for detail:", er)
	}
	secondBlock, err := self.GetBlockByID(second)
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

func (self *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *BlockState {
	b := BlockState{}
	b.BlockNum = n
	//TODO wait append
	//numIdx := fdb.db.Find("ByBlockNum",b)
	return &b
}

func (self *ForkDatabase) Remove(id *common.BlockIdType) {}

func (self *ForkDatabase) SetValidity(h *BlockState, valid bool) {
	if !valid {
		self.Remove(&h.ID)
	} else {
		h.Validated = true
	}
}
func (self *ForkDatabase) MarkInCurrentChain(b *BlockState, inCurrentChain bool) {}

func (self *ForkDatabase) Prune(b *BlockState) {}

func (self *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {}

/*func main(){

	db,err := eosiodb.NewDatabase("./","test.mat",true)
	if err != nil{
		fmt.Println("test")
		return
	}
	defer db.Close()
	fmt.Print("test")
}*/
