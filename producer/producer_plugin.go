package producer_plugin

import (
	"errors"
	"fmt"
	Chain "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/rlp"
	"gopkg.in/urfave/cli.v1"
	"time"
)

type EnumPendingBlockMode int

const (
	producing = EnumPendingBlockMode(iota)
	speculating
)

type EnumStartBlockRusult int

const (
	succeeded = EnumStartBlockRusult(iota)
	failed
	waiting
	exhausted
)

type signatureProviderType func(sha256 rlp.Sha256) ecc.Signature

type transactionIdWithExpireIndex map[common.TransactionIDType]common.TimePoint

type ProducerPlugin struct {
	timer              *common.Timer
	producers          map[common.AccountName]struct{}
	pendingBlockMode   EnumPendingBlockMode
	productionEnabled  bool
	productionPaused   bool
	signatureProviders map[ecc.PublicKey]signatureProviderType
	producerWatermarks map[common.AccountName]uint32

	persistentTransactions  transactionIdWithExpireIndex
	blacklistedTransactions transactionIdWithExpireIndex

	maxTransactionTimeMs      int32
	maxIrreversibleBlockAgeUs common.Microseconds
	produceTImeOffsetUs       int32
	lastBlockTimeOffsetUs     int32
	irreversibleBlockTime     common.TimePoint
	keosdProviderTimeoutUs    common.Microseconds

	lastSignedBlockTime common.TimePoint
	startTime           common.TimePoint
	lastSignedBlockNum  uint32

	confirmedBlock func(signature ecc.Signature) //TODO

	pendingIncomingTransactions []pendingIncomingTransaction

	// keep a expected ratio between defer txn and incoming txn
	incomingTrxWeight  float64
	incomingDeferRadio float64
}

type RuntimeOptions struct {
	MaxTransactionTime      int32
	MaxIrreversibleBlockAge int32
	ProduceTimeOffsetUs     int32
	LastBlockTimeOffsetUs   int32
	SubjectiveCpuLeewayUs   int32
	IncomingDeferRadio      float64
}

type WhitelistAndBlacklist struct {
	ActorWhitelist    map[common.AccountName]struct{}
	ActorBlacklist    map[common.AccountName]struct{}
	ContractWhitelist map[common.AccountName]struct{}
	ContractBlacklist map[common.AccountName]struct{}
	ActionBlacklist   map[[2]common.Name]struct{}
	KeyBlacklist      map[ecc.PublicKey]struct{}
}

type GreylistParams struct {
	Accounts []common.AccountName
}

func (pp *ProducerPlugin) init() {
	pp.timer = new(common.Timer)
	pp.producers = make(map[common.AccountName]struct{})
	pp.signatureProviders = make(map[ecc.PublicKey]signatureProviderType)
	pp.producerWatermarks = make(map[common.AccountName]uint32)

	pp.persistentTransactions = make(transactionIdWithExpireIndex)
	pp.blacklistedTransactions = make(transactionIdWithExpireIndex)

	pp.incomingTrxWeight = 0.0
	pp.incomingDeferRadio = 1.0 // 1:1
}

func (pp *ProducerPlugin) IsProducerKey(key ecc.PublicKey) bool {
	privateKey := pp.signatureProviders[key]
	if privateKey != nil {
		return true
	}
	return false
}

func (pp *ProducerPlugin) SignCompact(key *ecc.PublicKey, digest rlp.Sha256) ecc.Signature {
	if key != nil {
		privateKeyFunc := pp.signatureProviders[*key]
		if privateKeyFunc == nil {
			panic(ErrProducerPriKeyNotFound)
		}

		return privateKeyFunc(digest)
	}
	return ecc.Signature{}
}

