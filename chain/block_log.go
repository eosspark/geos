package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"math"
	"os"
)

type BlockLog struct {
	head   *types.SignedBlock
	headId common.BlockIdType

	blockStream *os.File
	blockFile   string
	//blockWrite bool

	indexStream *os.File
	indexFile   string
	//indexWrite bool

	genesisWriteToBlockLog bool
}

/// static field
const (
	nPos             = math.MaxUint64
	supportedVersion = uint32(1)
)

/// sizeof
const (
	SizeOfInt32 = 4
	SizeOfInt64 = 8
)

/// seek whence
const (
	beg = 0
	cur = 1
	end = 2
)

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
		//blockWrite:true,
		//indexWrite:true,
	}

	_, err := os.Stat(dataDir)
	if err != nil {
		os.Mkdir(dataDir, os.ModePerm)
	}

	blockLog.blockFile = dataDir + "/blocks.log"
	blockLog.indexFile = dataDir + "/blocks.index"

	blockLog.blockStream, _ = os.OpenFile(blockLog.blockFile, os.O_RDWR, os.ModePerm)
	blockLog.indexStream, _ = os.OpenFile(blockLog.indexFile, os.O_RDWR, os.ModePerm)

	if blockLog.blockStream == nil {
		blockLog.blockStream, _ = os.Create(blockLog.blockFile)
		blockLog.blockStream.Close()
		blockLog.blockStream, _ = os.OpenFile(blockLog.blockFile, os.O_RDWR, os.ModePerm)
	}

	if blockLog.indexStream == nil {
		blockLog.indexStream, _ = os.Create(blockLog.indexFile)
		blockLog.indexStream.Close()
		blockLog.indexStream, _ = os.OpenFile(blockLog.indexFile, os.O_RDWR, os.ModePerm)
	}

	logSize, _ := blockLog.blockStream.Seek(0, end)
	indexSize, _ := blockLog.indexStream.Seek(0, end)

	if logSize > 0 {

		blockLog.blockStream.Seek(0, 0)
		var version uint32 = 0
		bytes := make([]byte, SizeOfInt32)
		blockLog.blockStream.Read(bytes)
		rlp.DecodeBytes(bytes, &version)

		EosAssert(version > 0, &BlockLogAppendFail{}, "Block log was not setup properly with genesis information.")
		EosAssert(version == supportedVersion,
			&BlockLogUnsupportedVersion{},
			"Unsupported version of block log. Block log version is %d while code supports version %d",
			version, supportedVersion)

		blockLog.genesisWriteToBlockLog = true
		blockLog.head = blockLog.ReadHead()
		blockLog.headId = blockLog.head.BlockID()

		if indexSize > 0 {
			var blockPos int64 = 0
			bytes = make([]byte, SizeOfInt64)
			blockLog.blockStream.Seek(-SizeOfInt64, end) //sizeof(blockPos)
			blockLog.blockStream.Read(bytes)
			rlp.DecodeBytes(bytes, &blockPos)

			var indexPos int64 = 0
			bytes = make([]byte, SizeOfInt64)
			blockLog.indexStream.Seek(-SizeOfInt64, end) //sizeof(blockPos)
			blockLog.indexStream.Read(bytes)
			rlp.DecodeBytes(bytes, &indexPos)

			if blockPos < indexPos {
				//ilog("block_pos < index_pos, close and reopen index_stream")
				blockLog.ConstructIndex()
			} else if blockPos > indexPos {
				//ilog("Index is incomplete")
				blockLog.ConstructIndex()
			}
		} else {
			//ilog("Index is empty")
			blockLog.ConstructIndex()
		}
	} else if indexSize > 0 {
		//ilog("Index is nonempty, remove and recreate it")
		blockLog.indexStream.Close()
		os.Remove(blockLog.indexFile)
		blockLog.indexStream, _ = os.OpenFile(blockLog.indexFile, os.O_RDWR, os.ModePerm)
	}

	return blockLog
}
func (b *BlockLog) Append(block *types.SignedBlock) uint64 {

	EosAssert(b.genesisWriteToBlockLog, &BlockLogAppendFail{}, "Cannot append to block log until the genesis is first written")

	pos, err := b.blockStream.Seek(0, end)
	Throw(err)

	indexPos, err := b.indexStream.Seek(0, end)
	Throw(err)

	EosAssert(indexPos == int64(SizeOfInt64*(block.BlockNumber()-1)), &BlockLogAppendFail{},
		"Append to index file occuring at wrong position. position %d expected %d", indexPos, SizeOfInt64*(block.BlockNumber()-1))

	data, _ := rlp.EncodeToBytes(block)

	size := uint32(len(data))
	bytes, _ := rlp.EncodeToBytes(&size)
	b.blockStream.Write(bytes)
	b.blockStream.Write(data)

	posData, _ := rlp.EncodeToBytes(pos)
	b.blockStream.Write(posData)
	b.indexStream.Write(posData)

	b.head = block
	b.headId = block.BlockID()

	return uint64(pos)
}
func (b *BlockLog) flush() {
	b.blockStream.Sync()
	b.indexStream.Sync()
}
func (b *BlockLog) ResetToGenesis(gs *types.GenesisState, benesisBlock *types.SignedBlock) uint64 {
	var err error

	if b.blockStream != nil {
		b.blockStream.Close()
	}

	if b.indexStream != nil {
		b.indexStream.Close()
	}

	os.Remove(b.blockFile)
	os.Remove(b.indexFile)

	b.blockStream, err = os.Create(b.blockFile)
	b.indexStream, err = os.Create(b.indexFile)
	Throw(err)

	version := uint32(0) // version of 0 is invalid; it indicates that the genesis was not properly written to the block log
	bytes, _ := rlp.EncodeToBytes(version)
	b.blockStream.Write(bytes)

	bytes, _ = rlp.EncodeToBytes(gs)

	size := uint32(len(bytes))
	sizeBytes, _ := rlp.EncodeToBytes(&size)
	b.blockStream.Write(sizeBytes)
	b.blockStream.Write(bytes)
	b.genesisWriteToBlockLog = true

	ret := b.Append(benesisBlock)

	bytes, _ = rlp.EncodeToBytes(supportedVersion)
	b.blockStream.Write(bytes)

	b.flush()

	return ret
}
func (b *BlockLog) ReadBlock(pos uint64) (*types.SignedBlock, uint64) {
	//    s := size
	// if s == 0 {//the last block (head)
	// 	s = b.blockSteam.Seek(0,2)
	// }
	b.blockStream.Seek(int64(pos), beg)

	sizeBytes := make([]byte, SizeOfInt32)
	b.blockStream.Read(sizeBytes)
	var size uint32
	rlp.DecodeBytes(sizeBytes, &size)

	signedBlock := &types.SignedBlock{}
	bytes := make([]byte, size)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, signedBlock)

	var nextPos uint64
	bytes = make([]byte, SizeOfInt64)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &nextPos)

	return signedBlock, nextPos
}
func (b *BlockLog) ReadBlockByNum(blockNum uint32) *types.SignedBlock {
	returning, block := false, (*types.SignedBlock)(nil)
	Try(func() {
		pos := b.GetBlockPos(blockNum)
		if pos != nPos {
			block, _ = b.ReadBlock(pos)
			EosAssert(block.BlockNumber() == blockNum, &ReversibleBlocksException{}, "Wrong block was read from block log.")
		}
		returning = true
	}).FcLogAndRethrow()

	if returning {
		return block
	}

	return nil
}

