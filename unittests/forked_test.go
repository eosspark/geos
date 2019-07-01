package unittests

import (
	"io/ioutil"
	"testing"

	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"

	"github.com/stretchr/testify/assert"
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

func pushBlocks2(from *BaseTester, to *BaseTester) {

	for to.Control.ForkDbHeadBlockNum() < from.Control.ForkDbHeadBlockNum() {
		fb := from.Control.FetchBlockByNumber(to.Control.ForkDbHeadBlockNum() + 1)
		to.PushBlock2(fb, types.Irreversible)
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
			{ProducerName: dan, BlockSigningKey: ft.tester.getPublicKey(dan, "active")},
			{ProducerName: sam, BlockSigningKey: ft.tester.getPublicKey(sam, "active")},
			{ProducerName: pam, BlockSigningKey: ft.tester.getPublicKey(pam, "active")},
			{ProducerName: scott, BlockSigningKey: ft.tester.getPublicKey(scott, "active")},
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

		for bios.Control.PendingBlockState().Header.Producer.String() != "a" || bios.Control.HeadBlockState().Header.Producer.String() != "e" {
			bios.ProduceBlocks(1, false)
		}
		// sync remote node
		remote := newBaseTesterSecNode(true, chain.SPECULATIVE)
		pushBlocks(bios, remote)

		// produce 6 blocks on bios
		for i := 0; i < 6; i++ {
			bios.ProduceBlocks(1, false)
			assert.Equal(t, bios.Control.HeadBlockState().Header.Producer.String(), "a")
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
					copyB := types.NewSignedBlock1(&sk.SignedBlockHeader)
					if j == i {
						fork.blockMerkle = remote.Control.HeadBlockState().BlockrootMerkle
						//copyB.ActionMRoot.Hash[0] ^= 0x1ULL
						copyB.ActionMRoot.Hash[0] ^= 0x1
					} else if j < i {
						endB := fork.blocks[len(fork.blocks)-1]
						copyB.Previous = endB.BlockID()
					}

					headerBMRoot := crypto.Hash256(common.MakePair(copyB.Digest(), fork.blockMerkle.GetRoot()))
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
			for fidx := 0; fidx < len(fork.blocks)-1; fidx++ {
				ssk := fork.blocks[fidx]
				// push the block only if its not known already
				if bios.Control.FetchBlockById(ssk.BlockID()) == nil {
					bios.PushBlock(ssk)
				}
			}
			// push the block which should attempt the corrupted fork and fail

			var ex string
			try.Try(func() {
				bios.PushBlock(fork.blocks[len(fork.blocks)-1])
			}).Catch(func(e exception.Exception) {
				ex = e.DetailMessage()
			}).End()
			assert.True(t, inString(ex, "Block ID does not match"))
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
	dan := common.N("dan")
	sam := common.N("sam")
	pam := common.N("pam")
	accounts := []common.AccountName{dan, sam, pam}
	trace := c1.CreateAccounts(accounts, false, true)
	log.Info("trace:%v", trace)
	c1.ProduceBlocks(1, false)
	res := c1.SetProducers(&accounts)
	sch := []types.ProducerKey{
		{ProducerName: dan, BlockSigningKey: c1.getPublicKey(dan, "active")},
		{ProducerName: sam, BlockSigningKey: c1.getPublicKey(sam, "active")},
		{ProducerName: pam, BlockSigningKey: c1.getPublicKey(pam, "active")},
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
	log.Info("end push c2 blocks to c1")
	log.Info("c1 blocks:")
	c1.ProduceBlocks(24, false)
	// Switching active schedule to version 2 happens in this block.
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	assert.Equal(t, b.Producer.String(), "pam")
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	c1.ProduceBlocks(10, false)

	log.Info("push c1 blocks to c2")
	pushBlocks(c1, c2)
	log.Info("end push c1 blocks to c2")

	// Now with four block producers active and two identical chains (for now),
	// we can test out the case that would trigger the bug in the old fork db code:
	forkBlockNum = c1.Control.HeadBlockNum()
	log.Info("cam and dan go off on their own fork on c1 while sam and pam go off on their own fork on c2")
	log.Info("c1 blocks:")

	c1.ProduceBlocks(12, false) // dan produces 12 blocks
	b = c1.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*25), 0)
	// cam skips over sam and pam's blocks
	c1.ProduceBlocks(23, false) // cam finishes the remaining 11 blocks then dan produces his 12 blocks
	log.Info("c2 blocks:")
	c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*25), 0) // pam skips over dan and sam's blocks
	c2.ProduceBlocks(11, false)                                                      // pam finishes the remaining 11 blocks
	c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*25), 0) // sam skips over cam and dan's blocks
	c2.ProduceBlocks(11, false)                                                      // sam finishes the remaining 11 blocks

	log.Info("now cam and dan rejoin sam and pam on c2")
	c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*13), 0)  // cam skips over pam's blocks (this block triggers a block on this branch to become irreversible)
	c2.ProduceBlocks(11, false)                                                       // cam produces the remaining 11 blocks
	b = c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0) // dan produces a block

	// a node on chain 1 now gets all but the last block from chain 2 which should cause a fork switch
	log.Info("push c2 blocks (except for the last block by dan) to c1")

	start = forkBlockNum + 1
	end = c2.Control.HeadBlockNum()
	for ; start <= end; start++ {
		log.Info("c2 %v", start)
		fb := c2.Control.FetchBlockByNumber(start)
		c1.PushBlock(fb)
	}
	log.Info("end push c2 blocks to c1")
	log.Info("now push dan's block to c1 but first corrupt it so it is a bad block")
	badBlock := b
	badBlock.TransactionMRoot = badBlock.Previous
	c1.Control.AbortBlock()
	var ex string
	try.Try(func() {
		c1.Control.PushBlock(badBlock, types.Complete)
	}).Catch(func(e exception.Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "block not signed by expected key"))

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
	c.SetProducers(&accounts)
	log.Info("set producer schedule to [dan,sam,pam,scott]")
	c.ProduceBlocks(50, false)

	c2 := newBaseTesterSecNode(true, chain.SPECULATIVE)
	log.Info("push c blocks to c2")
	pushBlocks(c, c2)
	assert.Equal(t, uint32(61), c.Control.HeadBlockNum())
	assert.Equal(t, uint32(61), c2.Control.HeadBlockNum())
	forkNum := c.Control.HeadBlockNum()

	nextproducer := func(c *BaseTester, skipInterval int) common.Name {
		headTime := c.Control.HeadBlockTime()
		nextTime := headTime + common.TimePoint(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*int64(skipInterval)))
		return c.Control.HeadBlockState().GetScheduledProducer(types.NewBlockTimeStamp(nextTime)).ProducerName
	}
	// fork c: 2 producers: dan, sam
	// fork c2: 1 producer: scott
	skip1 := 1
	skip2 := 1
	for i := 0; i < 50; i++ {
		next1 := nextproducer(c, skip1)
		if next1 == dan || next1 == sam {
			c.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*int64(skip1)), 0)
			skip1 = 1
		} else {
			skip1++
		}
		next2 := nextproducer(c2, skip2)
		if next2 == scott {
			c2.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*int64(skip2)), 0)
			skip2 = 1
		} else {
			skip2++
		}
	}
	assert.Equal(t, uint32(87), c.Control.HeadBlockNum())
	assert.Equal(t, uint32(73), c2.Control.HeadBlockNum())

	p := forkNum
	// push fork from c2 => c
	for p < c2.Control.HeadBlockNum() {
		p++
		fb := c2.Control.FetchBlockByNumber(p)
		c.PushBlock(fb)
	}
	assert.Equal(t, uint32(73), c.Control.HeadBlockNum())
}

