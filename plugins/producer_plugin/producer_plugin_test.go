package producer_plugin

import (
					"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
	"testing"
	main "github.com/eosspark/eos-go/plugins/producer_plugin/testing"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"syscall"
	"fmt"
	"crypto/sha256"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/chain/types"
)

type ProducerPluginTester struct {
	*ProducerPlugin
	*ProducerPluginImpl

	io *asio.IoContext
	app *cli.App
	chain *main.ChainTester
}

const (
	BatchSize = 0
)

func BatchTest(f func(*testing.T), t *testing.T) {
	for i := 0; i < BatchSize; i++ {
		f(t)
	}
}

func (p *ProducerPluginTester) Exec() {
	sigint := asio.NewSignalSet(p.io, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		p.io.Stop()
		sigint.Cancel()
		p.PluginShutdown()
	})

	p.io.Run()
}

func (p *ProducerPluginTester) Stop() {
	p.io.Stop()
}



func NewProducerTester(t *testing.T) *ProducerPluginTester {
	ppt := new(ProducerPluginTester)
	ppt.io = asio.NewIoContext()
	ppt.app = cli.NewApp()

	ppt.chain = main.NewChainTester(0, common.AccountName(common.N("eosio")), common.AccountName(common.N("yuanc")))
	main.Control = ppt.chain.Control // use control instance from ChainTester

	ppt.ProducerPlugin = NewProducerPlugin(ppt.io)
	ppt.ProducerPluginImpl = ppt.my



	ppt.SetProgramOptions(&ppt.app.Flags)
	main.MakeTesterArguments("--enable-stale-production", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM\", \"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]",
	)
	ppt.app.Action = func(c *cli.Context) {
		ppt.PluginInitialize(c)
	}

	err := ppt.app.Run(os.Args)
	assert.NoError(t, err)

	assert.Equal(t, true, ppt.ProductionEnabled)
	assert.Equal(t, int32(30), ppt.MaxTransactionTimeMs)
	assert.Equal(t, common.Seconds(-1), ppt.MaxIrreversibleBlockAgeUs)
	//assert.Equal(t, struct{}{}, producerPlugin.my.Producers[common.AccountName(common.N("eosio"))])
	//assert.Equal(t, struct{}{}, producerPlugin.my.Producers[common.AccountName(common.N("yuanc"))])

	return ppt
}

//func produceone(when types.BlockTimeStamp) (b *types.SignedBlock) {
//	b = new(types.SignedBlock)
//	control := main.GetControllerInstance()
//	hbs := control.HeadBlockState()
//	newBs := hbs.GenerateNext(when)
//	nextProducer := hbs.GetScheduledProducer(when)
//	b.Timestamp = when
//	b.Producer = nextProducer.ProducerName
//	b.Previous = hbs.BlockId
//
//	//currentWatermark := plugin.my.ProducerWatermarks[nextProducer.AccountName]
//
//	if hbs.Header.Producer != newBs.Header.Producer {
//		b.Confirmed = 12
//	}
//
//	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
//	initPubKey := initPriKey.PublicKey()
//	initPriKey2, _ := ecc.NewPrivateKey("5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP")
//	initPubKey2 := initPriKey2.PublicKey()
//	signatureProviders := map[ecc.PublicKey]signatureProviderType{
//		initPubKey: func(sha256 crypto.Sha256) ecc.Signature {
//			sk, _ := initPriKey.Sign(sha256.Bytes())
//			return sk
//		},
//		initPubKey2: func(sha256 crypto.Sha256) ecc.Signature {
//			sk, _ := initPriKey2.Sign(sha256.Bytes())
//			return sk
//		},
//	}
//
//	newBs.Header = b.SignedBlockHeader
//	signatureProvider := signatureProviders[nextProducer.BlockSigningKey]
//	b.ProducerSignature = signatureProvider(newBs.SigDigest())
//
//	return
//}


