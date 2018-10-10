package producer_plugin

import (
	"crypto/sha256"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"testing"
	"time"
	Chain "github.com/eosspark/eos-go/plugins/producer_plugin/mock"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
	)

var plugin *ProducerPlugin

func initialize() {
	os.Args = []string{"--enable-stale-production", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM\", \"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]",
		//"--max-irreversible-block-age", "10",
	}

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	Chain.Initialize()
	producerPlugin := NewProducerPlugin()
	producerPlugin.PluginInitialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	plugin = &producerPlugin
}


func produceone(when common.BlockTimeStamp) (b *types.SignedBlock) {
	b = new(types.SignedBlock)
	control := Chain.GetControllerInstance()
	hbs := control.HeadBlockState()
	newBs := hbs.GenerateNext(&when)
	nextProducer := hbs.GetScheduledProducer(when)
	b.Timestamp = when
	b.Producer = nextProducer.AccountName
	b.Previous = hbs.ID

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
			sk,_ := initPriKey.Sign(sha256.Bytes())
			return sk
		},
		initPubKey2: func(sha256 crypto.Sha256) ecc.Signature {
			sk,_ := initPriKey2.Sign(sha256.Bytes())
			return sk
		},
	}

	newBs.Header = b.SignedBlockHeader
	signatureProvider := signatureProviders[nextProducer.BlockSigningKey]
	b.ProducerSignature = signatureProvider(newBs.SigDigest())

	return
}

func TestProducerPlugin_PluginInitialize(t *testing.T) {
	initialize()
	assert.Equal(t, true, plugin.my.ProductionEnabled)
	assert.Equal(t, int32(30), plugin.my.MaxTransactionTimeMs)
	assert.Equal(t, common.Seconds(-1), plugin.my.MaxIrreversibleBlockAgeUs)
	assert.Equal(t, struct{}{}, plugin.my.Producers[common.AccountName(common.N("eosio"))])
	assert.Equal(t, struct{}{}, plugin.my.Producers[common.AccountName(common.N("yuanc"))])
}

func TestProducerPlugin_PluginStartup(t *testing.T) {
	initialize()

	plugin.PluginStartup()

	chain := Chain.GetControllerInstance()
	for {
		time.Sleep(time.Second)
		if chain.LastIrreversibleBlockNum() > 0 && chain.HeadBlockNum() > 10 {
			plugin.PluginShutdown()
			break
		}
	}
}

func TestProducerPlugin_Pause(t *testing.T) {
	initialize()
	plugin.PluginStartup()
	once := false

	chain := Chain.GetControllerInstance()
	for {
		time.Sleep(time.Second)
		if chain.HeadBlockNum() == 10 && !once {
			println("pause")
			once = true
			plugin.Pause()
		}

		if plugin.Paused() {
			time.Sleep(3 * time.Second)
			println("rusume")
			plugin.Resume()
		}

		if chain.HeadBlockNum() > 20 {
			plugin.PluginShutdown()
			break
		}
	}

}

func TestProducerPlugin_SignCompact(t *testing.T) {
	initialize()
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
	initialize()
	pub1, _ := ecc.NewPublicKey("EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM")
	pub2, _ := ecc.NewPublicKey("EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu")
	assert.Equal(t, true, plugin.IsProducerKey(pub1))
	assert.Equal(t, true, plugin.IsProducerKey(pub2))
}

func Test_makeKeySignatureProvider(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()

	sp,_ := makeKeySignatureProvider(initPriKey)
	hash := crypto.Hash256("makeKeySignatureProvider")
	pk,_ := sp(hash).PublicKey(hash.Bytes())
	assert.Equal(t, initPubKey, pk)

}

func TestProducerPluginImpl_StartBlock(t *testing.T) {

}

func TestProducerPluginImpl_ScheduleProductionLoop(t *testing.T) {

}

func TestProducerPluginImpl_ScheduleDelayedProductionLoop(t *testing.T) {

}

func TestProducerPluginImpl_OnIncomingBlock(t *testing.T) {
	os.Args = []string{"--enable-stale-production", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM\", \"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss\"]",
		//"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]",
		//"--max-irreversible-block-age", "10",
	}

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	Chain.Initialize()
	producerPlugin := NewProducerPlugin()
	producerPlugin.PluginInitialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	plugin = &producerPlugin
	plugin.PluginStartup()

	chain := Chain.GetControllerInstance()

	for {
		time.Sleep(30 * time.Microsecond)

		if plugin.my.PendingBlockMode == EnumPendingBlockMode(speculating) {
			time.Sleep(500 * time.Millisecond)
			block := produceone(common.NewBlockTimeStamp(common.Now()))
			plugin.my.OnIncomingBlock(block)
		}

		if chain.HeadBlockNum() >= 24 {
			plugin.PluginShutdown()
			break
		}
	}
}

func TestProducerPluginImpl_OnIncomingTransactionAsync(t *testing.T) {

}

func TestProducerPluginImpl_OnBlock(t *testing.T) {

}

func TestProducerPluginImpl_CalculateNextBlockTime(t *testing.T) {
	initialize()
	chain := Chain.GetControllerInstance()
	pt := *plugin.my.CalculateNextBlockTime(common.AccountName(common.N("yuanc")), chain.HeadBlockState().Header.Timestamp)
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}

func TestProducerPluginImpl_CalculatePendingBlockTime(t *testing.T) {
	initialize()
	pt := plugin.my.CalculatePendingBlockTime()
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}

//func Test_Timer(t *testing.T) {
//	start := time.Now()
//	var apply = false
//	var timer = new(scheduleTimer)
//	var blockNum = 1
//
//	//sigs := make(chan os.Signal, 1)
//	//signal.Notify(sigs, syscall.SIGINT)
//
//	var scheduleProductionLoop func()
//
//	scheduleProductionLoop = func() {
//		timer.cancel()
//		base := time.Now()
//		minTimeToNextBlock := int64(common.DefaultConfig.BlockIntervalUs) - base.UnixNano()/1e3%int64(common.DefaultConfig.BlockIntervalUs)
//		wakeTime := base.Add(time.Microsecond * time.Duration(minTimeToNextBlock))
//
//		timer.expiresUntil(wakeTime)
//
//		// test after 12 block need to apply new block to continue
//		if blockNum%10 == 0 || (blockNum-1)%10 == 0 || (blockNum-2)%10 == 0 {
//			apply = true
//			return
//		}
//
//		timerCorelationId++
//		cid := timerCorelationId
//		timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
//			fmt.Println("exec async1...", time.Now())
//			blockNum++
//			fmt.Println("add.blockNum", blockNum)
//			scheduleProductionLoop()
//		})
//	}
//
//	applyBlock := func() {
//		for time.Now().Sub(start) <= KEEPTESTSEC*time.Second {
//			if apply {
//				apply = false
//				blockNum++
//				fmt.Println("exec apply...", time.Now(), "\n-----------apply block #.", blockNum)
//				scheduleProductionLoop()
//			}
//		}
//	}
//
//	naughty := func() {
//		for time.Now().Sub(start) <= KEEPTESTSEC*time.Second {
//			time.Sleep(666 * time.Millisecond)
//			scheduleProductionLoop()
//		}
//	}
//
//	//go func() {
//	//	sig := <-sigs
//	//	fmt.Println("sig: ", sig)
//	//}()
//
//	scheduleProductionLoop()
//	applyBlock()
//	naughty() //try to break the schedule timer
//}