func TestConfirmation(t *testing.T) {
	c := NewForkedTester().tester
	c.ProduceBlocks(10, false)
	dan := common.N("dan")
	sam := common.N("sam")
	pam := common.N("pam")
	scott := common.N("scott")
	invalid := common.N("invalid")
	accounts := []common.AccountName{dan, sam, pam, scott}
	c.CreateAccounts(accounts, false, true)
	/*res := */ c.SetProducers(&accounts)
	privSam := c.getPrivateKey(sam, "active")
	privDan := c.getPrivateKey(dan, "active")
	privPam := c.getPrivateKey(pam, "active")
	privScott := c.getPrivateKey(scott, "active")
	privInvalid := c.getPrivateKey(invalid, "active")

	log.Info("set producer schedule to [dan,sam,pam,scott]")
	c.ProduceBlocks(50, false)
	c.Control.AbortBlock() // discard pending block

	assert.Equal(t, uint32(61), c.Control.HeadBlockNum())

	blk := c.Control.ForkDB.GetBlockInCurrentChainByNum(55)
	blk61 := c.Control.ForkDB.GetBlockInCurrentChainByNum(61)
	blk50 := c.Control.ForkDB.GetBlockInCurrentChainByNum(50)

	assert.Equal(t, uint32(0), blk.BftIrreversibleBlocknum)
	assert.Equal(t, 0, len(blk.Confirmations))
	{
		var ex string
		try.Try(func() {
			h := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: sam, ProducerSignature: ecc.Signature{}}
			h.ProducerSignature, _ = privInvalid.Sign(blk.SigDigest().Bytes())
			c.Control.PushConfirmation(&h)
		}).Catch(func(e exception.Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "confirmation not signed by expected key"))
	}
	{
		var ex string
		try.Try(func() {
			h := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: invalid, ProducerSignature: ecc.Signature{}}
			h.ProducerSignature, _ = privInvalid.Sign(blk.SigDigest().Bytes())
			c.Control.PushConfirmation(&h)
		}).Catch(func(e exception.Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer not in current schedule"))
	}
	{
		// signed by sam
		h := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: sam, ProducerSignature: ecc.Signature{}}
		h.ProducerSignature, _ = privSam.Sign(blk.SigDigest().Bytes())
		c.Control.PushConfirmation(&h)
		assert.Equal(t, uint32(0), blk.BftIrreversibleBlocknum)
		assert.Equal(t, 1, len(blk.Confirmations))
		// double confirm not allowed
		var ex string
		try.Try(func() {
			h2 := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: sam, ProducerSignature: ecc.Signature{}}
			h2.ProducerSignature, _ = privSam.Sign(blk.SigDigest().Bytes())
			c.Control.PushConfirmation(&h2)
		}).Catch(func(e exception.Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "block already confirmed by this producer"))
	}

	{
		// signed by dan
		h := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: dan, ProducerSignature: ecc.Signature{}}
		h.ProducerSignature, _ = privDan.Sign(blk.SigDigest().Bytes())
		c.Control.PushConfirmation(&h)
		assert.Equal(t, uint32(0), blk.BftIrreversibleBlocknum)
		assert.Equal(t, 2, len(blk.Confirmations))

		//signed by pam
		h2 := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: pam, ProducerSignature: ecc.Signature{}}
		h2.ProducerSignature, _ = privPam.Sign(blk.SigDigest().Bytes())
		c.Control.PushConfirmation(&h2)
		// we have more than 2/3 of confirmations, bft irreversible number should be set
		assert.Equal(t, uint32(55), blk.BftIrreversibleBlocknum)
		assert.Equal(t, uint32(55), blk61.BftIrreversibleBlocknum)
		assert.Equal(t, uint32(0), blk50.BftIrreversibleBlocknum)
		assert.Equal(t, 3, len(blk.Confirmations))
	}
	{
		// signed by scott
		h := types.HeaderConfirmation{BlockId: blk.BlockId, Producer: scott, ProducerSignature: ecc.Signature{}}
		h.ProducerSignature, _ = privScott.Sign(blk.SigDigest().Bytes())
		c.Control.PushConfirmation(&h)
		assert.Equal(t, uint32(55), blk.BftIrreversibleBlocknum)
		assert.Equal(t, 4, len(blk.Confirmations))
	}

	{
		h := types.HeaderConfirmation{BlockId: blk50.BlockId, Producer: sam, ProducerSignature: ecc.Signature{}}
		h.ProducerSignature, _ = privSam.Sign(blk50.SigDigest().Bytes())
		c.Control.PushConfirmation(&h)

		h2 := types.HeaderConfirmation{BlockId: blk50.BlockId, Producer: dan, ProducerSignature: ecc.Signature{}}
		h2.ProducerSignature, _ = privDan.Sign(blk50.SigDigest().Bytes())
		c.Control.PushConfirmation(&h2)

		h3 := types.HeaderConfirmation{BlockId: blk50.BlockId, Producer: pam, ProducerSignature: ecc.Signature{}}
		h3.ProducerSignature, _ = privPam.Sign(blk50.SigDigest().Bytes())
		c.Control.PushConfirmation(&h3)
		assert.Equal(t, uint32(50), blk50.BftIrreversibleBlocknum)

		blk54 := c.Control.ForkDB.GetBlockInCurrentChainByNum(54)
		assert.Equal(t, uint32(50), blk54.BftIrreversibleBlocknum)
		assert.Equal(t, uint32(55), blk.BftIrreversibleBlocknum)
		assert.Equal(t, uint32(55), blk61.BftIrreversibleBlocknum)

		c.ProduceBlocks(20, false)
		blk81 := c.Control.ForkDB.GetBlockInCurrentChainByNum(81)
		assert.Equal(t, uint32(55), blk81.BftIrreversibleBlocknum)

	}

}

