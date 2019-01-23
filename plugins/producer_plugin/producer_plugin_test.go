package producer_plugin

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"testing"

	"crypto/sha256"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
)

func TestProducerPlugin_FindByApplication(t *testing.T) {
	plugin := app.App().FindPlugin(ProducerPlug).(*ProducerPlugin)
	assert.NotNil(t, plugin)
}

var eosio = common.N("eosio")
var yuanc = common.N("yuanc")

func producerPluginInitialize(arguments ...string) *ProducerPlugin {
	app := cli.NewApp()
	plugin := NewProducerPlugin(asio.NewIoContext())
	plugin.SetProgramOptions(&app.Flags)
	app.Action = func(option *cli.Context) {
		plugin.PluginInitialize(option)
	}
	app.Run(append(make([]string, 1, len(arguments)+1), arguments...))
	return plugin
}

func TestProducerPlugin_PluginInitialize(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]")

	assert.Equal(t, true, plugin.my.Producers.Contains(common.N("eosio")), common.N("yuanc"))
	assert.Equal(t, true, plugin.my.ProductionEnabled)

	pub, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	assert.NoError(t, err)
	assert.Contains(t, plugin.my.SignatureProviders, pub)

	pub, err = ecc.NewPublicKey("EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu")
	assert.NoError(t, err)
	assert.Contains(t, plugin.my.SignatureProviders, pub)

}

func TestProducerPlugin_PluginStartup(t *testing.T) {
	// startup by self
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	try.Try(func() {
		plugin.PluginStartup()
	}).Catch(func(e interface{}) {
		assert.Fail(t, "startup with exception")
	}).End()

	assert.Equal(t, producing, plugin.my.PendingBlockMode)

}

func TestProducerPlugin_Pause(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	plugin.PluginStartup()
	assert.Equal(t, producing, plugin.my.PendingBlockMode)
	assert.Equal(t, false, plugin.my.ProductionPaused)
	assert.Equal(t, false, plugin.Paused())

	plugin.Pause()
	assert.Equal(t, true, plugin.my.ProductionPaused)
	assert.Equal(t, true, plugin.Paused())

	result, _ := plugin.my.StartBlock()
	assert.Equal(t, speculating, plugin.my.PendingBlockMode)
	assert.Equal(t, waiting, result)

	plugin.Resume()
	assert.Equal(t, false, plugin.my.ProductionPaused)
	assert.Equal(t, false, plugin.Paused())

	result, _ = plugin.my.StartBlock()
	assert.Equal(t, producing, plugin.my.PendingBlockMode)
	assert.Equal(t, succeeded, result)
}

func TestProducerPlugin_SignCompact(t *testing.T) {
	data := "test producer_plugin SignCompact "

	dataByte, _ := rlp.EncodeToBytes(data)
	h := sha256.New()
	h.Write(dataByte)

	dataByteHash := h.Sum(nil)
	dataHash := crypto.Hash256(data)

	initPriKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	initPubKey := initPriKey.PublicKey()
	initPriKey2, _ := ecc.NewPrivateKey("5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP")
	initPubKey2 := initPriKey2.PublicKey()

	plugin := producerPluginInitialize("-e", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]")

	sign1, _ := initPriKey.Sign(dataByteHash)
	sign2, _ := initPriKey2.Sign(dataByteHash)

	assert.Equal(t, sign1, *plugin.SignCompact(&initPubKey, *dataHash))
	assert.Equal(t, sign2, *plugin.SignCompact(&initPubKey2, *dataHash))
	assert.NotEqual(t, sign1, sign2)
}

func TestProducerPlugin_IsProducerKey(t *testing.T) {
	pub1, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	pub2, _ := ecc.NewPublicKey("EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu")

	plugin := producerPluginInitialize("-e", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]")

	assert.Equal(t, true, plugin.IsProducerKey(pub1))
	assert.Equal(t, true, plugin.IsProducerKey(pub2))
}

