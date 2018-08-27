package database

import (
	"testing"
	"fmt"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/chain/types"
)

func Test_NewForkDatabase(t *testing.T) {
	forkdb,err := NewForkDatabase("./","forkdb.dat",true)
	if err != nil {
		t.Error(err)
	}
	log.Debug("forkdb block state:",forkdb)
	fmt.Println("Test_NewForkDatabase run seccuss")
	defer forkdb.database.Close()
}

func Test_AddBlockState(t *testing.T){
	var blockState =types.BlockState{}

	forkdb,err := NewForkDatabase("./","forkdb.dat",true)
	if err != nil {
		t.Error(err)
	}

	b,er := forkdb.AddBlockState(blockState)
	if er!=nil{
		t.Error(er)
	}
	log.Debug("AddBlockState return info:",b)
}