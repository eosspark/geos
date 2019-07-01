package unittests

import (
	"fmt"
	"io/ioutil"
	"math"

	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/log"
)

var producer = common.N("producer1111")
var producer1 = common.N("defproducera")
var producer2 = common.N("defproducerb")
var producer3 = common.N("defproducerc")
var producer4 = common.N("defproducerd")
var producer5 = common.N("defproducere")
var voter1 = common.N("producvotera")
var voter2 = common.N("producvoterb")
var voter3 = common.N("producvoterc")
var voter4 = common.N("producvoterd")

type EosioSystemTester struct {
	ValidatingTester
	abiSer      abi_serializer.AbiSerializer
	tokenAbiSer abi_serializer.AbiSerializer
}

func newEosioSystemTester(pushGenesis bool, readMode DBReadMode) *EosioSystemTester {
	e := &EosioSystemTester{}
	e.DefaultExpirationDelta = 6
	e.DefaultBilledCpuTimeUs = 2000
	e.AbiSerializerMaxTime = 1000 * 1000
	e.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	e.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)
	e.VCfg = *newConfig(readMode)

	e.ValidatingControl = NewController(&e.VCfg)
	e.ValidatingControl.Startup()
	e.init(true, readMode)
	return e
}

func initEosioSystemTester() *EosioSystemTester {
	e := newEosioSystemTester(true, SPECULATIVE)

	e.ProduceBlocks(2, false)
	e.CreateAccounts([]common.AccountName{eosioToken, eosioRam, eosioRamFee, eosioStake,
		eosioBpay, eosioVpay, eosioSaving, eosioName}, false, true)
	e.ProduceBlocks(100, false)

	//eosio.token
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(eosioToken, abi, nil)
	accnt := entity.AccountObject{Name: eosioToken}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	e.tokenAbiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)
	e.CreateCurrency(eosioToken, common.DefaultConfig.SystemAccountName, CoreFromString("10000000000.0000"))
	e.Issue(common.DefaultConfig.SystemAccountName, CoreFromString("1000000000.0000"), common.DefaultConfig.SystemAccountName)
	currencyBalance := e.GetBalance(eosio)
	expectedBalance := CoreFromString("1000000000.0000")
	if currencyBalance != expectedBalance {
		log.Error("error, initEosioSystemTester failed")
	}

	//eosio.system
	wasmName = "test_contracts/eosio.system.wasm"
	code, _ = ioutil.ReadFile(wasmName)
	e.SetCode(eosio, code, nil)
	abiName = "test_contracts/eosio.system.abi"
	abi, _ = ioutil.ReadFile(abiName)
	e.SetAbi(eosio, abi, nil)
	abiDef = abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)
	accnt = entity.AccountObject{Name: eosio}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef = abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	//TODO
	//e.abiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)

	e.ProduceBlocks(1, false)

	e.CreateAccountWithResources(
		alice,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		bob,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("0.4500"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		carol,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	if CoreFromString("1000000000.0000").Amount != e.GetBalance(eosio).Amount+e.GetBalance(eosioRamFee).Amount+
		e.GetBalance(eosioStake).Amount+e.GetBalance(eosioRam).Amount {
		log.Error("error")
	}
	return e
}

func (e EosioSystemTester) CreateAccountWithResources(name common.AccountName, creator common.AccountName,
	ramFunds common.Asset, multiSig bool, net common.Asset, cpu common.Asset) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{
			Threshold: 2,
			Keys:      []types.KeyWeight{{Key: e.getPublicKey(name, "owner"), Weight: 1}},
			Accounts:  []types.PermissionLevelWeight{{Permission: common.PermissionLevel{Actor: creator, Permission: common.DefaultConfig.ActiveName}, Weight: 1}},
		}
	} else {
		ownerAuth = types.NewAuthority(e.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(e.getPublicKey(name, "active"), 0)

	newAccount := NewAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(newAccount)
	act := &types.Action{
		Account:       newAccount.GetAccount(),
		Name:          newAccount.GetName(),
		Authorization: []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	buyRamData := common.Variants{
		"payer":    creator,
		"receiver": name,
		"quant":    ramFunds,
	}
	buyRam := e.GetAction(
		eosio,
		common.N("buyram"),
		[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		&buyRamData,
	)
	trx.Actions = append(trx.Actions, buyRam)

	delegateData := common.Variants{
		"from":               creator,
		"receiver":           name,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           0,
	}
	delegate := e.GetAction(
		eosio,
		common.N("delegatebw"),
		[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		&delegateData,
	)
	trx.Actions = append(trx.Actions, delegate)

	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) CreateAccountsWithResources(accounts []common.AccountName, creator common.AccountName) {
	for _, a := range accounts {
		e.CreateAccountWithResources2(a, creator, 8000)
	}
}

func (e EosioSystemTester) CreateAccountWithResources2(name common.AccountName, creator common.AccountName, ramBytes uint32) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)

	ownerAuth := types.NewAuthority(e.getPublicKey(name, "owner"), 0)
	activeAuth := types.NewAuthority(e.getPublicKey(name, "active"), 0)

	new := NewAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(new)
	act := &types.Action{
		Account:       new.GetAccount(),
		Name:          new.GetName(),
		Authorization: []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	buyRamBytesData := common.Variants{
		"payer":    creator,
		"receiver": name,
		"bytes":    ramBytes,
	}
	buyRam := e.GetAction(
		eosio,
		common.N("buyrambytes"),
		[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		&buyRamBytesData,
	)
	trx.Actions = append(trx.Actions, buyRam)

	delegateData := common.Variants{
		"from":               creator,
		"receiver":           name,
		"stake_net_quantity": CoreFromString("10.0000"),
		"stake_cpu_quantity": CoreFromString("10.0000"),
		"transfer":           0,
	}
	delegate := e.GetAction(
		eosio,
		common.N("delegatebw"),
		[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
		&delegateData,
	)
	trx.Actions = append(trx.Actions, delegate)

	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) SetupProducerAccounts(accounts []common.AccountName) *types.TransactionTrace {
	creator := eosio
	trx := types.SignedTransaction{}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	cpu := CoreFromString("80.0000")
	net := CoreFromString("80.0000")
	ram := CoreFromString("1.0000")
	for _, a := range accounts {
		ownerAuth := types.NewAuthority(e.getPublicKey(a, "owner"), 0)
		activeAuth := types.NewAuthority(e.getPublicKey(a, "active"), 0)
		newAccount := NewAccount{
			Creator: creator,
			Name:    a,
			Owner:   ownerAuth,
			Active:  activeAuth,
		}
		data, _ := rlp.EncodeToBytes(newAccount)
		newAccountAct := &types.Action{
			Account:       newAccount.GetAccount(),
			Name:          newAccount.GetName(),
			Authorization: []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
			Data:          data,
		}
		trx.Actions = append(trx.Actions, newAccountAct)

		buyRamData := common.Variants{
			"payer":    creator,
			"receiver": a,
			"quant":    ram,
		}
		buyRam := e.GetAction(
			eosio,
			common.N("buyram"),
			[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
			&buyRamData,
		)
		trx.Actions = append(trx.Actions, buyRam)

		delegateData := common.Variants{
			"from":               creator,
			"receiver":           a,
			"stake_net_quantity": net,
			"stake_cpu_quantity": cpu,
			"transfer":           0,
		}
		delegate := e.GetAction(
			eosio,
			common.N("delegatebw"),
			[]common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}},
			&delegateData,
		)
		trx.Actions = append(trx.Actions, delegate)
	}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) BuyRam(payer common.AccountName, receiver common.AccountName, eosin common.Asset) ActionResult {
	buyRam := common.Variants{
		"payer":    payer,
		"receiver": receiver,
		"quant":    eosin,
	}
	act := common.N("buyram")
	return e.EsPushAction(&payer, &act, &buyRam, true)
}

func (e EosioSystemTester) BuyRamBytes(payer common.AccountName, receiver common.AccountName, numBytes uint32) ActionResult {
	buyRamBytes := common.Variants{
		"payer":    payer,
		"receiver": receiver,
		"bytes":    numBytes,
	}
	act := common.N("buyrambytes")
	return e.EsPushAction(&payer, &act, &buyRamBytes, true)
}

func (e EosioSystemTester) SellRam(account common.AccountName, numBytes uint64) ActionResult {
	sellRam := common.Variants{
		"account": account,
		"bytes":   numBytes,
	}
	act := common.N("sellram")
	return e.EsPushAction(&account, &act, &sellRam, true)
}

func (e EosioSystemTester) EsPushAction(signer *common.AccountName, name *common.ActionName, data *common.Variants, auth bool) ActionResult {
	var authorizer common.AccountName
	if auth == true {
		authorizer = *signer
	} else {
		if *signer == bob {
			authorizer = alice
		} else {
			authorizer = bob
		}
	}
	act := e.GetAction(eosio, *name, []common.PermissionLevel{}, data)
	return e.PushAction(act, common.AccountName(authorizer))
}

func (e EosioSystemTester) Stake(from common.AccountName, to common.AccountName, net common.Asset, cpu common.Asset) ActionResult {
	stake := common.Variants{
		"from":               from,
		"receiver":           to,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           0,
	}
	act := common.N("delegatebw")
	return e.EsPushAction(&from, &act, &stake, true)
}

func (e EosioSystemTester) StakeWithTransfer(from common.AccountName, to common.AccountName, net common.Asset, cpu common.Asset) ActionResult {
	stake := common.Variants{
		"from":               from,
		"receiver":           to,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           true,
	}
	act := common.N("delegatebw")
	return e.EsPushAction(&from, &act, &stake, true)
}

func (e EosioSystemTester) UnStake(from common.AccountName, to common.AccountName, net common.Asset, cpu common.Asset) ActionResult {
	unStake := common.Variants{
		"from":                 from,
		"receiver":             to,
		"unstake_net_quantity": net,
		"unstake_cpu_quantity": cpu,
	}
	act := common.N("undelegatebw")
	return e.EsPushAction(&from, &act, &unStake, true)
}

func (e EosioSystemTester) BidName(bidder common.AccountName, newName common.AccountName, bid common.Asset) ActionResult {
	bidName := common.Variants{
		"bidder":  bidder,
		"newname": newName,
		"bid":     bid,
	}
	act := common.N("bidname")
	return e.EsPushAction(&bidder, &act, &bidName, true)
}

func (e EosioSystemTester) SetRam(signer common.AccountName, maxRamSize uint64) ActionResult {
	setRam := common.Variants{
		"max_ram_size": maxRamSize,
	}
	act := common.N("setram")
	return e.EsPushAction(&signer, &act, &setRam, true)
}

func (e EosioSystemTester) RmvProducer(signer common.AccountName, prodName common.AccountName) ActionResult {
	rmvProducer := common.Variants{
		"producer": prodName,
	}
	act := common.N("rmvproducer")
	return e.EsPushAction(&signer, &act, &rmvProducer, true)
}

func (e EosioSystemTester) RegProducer(acnt common.AccountName) ActionResult {
	regproducer := common.Variants{
		"producer":     acnt,
		"producer_key": e.getPublicKey(acnt, "active"),
		"url":          "",
		"location":     0,
	}
	act := common.N("regproducer")
	r := e.EsPushAction(&acnt, &act, &regproducer, true)
	if r != e.Success() {
		log.Error("Wrong: Push action regproducer failed.")
	}
	return r
}

func (e EosioSystemTester) RegProxy(acnt common.AccountName) ActionResult {
	act := common.N("regproxy")
	regproxy := common.Variants{
		"proxy":   acnt,
		"isproxy": true,
	}
	return e.EsPushAction(&acnt, &act, &regproxy, true)
}

func (e EosioSystemTester) UnRegProxy(acnt common.AccountName) ActionResult {
	act := common.N("regproxy")
	regproxy := common.Variants{
		"proxy":   acnt,
		"isproxy": false,
	}
	return e.EsPushAction(&acnt, &act, &regproxy, true)
}

func (e EosioSystemTester) Vote(voter common.AccountName, producers []common.AccountName, proxy common.AccountName) ActionResult {
	vote := common.Variants{
		"voter":     voter,
		"proxy":     proxy,
		"producers": producers,
	}
	act := common.N("voteproducer")
	r := e.EsPushAction(&voter, &act, &vote, true)
	return r
}

func (e EosioSystemTester) ClaimRewards(producer common.AccountName) ActionResult {
	claimrewards := common.Variants{
		"owner": producer,
	}
	act := common.N("claimrewards")
	r := e.EsPushAction(&producer, &act, &claimrewards, true)
	return r
}

func (e EosioSystemTester) GetBalance(act common.AccountName) common.Asset {
	PrimaryKey := uint64(CORE_SYMBOL.ToSymbolCode())
	data := e.GetRowByAccount(uint64(eosioToken), uint64(act), uint64(common.N("accounts")), PrimaryKey)
	if len(data) == 0 {
		return common.Asset{Amount: 0, Symbol: CORE_SYMBOL}
	} else {
		asset := common.Asset{}
		rlp.DecodeBytes(data, &asset)
		return asset
	}
}

func (e EosioSystemTester) GetTotalStake(act common.AccountName) common.Variants {
	type UserResources struct {
		Owner     common.AccountName
		NetWeight common.Asset
		CpuWeight common.Asset
		RamBytes  uint64
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(act), uint64(common.N("userres")), uint64(act))
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := UserResources{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"owner":      res.Owner,
			"net_weight": res.NetWeight,
			"cpu_weight": res.CpuWeight,
			"ram_bytes":  res.RamBytes,
		}
	}
}

func (e EosioSystemTester) GetVoterInfo(act common.AccountName) common.Variants {
	type VoterInfo struct {
		Owner             common.AccountName
		Proxy             common.AccountName
		Producers         []common.AccountName
		Staked            int64
		LastVoteWeight    float64
		ProxiedVoteWeight float64
		IsProxy           bool
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(eosio), uint64(common.N("voters")), uint64(act))
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := VoterInfo{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"owner":     res.Owner,
			"proxy":     res.Proxy,
			"producers": res.Producers,
			"staked":    res.Staked,
			//"last_vote_weight":    res.LastVoteWeight,
			"proxied_vote_weight": res.ProxiedVoteWeight,
			"is_proxy":            res.IsProxy,
		}
	}
}

func (e EosioSystemTester) GetProducerInfo(act common.AccountName) common.Variants {
	type ProducerInfo struct {
		Owner         common.AccountName
		TotalVotes    float64
		ProducerKey   ecc.PublicKey
		IsActive      bool
		Url           string
		UnpaidBlocks  uint32
		LastClaimTime uint64
		Location      uint16
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(eosio), uint64(common.N("producers")), uint64(act))
	res := ProducerInfo{}
	err := rlp.DecodeBytes(data, &res)
	if err != nil {
		fmt.Println(err)
	}
	return common.Variants{
		"owner":           res.Owner,
		"total_votes":     res.TotalVotes,
		"producer_key":    res.ProducerKey,
		"is_active":       res.IsActive,
		"url":             res.Url,
		"unpaid_blocks":   res.UnpaidBlocks,
		"last_claim_time": res.LastClaimTime,
		"location":        res.Location,
	}
}

func (e EosioSystemTester) CreateCurrency(contract common.Name, manager common.Name, maxSupply common.Asset) {
	act := common.Variants{
		"issuer":         manager,
		"maximum_supply": maxSupply,
	}
	acttype := common.N("create")
	e.PushAction2(
		&contract,
		&acttype,
		contract,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}

func (e EosioSystemTester) Issue(to common.Name, amount common.Asset, manager common.Name) {
	act := common.Variants{
		"to":       to,
		"quantity": amount,
		"memo":     "issue",
	}
	acttype := common.N("issue")
	contract := eosioToken
	e.PushAction2(
		&contract,
		&acttype,
		manager,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}

func (e EosioSystemTester) Transfer(from common.Name, to common.Name, amount common.Asset, manager common.Name) {
	act := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": amount,
		"memo":     "transfer",
	}
	acttype := common.N("transfer")
	contract := eosioToken
	e.PushAction2(
		&contract,
		&acttype,
		manager,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}

func (e EosioSystemTester) Stake2Votes(stake common.Asset) float64 {
	now := e.Control.PendingBlockTime().TimeSinceEpoch().Count() / 1000000
	return float64(stake.Amount) * math.Pow(float64(2), float64((now-common.DefaultConfig.BlockTimestampEpochMs/1000)/(86400*7))/float64(52))
}

func (e EosioSystemTester) GetStats(symbol common.Symbol) common.Variants {
	type CurrencyStats struct {
		Supply    common.Asset
		MaxSupply common.Asset
		Issuer    common.AccountName
	}
	symbolCode := symbol.ToSymbolCode()
	data := e.GetRowByAccount(uint64(eosioToken), uint64(symbolCode), uint64(common.N("stat")), uint64(symbolCode))
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := CurrencyStats{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"supply":     res.Supply,
			"max_supply": res.MaxSupply,
			"issuer":     res.Issuer,
		}
	}
}

func (e EosioSystemTester) GetTokenSupply() common.Asset {
	return e.GetStats(CORE_SYMBOL)["supply"].(common.Asset)
}

func (e EosioSystemTester) GetGlobalState() common.Variants {
	type EosioGlobalState struct {
		types.ChainConfig
		MaxRamSize                 uint64
		TotalRamBytesReserved      uint64
		TotalRamStake              int64
		LastProducerScheduleUpdate types.BlockTimeStamp
		LastPervoteBucketFill      uint64
		PervoteBucket              int64
		PerblockBucket             int64
		TotalUnpaidBlocks          uint32
		TotalActivatedStake        int64
		ThreshActivatedStakeTime   uint64
		LastProducerScheduleSize   uint16
		TotalProducerVoteWeight    float64
		LastNameClose              types.BlockTimeStamp
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(eosio), uint64(common.N("global")), uint64(common.N("global")))
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := EosioGlobalState{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"max_ram_size":                  res.MaxRamSize,
			"total_ram_bytes_reserved":      res.TotalRamBytesReserved,
			"total_ram_stake":               res.TotalRamStake,
			"last_producer_schedule_update": res.LastProducerScheduleUpdate,
			"last_pervote_bucket_fill":      res.LastPervoteBucketFill,
			"pervote_bucket":                res.PervoteBucket,
			"perblock_bucket":               res.PerblockBucket,
			"total_unpaid_blocks":           res.TotalUnpaidBlocks,
			"total_activated_stake":         res.TotalActivatedStake,
			"thresh_activated_stake_time":   res.ThreshActivatedStakeTime,
			"last_producer_schedule_size":   res.LastProducerScheduleSize,
			"total_producer_vote_weight":    res.TotalProducerVoteWeight,
			"last_name_close":               res.LastNameClose,
		}
	}
}

func (e EosioSystemTester) GetRefundRequest(account common.AccountName) common.Variants {
	type RefundRequest struct {
		Owner       common.AccountName
		RequestTime common.TimePointSec
		NetAmount   common.Asset
		CpuAmount   common.Asset
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(account), uint64(common.N("refunds")), uint64(account))
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := RefundRequest{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"owner":        res.Owner,
			"request_time": res.RequestTime,
			"net_amount":   res.NetAmount,
			"cpu_amount":   res.CpuAmount,
		}
	}
}

func (e EosioSystemTester) InitializeMultisig() abi_serializer.AbiSerializer {
	var msigAbiSer abi_serializer.AbiSerializer
	e.CreateAccountWithResources2(eosioMsig, eosio, 8000)
	e.BuyRam(eosio, eosioMsig, CoreFromString("5000.0000"))
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	data := common.Variants{
		"account": eosioMsig,
		"is_priv": 1,
	}
	acttype := common.N("setpriv")
	e.PushAction2(&eosio, &acttype, eosio, &data, e.DefaultExpirationDelta, 0)
	wasmName := "test_contracts/eosio.msig.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(eosioMsig, code, nil)
	abiName := "test_contracts/eosio.msig.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(eosioMsig, abi, nil)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	accnt := entity.AccountObject{Name: eosioMsig}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::InitializeMultisig failed with ToAbi")
	}
	msigAbiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)
	return msigAbiSer
}

