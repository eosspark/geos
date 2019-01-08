package unittests

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type ForkedTester struct {
	tester *BaseTester
}

type ForkTracker struct {
	blocks      []*types.SignedBlock
	blockMerkle types.IncrementalMerkle
}

func NewForkedTester() *ForkedTester {
	ft := ForkedTester{}
	ft.tester = newBaseTester(true, chain.SPECULATIVE)
	return &ft
}

func pushBlocks(from *BaseTester, to *BaseTester) {
	for to.Control.ForkDbHeadBlockNum() < from.Control.ForkDbHeadBlockNum() {
		fb := from.Control.FetchBlockByNumber(to.Control.ForkDbHeadBlockNum() + 1)
		to.PushBlock(fb)
	}
}

func TestIrrblock(t *testing.T) {
	try.Try(func() {
		ft := NewForkedTester()
		ft.tester.ProduceBlocks(10, false)
		dan := common.N("dan")
		sam := common.N("sam")
		pam := common.N("pam")
		scott := common.N("scott")
		accounts := []common.AccountName{dan, sam, pam, scott}
		r := ft.tester.CreateAccounts(accounts, false, true)
		res := ft.tester.SetProducers(&accounts)
		sch := []types.ProducerKey{
			{dan, ft.tester.getPublicKey(dan, "active")},
			{sam, ft.tester.getPublicKey(sam, "active")},
			{pam, ft.tester.getPublicKey(pam, "active")},
			{scott, ft.tester.getPublicKey(scott, "active")},
		}
		log.Info("set producer schedule to [dan,sam,pam] trace:%v,setProducers result:%v,%v", r, res, sch)
		ft.tester.ProduceBlocks(50, false)

		ft.tester.close()
	}).FcLogAndRethrow().End()
}

func TestForkWithBadBlock(t *testing.T) {
	try.Try(func() {
		ft := NewForkedTester()
		bios := ft.tester
		bios.ProduceBlocks(1, false)
		bios.ProduceBlocks(1, false)
		accounts := []common.AccountName{common.N("a"), common.N("b"), common.N("c"), common.N("d"), common.N("e")}
		bios.CreateAccounts(accounts, false, true)
		bios.ProduceBlocks(1, false)
		bios.SetProducers(&accounts)

		iii := 3
		for bios.Control.PendingBlockState().Header.Producer.String() != "a" || bios.Control.HeadBlockState().Header.Producer.String() != "e" {
			bios.ProduceBlocks(1, false)
			iii++
		}
		fmt.Println(iii)
		// sync remote node
		remote := newBaseTesterSecNode(true, chain.SPECULATIVE)
		pushBlocks(bios, remote)

		// produce 6 blocks on bios
		for i := 0; i < 6; i++ {
			bios.ProduceBlocks(1, false)
			assert.Equal(t, bios.Control.HeadBlockState().Header.Producer.String(), "a")
			iii++
		}
		forks := make([]ForkTracker, 7)

		// enough to skip A's blocks
		offset := common.Milliseconds(common.DefaultConfig.BlockIntervalMs * 13)
		// skip a's blocks on remote
		// create 7 forks of 7 blocks so this fork is longer where the ith block is corrupted
		for i := 0; i < 7; i++ {

			sk := remote.produceBlock(offset, false, 0)
			assert.Equal(t, sk.Producer.String(), "b")
			for j := 0; j < 7; j++ {
				fork := forks[j]
				if j <= i {
					copyB := sk
					if j == i {
						fork.blockMerkle = remote.Control.HeadBlockState().BlockrootMerkle
						//copyB.ActionMRoot.Hash[0] ^= 0x1ULL
						copyB.ActionMRoot.Hash[0] ^= 0x1
					} else if j < i {
						endB := fork.blocks[len(fork.blocks)-1]
						copyB.Previous = endB.BlockID()
					}
					p := common.MakePair(copyB.Digest(), fork.blockMerkle.GetRoot())
					headerBMRoot := crypto.Hash256(p)
					sigDigest := crypto.Hash256(common.MakePair(headerBMRoot, remote.Control.HeadBlockState().PendingScheduleHash))
					pk := remote.getPrivateKey(common.N("b"), "active")
					copyB.ProducerSignature, _ = pk.Sign(sigDigest.Bytes())

					// add this new block to our corrupted block merkle
					fork.blockMerkle.Append(copyB.BlockID())
					fork.blocks = append(fork.blocks, copyB)
				} else {
					fork.blocks = append(fork.blocks, sk)
				}

				forks[j] = fork
			}
			offset = common.Milliseconds(common.DefaultConfig.BlockIntervalMs)
		}

		// go from most corrupted fork to least
		for i := 0; i < len(forks); i++ {
			//BOOST_TEST_CONTEXT("Testing Fork: " << i) {
			fork := forks[i]
			for fidx := 0; fidx < len(fork.blocks); fidx++ {
				ssk := fork.blocks[fidx]
				// push the block only if its not known already
				if bios.Control.FetchBlockById(ssk.BlockID()) != nil {
					bios.PushBlock(ssk)
				}
			}
			try.Try(func() {
				bios.PushBlock(fork.blocks[len(fork.blocks)-1])
			}).Catch(func(e exception.FcException) {
				try.Throw(e)
			})
			//}
		}
		lib := bios.Control.HeadBlockState().DposIrreversibleBlocknum
		for tries := 0; bios.Control.HeadBlockState().DposIrreversibleBlocknum == lib && tries < 10000; tries++ {
			//++<10000
			bios.ProduceBlocks(1, false)
		}

	}).FcLogAndRethrow().End()
}

