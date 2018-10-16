package producer_plugin

import (
	"fmt"
	Chain "github.com/eosspark/eos-go/plugins/producer_plugin/mock" /*Debug model*/
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"gopkg.in/urfave/cli.v1"
	"time"
	"encoding/json"
	"strings"
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
	KeyBlacklist      *map[ecc.PublicKey]struct{}
}

type GreylistParams struct {
	Accounts []common.AccountName
}

func init() {
	//TODO: initialize plugin for appbase
	fmt.Println("register plugin")
}

func NewProducerPlugin() ProducerPlugin {
	pp := new(ProducerPlugin)
	my := new(ProducerPluginImpl)

	my.Timer = new(common.Timer)
	my.Producers = make(map[common.AccountName]struct{})
	my.SignatureProviders = make(map[ecc.PublicKey]signatureProviderType)
	my.ProducerWatermarks = make(map[common.AccountName]uint32)

	my.PersistentTransactions = make(transactionIdWithExpireIndex)
	my.BlacklistedTransactions = make(transactionIdWithExpireIndex)

	my.IncomingTrxWeight = 0.0
	my.IncomingDeferRadio = 1.0 // 1:1
	my.Self = pp

	pp.my = my

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
		EosAssert(privateKeyFunc != nil, &ProducerPrivKeyNotFound{}, "Local producer has no private key in config.ini corresponding to public key %s", key)

		return privateKeyFunc(digest)
	}
	return ecc.Signature{}
}

func (pp *ProducerPlugin) PluginInitialize(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "enable-stale-production, e",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.BoolFlag{
			Name:  "pause-on-startup, x",
			Usage: "Start this node in a state where production is paused.",
		},
		cli.IntFlag{
			Name:  "max-transaction-age",
			Usage: "Limits the maximum time (in milliseconds) that is allowed a pushed transaction's code to execute before being considered invalid",
			Value: 30,
		},
		cli.IntFlag{
			Name:  "max-irreversible-block-age",
			Usage: "Limits the maximum age (in seconds) of the DPOS Irreversible Block for a chain this node will produce blocks on (use negative value to indicate unlimited)",
			Value: -1,
		},
		cli.StringSliceFlag{
			Name:  "producer-name, p",
			Usage: "ID of producer controlled by this node(e.g. inita; may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "private-key",
			Usage: "(DEPRECATED - Use signature-provider instead) Tuple of [public key, WIF private key] (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "signature-provider",
			Usage:
			"Key=Value pairs in the form <public-key>=<provider-spec>\n" 			   +
			"Where:\n" 																   +
			"   <public-key>    \tis a string form of a vaild EOSIO public key\n\n"    +
			"   <provider-spec> \tis a string in the form <provider-type>:<data>\n\n"  +
			"   <provider-type> \tis KEY, or KEOSD\n\n"                                +
			"   KEY:<data>      \tis a string form of a valid EOSIO private key which maps to the provided public key\n\n",
		},
		cli.IntFlag{
			Name:  "keosd-provider-timeout",
			Usage: "Limits the maximum time (in milliseconds) that is allowd for sending blocks to a keosd provider for signing",
			Value: 5,
		},
		cli.StringSliceFlag{
			Name:  "greylist-account",
			Usage: "account that can not access to extended CPU/NET virtual resources",
		},
		cli.IntFlag{
			Name:  "produce-time-offset-us",
			Usage: "offset of non last block producing time in micro second. Negative number results in blocks to go out sooner, and positive number results in blocks to go out later",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "last-block-time-offset-us",
			Usage: "offset of last block producing time in micro second. Negative number results in blocks to go out sooner, and positive number results in blocks to go out later",
			Value: 0,
		},
		cli.Float64Flag{
			Name:  "incoming-defer-ratio",
			Usage: "ratio between incoming transations and deferred transactions when both are exhausted",
			Value: 1.0,
		},
	}

	app.Action = func(c *cli.Context) {
		for _, p := range c.StringSlice("producer-name") {
			pp.my.Producers[common.AccountName(common.N(p))] = struct{}{}
		}

		for _, keyIdToWifPairString := range c.StringSlice("private-key") {
			var keyIdToWifPair [2]string
			json.Unmarshal([]byte(keyIdToWifPairString), &keyIdToWifPair)
			pubKey, err := ecc.NewPublicKey(keyIdToWifPair[0])
			if err != nil {
				panic(err)
			}
			priKey, err2 := ecc.NewPrivateKey(keyIdToWifPair[1])
			if err2 != nil {
				panic(err2)
			}
			pp.my.SignatureProviders[pubKey] = makeKeySignatureProvider(priKey)
		}

		for _, keySpecPair := range c.StringSlice("signature-provider") {
			delim := strings.Index(keySpecPair, "=")
			EosAssert(delim >= 0, &PluginConfigException{}, "Missing \"=\" in the key spec pair")
			pubKeyStr := keySpecPair[0:delim]
			specStr := keySpecPair[delim+1:]

			specDelim := strings.Index(specStr, ":")
			EosAssert(specDelim >= 0, &PluginConfigException{}, "Missing \":\" in the key spec pair")
			specTypeStr := specStr[0:specDelim]
			specData := specStr[specDelim+1:]

			pubKey, err := ecc.NewPublicKey(pubKeyStr)
			if err != nil { panic(err) }

			if specTypeStr == "KEY" {
				priKey, e := ecc.NewPrivateKey(specData)
				if e != nil { panic(nil) }
				pp.my.SignatureProviders[pubKey] = makeKeySignatureProvider(priKey)
			} else if specTypeStr == "KEOSD" {
				pp.my.SignatureProviders[pubKey] = makeKeosdSignatureProvider(pp.my, specData, pubKey)
			}
		}

		pp.my.ProductionEnabled = c.Bool("enable-stale-production")

		pp.my.ProductionPaused = c.Bool("pause-on-startup")

		pp.my.KeosdProviderTimeoutUs = common.Milliseconds(int64(c.Int("keosd-provider-timeout")))

		pp.my.ProduceTimeOffsetUs = int32(c.Int("produce-time-offset-us"))

		pp.my.MaxTransactionTimeMs = int32(c.Int("max-transaction-age"))

		pp.my.MaxIrreversibleBlockAgeUs = common.Seconds(int64(c.Int("max-irreversible-block-age")))

		pp.my.IncomingDeferRadio = c.Float64("incoming-defer-ratio")

		greylist := c.StringSlice("greylist-account")
		if len(greylist) > 0 {
			param := GreylistParams{}
			for _, a := range greylist {
				param.Accounts = append(param.Accounts, common.AccountName(common.N(a)))
			}
			pp.AddGreylistAccounts(param)
		}

	}
}

