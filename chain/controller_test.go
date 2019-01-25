package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

var path string = "/tmp/data/"

func TestController_ProduceProcess(t *testing.T) {
	//timer := time.NewTicker(1 * time.Second)
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Startup()
	for i := 0; i < 100; i++ {
		produceProcess(con)
	}
	con.Close()
}

func produceProcess(con *Controller) {

	signatureProviders := make(map[ecc.PublicKey]signatureProviderType)
	//con := NewController(NewConfig())
	con.AbortBlock()
	now := common.Now()
	var base common.TimePoint
	if now > con.HeadBlockTime() {
		base = now
	} else {
		base = con.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}
	con.StartBlock(types.NewBlockTimeStamp(blockTime), 0)
	unappliedTrxs := con.GetUnappliedTransactions()
	if len(unappliedTrxs) > 0 {
		for _, trx := range unappliedTrxs {
			trace := con.PushTransaction(trx, common.MaxTimePoint(), 0)
			if trace.Except != nil {
				log.Error("produce is failed isExhausted=true")
			} else {
				con.DropUnappliedTransaction(trx)
			}
		}
	}
	con.Pending.PendingBlockState = con.PendingBlockState()
	con.FinalizeBlock()
	pubKey, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	if err != nil {
		log.Error("produceLoop NewPublicKey is error :%s", err.Error())
	}
	priKey, err2 := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	if err2 != nil {
		log.Error("produceLoop NewPrivateKey is error :%s", err.Error())
	}
	pbs := con.PendingBlockState()

	signatureProviders[pubKey] = makeKeySignatureProvider(priKey)
	a := signatureProviders[pbs.BlockSigningKey]
	con.SignBlock(func(d crypto.Sha256) ecc.Signature {
		return a(d)
	})

	con.CommitBlock(true)
}

type signatureProviderType = func(sha256 crypto.Sha256) ecc.Signature

func makeKeySignatureProvider(key *ecc.PrivateKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) ecc.Signature {
		sign, err := key.Sign(digest.Bytes())
		if err != nil {
			try.Throw(err)
		}
		return sign
	}
	return signFunc
}

func CallBackApplayHandler(p *ApplyContext) {
	fmt.Println("SetApplyHandler CallBack")
}

func CallBackApplayHandler2(p *ApplyContext) {
	fmt.Println("SetApplyHandler CallBack2")
}
func TestSetApplyHandler(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Startup()
	//applyCon := ApplyContext{}
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("newaccount")), CallBackApplayHandler)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("setcode")), CallBackApplayHandler2)

	handler1 := con.FindApplyHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("newaccount")))
	handler1(nil)

	handler2 := con.FindApplyHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("setcode")))
	handler2(nil)

	fmt.Println(len(con.ApplyHandlers))
	con.Close()
}

var IrreversibleBlock chan types.BlockState = make(chan types.BlockState)

func TestController_CreateNativeAccount(t *testing.T) {
	//CreateNativeAccount(name common.AccountName,owner types.Authority,active types.Authority,isPrivileged bool)
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	name := common.AccountName(common.N("eos"))

	owner := types.Authority{}
	owner.Threshold = 2
	active := types.Authority{}
	active.Threshold = 1
	con.CreateNativeAccount(name, owner, active, false)
	fmt.Println(name)
	result := entity.AccountObject{}
	result.Name = name
	//control.DB.Find("name", result)

	//fmt.Println("check account name:", strings.Compare(name.String(), "eos"))
	assert.Equal(t, "eos", name.String())
	con.Close()
}

func TestController_GetGlobalProperties(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	result := con.GetGlobalProperties()
	gp := entity.GlobalPropertyObject{}
	gp.ID = common.IdType(1)
	err := con.DB.Find("ID", gp, &gp)
	if err != nil {
		assert.Error(t, err, gp)
	}
	assert.Equal(t, false, common.Empty(result)) //GlobalProperties not initialized
	assert.Equal(t, false, result == &gp)
	con.Close()
}

func TestController_GetDynamicGlobalProperties(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Startup()
	con.GetDynamicGlobalProperties()
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 0
	con.Close()
}
func inString(s1, s2 string) bool {
	return strings.Contains(s1, s2)
}
func TestController_GetBlockIdForNum_NotFound(t *testing.T) {
	os.RemoveAll(path)
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Startup()
	var ex string
	try.Try(func() {
		con.GetBlockIdForNum(10)
	}).Catch(func(e exception.UnknownBlockException) {
		ex = e.DetailMessage()
	}).End()
	//fmt.Println("--A--",ex)
	assert.True(t, inString(ex, "Could not find block: 10"))
	con.Close()
}

func TestController_StartBlock(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Startup()
	s := types.Irreversible
	now := common.Now()
	var base common.TimePoint
	if now > con.HeadBlockTime() {
		base = now
	} else {
		base = con.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}
	con.StartBlock(types.NewBlockTimeStamp(blockTime), uint16(s))
	con.Close()
}

func TestController_Close(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	con := NewController(cfg)
	con.Close()
}

func TestController_UpdateProducersAuthority(t *testing.T) {
	cfg := NewConfig()
	cfg.BlocksDir = path + cfg.BlocksDir
	cfg.StateDir = path + cfg.StateDir
	c := NewController(cfg)
	c.Startup()
	c.AbortBlock()
	now := common.Now()
	var base common.TimePoint
	if now > c.HeadBlockTime() {
		base = now
	} else {
		base = c.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}
	c.StartBlock(types.NewBlockTimeStamp(blockTime), 0)
	c.updateProducersAuthority()
	c.Close()
}