func TestReadModes(t *testing.T) {
	c := NewForkedTester().tester
	dan := common.N("dan")
	sam := common.N("sam")
	pam := common.N("pam")
	c.ProduceBlocks(1, false)
	c.ProduceBlocks(1, false)

	accounts := []common.AccountName{dan, sam, pam}
	c.CreateAccounts(accounts, false, true)
	c.ProduceBlocks(1, false)
	c.SetProducers(&accounts)
	c.ProduceBlocks(200, false)
	headBlockNum := c.Control.HeadBlockNum()

	head := newBaseTester(true, chain.HEADER)
	pushBlocks(c, head)
	assert.Equal(t, headBlockNum, head.Control.ForkDbHeadBlockNum())
	assert.Equal(t, headBlockNum, head.Control.HeadBlockNum())

	readOnly := newBaseTester(false, chain.READONLY)
	pushBlocks(c, readOnly)
	assert.Equal(t, headBlockNum, readOnly.Control.ForkDbHeadBlockNum())
	assert.Equal(t, headBlockNum, readOnly.Control.HeadBlockNum())

	irreversible := newBaseTester(true, chain.IRREVERSIBLE)
	pushBlocks2(c, irreversible)
	/*fmt.Println("**************irreversible**************：", c.Control.LastIrreversibleBlockNum(),irreversible.Control.LastIrreversibleBlockNum(), irreversible.Control.HeadBlockNum(), headBlockNum)
	fmt.Println("**************irreversible**************：", irreversible.Control.LastIrreversibleBlockNum())*/
	assert.Equal(t, headBlockNum, irreversible.Control.ForkDbHeadBlockNum())
	assert.Equal(t, headBlockNum-49, irreversible.Control.HeadBlockNum())

}
