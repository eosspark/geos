package producer_plugin

import (
	"crypto/sha256"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
	"os"
	"testing"
	Chain "github.com/eosspark/eos-go/plugins/producer_plugin/mock"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"syscall"
	"fmt"
)

var plugin *ProducerPlugin
var io *asio.IoContext

func BatchTest(f func(*testing.T), t *testing.T) {
	const Batch = 0
	for i:=0; i<Batch; i++ {
		f(t)
	}
}

func makeArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

func initialize(t *testing.T) {
	makeArguments("--enable-stale-production", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM\", \"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]",
	)

	app := cli.NewApp()
	io = asio.NewIoContext()

	Chain.Initialize(
		0,
		common.AccountName(common.N("eosio")),
		common.AccountName(common.N("yuanc")),


		)
	plugin = NewProducerPlugin(io)
	plugin.SetProgramOptions(&app.Flags)

	app.Action = func(c *cli.Context) {
		plugin.PluginInitialize(c)
	}

	err := app.Run(os.Args)
	assert.NoError(t, err)

	assert.Equal(t, true, plugin.my.ProductionEnabled)
	assert.Equal(t, int32(30), plugin.my.MaxTransactionTimeMs)
	assert.Equal(t, common.Seconds(-1), plugin.my.MaxIrreversibleBlockAgeUs)
	//assert.Equal(t, struct{}{}, producerPlugin.my.Producers[common.AccountName(common.N("eosio"))])
	//assert.Equal(t, struct{}{}, producerPlugin.my.Producers[common.AccountName(common.N("yuanc"))])
}

func Test_commend(t *testing.T) {

	makeArguments("-e", "-n", "eos", "-c", "18", "-p", "eosio", "-p", "yc")

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "enable, e",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.IntFlag{
			Name:  "count, c",
			Usage: "Enable block production, even if the chain is stale.",
			Value: -1,
		},
		cli.StringSliceFlag{
			Name:  "producers, p",
			Usage: "Enable block production, even if the chain is stale.",
		},
	}

	app.Action = func(c *cli.Context) {
		assert.Equal(t, true, c.Bool("enable"))
		assert.Equal(t, "eos", c.String("name"))
		assert.Equal(t, 18, c.Int("count"))

	}

	app.Action = func(c *cli.Context) {
		assert.Equal(t, "eosio", c.StringSlice("producers")[0])
		assert.Equal(t, "yc", c.StringSlice("producers")[1])
	}

	err := app.Run(os.Args)
	if err != nil {
		t.Fatal(err)
	}

}

func produceone(when common.BlockTimeStamp) (b *types.SignedBlock) {
	b = new(types.SignedBlock)
	control := Chain.GetControllerInstance()
	hbs := control.HeadBlockState()
	newBs := hbs.GenerateNext(when)
	nextProducer := hbs.GetScheduledProducer(when)
	b.Timestamp = when
	b.Producer = nextProducer.ProducerName
	b.Previous = hbs.BlockId

	//currentWatermark := plugin.my.ProducerWatermarks[nextProducer.AccountName]

	if hbs.Header.Producer != newBs.Header.Producer {
		b.Confirmed = 12
	}

	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()
	initPriKey2, _ := ecc.NewPrivateKey("5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP")
	initPubKey2 := initPriKey2.PublicKey()
	signatureProviders := map[ecc.PublicKey]signatureProviderType{
		initPubKey: func(sha256 crypto.Sha256) ecc.Signature {
			sk, _ := initPriKey.Sign(sha256.Bytes())
			return sk
		},
		initPubKey2: func(sha256 crypto.Sha256) ecc.Signature {
			sk, _ := initPriKey2.Sign(sha256.Bytes())
			return sk
		},
	}

	newBs.Header = b.SignedBlockHeader
	signatureProvider := signatureProviders[nextProducer.BlockSigningKey]
	b.ProducerSignature = signatureProvider(newBs.SigDigest())

	return
}

func exec() {
	sigint := asio.NewSignalSet(io, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		io.Stop()
		sigint.Cancel()
		plugin.PluginShutdown()
	})

	io.Run()
}

func TestProducerPlugin_PluginStartup(t *testing.T) {
	initialize(t)
	chain := Chain.GetControllerInstance()

	timer := common.NewTimer(io)
	var delayStop func()
	delayStop = func() {
		timer.ExpiresFromNow(common.Seconds(1))
		timer.AsyncWait(func(error) {
			if chain.LastIrreversibleBlockNum() > 0 && chain.HeadBlockNum() > 20 {
				io.Stop()
				plugin.PluginShutdown()
				return
			}
			delayStop()
		})
	}

	plugin.PluginStartup()

	delayStop()

	exec()
}

func TestBatch_PluginStartup(t *testing.T) {
	BatchTest(TestProducerPlugin_PluginStartup, t)
}