func (b *BlockLog) ReadBlockById(id *common.BlockIdType) *types.SignedBlock {
	return b.ReadBlockByNum(types.NumFromID(id))
}

func (b *BlockLog) GetBlockPos(blockNum uint32) uint64 {

	if !(b.head != nil && blockNum <= types.NumFromID(&b.headId) && blockNum > 0) {
		return nPos
	}

	var pos uint64
	bytes := make([]byte, SizeOfInt64)
	b.indexStream.Seek(SizeOfInt64*(int64(blockNum)-1), beg)
	b.indexStream.Read(bytes)
	rlp.DecodeBytes(bytes, &pos)

	return pos //, nextPos - pos
}
func (b *BlockLog) ReadHead() *types.SignedBlock {

	s, _ := b.blockStream.Seek(0, end)
	if s <= SizeOfInt64 {
		return nil
	}

	var pos uint64
	bytes := make([]byte, end)
	b.blockStream.Seek(-SizeOfInt64, end)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &pos)

	block, _ := b.ReadBlock(pos)
	return block
}
func (b *BlockLog) Head() *types.SignedBlock {
	return b.head
}

func (b *BlockLog) ConstructIndex() {
	//ilog("Reconstructing Block Log Index...")
	b.indexStream.Close()
	os.Remove(b.indexFile)
	b.indexStream, _ = os.Create(b.indexFile)
	b.indexStream.Close()
	b.indexStream, _ = os.OpenFile(b.indexFile, os.O_RDWR, os.ModePerm)

	var gsSize uint32
	b.blockStream.Seek(4, beg)
	gsSizeBytes := make([]byte, SizeOfInt32)
	b.blockStream.Read(gsSizeBytes)
	rlp.DecodeBytes(gsSizeBytes, &gsSize)
	pos, _ := b.blockStream.Seek(int64(gsSize), cur)

	bytes, _ := rlp.EncodeToBytes(pos)

	if pos == 0 {
		return
	}

	for pos > 0 {
		b.indexStream.Write(bytes)

		var size uint32
		sizeBytes := make([]byte, SizeOfInt32)
		b.blockStream.Read(sizeBytes)
		rlp.DecodeBytes(sizeBytes, &size)
		if size == 0 {
			break
		}

		pos, _ = b.blockStream.Seek(int64(size), cur)
		if pos == 0 {
			break
		}

		pos += 8 //8 bytes pos
		bytes, _ = rlp.EncodeToBytes(pos)
	}

	return
}
func repairLog(dataDir string, truncateAtBlock uint32) string { return "" }

func ExtractGenesisState(dataDir string) types.GenesisState {

	blockStream, _ := os.OpenFile(dataDir+"/blocks.log", os.O_RDWR, os.ModePerm)
	blockStream.Seek(0, beg)
	var version uint32
	bytes := make([]byte, SizeOfInt32)
	blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &version)

	EosAssert(version > 0, &BlockLogAppendFail{}, "Block log was not setup properly with genesis information.")
	EosAssert(version == supportedVersion,
		&BlockLogUnsupportedVersion{},
		"Unsupported version of block log. Block log version is %d while code supports version %d",
		version, supportedVersion)

	var gsSize uint32
	gsSizeBytes := make([]byte, SizeOfInt32)
	blockStream.Read(gsSizeBytes)
	rlp.DecodeBytes(gsSizeBytes, &gsSize)

	gsBytes := make([]byte, gsSize)
	blockStream.Read(gsBytes)

	gs := types.GenesisState{}
	rlp.DecodeBytes(gsBytes, &gs)

	return gs
}