func Test_makeKeySignatureProvider(t *testing.T) {
	initPriKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	initPubKey := initPriKey.PublicKey()

	sp := makeKeySignatureProvider(initPriKey)
	hash := crypto.Hash256("makeKeySignatureProvider")
	pk, _ := sp(*hash).PublicKey(hash.Bytes())
	assert.Equal(t, initPubKey, pk)

}

func TestProducerPluginImpl_StartBlock(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	plugin.PluginStartup()
	result, last := plugin.my.StartBlock()
	assert.Equal(t, last, uint32(types.NewBlockTimeStamp(plugin.my.CalculatePendingBlockTime()))%uint32(common.DefaultConfig.ProducerRepetitions) == uint32(common.DefaultConfig.ProducerRepetitions-1))
	assert.Equal(t, succeeded, result)
	assert.Equal(t, producing, plugin.my.PendingBlockMode)

	plugin.my.ProducerWatermarks[common.N("eosio")] = 100
	result, _ = plugin.my.StartBlock()
	assert.Equal(t, waiting, result)
	assert.Equal(t, speculating, plugin.my.PendingBlockMode)

	plugin = producerPluginInitialize()
	plugin.PluginStartup()
	result, _ = plugin.my.StartBlock()
	assert.Equal(t, waiting, result)
	assert.Equal(t, speculating, plugin.my.PendingBlockMode)

	plugin = producerPluginInitialize("-e")
	plugin.PluginStartup()
	result, _ = plugin.my.StartBlock()
	assert.Equal(t, waiting, result)
	assert.Equal(t, speculating, plugin.my.PendingBlockMode)

	plugin = producerPluginInitialize("-e", "-p", "eosio")
	plugin.PluginStartup()
	result, _ = plugin.my.StartBlock()
	assert.Equal(t, waiting, result)
	assert.Equal(t, speculating, plugin.my.PendingBlockMode)
}

func TestProducerPluginImpl_ScheduleProductionLoop(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	hbn := plugin.my.Chain.HeadBlockNum()

	plugin.PluginStartup()
	assert.Equal(t, true, plugin.my.MaybeProduceBlock(), "produce failed in startup ")
	assert.EqualValues(t, hbn+1, plugin.my.Chain.HeadBlockNum())

	plugin.my.ScheduleProductionLoop()
	assert.Equal(t, true, plugin.my.MaybeProduceBlock(), "produce failed in schedule loop")
	assert.EqualValues(t, hbn+2, plugin.my.Chain.HeadBlockNum())
}

func TestProducerPluginImpl_MaybeProduceBlock(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	hbn := plugin.my.Chain.HeadBlockNum()
	plugin.PluginStartup()

	assert.Equal(t, true, plugin.my.MaybeProduceBlock())
	assert.EqualValues(t, hbn+1, plugin.my.Chain.HeadBlockNum())
	assert.Equal(t, true, plugin.my.MaybeProduceBlock())
	assert.EqualValues(t, hbn+2, plugin.my.Chain.HeadBlockNum())
}

func TestProducerPluginImpl_CalculateNextBlockTime(t *testing.T) {

	account1 := common.Name(common.N("eosio"))
	account2 := common.Name(common.N("yuanc"))

	plugin := producerPluginInitialize("-e", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]")

	pt := *plugin.my.CalculateNextBlockTime(&account1, 0)
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms

	time := plugin.my.CalculateNextBlockTime(&account1, 100)
	assert.Equal(t, "2000-01-01T00:00:50.5", time.String())

	time = plugin.my.CalculateNextBlockTime(&account2, 100)
	assert.Nil(t, time)
}

func TestProducerPluginImpl_CalculatePendingBlockTime(t *testing.T) {
	plugin := producerPluginInitialize("-e", "-p", "eosio", "-p", "yuanc",
		"--private-key", "[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]",
		"--private-key", "[\"EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu\", \"5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP\"]")

	pt := plugin.my.CalculatePendingBlockTime()
	assert.Equal(t, pt/1e3, (pt/1e3/500)*500) // make sure pt can be divisible by 500ms
}

