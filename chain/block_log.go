package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"math"
	"os"
)

type BlockLog struct {
	nPos             uint64
	supportedVersion uint32

	head   *types.SignedBlock
	headId common.BlockIdType

	blockSteam *os.File
	blockFile  string
	//blockWrite bool

	indexStream *os.File
	indexFile   string
	//indexWrite bool

	genesisWriteToBlockLog bool
}

// func (b *BlockLog)checkBlockRead()  {
// 	if b.blockWrite {
// 		b.blockSteam.Close()
// 		b.blockSteam,_ = os.OpenFile(blockLog.blockFile, os.O_RDONLY)
// 		b.blockWrite = false
// 	}
// }
// func (b *BlockLog)checkBlockWrite()  {
// 	if b.blockWrite {
// 		b.blockSteam.Close()
// 		b.blockSteam,_ = os.OpenFile(blockLog.blockFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE)
// 		b.blockWrite = true
// 	}
// }
// func (b *BlockLog)checkIndexRead()  {
// 	if b.indexWrite {
// 		b.indexStream.Close()
// 		b.indexStream,_ = os.OpenFile(blockLog.indexFile, os.O_RDONLY)
// 		b.indexWrite = false
// 	}
// }
// func (b *BlockLog)checkIndexWrite()  {
// 	if b.indexWrite {
// 		b.indexStream.Close()
// 		b.indexStream,_ = os.OpenFile(blockLog.indexFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE)
// 		b.indexWrite = true
// 	}
// }

func NewBlockLog(dataDir string) *BlockLog {

	blockLog := &BlockLog{
		nPos:             math.MaxUint64,
		supportedVersion: 1,

		//blockWrite:true,
		//indexWrite:true,
	}

	_, err := os.Stat(dataDir)
	if err != nil {
		os.Mkdir(dataDir, os.ModePerm)
	}

	blockLog.blockFile = dataDir + "/blocks.log"
	blockLog.indexFile = dataDir + "/blocks.index"

	//if blockLog.blockSteam.IsOpen() {blockLog.blockSteam.Close()}
	//if blockLog.indexSteam.IsOpen() {blockLog.indexSteam.Close()}

	// blockLog.blockSteam,_ = os.OpenFile(blockLog.blockFile, os.O_WRONLY|O_APPEND|O_CREATE)
	// blockLog.indexStream,_ = os.Open(blockLog.indexFile, os.O_WRONLY|O_APPEND|O_CREATE)

	blockLog.blockSteam, _ = os.OpenFile(blockLog.blockFile, os.O_RDWR, os.ModePerm)
	blockLog.indexStream, _ = os.OpenFile(blockLog.indexFile, os.O_RDWR, os.ModePerm)

	logSize, _ := blockLog.blockSteam.Seek(0, 1)
	indexSize, _ := blockLog.indexStream.Seek(0, 1)

	if logSize > 0 {

		blockLog.blockSteam.Seek(0, os.SEEK_CUR)

	} else if indexSize > 0 {

	}

	return blockLog
}
func (b *BlockLog) Append(block *types.SignedBlock) uint64                                 { return 0 }
func (b *BlockLog) flush()                                                                 {}
func (b *BlockLog) ResetToGenesis(gs *types.GenesisState, benesisBlock *types.SignedBlock) { return }
func (b *BlockLog) ReadBlock(pos uint64) (*types.SignedBlock, uint64) {
	b.blockSteam.Seek(int64(pos), 1)
	return nil, 0
}
func (b *BlockLog) ReadBlockByNum(blockNum uint32) *types.SignedBlock       { return nil }
func (b *BlockLog) GetBlockPos(blockNum uint32) uint64                      { return 0 }
func (b *BlockLog) ReadHead() *types.SignedBlock                            { return nil }
func (b *BlockLog) Head() *types.SignedBlock                                { return nil }
func (b *BlockLog) ConstructIndex()                                         { return }
func (b *BlockLog) repairLog(dataDir string, truncateAtBlock uint32) string { return "" }
func (b *BlockLog) ExtractGenesisState(dataDir string) types.GenesisState   { return types.GenesisState{} }
