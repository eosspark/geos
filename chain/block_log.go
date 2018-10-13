package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"math"
	"os"
)

type BlockLog struct {
	nPos             uint64
	nSize            uint64
	supportedVersion uint32

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

	blockLog.blockStream, _ = os.OpenFile(blockLog.blockFile, os.O_RDWR, os.ModePerm)
	blockLog.indexStream, _ = os.OpenFile(blockLog.indexFile, os.O_RDWR, os.ModePerm)

	logSize, _ := blockLog.blockStream.Seek(0, 2)
	indexSize, _ := blockLog.indexStream.Seek(0, 2)

	if logSize > 0 {

		blockLog.blockStream.Seek(0, 0)
		var version uint32 = 0
		bytes := make([]byte, 4)
		blockLog.blockStream.Read(bytes)
		rlp.DecodeBytes(bytes, &version)

		//assert(version > 0, block_log_exception, "Block log was not setup properly with genesis information." )
		//assert(version == blockLog.supportedVersion, block_log_unsupported_version,
		//            "Unsupported version of block log. Block log version is ${version} while code supports version ${supported}",
		//            ("version", version)("supported", block_log::supported_version) )

		blockLog.genesisWriteToBlockLog = true
		blockLog.head = blockLog.Head()
		blockLog.headId = blockLog.head.BlockID()

		if indexSize > 0 {
			var blockPos uint64 = 0
			bytes = make([]byte, 8)
			blockLog.blockStream.Seek(-8, 2) //sizeof(blockPos)
			blockLog.blockStream.Read(bytes)
			rlp.DecodeBytes(bytes, &blockPos)

			var indexPos uint64 = 0
			bytes = make([]byte, 8)
			blockLog.indexStream.Seek(-8, 2) //sizeof(blockPos)
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
	//assert(b.genesisWriteToBlockLog,block_log_append_fail, "Cannot append to block log until the genesis is first written" )
	pos, _ := b.blockStream.Seek(0, 2)
	//indexPos, _ := b.indexStream.Seek(0, 2)

	// assert(indexPos == 8 * (block.BlockNumber() - 1),
	//                block_log_append_fail,
	//                "Append to index file occuring at wrong position.",
	//                ("position", (uint64_t) indexPos
	//                ("expected", 8 * (block.BlockNumber() - 1));

	data, _ := rlp.EncodeToBytes(block)
	b.blockStream.Write(data)

	posData, _ := rlp.EncodeToBytes(pos)
	b.blockStream.Write(posData)
	b.indexStream.Write(posData)

	b.head = block
	b.headId = block.BlockID()

	return 0
}
func (b *BlockLog) flush() {
	b.blockStream.Close()
	b.indexStream.Close()
}
func (b *BlockLog) ResetToGenesis(gs *types.GenesisState, benesisBlock *types.SignedBlock) { return }
func (b *BlockLog) ReadBlock(pos uint64, size uint64) (*types.SignedBlock, uint64) {
	//    s := size
	// if s == 0 {//the last block (head)
	// 	s = b.blockSteam.Seek(0,2)
	// }
	b.blockStream.Seek(int64(pos), 0)

	signedBlock := &types.SignedBlock{}
	bytes := make([]byte, size)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, signedBlock)

	var nextPos uint64
	bytes = make([]byte, 8)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &nextPos)

	return signedBlock, nextPos
}
func (b *BlockLog) ReadBlockByNum(blockNum uint32) *types.SignedBlock {

	//signedBlock := &types.SignedBlock{}
	var block *types.SignedBlock
	pos, size := b.GetBlockPos(blockNum)

	if pos != b.nPos {
		block, _ = b.ReadBlock(pos, size)
	}

	// assert(b.BlockNum() == block_num, reversible_blocks_exception,
	//       "Wrong block was read from block log.", ("returned", b.BlockNum())("expected", blockNum))

	return block
}
func (b *BlockLog) GetBlockPos(blockNum uint32) (uint64, uint64) {

	if !(b.head != nil && blockNum <= types.NumFromID(b.headId) && blockNum > 0) {
		return b.nPos, b.nSize
	}

	var pos, nextPos uint64
	bytes := make([]byte, 8)
	b.indexStream.Seek(8*(int64(blockNum)-1), 0)
	b.indexStream.Read(bytes)
	rlp.DecodeBytes(bytes, &pos)

	b.indexStream.Read(bytes)
	rlp.DecodeBytes(bytes, &nextPos)

	return pos, nextPos - pos
}
func (b *BlockLog) ReadHead() *types.SignedBlock {

	s, _ := b.blockStream.Seek(0, 2)
	if s == 0 {
		return &types.SignedBlock{}
	}

	var pos uint64
	bytes := make([]byte, 8)
	b.blockStream.Seek(-8, 2)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &pos)

	size, _ := b.blockStream.Seek(0, 2)
	block, _ := b.ReadBlock(pos, uint64(size))

	return block
}
func (b *BlockLog) Head() *types.SignedBlock {
	return b.head
}

func (b *BlockLog) ConstructIndex() {
	//ilog("Reconstructing Block Log Index...")
	b.indexStream.Close()
	os.Remove(b.indexFile)
	b.indexStream, _ = os.OpenFile(b.indexFile, os.O_RDWR, os.ModePerm)

	tmpfilename := ""
	tmp, _ := os.OpenFile(tmpfilename, os.O_RDWR, os.ModePerm)

	var endPos int64
	bytes := make([]byte, 8)
	b.blockStream.Seek(-8, 2)
	b.blockStream.Read(bytes)
	rlp.DecodeBytes(bytes, &endPos)

	var pos int64 = 4
	for {
		tmp.Write(bytes)

		endPos -= 8
		if pos < endPos {
			break
		}
		b.blockStream.Seek(endPos, 0)
		b.blockStream.Read(bytes)
		rlp.DecodeBytes(bytes, &endPos)
	}

	count, _ := tmp.Seek(0, 2)
	count /= 8
	for i := int64(0); i < count; i++ {

		tmp.Seek(-i*8, 2)
		tmp.Read(bytes)

		b.indexStream.Write(bytes)
	}

	return
}
func (b *BlockLog) repairLog(dataDir string, truncateAtBlock uint32) string { return "" }
func (b *BlockLog) ExtractGenesisState(dataDir string) types.GenesisState   { return types.GenesisState{} }