func TestProducerPlugin_Pause(t *testing.T) {
	initialize(t)
	plugin.PluginStartup()

	chain := Chain.GetControllerInstance()

	pause := common.NewTimer(io)
	var pauses func()
	pauses = func() {
		pause.ExpiresFromNow(common.Microseconds(300))
		pause.AsyncWait(func(error) {
			//fmt.Println("pauses")
			if chain.HeadBlockNum() == 5 {
				fmt.Println("pause")
				plugin.Pause()
				return
			}
			pauses()
		})
	}
	pauses()

	resume := common.NewTimer(io)
	var resumes func()
	resumes = func() {
		resume.ExpiresFromNow(common.Microseconds(300))
		resume.AsyncWait(func(error) {
			if plugin.Paused() {
				fmt.Println("resume")
				assert.Equal(t, uint32(5), chain.HeadBlockNum())
				plugin.Resume()
				return
			}
			resumes()
		})
	}

	resumes()

	stop := common.NewTimer(io)
	var stops func()
	stops = func() {
		stop.ExpiresFromNow(common.Microseconds(300))
		stop.AsyncWait(func(error) {
			if chain.HeadBlockNum() >= 10 {
				io.Stop()
				plugin.PluginShutdown()
				return
			}
			stops()
		})
	}
	stops()

	exec()

}

func TestBatch_Pause(t *testing.T) {
	BatchTest(TestProducerPlugin_Pause, t)
}

func TestProducerPlugin_SignCompact(t *testing.T) {
	initialize(t)
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
	sign2 := plugin.SignCompact(&initPubKey, dataHash)
	sign3 := plugin.SignCompact(&initPubKey2, dataHash)

	assert.Equal(t, sign1, sign2)
	assert.NotEqual(t, sign1, sign3)
}

func TestProducerPlugin_IsProducerKey(t *testing.T) {
	initialize(t)
	pub1, _ := ecc.NewPublicKey("EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM")
	pub2, _ := ecc.NewPublicKey("EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu")
	assert.Equal(t, true, plugin.IsProducerKey(pub1))
	assert.Equal(t, true, plugin.IsProducerKey(pub2))
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
	makeArguments("--enable-stale-production")

	app := cli.NewApp()
	io = asio.NewIoContext()

	plugin = NewProducerPlugin(io)
	plugin.SetProgramOptions(&app.Flags)
	app.Action = func(c *cli.Context) {
		plugin.PluginInitialize(c)
	}

	err := app.Run(os.Args)
	assert.NoError(t, err)

	//receive blocks 1 minutes before
	now := common.NewBlockTimeStamp(common.Now().SubUs(common.Minutes(1)))

	Chain.Initialize(now)
	mock := Chain.NewMockChain(now)

	chain := Chain.GetControllerInstance()

	for i:=1; i<100; i++ {
		oldBs := chain.HeadBlockState()
		now ++
		block := mock.ProduceOne(now)
		plugin.my.OnIncomingBlock(block)

		newBs := chain.HeadBlockState()

		assert.Equal(t, oldBs.BlockNum+1, newBs.BlockNum)
		assert.Equal(t, oldBs.BlockId, newBs.Header.Previous)
	}

}

func TestProducerPlugin_ReceiveBlockInSchedule(t *testing.T) {
	makeArguments("--enable-stale-production")

	app := cli.NewApp()
	io = asio.NewIoContext()

	plugin := NewProducerPlugin(io)
	plugin.SetProgramOptions(&app.Flags)

	app.Action = func(c *cli.Context) {
		plugin.PluginInitialize(c)
	}


	err := app.Run(os.Args)
	assert.NoError(t, err)

	//receive blocks 1 minutes before
	now := common.NewBlockTimeStamp(common.Now().SubUs(common.Minutes(1)))

	Chain.Initialize(now)
	chain := Chain.GetControllerInstance()
	mock := Chain.NewMockChain(now)

	plugin.PluginStartup()

	timer := common.NewTimer(io)

	var asyncApply func()
	asyncApply = func() {
		timer.ExpiresFromNow(common.Milliseconds(1))
		timer.AsyncWait(func(err error) {
			now ++
			block := mock.ProduceOne(now)

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
	exec()
}

func TestProducerPluginImpl_OnIncomingTransactionAsync(t *testing.T) {

}

func TestProducerPluginImpl_OnBlock(t *testing.T) {

}

func TestProducerPluginImpl_CalculateNextBlockTime(t *testing.T) {
	initialize(t)
	chain := Chain.GetControllerInstance()
	account := common.Name(common.N("yuanc"))
	pt := *plugin.my.CalculateNextBlockTime(&account, chain.HeadBlockState().Header.Timestamp)
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}

func TestProducerPluginImpl_CalculatePendingBlockTime(t *testing.T) {
	initialize(t)
	pt := plugin.my.CalculatePendingBlockTime()
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}
