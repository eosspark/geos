package producer_plugin

import (
	"errors"
	"fmt"
	Chain "github.com/eosspark/eos-go/plugins/producer_plugin/mock" /*Debug model*/
	//Chain "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/log"
	"gopkg.in/urfave/cli.v1"
	"time"
	"encoding/json"
	. "github.com/eosspark/eos-go/exception"
)

type ProducerPlugin struct {
	my *ProducerPluginImpl
}

type RuntimeOptions struct {
	MaxTransactionTime      *int32
	MaxIrreversibleBlockAge *int32
	ProduceTimeOffsetUs     *int32
	LastBlockTimeOffsetUs   *int32
	SubjectiveCpuLeewayUs   *int32
	IncomingDeferRadio      *float64
}

type WhitelistAndBlacklist struct {
	ActorWhitelist    *map[common.AccountName]struct{}
	ActorBlacklist    *map[common.AccountName]struct{}
	ContractWhitelist *map[common.AccountName]struct{}
	ContractBlacklist *map[common.AccountName]struct{}
	ActionBlacklist   *map[[2]common.AccountName]struct{}
	KeyBlacklist      *map[common.PublicKeyType]struct{}
}

type GreylistParams struct {
	Accounts []common.AccountName
}

func NewProducerPlugin() ProducerPlugin {
	pp := new(ProducerPlugin)

	impl := new(ProducerPluginImpl)
	impl.Timer = new(common.Timer)
	impl.Producers = make(map[common.AccountName]struct{})
	impl.SignatureProviders = make(map[ecc.PublicKey]signatureProviderType)
	impl.ProducerWatermarks = make(map[common.AccountName]uint32)

	impl.PersistentTransactions = make(transactionIdWithExpireIndex)
	impl.BlacklistedTransactions = make(transactionIdWithExpireIndex)

	impl.IncomingTrxWeight = 0.0
	impl.IncomingDeferRadio = 1.0 // 1:1
	impl.Self = pp

	pp.my = impl
	return *pp
}

func (pp *ProducerPlugin) IsProducerKey(key ecc.PublicKey) bool {
	privateKey := pp.my.SignatureProviders[key]
	if privateKey != nil {
		return true
	}
	return false
}

func (pp *ProducerPlugin) SignCompact(key *ecc.PublicKey, digest crypto.Sha256) ecc.Signature {
	if key != nil {
		privateKeyFunc := pp.my.SignatureProviders[*key]
		if privateKeyFunc == nil {
			panic(ErrProducerPriKeyNotFound)
		}

		return privateKeyFunc(digest)
	}
	return ecc.Signature{}
}

func (pp *ProducerPlugin) PluginInitialize(app *cli.App) {
	//pp.signatureProviders[initPubKey] = func(hash []byte) ecc.Signature {
	//	sig, _ := initPriKey.Sign(hash)
	//	return sig
	//}

	// pp.my.SignatureProviders[initPubKey], _ = makeKeosdSignatureProvider(pp, "http://", initPubKey)

	var maxTransactionTimeMs int
	var maxIrreversibleBlockAgeUs int
	var privateKeys cli.StringSlice
	var producerNames cli.StringSlice

	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:        "enable-stale-production, e",
			Usage:       "Enable block production, even if the chain is stale.",
			Destination: &pp.my.ProductionEnabled,
		},
		cli.IntFlag{
			Name:        "max-transaction-age",
			Usage:       "Limits the maximum time (in milliseconds) that is allowed a pushed transaction's code to execute before being considered invalid",
			Value:       30,
			Destination: &maxTransactionTimeMs,
		},
		cli.IntFlag{
			Name:        "max-irreversible-block-age",
			Usage:       "Limits the maximum age (in seconds) of the DPOS Irreversible Block for a chain this node will produce blocks on (use negative value to indicate unlimited)",
			Value:       -1,
			Destination: &maxIrreversibleBlockAgeUs,
		},
		cli.StringSliceFlag{
			Name:  "producer-name, p",
			Usage: "ID of producer controlled by this node(e.g. inita; may specify multiple times)",
			Value: &producerNames,
		},
		cli.StringSliceFlag{
			Name:  "private-key",
			Usage: "ID of producer controlled by this node(e.g. inita; may specify multiple times)",
			Value: &privateKeys,
		},
	}

	app.Action = func(c *cli.Context) {
		pp.my.MaxTransactionTimeMs = int32(maxTransactionTimeMs)
		pp.my.MaxIrreversibleBlockAgeUs = common.Seconds(int64(maxIrreversibleBlockAgeUs))

		if len(producerNames) > 0 {
			for _, p := range producerNames {
				pp.my.Producers[common.AccountName(common.N(p))] = struct{}{}
			}
		}

		if len(privateKeys) > 0 {
			for _, pairString := range privateKeys {
				var pair [2]string
				json.Unmarshal([]byte(pairString), &pair)
				pubKey, err := ecc.NewPublicKey(pair[0])
				if err != nil {
					panic(err)
				}
				priKey, err2 := ecc.NewPrivateKey(pair[1])
				if err2 != nil {
					panic(err2)
				}
				pp.my.SignatureProviders[pubKey], _ = makeKeySignatureProvider(priKey)
			}
		}

		fmt.Println("max-transaction-age:", pp.my.MaxTransactionTimeMs)
		fmt.Println("max-irreversible-block-age:", pp.my.MaxIrreversibleBlockAgeUs)
		fmt.Println("producer-name:", pp.my.Producers)
	}
}