func TestProducerPlugin_PluginStartup(t *testing.T) {
	tester := NewProducerTester(t)

	chain := main.GetControllerInstance()

	timer := common.NewTimer(tester.io)
	var delayStop func()
	delayStop = func() {
		timer.ExpiresFromNow(common.Seconds(1))
		timer.AsyncWait(func(error) {
			if chain.LastIrreversibleBlockNum() > 0 && chain.HeadBlockNum() > 20 {
				tester.Stop()
				tester.PluginShutdown()
				return
			}
			delayStop()
		})
	}

	tester.PluginStartup()

	delayStop()

	tester.Exec()
}


func TestBatch_PluginStartup(t *testing.T) {
	BatchTest(TestProducerPlugin_PluginStartup, t)
}

func TestProducerPlugin_Pause(t *testing.T) {
	tester := NewProducerTester(t)
	tester.PluginStartup()

	chain := main.GetControllerInstance()

	pause := common.NewTimer(tester.io)
	var pauses func()
	pauses = func() {
		pause.ExpiresFromNow(common.Microseconds(300))
		pause.AsyncWait(func(error) {
			//fmt.Println("pauses")
			if chain.HeadBlockNum() == 5 {
				fmt.Println("pause")
				tester.Pause()
				return
			}
			pauses()
		})
	}
	pauses()

	resume := common.NewTimer(tester.io)
	var resumes func()
	resumes = func() {
		resume.ExpiresFromNow(common.Microseconds(300))
		resume.AsyncWait(func(error) {
			if tester.Paused() {
				fmt.Println("resume")
				assert.Equal(t, uint32(5), chain.HeadBlockNum())
				tester.Resume()
				return
			}
			resumes()
		})
	}

	resumes()

	stop := common.NewTimer(tester.io)
	var stops func()
	stops = func() {
		stop.ExpiresFromNow(common.Microseconds(300))
		stop.AsyncWait(func(error) {
			if chain.HeadBlockNum() >= 10 {
				tester.io.Stop()
				tester.PluginShutdown()
				return
			}
			stops()
		})
	}
	stops()

	tester.Exec()

}

func TestBatch_Pause(t *testing.T) {
	BatchTest(TestProducerPlugin_Pause, t)
}

func TestProducerPlugin_SignCompact(t *testing.T) {
	tester := NewProducerTester(t)
	data := "test producer_plugin's is_producer_key "

	dataByte, _ := rlp.EncodeToBytes(data)
	h := sha256.New()
	h.Write(dataByte)

	dataByteHash := h.Sum(nil)

	dataHash := crypto.Hash256(data)

	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()
	initPriKey2, _ := ecc.NewPrivateKey("5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP")
	initPubKey2 := initPriKey2.PublicKey()

	sign1, _ := initPriKey.Sign(dataByteHash)
	sign2 := tester.SignCompact(&initPubKey, dataHash)
	sign3 := tester.SignCompact(&initPubKey2, dataHash)

	assert.Equal(t, sign1, sign2)
	assert.NotEqual(t, sign1, sign3)
}

func TestProducerPlugin_IsProducerKey(t *testing.T) {
	tester := NewProducerTester(t)
	pub1, _ := ecc.NewPublicKey("EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM")
	pub2, _ := ecc.NewPublicKey("EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu")
	assert.Equal(t, true, tester.IsProducerKey(pub1))
	assert.Equal(t, true, tester.IsProducerKey(pub2))
}

func Test_makeKeySignatureProvider(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()

	sp := makeKeySignatureProvider(initPriKey)
	hash := crypto.Hash256("makeKeySignatureProvider")
	pk, _ := sp(hash).PublicKey(hash.Bytes())
	assert.Equal(t, initPubKey, pk)

}

func TestProducerPluginImpl_StartBlock(t *testing.T) {

}

func TestProducerPluginImpl_ScheduleProductionLoop(t *testing.T) {

}