func TestProducerPluginImpl_OnIncomingBlock(t *testing.T) {
	plugin := producerPluginInitialize("-p", "eosio")
	//
	chain := plugin.my.Chain
	block := &types.SignedBlock{}
	block.Timestamp = types.NewBlockTimeStamp(*plugin.my.CalculateNextBlockTime(&eosio, chain.HeadBlockState().SignedBlock.Timestamp))
	block.Producer = common.N("eosio")
	block.Previous = chain.HeadBlockState().BlockId

	priKey, err := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	assert.NoError(t, err)

	blockrootMerkle := chain.HeadBlockState().BlockrootMerkle
	blockrootMerkle.Append(chain.HeadBlockState().BlockId)

	blockHash := block.Digest()
	scheduleHash := crypto.Hash256(types.ProducerScheduleType{Version: 0, Producers: []types.ProducerKey{{eosio, priKey.PublicKey()}}})
	headerBmroot := crypto.Hash256(common.MakePair(blockHash, blockrootMerkle.GetRoot()))
	digest := crypto.Hash256(common.MakePair(headerBmroot, scheduleHash))

	block.ProducerSignature, err = priKey.Sign(digest.Bytes())
	assert.NoError(t, err)

	plugin.my.OnIncomingBlock(block)

	assert.Equal(t, block.BlockNumber(), chain.HeadBlockNum())
}

func BenchmarkProducerPluginImpl_OnIncomingBlock(b *testing.B) {
	b.StopTimer()
	log.Root().SetHandler(log.DiscardHandler())
	plugin := producerPluginInitialize("-p", "eosio")
	chain := plugin.my.Chain

	for i := 0; i < b.N; i++ {
		block := &types.SignedBlock{}
		block.Timestamp = types.NewBlockTimeStamp(*plugin.my.CalculateNextBlockTime(&eosio, chain.HeadBlockState().SignedBlock.Timestamp))
		block.Producer = common.N("eosio")
		block.Previous = chain.HeadBlockState().BlockId

		priKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")

		blockrootMerkle := chain.HeadBlockState().BlockrootMerkle
		blockrootMerkle.Append(chain.HeadBlockState().BlockId)

		blockHash := block.Digest()
		scheduleHash := crypto.Hash256(types.ProducerScheduleType{Version: 0, Producers: []types.ProducerKey{{eosio, priKey.PublicKey()}}})
		headerBmroot := crypto.Hash256(common.MakePair(blockHash, blockrootMerkle.GetRoot()))
		digest := crypto.Hash256(common.MakePair(headerBmroot, scheduleHash))

		block.ProducerSignature, _ = priKey.Sign(digest.Bytes())

		b.StartTimer()
		plugin.my.OnIncomingBlock(block)
		b.StopTimer()
	}
}

func BenchmarkProducerPluginImpl_StartBlock(b *testing.B) {
	b.StopTimer()
	log.Root().SetHandler(log.DiscardHandler())

	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	plugin.PluginStartup()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		plugin.my.StartBlock()
	}
}

func BenchmarkProducerPluginImpl_ScheduleProductionLoop(b *testing.B) {
	b.StopTimer()
	log.Root().SetHandler(log.DiscardHandler())
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", \"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")

	plugin.PluginStartup()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		plugin.my.ScheduleProductionLoop()
		plugin.my.MaybeProduceBlock()
	}
}

func BenchmarkProducerPluginImpl_MaybeProduceBlock(b *testing.B) {
	b.StopTimer()
	log.Root().SetHandler(log.DiscardHandler())
	plugin := producerPluginInitialize("-e", "-p", "eosio", "--private-key",
		"[\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\", "+
			"\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]")
	plugin.PluginStartup()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		plugin.my.MaybeProduceBlock()
	}
}