func (pp *ProducerPlugin) PluginStartup() {
	log.Info("producer plugin:  plugin_startup() begin")

	chain := Chain.GetControllerInstance()
	if !(len(pp.my.Producers) == 0 || chain.GetReadMode() == Chain.DBReadMode(Chain.SPECULATIVE)) {
		panic("node cannot have any producer-name configured because block production is impossible when read_mode is not \"speculative\"")
	}

	//TODO if
	// my->_accepted_block_connection.emplace(chain.accepted_block.connect( [this]( const auto& bsp ){ my->on_block( bsp ); } ));
	// my->_irreversible_block_connection.emplace(chain.irreversible_block.connect( [this]( const auto& bsp ){ my->on_irreversible_block( bsp->block ); } ));

	libNum := chain.LastIrreversibleBlockNum()
	lib := chain.FetchBlockByNumber(libNum)
	if lib != nil {
		pp.my.OnIrreversibleBlock(lib)
	} else {
		pp.my.IrreversibleBlockTime = common.MaxTimePoint()
	}

	if len(pp.my.Producers) > 0 {
		log.Info(fmt.Sprintf("Launching block production for %d producers at %s.", len(pp.my.Producers), common.Now()))

		if pp.my.ProductionEnabled {
			if chain.HeadBlockNum() == 0 {
				newChainBanner(chain)
			}
		}
	}

	pp.my.ScheduleProductionLoop()

	log.Info("producer plugin:  plugin_startup() end")
}

func (pp *ProducerPlugin) PluginShutdown() {
	pp.my.Timer.Cancel()
}

func (pp *ProducerPlugin) Pause() {
	pp.my.ProductionPaused = true
}

func (pp *ProducerPlugin) Resume() {
	pp.my.ProductionPaused = false
	// it is possible that we are only speculating because of this policy which we have now changed
	// re-evaluate that now
	//
	if pp.my.PendingBlockMode == EnumPendingBlockMode(speculating) {
		chain := Chain.GetControllerInstance()
		chain.AbortBlock()
		pp.my.ScheduleProductionLoop()
	}
}

func (pp *ProducerPlugin) Paused() bool {
	return pp.my.ProductionPaused
}

func (pp *ProducerPlugin) UpdateRuntimeOptions(options RuntimeOptions) {
	checkSpeculation := false

	if options.MaxTransactionTime != nil {
		pp.my.MaxTransactionTimeMs = *options.MaxTransactionTime
	}

	if options.MaxIrreversibleBlockAge != nil {
		pp.my.MaxIrreversibleBlockAgeUs = common.Seconds(int64(*options.MaxIrreversibleBlockAge))
		checkSpeculation = true
	}

	if options.ProduceTimeOffsetUs != nil {
		pp.my.ProduceTimeOffsetUs = *options.ProduceTimeOffsetUs
	}

	if options.LastBlockTimeOffsetUs != nil {
		pp.my.LastBlockTimeOffsetUs = *options.LastBlockTimeOffsetUs
	}

	if options.IncomingDeferRadio != nil {
		pp.my.IncomingDeferRadio = *options.IncomingDeferRadio
	}

	if checkSpeculation && pp.my.PendingBlockMode == EnumPendingBlockMode(speculating) {
		chain := Chain.GetControllerInstance()
		chain.AbortBlock()
		pp.my.ScheduleProductionLoop()
	}

	if options.SubjectiveCpuLeewayUs != nil {
		chain := Chain.GetControllerInstance()
		chain.SetSubjectiveCpuLeeway(common.Microseconds(*options.SubjectiveCpuLeewayUs))
	}
}

