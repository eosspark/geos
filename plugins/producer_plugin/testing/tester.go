package testing

import (
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/crypto"
	"os"
	"github.com/eosspark/container/maps/treemap"
)

type ChainTester struct {
	Control  *Controller
	KeyPairs map[common.AccountName]common.Pair //[]<pubKey, priKey>
}

//var chain *Controller

//var initPriKey, _ = ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
//var initPubKey = initPriKey.PublicKey()
//var initPriKey2, _ = ecc.NewPrivateKey("5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP")
//var initPubKey2 = initPriKey2.PublicKey()
//var eosio = common.AccountName(common.N("eosio"))
//var yuanc = common.AccountName(common.N("yuanc"))

func NewChainTester(when types.BlockTimeStamp, names ...common.AccountName) *ChainTester {
	tester := new(ChainTester)
	priKey, err := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	maythrow(err)
	pubKey := priKey.PublicKey()

	tester.KeyPairs = make(map[common.AccountName]common.Pair)
	tester.KeyPairs[common.AccountName(common.N("eosio"))] = common.MakePair(pubKey, priKey)
	tester.KeyPairs[common.AccountName(common.N("yuanc"))] = common.MakePair(pubKey, priKey)

	hbs := tester.NewHeaderStateTester(when)
	sbk := tester.NewSignedBlockTester(hbs)
	sch := tester.NewProducerScheduleTester(names...)

	tester.Control = newController()
	tester.Control.head = types.NewBlockState(hbs)
	tester.Control.head.SignedBlock = sbk

	tester.Control.head.ActiveSchedule = sch
	tester.Control.head.PendingSchedule = sch

	tester.Control.forkDb.add(tester.Control.head)

	return tester
}


func (t *ChainTester) NewProducerScheduleTester(names ...common.AccountName) types.ProducerScheduleType {
	if len(names) == 0 {
		names = append(names, common.AccountName(common.N("eosio")))
	}

	initSchedule := types.ProducerScheduleType{Version: 0, Producers: []types.ProducerKey{}}

	for _, n := range names {
		pk := types.ProducerKey{ProducerName: n, BlockSigningKey: t.KeyPairs[n].First.(ecc.PublicKey)}
		initSchedule.Producers = append(initSchedule.Producers, pk)
	}

	return initSchedule
}

func (t *ChainTester) NewSignedBlockTester(bhs *types.BlockHeaderState) *types.SignedBlock {
	genSigned := new(types.SignedBlock)
	genSigned.SignedBlockHeader = bhs.Header
	return genSigned
}

func (t *ChainTester) NewHeaderStateTester(when types.BlockTimeStamp) *types.BlockHeaderState {
	if when == 0 {
		when = types.NewBlockTimeStamp(common.Now())
	}
	genHeader := new(types.BlockHeaderState)
	genHeader.Header.Timestamp = when
	genHeader.Header.Confirmed = 1
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()
	genHeader.ProducerToLastProduced = *treemap.NewWith(common.NameComparator)
	genHeader.ProducerToLastImpliedIrb = *treemap.NewWith(common.NameComparator)

	return genHeader
}



func (t *ChainTester) ProduceBlock(when types.BlockTimeStamp) *types.SignedBlock {

	t.Control.AbortBlock()
	t.Control.StartBlock(when, 0)
	t.Control.FinalizeBlock()

	producer := t.Control.HeadBlockState().GetScheduledProducer(when).ProducerName

	s := t.Control.SignBlock(func(digest crypto.Sha256) ecc.Signature {
		sign, err := t.KeyPairs[producer].Second.(*ecc.PrivateKey).Sign(digest.Bytes())
		if err != nil {
			panic(err)
		}
		return sign
	})

	t.Control.CommitBlock(true)

	return s
}

func MakeTesterArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

func maythrow(err error) {
	if err != nil {
		panic(err)
	}
}