func (pp *ProducerPlugin) Initialize(app *cli.App) {
	pp.init()

	//pp.signatureProviders[initPubKey] = func(hash []byte) ecc.Signature {
	//	sig, _ := initPriKey.Sign(hash)
	//	return sig
	//}

	pp.signatureProviders[initPubKey], _ = makeKeySignatureProvider(*initPriKey)
	pp.signatureProviders[initPubKey], _ = makeKeosdSignatureProvider(pp, "http://", initPubKey)

	var maxTransactionTimeMs int
	var maxIrreversibleBlockAgeUs int
	var producerName cli.StringSlice

	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:        "enable-stale-production, e",
			Usage:       "Enable block production, even if the chain is stale.",
			Destination: &pp.productionEnabled,
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
			Value: &producerName,
		},
	}

	app.Action = func(c *cli.Context) {
		pp.maxTransactionTimeMs = int32(maxTransactionTimeMs)
		pp.maxIrreversibleBlockAgeUs = common.Seconds(int64(maxIrreversibleBlockAgeUs))

		if len(producerName) > 0 {
			for _, p := range producerName {
				pp.producers[common.AccountName(common.StringToName(p))] = struct{}{}
			}
		}

		fmt.Println("max-transaction-age:", pp.maxTransactionTimeMs)
		fmt.Println("max-irreversible-block-age:", pp.maxIrreversibleBlockAgeUs)
		fmt.Println("producer-name:", pp.producers)
	}
}

func (pp *ProducerPlugin) Startup() {
	log.Info("producer plugin:  plugin_startup() begin")

	if !(len(pp.producers) == 0 || chain.GetReadMode() == Chain.DBReadMode(Chain.SPECULATIVE)) {
		panic("node cannot have any producer-name configured because block production is impossible when read_mode is not \"speculative\"")
	}

	//TODO if
	// my->_accepted_block_connection.emplace(chain.accepted_block.connect( [this]( const auto& bsp ){ my->on_block( bsp ); } ));
	// my->_irreversible_block_connection.emplace(chain.irreversible_block.connect( [this]( const auto& bsp ){ my->on_irreversible_block( bsp->block ); } ));

	libNum := chain.LastIrreversibleBlockNum()
	lib := chain.FetchBlockByNumber(libNum)
	if lib != nil {
		pp.onIrreversibleBlock(lib)
	} else {
		pp.irreversibleBlockTime = common.MaxTimePoint()
	}

	if len(pp.producers) > 0 {
		log.Info(fmt.Sprintf("Launching block production for %d producers at %s.", len(pp.producers), common.Now()))

		if pp.productionEnabled {
			if chain.HeadBlockNum() == 0 {
				newChainBanner(Chain.Controller{}) //TODO
			}
		}
	}

	pp.scheduleProductionLoop()

	log.Info("producer plugin:  plugin_startup() end")
}

func (pp *ProducerPlugin) Shutdown() {
	pp.timer.Cancel()
}

func (pp *ProducerPlugin) Pause() {
	pp.productionPaused = true
}

func (pp *ProducerPlugin) Resume() {
	pp.productionPaused = false
	// it is possible that we are only speculating because of this policy which we have now changed
	// re-evaluate that now
	//
	if pp.pendingBlockMode == EnumPendingBlockMode(speculating) {
		chain.AbortBlock()
		pp.scheduleProductionLoop()
	}
}

func (pp *ProducerPlugin) Paused() bool {
	return pp.productionPaused
}

func failureIsSubjective(e error, deadlineIsSubjective bool) bool {
	//TODO wait for error definition
	return false
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key ecc.PrivateKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest rlp.Sha256) (sign ecc.Signature) {
		sign, err = key.Sign(digest.Bytes())
		return
	}
	return
}

func makeKeosdSignatureProvider(produce *ProducerPlugin, url string, pubKey ecc.PublicKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest rlp.Sha256) ecc.Signature {
		if produce != nil {
			//TODO
			return ecc.Signature{}
		} else {
			return ecc.Signature{}
		}
	}
	return
}

func newChainBanner(db Chain.Controller) {
	fmt.Print("\n" +
		"*******************************\n" +
		"*                             *\n" +
		"*   ------ NEW CHAIN ------   *\n" +
		"*   -  Welcome to EOSIO!  -   *\n" +
		"*   -----------------------   *\n" +
		"*                             *\n" +
		"*******************************\n" +
		"\n")

	//TODO
}

//errors
var (
	ErrProducerFail             = errors.New("called produce_block while not actually producing")
	ErrMissingPendingBlockState = errors.New("pending_block_state does not exist but it should, another plugin may have corrupted it")
	ErrProducerPriKeyNotFound   = errors.New("attempting to produce a block for which we don't have the private key")
	ErrBlockFromTheFuture       = errors.New("received a block from the future, ignoring it")
)