func TestForking(t *testing.T) {
	c1 := NewForkedTester().tester
	c1.ProduceBlocks(1, false)
	c1.ProduceBlocks(1, false)
	//c1.ProduceBlocks(1,false)
	dan := common.N("dan")
	sam := common.N("sam")
	pam := common.N("pam")
	accounts := []common.AccountName{dan, sam, pam}
	trace := c1.CreateAccounts(accounts, false, true)
	log.Info("trace:%v", trace)
	c1.ProduceBlocks(1, false)
	res := c1.SetProducers(&accounts)
	sch := []types.ProducerKey{
		{dan, c1.getPublicKey(dan, "active")},
		{sam, c1.getPublicKey(sam, "active")},
		{pam, c1.getPublicKey(pam, "active")},
	}
	log.Info("set producer schedule to [dan,sam,pam] trace:%v,setProducers result:%v,", res, sch)
	c1.ProduceBlocks(30, false)
	r2 := c1.CreateAccounts([]common.AccountName{eosioToken}, false, true)
	log.Info("create eosio.token:%v", r2)
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	c1.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("eosio.token.abi is err : %v", err)
	}
	c1.SetAbi(common.AccountName(eosioToken), abi, nil)

	c1.ProduceBlocks(10, false)

	data := common.Variants{
		"issuer":         "eosio",
		"maximum_supply": CoreFromString("10000000.0000"),
	}
	actionName := common.N("create")
	cr := c1.PushAction2(&eosioToken, &actionName, eosioToken, &data, c1.DefaultExpirationDelta, 0)

	log.Info("create action :%v", cr)
	an := common.N("issue")
	data2 := common.Variants{
		"to":       "dan",
		"quantity": CoreFromString("100.0000"),
		"memo":     "",
	}
	cr = c1.PushAction2(&eosioToken, &an, common.DefaultConfig.SystemAccountName, &data2, c1.DefaultExpirationDelta, 0)
	log.Info("create action :%v", cr)

	c2 := newBaseTesterSecNode(true, chain.SPECULATIVE)
	log.Info("push c1 blocks to c2")
	pushBlocks(c1, c2)
	log.Info("end push c1 blocks to c2")

	log.Info("c1 blocks:")
	c1.ProduceBlocks(3, false)

	b := c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	fmt.Println("c1 produce blockNum:", b.BlockNumber())
	assert.Equal(t, b.Producer.String(), "dan")
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	assert.Equal(t, b.Producer.String(), "sam")
	c1.ProduceBlocks(10, false)
	c1.CreateAccounts([]common.AccountName{common.N("cam")}, false, true)
	prods := []common.AccountName{common.N("dan"), common.N("sam"), common.N("pam"), common.N("cam")}
	c1.SetProducers(&prods)
	log.Info("set producer schedule to [dan,sam,pam,cam]")
	c1.ProduceBlocks(1, false)
	// The next block should be produced by pam.

	// Sync second chain with first chain.
	log.Info("push c1 blocks to c2")
	pushBlocks(c1, c2)

	// Now sam and pam go on their own fork while dan is producing blocks by himself.
	log.Info("sam and pam go off on their own fork on c2 while dan produces blocks by himself in c1")
	forkBlockNum := c1.Control.HeadBlockNum()
	log.Info("c2 blocks:")

	c2.ProduceBlocks(12, false) // pam produces 12 blocks
	b = c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*13), 0)

	assert.Equal(t, b.Producer.String(), "sam")
	c2.ProduceBlocks(11+12, false)

	log.Info("c1 blocks:")
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*13), 0)
	assert.Equal(t, b.Producer.String(), "dan")
	c1.ProduceBlocks(11, false)
	// dan on chain 1 now gets all of the blocks from chain 2 which should cause fork switch
	log.Info("push c2 blocks to c1")
	start := forkBlockNum + 1
	end := c2.Control.HeadBlockNum()
	for ; start <= end; start++ {
		log.Info("c2 %v", start)
		fb := c2.Control.FetchBlockByNumber(start)
		c1.PushBlock(fb)
	}
	/*log.Info("end push c2 blocks to c1")
	log.Info("c1 blocks:")
	c1.ProduceBlocks(24,false)
	// Switching active schedule to version 2 happens in this block.
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs),0)
	assert.Equal(t,b.Producer.String(),"pam")
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs),0)

	c1.ProduceBlocks(10,false)

	log.Info( "push c1 blocks to c2" )
	pushBlocks(c1, c2);
	log.Info( "end push c1 blocks to c2" )

	// Now with four block producers active and two identical chains (for now),
	// we can test out the case that would trigger the bug in the old fork db code:
	forkBlockNum = c1.Control.HeadBlockNum()
	log.Info( "cam and dan go off on their own fork on c1 while sam and pam go off on their own fork on c2" );
	log.Info( "c1 blocks:" )

	c1.ProduceBlocks(12,false)// dan produces 12 blocks
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs * 25),0)
	// cam skips over sam and pam's blocks
	c1.ProduceBlocks(23,false) // cam finishes the remaining 11 blocks then dan produces his 12 blocks
	log.Info( "c2 blocks:" )
	c2.ProduceBlock( common.Milliseconds(common.DefaultConfig.BlockIntervalMs * 25),0 ) // pam skips over dan and sam's blocks
	c2.ProduceBlocks(11,false) // pam finishes the remaining 11 blocks
	c2.ProduceBlock( common.Milliseconds(common.DefaultConfig.BlockIntervalMs * 25),0 ) // sam skips over cam and dan's blocks
	c2.ProduceBlocks(11,false) // sam finishes the remaining 11 blocks

	log.Info( "now cam and dan rejoin sam and pam on c2" )
	c2.ProduceBlock( common.Milliseconds(common.DefaultConfig.BlockIntervalMs * 13),0 ) // cam skips over pam's blocks (this block triggers a block on this branch to become irreversible)
	c2.ProduceBlocks(11,false) // cam produces the remaining 11 blocks
	b = c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs),0) // dan produces a block

	// a node on chain 1 now gets all but the last block from chain 2 which should cause a fork switch
	log.Info( "push c2 blocks (except for the last block by dan) to c1" )

	start=forkBlockNum+1
	end=c2.Control.HeadBlockNum()
	for ;start<=end;start++{
		log.Info("c2 ",start)
		fb := c2.Control.FetchBlockByNumber(start)
		c1.PushBlock(fb)
	}
	log.Info( "end push c2 blocks to c1" )
	log.Info( "now push dan's block to c1 but first corrupt it so it is a bad block" )
	badBlock :=b
	badBlock.TransactionMRoot = badBlock.Previous
	c1.Control.AbortBlock()
	returning := false
	try.Try(func() {
		c1.Control.PushBlock(badBlock,types.Complete)
	}).Catch(func(e exception.Exception) {
		log.Error("test forking is error:%s",e.DetailMessage())
		returning = true
	}).End()
	assert.Equal(t,false,returning)*/

}

func TestPruneRemoveBranch(t *testing.T) {

	c := NewForkedTester().tester
	c.ProduceBlocks(10, false)
	dan := common.N("dan")
	sam := common.N("sam")
	pam := common.N("pam")
	scott := common.N("scott")
	accounts := []common.AccountName{dan, sam, pam, scott}
	c.CreateAccounts(accounts, false, true)
	log.Info("set producer schedule to [dan,sam,pam,scott]")
	c.ProduceBlocks(50, false)

	c2 := newBaseTesterSecNode(true, chain.SPECULATIVE)
	log.Info("push c blocks to c2")
	pushBlocks(c, c2)
	assert.Equal(t, uint32(61), c.Control.HeadBlockNum())
	assert.Equal(t, uint32(61), c2.Control.HeadBlockNum())

}
