package testing

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type DBReadMode int8

const (
	SPECULATIVE = DBReadMode(iota)
	HEADER       //HEAD
	READONLY
	IRREVERSIBLE
)

type ValidationMode int8

const (
	FULL = ValidationMode(iota)
	LIGHT
)


type Controller struct {
	head    *types.BlockState
	pending *types.BlockState
	forkDb  *forkDatabase
}

func newController() *Controller {
	c := new(Controller)
	c.forkDb = new(forkDatabase)
	c.forkDb.index = make([]*types.BlockState, 0)
	return c
}

var Control *Controller
func GetControllerInstance() *Controller {
	if Control == nil {
		Control = newController()

	}
	return Control
}

type forkDatabase struct {
	head  *types.BlockState
	index []*types.BlockState
}

func (db *forkDatabase) find(id common.BlockIdType) *types.BlockState {
	for _, n := range db.index {
		if n.BlockId == id {
			return n
		}
	}
	return nil
}

func (db *forkDatabase) findByNum(num uint32) *types.BlockState {
	for _, n := range db.index {
		if n.BlockNum == num {
			return n
		}
	}
	return nil
}

func (db *forkDatabase) add(n *types.BlockState) *types.BlockState {
	db.index = append(db.index, n)
	db.head = n
	return n
}

func (db *forkDatabase) add2(b *types.SignedBlock, trust bool) *types.BlockState {
	prior := db.find(b.Previous)
	result := types.NewBlockState3(&prior.BlockHeaderState, b, trust)
	return db.add(result)
}

func (c *Controller) LastIrreversibleBlockNum() uint32 {
	return c.head.DposIrreversibleBlocknum
}

func (c *Controller) HeadBlockState() *types.BlockState {
	return c.head

}
func (c *Controller) HeadBlockTime() common.TimePoint {
	return c.head.Header.Timestamp.ToTimePoint()
}

func (c *Controller) PendingBlockTime() common.TimePoint {
	return c.pending.Header.Timestamp.ToTimePoint()
}

func (c *Controller) HeadBlockNum() uint32 {
	return c.head.BlockNum
}

func (c *Controller) PendingBlockState() *types.BlockState {
	return c.pending
}

func (c *Controller) GetUnappliedTransactions() []*types.TransactionMetadata {
	return make([]*types.TransactionMetadata, 0)
}

func (c *Controller) DropUnappliedTransaction(trx *types.TransactionMetadata) {}

func (c *Controller) GetScheduledTransactions() []common.TransactionIdType {
	return make([]common.TransactionIdType, 0)
}

func (c *Controller) AbortBlock() {
	//fmt.Println("abort block...")
	if c.pending != nil {
		c.pending = nil
	}
}

func (c *Controller) StartBlock(when types.BlockTimeStamp, confirmBlockCount uint16) {
	//fmt.Println("start block...")
	c.pending = types.NewBlockState2(&c.head.BlockHeaderState, when)
	c.pending.SetConfirmed(confirmBlockCount)

}
func (c *Controller) FinalizeBlock() {
	//fmt.Println("finalize block...")
	c.pending.BlockId = c.pending.Header.BlockID()
}

func (c *Controller) SignBlock(callback func(sha256 crypto.Sha256) ecc.Signature) *types.SignedBlock {
	//fmt.Println("sign block...")
	p := c.pending
	p.Sign(callback)
	//println("after signer")
	p.SignedBlock.SignedBlockHeader = p.Header
	return p.SignedBlock
}

func (c *Controller) CommitBlock(addToForkDb bool) {
	//fmt.Println("commit block...")

	if addToForkDb {
		c.pending.Validated = true
		c.forkDb.add(c.pending)
		c.head = c.forkDb.head
	}

	//c.pending = nil
}

func (c *Controller) PushTransaction(trx *types.TransactionMetadata, deadline common.TimePoint, billedCpuTimeUs uint32) *types.TransactionTrace {
	//c.pending.SignedBlock.Transactions = append(c.pending.SignedBlock.Transactions, )
	c.pending.Trxs = append(c.pending.Trxs, trx)
	return nil
}
func (c *Controller) PushScheduledTransaction(trx *common.TransactionIdType, deadline common.TimePoint, billedCpuTimeUs uint32) *types.TransactionTrace {
	return nil
}

func (c *Controller) PushReceipt(trx interface{}) types.TransactionReceipt {
	//c.pending.SignedBlock.Transactions = append(c.pending.SignedBlock.Transactions, )
	return types.TransactionReceipt{}
}

func (c *Controller) PushBlock(b *types.SignedBlock, status types.BlockStatus) {
	c.forkDb.add2(b, false)
	if c.GetReadMode() != DBReadMode(IRREVERSIBLE) {
		c.MaybeSwitchForks()
	}

}

func (c *Controller) MaybeSwitchForks() {
	newHead := c.forkDb.head

	if newHead.Header.Previous == c.head.BlockId {
		c.ApplyBlock(newHead.SignedBlock)
		c.head = newHead

	} else if newHead.ID != c.head.ID {
		fmt.Println(" newHead.ID != c.head.ID ")
	}
}

func (c *Controller) ApplyBlock(b *types.SignedBlock) {
	c.StartBlock(b.Timestamp, b.Confirmed)
	c.FinalizeBlock()

	c.pending.Header.ProducerSignature = b.ProducerSignature
	c.pending.SignedBlock.SignedBlockHeader = c.pending.Header

	c.CommitBlock(false)
}

func (c *Controller) FetchBlockById(id common.BlockIdType) *types.SignedBlock {
	state := c.forkDb.find(id)
	if state != nil {
		return state.SignedBlock
	}
	bptr := c.FetchBlockByNumber(types.NumFromID(&id))
	if bptr != nil && bptr.BlockID() == id {
		return bptr
	}
	return nil
}

func (c *Controller) FetchBlockByNumber(num uint32) *types.SignedBlock {
	state := c.forkDb.findByNum(num)
	if state != nil {
		return state.SignedBlock
	}
	return nil
}

func (c *Controller) IsKnownUnexpiredTransaction(id *common.TransactionIdType) bool {
	return false
}

func (c *Controller) GetReadMode() DBReadMode {
	return DBReadMode(SPECULATIVE)
}

func (c *Controller) GetValidationMode() ValidationMode { return ValidationMode(FULL) }

func (c *Controller) SetSubjectiveCpuLeeway(leeway common.Microseconds) {}

func (c *Controller) AddResourceGreyList(name *common.AccountName)    {}
func (c *Controller) RemoveResourceGreyList(name *common.AccountName) {}

func (c *Controller) GetResourceGreyList() *common.FlatSet { return nil }

func (c *Controller) GetActorWhiteList() *common.FlatSet { return nil }

func (c *Controller) GetActorBlackList() *common.FlatSet { return nil }

func (c *Controller) GetContractWhiteList() *common.FlatSet { return nil }

func (c *Controller) GetContractBlackList() *common.FlatSet { return nil }

func (c *Controller) GetActionBlockList() *common.FlatSet { return nil }

func (c *Controller) GetKeyBlackList() *common.FlatSet { return nil }

func (c *Controller) SetActorWhiteList(params *common.FlatSet) {}
func (c *Controller) SetActorBlackList(params *common.FlatSet) {}

func (c *Controller) SetContractWhiteList(params *common.FlatSet) {}
func (c *Controller) SetContractBlackList(params *common.FlatSet) {}

func (c *Controller) SetActionBlackList(params *common.FlatSet) {}

func (c *Controller) SetKeyBlackList(params *common.FlatSet) {}
