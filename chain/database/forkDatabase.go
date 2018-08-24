package database

import (
	eosiodb "github.com/db"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type ForkDatabase struct {
	database   eosiodb.Database
	Index forkMultiIndexType
}

func NewForkDatabase(path string, fileName string) (*ForkDatabase,error) {
	db, err := eosiodb.NewDatabase(path, fileName,true)
	if err != nil{
		return nil,err
	}
	return &ForkDatabase{database : *db,},nil
}

type forkMultiIndexType struct {
	Head             types.BlockState
	BlockHeaderState types.BlockHeaderState
}

func (fdb *ForkDatabase) GetBlock(blockId common.BlockIDType) /**BlockState*/ {
	blockId   = fdb.Index.BlockHeaderState.ID

	//fdb.database.ByIndex(blockId, to)
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