func (pp *ProducerPlugin) GetRuntimeOptions() RuntimeOptions {
	var maxIrreversibleBlockAge int32 = -1
	if pp.my.MaxIrreversibleBlockAgeUs.Count() >= 0 {
		maxIrreversibleBlockAge = int32(pp.my.MaxIrreversibleBlockAgeUs.Count() / 1e6)
	}
	return RuntimeOptions{
		&pp.my.MaxTransactionTimeMs,
		&maxIrreversibleBlockAge,
		&pp.my.ProduceTimeOffsetUs,
		&pp.my.LastBlockTimeOffsetUs,
		nil, nil,
	}
}

func (pp *ProducerPlugin) AddGreylistAccounts(params GreylistParams) {
	chain := Chain.GetControllerInstance()
	for _, acc := range params.Accounts {
		chain.AddResourceGreyList(&acc)
	}
}

func (pp *ProducerPlugin) RemoveGreylistAccounts(params GreylistParams) {
	chain := Chain.GetControllerInstance()
	for _, acc := range params.Accounts {
		chain.RemoveResourceGreyList(&acc)
	}
}

func (pp *ProducerPlugin) GetGreylist() GreylistParams {
	chain := Chain.GetControllerInstance()
	result := GreylistParams{}
	list := chain.GetResourceGreyList()
	result.Accounts = make([]common.AccountName, 0, len(*list))
	for acc := range *list {
		result.Accounts = append(result.Accounts, acc)
	}
	return result
}

func (pp *ProducerPlugin) GetWhitelistBlacklist() WhitelistAndBlacklist {
	chain := Chain.GetControllerInstance()
	return WhitelistAndBlacklist{
		chain.GetActorWhiteList(),
		chain.GetActorBlackList(),
		chain.GetContractWhiteList(),
		chain.GetContractBlackList(),
		chain.GetActionBlockList(),
		chain.GetKeyBlackList(),
	}
}

func (pp *ProducerPlugin) SetWhitelistBlacklist(params WhitelistAndBlacklist) {
	chain := Chain.GetControllerInstance()
	if params.ActorWhitelist != nil {
		chain.SetActorWhiteList(params.ActorWhitelist)
	}
	if params.ActorBlacklist != nil {
		chain.SetActorBlackList(params.ActorBlacklist)
	}
	if params.ContractWhitelist != nil {
		chain.SetContractWhiteList(params.ContractWhitelist)
	}
	if params.ContractBlacklist != nil {
		chain.SetContractBlackList(params.ContractBlacklist)
	}
	if params.ActionBlacklist != nil {
		chain.SetActionBlackList(params.ActionBlacklist)
	}
	if params.KeyBlacklist != nil {
		chain.SetKeyBlackList(params.KeyBlacklist)
	}
}

func failureIsSubjective(e Exception, deadlineIsSubjective bool) bool {
	//TODO wait for error definition
	return false
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key *ecc.PrivateKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest crypto.Sha256) (sign ecc.Signature) {
		sign, err = key.Sign(digest.Bytes())
		return
	}
	return
}

func makeKeosdSignatureProvider(produce *ProducerPlugin, url string, pubKey ecc.PublicKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest crypto.Sha256) ecc.Signature {
		if produce != nil {
			//TODO
			return ecc.Signature{}
		} else {
			return ecc.Signature{}
		}
	}
	return
}

func newChainBanner(db *Chain.Controller) {
	fmt.Print("\n" +
		"*******************************\n" +
		"*                             *\n" +
		"*   ------ NEW CHAIN ------   *\n" +
		"*   -  Welcome to EOSIO!  -   *\n" +
		"*   -----------------------   *\n" +
		"*                             *\n" +
		"*******************************\n" +
		"\n")

	if db.HeadBlockState().Header.Timestamp.ToTimePoint() < common.Now().SubUs(common.Microseconds(200 * common.DefaultConfig.BlockIntervalMs)) {
		fmt.Print("Your genesis seems to have an old timestamp\n" +
			"Please consider using the --genesis-timestamp option to give your genesis a recent timestamp\n\n" )
	}
}

//errors
var (
	ErrProducerFail             = errors.New("called produce_block while not actually producing")
	ErrMissingPendingBlockState = errors.New("pending_block_state does not exist but it should, another plugin may have corrupted it")
	ErrProducerPriKeyNotFound   = errors.New("attempting to produce a block for which we don't have the private key")
	ErrBlockFromTheFuture       = errors.New("received a block from the future, ignoring it")
)