func (e EosioSystemTester) ActiveAndVoteProducers() []common.AccountName {
	//stake more than 15% of total EOS supply to activate chain
	e.Transfer(eosio, alice, CoreFromString("650000000.0000"), eosio)
	e.Stake(alice, alice, CoreFromString("300000000.0000"), CoreFromString("300000000.0000"))

	// create accounts {defproducera, defproducerb, ..., defproducerz} and register as producers
	var producersNames []common.AccountName
	root := "defproducer"
	for c := 'a'; c < 'a'+21; c++ {
		acc := common.N(root + string(c))
		producersNames = append(producersNames, acc)
	}
	e.SetupProducerAccounts(producersNames)

	for _, a := range producersNames {
		e.RegProducer(a)
	}
	e.ProduceBlocks(250, false)
	auth := types.Authority{
		Threshold: 1,
		Keys:      []types.KeyWeight{{Key: e.getPublicKey(eosio, "active"), Weight: 1}},
		Accounts: []types.PermissionLevelWeight{
			{Permission: common.PermissionLevel{Actor: eosio, Permission: common.DefaultConfig.EosioCodeName}, Weight: 1},
			{Permission: common.PermissionLevel{Actor: common.DefaultConfig.ProducersAccountName, Permission: common.DefaultConfig.ActiveName}, Weight: 1},
		},
		Waits: []types.WaitWeight{},
	}
	data := common.Variants{
		"account":    eosio,
		"permission": common.N("active"),
		"parent":     common.N("owner"),
		"auth":       auth,
	}

	actName := UpdateAuth{}.GetName()
	traceAuth := e.PushAction2(
		&eosio,
		&actName,
		eosio,
		&data,
		e.DefaultExpirationDelta,
		0,
	)
	if traceAuth.Receipt.Status != types.TransactionStatusExecuted {
		log.Error("wrong: updateAuth failed.")
	}

	//vote for producers
	e.Transfer(eosio, alice, CoreFromString("100000000.0000"), eosio)
	e.Stake(alice, alice, CoreFromString("30000000.0000"), CoreFromString("30000000.0000"))
	e.BuyRam(alice, alice, CoreFromString("30000000.0000"))
	e.Vote(alice, producersNames[0:21], common.AccountName(0))
	e.ProduceBlocks(250, false)

	producerKeys := e.Control.HeadBlockState().ActiveSchedule.Producers
	if len(producerKeys) != 21 || producerKeys[0].ProducerName != producer1 {
		log.Error("wrong: update producers failed.")
	}
	return producersNames
}

