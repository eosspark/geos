package unittests

import (
	"testing"

	//"github.com/eosspark/eos-go/chain"
	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/database"

	"github.com/stretchr/testify/assert"
)

func TestUndo(t *testing.T) {
	vt := newValidatingTester(true, chain.SPECULATIVE)
	db := vt.Control.DB
	sec := db.StartSession()

	billy := DbHouse{Id: uint64(1), Area: uint64(7), Name: "billy", Carnivore: Carnivore{Lion: 7, Tiger: 7}}

	db.Insert(&billy) //insert success

	// Make sure we can retrieve that account by name
	tmp := DbHouse{}
	db.Find("Name", billy, &tmp) //find success

	LogObj(tmp)
	//Undo creation of the account
	sec.Undo()

	// Make sure we can no longer find the account
	tmp2 := DbHouse{}
	db.Find("Name", billy, &tmp2) //can't find after undo
	LogObj(tmp2)

	vt.close()
}

// Test the block fetching methods on database, fetch_bock_by_id, and fetch_block_by_number
func TestGetBlocks(t *testing.T) {
	test := newValidatingTester(true, chain.SPECULATIVE)
	blockIds := make([]common.BlockIdType, 0)

	var NumOfBlocksToProd uint32 = 200
	// Produce 200 blocks and check their IDs should match the above
	test.ProduceBlocks(NumOfBlocksToProd, false)
	for i := 0; i < int(NumOfBlocksToProd); i++ {
		blockIds = append(blockIds, test.Control.FetchBlockByNumber(uint32(i+1)).BlockID())
		assert.Equal(t, types.NumFromID(&blockIds[len(blockIds)-1]), uint32(i+1))
		assert.Equal(t, test.Control.FetchBlockByNumber(uint32(i+1)).BlockID(), blockIds[len(blockIds)-1])
	}

	// Utility function to check expected irreversible block
	CalcExpLastIrrBlockNum := func(headBlockNum uint32) uint32 {
		producerSize := len(test.Control.HeadBlockState().ActiveSchedule.Producers)
		maxReversibleRounds := common.EosPercent(uint64(producerSize), uint32(common.DefaultConfig.Percent_100-common.DefaultConfig.IrreversibleThresholdPercent))
		if maxReversibleRounds == 0 {
			return headBlockNum
		} else {
			currentRound := headBlockNum / uint32(common.DefaultConfig.ProducerRepetitions)
			irreversibleRound := currentRound - uint32(maxReversibleRounds)
			return (irreversibleRound+1)*uint32(common.DefaultConfig.ProducerRepetitions) - 1
		}
	}
	// Check the last irreversible block number is set correctly
	expectedLastIrreversibleBlockNumber := CalcExpLastIrrBlockNum(NumOfBlocksToProd)
	assert.Equal(t, test.Control.HeadBlockState().DposIrreversibleBlocknum, expectedLastIrreversibleBlockNumber)
	// Check that block 201 cannot be found (only 20 blocks exist)
	assert.Equal(t, test.Control.FetchBlockByNumber(NumOfBlocksToProd+1+1), (*types.SignedBlock)(nil))

	var NextNumOfBlocksToProd uint32 = 100
	// Produce 100 blocks and check their IDs should match the above
	test.ProduceBlocks(NextNumOfBlocksToProd, false)

	nextExpectedLastIrreversibleBlockNumber := CalcExpLastIrrBlockNum(NumOfBlocksToProd + NextNumOfBlocksToProd)
	// Check the last irreversible block number is updated correctly
	assert.Equal(t, test.Control.HeadBlockState().DposIrreversibleBlocknum, nextExpectedLastIrreversibleBlockNumber)
	// Check that block 201 can now be found
	CheckNoThrow(t, func() { test.Control.FetchBlockByNumber(NumOfBlocksToProd + 1) })
	// Check the latest head block match
	assert.Equal(t, test.Control.FetchBlockByNumber(NumOfBlocksToProd+NextNumOfBlocksToProd+1).BlockID(), test.Control.HeadBlockId())

	test.close()

}

func Objects() ([]DbTableIdObject, []DbHouse) {

	objs := []DbTableIdObject{}
	DbHouses := []DbHouse{}
	for i := 1; i <= 3; i++ {
		number := i * 10
		obj := DbTableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 1), Payer: AccountName(number + 4 + i + 1), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house := DbHouse{Area: uint64(number + 7), Carnivore: Carnivore{Lion: number + 7, Tiger: number + 7}}
		DbHouses = append(DbHouses, house)
		obj = DbTableIdObject{Code: AccountName(number + 2), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 2), Payer: AccountName(number + 4 + i + 2), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = DbHouse{Area: uint64(number + 8), Carnivore: Carnivore{Lion: number + 8, Tiger: number + 8}}
		DbHouses = append(DbHouses, house)

		obj = DbTableIdObject{Code: AccountName(number + 3), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 3), Payer: AccountName(number + 4 + i + 3), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = DbHouse{Area: uint64(number + 9), Carnivore: Carnivore{Lion: number + 9, Tiger: number + 9}}
		DbHouses = append(DbHouses, house)
	}
	return objs, DbHouses
}