func (pp *ProducerPlugin) PluginStartup() {
	log.Info("producer plugin:  plugin_startup() begin")

	chain := Chain.GetControllerInstance()
	EosAssert(len(pp.my.Producers) == 0 || chain.GetReadMode() == Chain.DBReadMode(Chain.SPECULATIVE), &PluginConfigException{},
		"node cannot have any producer-name configured because block production is impossible when read_mode is not \"speculative\"" )

	EosAssert(len(pp.my.Producers) == 0 || chain.GetValidationMode() == Chain.ValidationMode(Chain.FULL), &PluginConfigException{},
		"node cannot have any producer-name configured because block production is not safe when validation_mode is not \"full\"" )

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
	code := e.Code()
	return (code == BlockCpuUsageExceeded{}.Code()) ||
		   (code == BlockCpuUsageExceeded{}.Code()) ||
		   (code == DeadlineException{}.Code() && deadlineIsSubjective)
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key *ecc.PrivateKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) ecc.Signature {
		sign, err := key.Sign(digest.Bytes())
		if err != nil {
			panic(err)
		}
		return sign
	}
	return signFunc
}

func makeKeosdSignatureProvider(produce *ProducerPluginImpl, url string, publicKey ecc.PublicKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) ecc.Signature {
		if produce != nil {
			//TODO
			return ecc.Signature{}
		} else {
			return ecc.Signature{}
		}
	}
	return signFunc
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

	if db.HeadBlockState().Header.Timestamp.ToTimePoint() < common.Now().SubUs(common.Microseconds(200*common.DefaultConfig.BlockIntervalMs)) {
		fmt.Print("Your genesis seems to have an old timestamp\n" +
			"Please consider using the --genesis-timestamp option to give your genesis a recent timestamp\n\n")
	}
}