func (e EosioSystemTester) Cross15PercentThreshold() {
	e.SetupProducerAccounts([]common.AccountName{producer})
	e.RegProducer(producer)
	{
		trx := types.SignedTransaction{}
		e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
		delegatebwData := common.Variants{
			"from":               eosio,
			"receiver":           producer,
			"stake_net_quantity": CoreFromString("150000000.0000"),
			"stake_cpu_quantity": CoreFromString("0.0000"),
			"transfer":           1,
		}
		delegate := e.GetAction(
			eosio,
			common.N("delegatebw"),
			[]common.PermissionLevel{{Actor: eosio, Permission: common.DefaultConfig.ActiveName}},
			&delegatebwData,
		)
		trx.Actions = append(trx.Actions, delegate)
		voteproducerData := common.Variants{
			"voter":     producer,
			"proxy":     common.AccountName(0),
			"producers": []common.AccountName{ /*common.AccountName(1),*/ producer},
		}
		voteproducer := e.GetAction(
			eosio,
			common.N("voteproducer"),
			[]common.PermissionLevel{{Actor: producer, Permission: common.DefaultConfig.ActiveName}},
			&voteproducerData,
		)
		trx.Actions = append(trx.Actions, voteproducer)
		undelegatebwData := common.Variants{
			"from":                 producer,
			"receiver":             producer,
			"unstake_net_quantity": CoreFromString("150000000.0000"),
			"unstake_cpu_quantity": CoreFromString("0.0000"),
		}
		undelegate := e.GetAction(
			eosio,
			common.N("undelegatebw"),
			[]common.PermissionLevel{{Actor: producer, Permission: common.DefaultConfig.ActiveName}},
			&undelegatebwData,
		)
		trx.Actions = append(trx.Actions, undelegate)
		e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
		eosioPriKey := e.getPrivateKey(eosio, "active")
		producerPriKey := e.getPrivateKey(producer, "active")
		chainId := e.Control.GetChainId()
		trx.Sign(&eosioPriKey, &chainId)
		trx.Sign(&producerPriKey, &chainId)
		e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
	}

}

func (e EosioSystemTester) Voter(acnt common.AccountName) common.Variants {
	return common.Variants{
		"owner":     acnt,
		"proxy":     common.AccountName(0),
		"producers": []common.AccountName{},
		"staked":    int64(0),
		//"last_vote_weight":    float64(0),
		"proxied_vote_weight": float64(0),
		"is_proxy":            false,
	}
}

func (e EosioSystemTester) VoterAccountAsset(acnt common.AccountName, voteStake common.Asset) common.Variants {
	voter := e.Voter(acnt)
	voter["staked"] = voteStake.Amount
	return voter
}

func (e EosioSystemTester) Proxy(acnt common.AccountName) common.Variants {
	voter := e.Voter(acnt)
	voter["is_proxy"] = true
	return voter
}

func (e EosioSystemTester) ProxyStake(acnt common.AccountName, voteStake common.Asset) common.Variants {
	voter := e.Proxy(acnt)
	voter["staked"] = voteStake.Amount
	return voter
}