func TestProducerPluginImpl_ScheduleDelayedProductionLoop(t *testing.T) {

}

func TestProducerPluginImpl_OnIncomingBlock(t *testing.T) {


	app := cli.NewApp()
	io := asio.NewIoContext()

	plugin := NewProducerPlugin(io)
	plugin.SetProgramOptions(&app.Flags)
	main.MakeTesterArguments("--enable-stale-production")
	app.Action = func(c *cli.Context) {
		plugin.PluginInitialize(c)
	}

	err := app.Run(os.Args)
	assert.NoError(t, err)

	//receive blocks 1 minutes before
	ago := types.NewBlockTimeStamp(common.Now().SubUs(common.Minutes(1)))

	tester := main.NewChainTester(ago)

	main.Control = main.NewChainTester(ago).Control
	chain := main.GetControllerInstance()

	for i := 1; i < 100; i++ {
		oldBs := chain.HeadBlockState()
		ago ++
		block := tester.ProduceBlock(ago)
		plugin.my.OnIncomingBlock(block)

		newBs := chain.HeadBlockState()

		assert.Equal(t, oldBs.BlockNum+1, newBs.BlockNum)
		assert.Equal(t, oldBs.BlockId, newBs.Header.Previous)
	}

}

func TestProducerPlugin_ReceiveBlockInSchedule(t *testing.T) {
	app := cli.NewApp()
	io := asio.NewIoContext()

	plugin := NewProducerPlugin(io)
	plugin.SetProgramOptions(&app.Flags)
	main.MakeTesterArguments("--enable-stale-production")
	app.Action = func(c *cli.Context) {
		plugin.PluginInitialize(c)
	}

	err := app.Run(os.Args)
	assert.NoError(t, err)

	//receive blocks 1 minutes before
	ago := types.NewBlockTimeStamp(common.Now().SubUs(common.Minutes(1)))

	tester := main.NewChainTester(ago)

	main.Control = main.NewChainTester(ago).Control
	chain := main.GetControllerInstance()

	plugin.PluginStartup()

	timer := common.NewTimer(io)

	var asyncApply func()
	asyncApply = func() {
		timer.ExpiresFromNow(common.Milliseconds(1))
		timer.AsyncWait(func(err error) {
			ago ++
			block := tester.ProduceBlock(ago)

			oldBs := chain.HeadBlockState()
			plugin.my.OnIncomingBlock(block)
			newBs := chain.HeadBlockState()

			assert.Equal(t, oldBs.BlockNum+1, newBs.BlockNum)
			assert.Equal(t, oldBs.BlockId, newBs.Header.Previous)

			if chain.HeadBlockNum() >= 100 {
				io.Stop()
				plugin.PluginShutdown()
				return
			}

			asyncApply()
		})
	}

	asyncApply()
	io.Run()
}

func TestProducerPluginImpl_OnIncomingTransactionAsync(t *testing.T) {

}

func TestProducerPluginImpl_OnBlock(t *testing.T) {

}

func TestProducerPluginImpl_CalculateNextBlockTime(t *testing.T) {
	tester := NewProducerTester(t)
	chain := main.GetControllerInstance()

	account1 := common.Name(common.N("eosio"))
	account2 := common.Name(common.N("yuanc"))

	pt := *tester.CalculateNextBlockTime(&account1, chain.HeadBlockState().Header.Timestamp)
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms

	time := tester.CalculateNextBlockTime(&account1, 100)
	assert.Equal(t, "2000-01-01 00:00:50.5 +0000 UTC", time.String())

	tester.CalculateNextBlockTime(&account2, 100)
	time = tester.CalculateNextBlockTime(&account2, 100)
	assert.Equal(t, "2000-01-01 00:00:54 +0000 UTC", time.String())
}

func TestProducerPluginImpl_CalculatePendingBlockTime(t *testing.T) {
	tester := NewProducerTester(t)
	pt :=tester.CalculatePendingBlockTime()
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}
/**/