package unittests

import (
	"fmt"
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"strings"
	"testing"
)

func TestBuySell(t *testing.T) {
	e := initEosioSystemTester()
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(eosio, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(alice)
	initBytes := total["ram_bytes"].(uint64)
	initialRamBalance := e.GetBalance(eosioRam)
	initialRamFeeBalance := e.GetBalance(eosioRamFee)
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("200.0000")))
	assert.Equal(t, CoreFromString("800.0000"), e.GetBalance(alice))
	assert.Equal(t, initialRamBalance.Add(CoreFromString("199.0000")), e.GetBalance(eosioRam))
	assert.Equal(t, initialRamFeeBalance.Add(CoreFromString("1.0000")), e.GetBalance(eosioRamFee))

	total = e.GetTotalStake(alice)
	bytes := total["ram_bytes"].(uint64)

	boughtBytes := bytes - initBytes

	assert.Equal(t, true, 0 < boughtBytes)
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("998.0049"), e.GetBalance(alice))
	total = e.GetTotalStake(alice)
	assert.Equal(t, initBytes, total["ram_bytes"].(uint64))

	e.Transfer(eosio, alice, CoreFromString("100000000.0000"), eosio)
	assert.Equal(t, CoreFromString("100000998.0049"), e.GetBalance(alice))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, CoreFromString("90000998.0049"), e.GetBalance(alice))

	total = e.GetTotalStake(alice)
	bytes = total["ram_bytes"].(uint64)
	boughtBytes = bytes - initBytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))

	total = e.GetTotalStake(alice)
	bytes = total["ram_bytes"].(uint64)
	boughtBytes = bytes - initBytes
	assert.Equal(t, initBytes, total["ram_bytes"].(uint64))
	assert.Equal(t, CoreFromString("99901248.0041"), e.GetBalance(alice))

	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("30.0000")))
	assert.Equal(t, CoreFromString("99900688.0041"), e.GetBalance(alice))

	newTotal := e.GetTotalStake(alice)
	newBytes := newTotal["ram_bytes"].(uint64)
	boughtBytes = newBytes - bytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("99901242.4179"), e.GetBalance(alice))

	newTotal = e.GetTotalStake(alice)
	startBytes := total["ram_bytes"].(uint64)

	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("100000.0000")))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("300000.0000")))
	assert.Equal(t, CoreFromString("49301242.4179"), e.GetBalance(alice))

	finalTotal := e.GetTotalStake(alice)
	endBytes := finalTotal["ram_bytes"].(uint64)
	boughtBytes = endBytes - startBytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("99396507.4142"), e.GetBalance(alice))
	e.close()
}

func TestStakeUnstake(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.ProduceBlocks(10, false)
	e.ProduceBlock(common.Hours(3*24), 0)

	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)

	assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))
	assert.Equal(t, e.Success(), e.Stake(eosio, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	initEosioStakeBalance := e.GetBalance(eosioStake)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	assert.Equal(t, initEosioStakeBalance.Add(CoreFromString("300.0000")), e.GetBalance(eosioStake))
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	e.ProduceBlock(common.Hours(3*24-1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	assert.Equal(t, initEosioStakeBalance.Add(CoreFromString("300.0000")), e.GetBalance(eosioStake))

	//after 3 days funds should be released
	e.ProduceBlock(common.Hours(1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))
	assert.Equal(t, initEosioStakeBalance, e.GetBalance(eosioStake))

	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))
	bytes := total["ram_bytes"].(uint64)
	assert.True(t, 0 < bytes)

	//unstake from bob111111111
	assert.Equal(t, e.Success(), e.UnStake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("10.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("10.0000"), total["cpu_weight"].(common.Asset))

	e.ProduceBlock(common.Hours(3*24-1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//after 3 days funds should be released
	e.ProduceBlock(common.Hours(1),0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("0.0000")), e.GetVoterInfo(alice))
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))

	e.close()
}

func TestStakeUnstakeWithTransfer(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(eosio, CoreFromString("1000.0000"), eosio)
	e.Issue(eosioStake, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))

	//eosio stakes for alice with transfer flag
	e.Transfer(eosio, bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.StakeWithTransfer(bob, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	//check that alice has both bandwidth and voting power
	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))

	//alice stakes for herself
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	//now alice's stake should be equal to transfered from eosio + own stake
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	assert.Equal(t, CoreFromString("410.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("210.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("600.0000")), e.GetVoterInfo(alice))

	//alice can unstake everything (including what was transfered)
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("400.0000"), CoreFromString("200.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	e.ProduceBlock(common.Hours(3*24-1),0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//after 3 days funds should be released
	e.ProduceBlock(common.Hours(1),0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("1300.0000"), e.GetBalance(alice))

	//stake should be equal to what was staked in constructor, voting power should be 0
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("10.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("10.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("0.0000")), e.GetVoterInfo(alice))

	// Now alice stakes to bob with transfer flag
	assert.Equal(t, e.Success(), e.StakeWithTransfer(alice, bob, CoreFromString("100.0000"), CoreFromString("100.0000")))
	e.close()
}

func TestStakeToSelfWithTransfer(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	stakeWithTransfer := func() { e.StakeWithTransfer(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")) }
	CatchThrowMsg(t, "cannot use transfer flag if delegating to self", stakeWithTransfer)
	e.close()
}

func TestStakeWhilePendingRefund(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(eosio, CoreFromString("1000.0000"), eosio)
	e.Issue(eosioStake, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))

	//eosio stakes for alice with transfer flag
	e.Transfer(eosio, bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.StakeWithTransfer(bob, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	//check that alice has both bandwidth and voting power
	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))

	//alice stakes for herself
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"),CoreFromString("100.0000")))

	//now alice's stake should be equal to transfered from eosio + own stake
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	assert.Equal(t, CoreFromString("410.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("210.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("600.0000")), e.GetVoterInfo(alice))

	//alice can unstake everything (including what was transfered)
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("400.0000"),CoreFromString("200.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	e.ProduceBlock(common.Hours(3*24-1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//after 3 days funds should be released
	e.ProduceBlock(common.Hours(1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("1300.0000"), e.GetBalance(alice))

	//stake should be equal to what was staked in constructor, voting power should be 0
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("10.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("10.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("0.0000")), e.GetVoterInfo(alice))

	// Now alice stakes to bob with transfer flag
	assert.Equal(t, e.Success(), e.StakeWithTransfer(alice, bob, CoreFromString("100.0000"), CoreFromString("100.0000")))
	e.close()
}

func TestFailWithoutAuth(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(eosio, alice, CoreFromString("2000.0000"), CoreFromString("1000.0000")))
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("10.0000"), CoreFromString("10.0000")))
	var ex string
	Try(func(){
		act := common.N("delegatebw")
		data := common.Variants{
			"from": 			  alice,
			"receiver":           bob,
			"stake_net_quantity": CoreFromString("10.0000"),
			"stake_cpu_quantity": CoreFromString("10.0000"),
			"transfer": 		  0,
		}
		e.EsPushAction(&alice, &act, &data, false)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "missing authority of alice1111111"))

	Try(func(){
		act := common.N("undelegatebw")
		data := common.Variants{
			"from": 			    alice,
			"receiver":             bob,
			"unstake_net_quantity": CoreFromString("200.0000"),
			"unstake_cpu_quantity": CoreFromString("100.0000"),
			"transfer": 		    0,
		}
		e.EsPushAction(&alice, &act, &data, false)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "missing authority of alice1111111"))
	e.close()
}

func TestStakeNegative(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)

	var ex string
	Try(func(){
		e.Stake(alice, alice, CoreFromString("-0.0001"), CoreFromString("0.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must stake a positive amount"))

	Try(func(){
		e.Stake(alice, alice, CoreFromString("0.0000"), CoreFromString("-0.0001"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must stake a positive amount"))

	Try(func(){
		e.Stake(alice, alice, CoreFromString("00.0000"), CoreFromString("00.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must stake a positive amount"))

	Try(func(){
		e.Stake(alice, alice, CoreFromString("0.0000"), CoreFromString("00.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must stake a positive amount"))

	e.close()
}

func TestUnstakeNegative(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0001"), CoreFromString("100.0001")))

	total := e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("210.0001"), total["net_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0002")), e.GetVoterInfo(alice))

	var ex string
	Try(func(){
		e.UnStake(alice, bob, CoreFromString("-1.0000"), CoreFromString("0.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must unstake a positive amount"))

	Try(func(){
		e.UnStake(alice, bob, CoreFromString("0.0000"), CoreFromString("-1.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must unstake a positive amount"))

	Try(func(){
		e.UnStake(alice, alice, CoreFromString("0.0000"), CoreFromString("0.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "must unstake a positive amount"))

	e.close()
}

func TestUnstakeMoreThanAtStake(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//trying to unstake more net bandwith than at stake
	var ex string
	Try(func(){
		e.UnStake(alice, alice, CoreFromString("200.0001"), CoreFromString("0.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked net bandwidth"))

	//trying to unstake more cpu bandwith than at stake
	Try(func(){
		e.UnStake(alice, alice, CoreFromString("0.0000"), CoreFromString("100.0001"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked cpu bandwidth"))

	//check that nothing has changed
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	e.close()
}

func TestDelegateToAnotherUser(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//all voting power goes to alice1111111
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))

	//but not to bob111111111
	assert.Equal(t, common.Variants{}, e.GetVoterInfo(bob))

	//bob111111111 should not be able to unstake what was staked by alice1111111
	var ex string
	Try(func(){
		e.UnStake(bob, bob, CoreFromString("10.0000"), CoreFromString("0.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked net bandwidth"))

	Try(func(){
		e.UnStake(bob, bob, CoreFromString("0.0000"), CoreFromString("10.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked cpu bandwidth"))

	e.Issue(carol, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(carol, bob, CoreFromString("20.0000"), CoreFromString("10.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("230.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("120.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("970.0000"), e.GetBalance(carol))
	assert.Equal(t, e.VoterAccountAsset(carol, CoreFromString("30.0000")), e.GetVoterInfo(carol))

	//alice1111111 should not be able to unstake money staked by carol1111111
	Try(func(){
		e.UnStake(alice, bob, CoreFromString("201.0000"), CoreFromString("1.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked net bandwidth"))

	Try(func(){
		e.UnStake(alice, bob, CoreFromString("1.0000"), CoreFromString("101.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "insufficient staked cpu bandwidth"))

	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("230.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("120.0000"), total["cpu_weight"].(common.Asset))

	//balance should not change after unsuccessfull attempts to unstake
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//voting power too
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, e.VoterAccountAsset(carol, CoreFromString("30.0000")), e.GetVoterInfo(carol))
	assert.Equal(t, common.Variants{}, e.GetVoterInfo(bob))
	e.close()
}

func TestStakeUnstakeSeparate(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))

	//everything at once
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("10.0000"), CoreFromString("20.0000")))
	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("20.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("30.0000"), total["cpu_weight"].(common.Asset))

	//net
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("100.0000"), CoreFromString("0.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("120.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("30.0000"), total["cpu_weight"].(common.Asset))

	//cpu
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("0.0000"), CoreFromString("200.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("120.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("230.0000"), total["cpu_weight"].(common.Asset))

	//unstake net
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("100.0000"), CoreFromString("0.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("20.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("230.0000"), total["cpu_weight"].(common.Asset))

	//unstake cpu
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("0.0000"), CoreFromString("200.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("20.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("30.0000"), total["cpu_weight"].(common.Asset))

	e.close()
}

func TestAddingStakePartialUnstake(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("300.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("100.0000"), CoreFromString("50.0000")))

	total := e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("310.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("160.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("450.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("550.0000"), e.GetBalance(alice))

	//unstake a share
	assert.Equal(t, e.Success(), e.UnStake(alice, bob, CoreFromString("150.0000"), CoreFromString("75.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("160.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("85.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("225.0000")), e.GetVoterInfo(alice))

	//unstake more
	assert.Equal(t, e.Success(), e.UnStake(alice, bob, CoreFromString("50.0000"), CoreFromString("25.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("110.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("60.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("150.0000")), e.GetVoterInfo(alice))

	//combined amount should be available only in 3 days
	e.ProduceBlock(common.Days(2), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("550.0000"), e.GetBalance(alice))
	e.ProduceBlock(common.Days(1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("850.0000"), e.GetBalance(alice))

	e.close()
}

func TestStakeFromRefund(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))
	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("50.0000"), CoreFromString("50.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("60.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("60.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("400.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("600.0000"), e.GetBalance(alice))

	//unstake a share
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("100.0000"), CoreFromString("50.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("110.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("60.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("250.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("600.0000"), e.GetBalance(alice))
	refund := e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("100.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("50.0000"), refund["cpu_amount"].(common.Asset))

	//alice delegates to bob, should pull from liquid balance not refund
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("50.0000"), CoreFromString("50.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("110.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("60.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("350.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, CoreFromString("500.0000"), e.GetBalance(alice))
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("100.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("50.0000"), refund["cpu_amount"].(common.Asset))

	//stake less than pending refund, entire amount should be taken from refund
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("50.0000"), CoreFromString("25.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("160.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("85.0000"), total["cpu_weight"].(common.Asset))
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("50.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("25.0000"), refund["cpu_amount"].(common.Asset))

	//balance should stay the same
	assert.Equal(t, CoreFromString("500.0000"), e.GetBalance(alice))

	//stake exactly pending refund amount
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("50.0000"), CoreFromString("25.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	//pending refund should be removed
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, common.Variants{}, refund)

	//balance should stay the same
	assert.Equal(t, CoreFromString("500.0000"), e.GetBalance(alice))

	//create pending refund again
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("10.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("10.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("500.0000"), e.GetBalance(alice))
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("200.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("100.0000"), refund["cpu_amount"].(common.Asset))

	//stake more than pending refund
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("300.0000"), CoreFromString("200.0000")))
	total = e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("310.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("210.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, e.VoterAccountAsset(alice, CoreFromString("700.0000")), e.GetVoterInfo(alice))
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, common.Variants{}, refund)

	//200 core tokens should be taken from alice's account
	assert.Equal(t, CoreFromString("300.0000"), e.GetBalance(alice))

	e.close()
}

func TestStakeToAnotherUserNotFromRefund(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(alice)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	//unstake
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))
	refund := e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("200.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("100.0000"), refund["cpu_amount"].(common.Asset))

	//stake to another user
	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	total = e.GetTotalStake(bob)
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("400.0000"), e.GetBalance(alice))
	refund = e.GetRefundRequest(alice)
	assert.Equal(t, CoreFromString("200.0000"), refund["net_amount"].(common.Asset))
	assert.Equal(t, CoreFromString("100.0000"), refund["cpu_amount"].(common.Asset))

	e.close()
}

func TestProducerRegisterUnregister(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)

	key, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key,
			"url":          "http://block.one",
			"location":     1,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	info := e.GetProducerInfo(alice)
	assert.Equal(t, alice, info["owner"].(common.AccountName))
	assert.Equal(t, float64(0), info["total_votes"].(float64))
	assert.Equal(t, "http://block.one", info["url"].(string))

	//change parameters one by one to check for things like #3783
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key,
			"url":          "http://block.two",
			"location":     1,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}
	info = e.GetProducerInfo(alice)
	assert.Equal(t, alice, info["owner"].(common.AccountName))
	assert.Equal(t, key, info["producer_key"].(ecc.PublicKey))
	assert.Equal(t, string("http://block.two"), info["url"].(string))
	assert.Equal(t, uint16(1), info["location"].(uint16))

	key2, _ := ecc.NewPublicKey("EOS5jnmSKrzdBHE9n8hw58y7yxFWBC8SNiG7m8S1crJH3KvAnf9o6")
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key2,
			"url":          "http://block.two",
			"location":     2,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}
	info = e.GetProducerInfo(alice)
	assert.Equal(t, alice, info["owner"].(common.AccountName))
	assert.Equal(t, key2, info["producer_key"].(ecc.PublicKey))
	assert.Equal(t, string("http://block.two"), info["url"].(string))
	assert.Equal(t, uint16(2), info["location"].(uint16))

	//unregister producer
	{
		act := common.N("unregprod")
		data := common.Variants{
			"producer": alice,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//key should be empty
	info = e.GetProducerInfo(alice)
	assert.Equal(t, *ecc.NewPublicKeyNil(), info["producer_key"].(ecc.PublicKey))

	//everything else should stay the same
	assert.Equal(t, alice, info["owner"].(common.AccountName))
	assert.Equal(t, string("http://block.two"), info["url"].(string))
	assert.Equal(t, uint16(2), info["location"].(uint16))

	//unregister bob111111111 who is not a producer
	{
		var ex string
		Try(func() {
			act := common.N("unregprod")
			data := common.Variants{
				"producer": bob,
			}
			e.EsPushAction(&bob, &act, &data, true)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer not found"))
	}

	e.close()
}

func TestVoteForProducer(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": e.getPublicKey(alice, "active"),
			"url":          "http://block.one",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}
	prod := e.GetProducerInfo(alice)
	assert.Equal(t, alice, prod["owner"].(common.AccountName))
	assert.Equal(t, float64(0), prod["total_votes"].(float64))
	assert.Equal(t, string("http://block.one"), prod["url"].(string))

	e.Issue(bob, CoreFromString("2000.0000"), eosio)
	e.Issue(carol, CoreFromString("3000.0000"), eosio)

	//bob111111111 makes stake
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("11.0000"), CoreFromString("0.1111")))
	assert.Equal(t, CoreFromString("1988.8889"), e.GetBalance(bob))
	assert.Equal(t, e.VoterAccountAsset(bob, CoreFromString("11.1111")), e.GetVoterInfo(bob))

	//bob111111111 votes for alice1111111
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{alice}, common.AccountName(0)))

	//check that producer parameters stay the same after voting
	prod = e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("11.1111")), prod["total_votes"].(float64))
	assert.Equal(t, alice, prod["owner"].(common.AccountName))
	assert.Equal(t, string("http://block.one"), prod["url"].(string))

	//carol1111111 makes stake
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("22.0000"), CoreFromString("0.2222")))
	assert.Equal(t, e.VoterAccountAsset(carol, CoreFromString("22.2222")), e.GetVoterInfo(carol))
	assert.Equal(t, CoreFromString("2977.7778"), e.GetBalance(carol))

	//carol1111111 votes for alice1111111
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{alice}, common.AccountName(0)))

	//new stake votes be added to alice1111111's total_votes
	prod = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("33.3333"))) <= EPSINON)

	//bob111111111 increases his stake
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("33.0000"), CoreFromString("0.3333")))

	//alice1111111 stake with transfer to bob111111111
	assert.Equal(t, e.Success(), e.StakeWithTransfer(alice, bob, CoreFromString("22.0000"), CoreFromString("0.2222")))

	//should increase alice1111111's total_votes
	prod = e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("88.8888")), prod["total_votes"].(float64))

	//carol1111111 unstakes part of the stake
	assert.Equal(t, e.Success(), e.UnStake(carol, carol, CoreFromString("2.0000"), CoreFromString("0.0002")))

	//should decrease alice1111111's total_votes
	prod = e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("86.8886")), prod["total_votes"].(float64))

	//bob111111111 revokes his vote
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{}, common.AccountName(0)))

	//should decrease alice1111111's total_votes
	prod = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.2220"))) <= EPSINON)

	//but eos should still be at stake
	assert.Equal(t, CoreFromString("1955.5556"), e.GetBalance(bob))

	//carol1111111 unstakes rest of eos
	assert.Equal(t, e.Success(), e.UnStake(carol, carol, CoreFromString("20.0000"), CoreFromString("0.2220")))

	//should decrease alice1111111's total_votes to zero
	prod = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(prod["total_votes"].(float64) - float64(0)) <= EPSINON)

	//carol1111111 should receive funds in 3 days
	e.ProduceBlock(common.Days(3), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("3000.0000"), e.GetBalance(carol))

	e.close()
}

func TestUnregisteredProducerVoting(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(bob, CoreFromString("2000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("13.0000"), CoreFromString("0.5791")))

	//bob111111111 should not be able to vote for alice1111111 who is not a producer
	{
		var ex string
		Try(func(){
			e.Vote(bob, []common.AccountName{alice}, common.AccountName(0))
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer is not registered"))
	}

	//alice1111111 registers as a producer
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": e.getPublicKey(alice, "active"),
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//and then unregisters
	{
		act := common.N("unregprod")
		data := common.Variants{
			"producer": alice,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//key should be empty
	prod := e.GetProducerInfo(alice)
	assert.Equal(t, *ecc.NewPublicKeyNil(), prod["producer_key"].(ecc.PublicKey))

	//bob111111111 should not be able to vote for alice1111111 who is an unregistered producer
	{
		var ex string
		Try(func(){
			e.Vote(bob, []common.AccountName{alice}, common.AccountName(0))
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer is not currently registered"))
	}

	e.close()
}

func TestMoreThan30ProducerVoting(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(bob, CoreFromString("2000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("13.0000"), CoreFromString("0.5791")))
	assert.Equal(t, e.VoterAccountAsset(bob, CoreFromString("13.5791")), e.GetVoterInfo(bob))

	//bob111111111 should not be able to vote for alice1111111 who is not a producer
	var producers []common.AccountName
	for i := 0; i < 31; i ++ {
		producers = append(producers, alice)
	}
	var ex string
	Try(func() {
		e.Vote(bob, producers, common.AccountName(0))
	}).Catch(func(e Exception){
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "attempt to vote for too many producers"))
	e.close()
}

func TestVoteSameProducer30Times(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(bob, CoreFromString("2000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("50.0000"), CoreFromString("50.0000")))
	assert.Equal(t, e.VoterAccountAsset(bob, CoreFromString("100.0000")), e.GetVoterInfo(bob))

	//alice1111111 becomes a producer
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": e.getPublicKey(alice, "active"),
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//bob111111111 should not be able to vote for alice1111111 for 30 times
	var producers []common.AccountName
	for i := 0; i < 30; i ++ {
		producers = append(producers, alice)
	}
	var ex string
	Try(func(){
		e.Vote(bob, producers, common.AccountName(0))
	}).Catch(func(e Exception){
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "producer votes must be unique and sorted"))

	prod := e.GetProducerInfo(alice)
	assert.True(t, math.Abs(prod["total_votes"].(float64) - float64(0)) <= EPSINON)

	e.close()
}

func TestProducerKeepVotes(t *testing.T) {
	e := initEosioSystemTester()
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	key := e.getPublicKey(alice, "active")
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key,
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//bob111111111 makes stake
	e.Issue(bob, CoreFromString("2000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("13.0000"), CoreFromString("0.5791")))
	assert.Equal(t, e.VoterAccountAsset(bob, CoreFromString("13.5791")), e.GetVoterInfo(bob))

	//bob111111111 votes for alice1111111
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{alice}, common.AccountName(0)))

	prod := e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("13.5791")), prod["total_votes"].(float64))

	//unregister producer
	{
		act := common.N("unregprod")
		data := common.Variants{
			"producer": alice,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	prod = e.GetProducerInfo(alice)
	//key should be empty && votes should stay the same
	assert.Equal(t, *ecc.NewPublicKeyNil(), prod["producer_key"].(ecc.PublicKey))
	assert.Equal(t, e.Stake2Votes(CoreFromString("13.5791")), prod["total_votes"].(float64))

	//register the same producer again
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key,
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	prod = e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("13.5791")), prod["total_votes"].(float64))

	//the same producer again
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": key,
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	prod = e.GetProducerInfo(alice)
	assert.Equal(t, e.Stake2Votes(CoreFromString("13.5791")), prod["total_votes"].(float64))

	e.close()
}

func TestVoteForTwoProducers(t *testing.T) {
	e := initEosioSystemTester()

	//alice1111111 becomes a producer
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     alice,
			"producer_key": e.getPublicKey(alice, "active"),
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&alice, &act, &data, true))
	}

	//bob111111111 becomes a producer
	{
		act := common.N("regproducer")
		data := common.Variants{
			"producer":     bob,
			"producer_key": e.getPublicKey(bob, "active"),
			"url":          "",
			"location":     0,
		}
		assert.Equal(t, e.Success(), e.EsPushAction(&bob, &act, &data, true))
	}

	//carol1111111 votes for alice1111111 and bob111111111
	e.Issue(carol, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("15.0005"), CoreFromString("5.0000")))
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{alice, bob}, common.AccountName(0)))

	aliceInfo := e.GetProducerInfo(alice)
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= EPSINON)
	bobInfo := e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= EPSINON)

	//carol1111111 votes for alice1111111 (but revokes vote for bob111111111)
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{alice}, common.AccountName(0)))
	aliceInfo = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= EPSINON)
	bobInfo = e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//alice1111111 votes for herself and bob111111111
	e.Issue(alice, CoreFromString("2.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("1.0000"), CoreFromString("1.0000")))
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{alice, bob}, common.AccountName(0)))
	aliceInfo = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("22.0005"))) <= EPSINON)
	bobInfo = e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("2.0000"))) <= EPSINON)

	e.close()
}

func TestProxyRegisterUnregisterKeepsStake(t *testing.T) {
	e := initEosioSystemTester()

	//register proxy by first action for this user ever
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.Equal(t, e.Proxy(alice), e.GetVoterInfo(alice))

	//unregister proxy
	assert.Equal(t, e.Success(), e.UnRegProxy(alice))
	assert.Equal(t, e.Voter(alice), e.GetVoterInfo(alice))

	//stake and then register as a proxy
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("200.0002"), CoreFromString("100.0001")))
	assert.Equal(t, e.Success(), e.RegProxy(bob))
	assert.Equal(t, e.ProxyStake(bob, CoreFromString("300.0003")), e.GetVoterInfo(bob))

	//unregister and check that stake is still in place
	assert.Equal(t, e.Success(), e.UnRegProxy(bob))
	assert.Equal(t, e.VoterAccountAsset(bob, CoreFromString("300.0003")), e.GetVoterInfo(bob))

	//register as a proxy and then stake
	assert.Equal(t, e.Success(), e.RegProxy(carol))
	e.Issue(carol, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("246.0002"), CoreFromString("531.0001")))

	//check that both proxy flag and stake a correct
	assert.Equal(t, e.ProxyStake(carol, CoreFromString("777.0003")), e.GetVoterInfo(carol))

	//unregister
	assert.Equal(t, e.Success(), e.UnRegProxy(carol))
	assert.Equal(t, e.VoterAccountAsset(carol, CoreFromString("777.0003")), e.GetVoterInfo(carol))

	e.close()
}

func TestProxyStakeUnstakeKeepsProxyFlag(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Proxy(alice), e.GetVoterInfo(alice))

	//stake
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("100.0000"), CoreFromString("50.0000")))
	assert.Equal(t, e.ProxyStake(alice, CoreFromString("150.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("30.0000"), CoreFromString("20.0000")))
	assert.Equal(t, e.ProxyStake(alice, CoreFromString("200.0000")), e.GetVoterInfo(alice))

	//unstake
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("65.0000"), CoreFromString("35.0000")))
	assert.Equal(t, e.ProxyStake(alice, CoreFromString("100.0000")), e.GetVoterInfo(alice))
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("65.0000"), CoreFromString("35.0000")))
	assert.Equal(t, e.ProxyStake(alice, CoreFromString("0.0000")), e.GetVoterInfo(alice))

	e.close()
}

func TestProxyActionsAffectProducers(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.CreateAccountsWithResources([]common.AccountName{producer1, producer2, producer3}, eosio)
	assert.Equal(t, e.Success(), e.RegProducer(producer1))
	assert.Equal(t, e.Success(), e.RegProducer(producer2))
	assert.Equal(t, e.Success(), e.RegProducer(producer3))

	//register as a proxy
	assert.Equal(t, e.Success(), e.RegProxy(alice))

	//accumulate proxied votes
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{}, alice))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//vote for producers
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1, producer2}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//vote for another producers
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1, producer3}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//unregister proxy
	assert.Equal(t, e.Success(), e.UnRegProxy(alice))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//register proxy again
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//stake increase by proxy itself affects producers
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("30.0001"), CoreFromString("20.0001")))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= EPSINON)

	//stake decrease by proxy itself affects producers
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("10.0001"), CoreFromString("10.0001")))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("180.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("180.0003"))) <= EPSINON)

	e.close()
}

func TestProducerPay(t *testing.T) {
	e := initEosioSystemTester()
	continueRate := float64(4.879 / 100)
	usecsPerYear := float64(52 * 7 * 24 * 3600 * 1000000)
	secsPerYear := float64(52 * 7 * 24 * 3600)
	largeAsset := CoreFromString("80.0000")
	e.CreateAccountWithResources(producer1, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(producer2, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(producer3, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(voter1, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(voter2, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)

	assert.Equal(t, e.Success(), e.RegProducer(producer1))
	prod := e.GetProducerInfo(producer1)
	assert.Equal(t, producer1, prod["owner"].(common.AccountName))
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	e.Transfer(eosio, voter1, CoreFromString("400000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(voter1, voter1, CoreFromString("100000000.0000"), CoreFromString("100000000.0000")))
	assert.Equal(t, e.Success(), e.Vote(voter1, []common.AccountName{producer1}, common.AccountName(0)))

	// defproducera is the only active producer
	// produce enough blocks so new schedule kicks in and defproducera produces some blocks
	{
		e.ProduceBlocks(50, false)
		initialGlobalState := e.GetGlobalState()
		initialClaimTime := initialGlobalState["last_pervote_bucket_fill"].(uint64)
		initialPervoteBucket := initialGlobalState["pervote_bucket"].(int64)
		initialPerblockBucket := initialGlobalState["perblock_bucket"].(int64)
		initialSavings := e.GetBalance(eosioSaving).Amount
		initialTotUnpaidBlocks := initialGlobalState["total_unpaid_blocks"].(uint32)
		prod = e.GetProducerInfo(producer1)
		unpaidBlocks := prod["unpaid_blocks"].(uint32)
		assert.True(t, 1 < unpaidBlocks)
		assert.Equal(t, uint64(0), prod["last_claim_time"].(uint64))
		assert.Equal(t, initialTotUnpaidBlocks, unpaidBlocks)
		initialSupply := e.GetTokenSupply()
		initialBalance := e.GetBalance(producer1)
		e.ClaimRewards(producer1)
		globalState := e.GetGlobalState()
		claimTime := globalState["last_pervote_bucket_fill"].(uint64)
		pervoteBucket := globalState["pervote_bucket"].(int64)
		//perblockBucket := globalState["perblock_bucket"].(int64)
		savings := e.GetBalance(eosioSaving).Amount
		totUnpaidBlocks := globalState["total_unpaid_blocks"].(uint32)

		prod = e.GetProducerInfo(producer1)
		assert.Equal(t, uint32(1), prod["unpaid_blocks"].(uint32))
		assert.Equal(t, uint32(1), totUnpaidBlocks)
		supply := e.GetTokenSupply()
		balance := e.GetBalance(producer1)

		assert.Equal(t, claimTime, prod["last_claim_time"].(uint64))
		usecsBetweenFills := claimTime - initialClaimTime
		secsBetweenFills := usecsBetweenFills / 1000000

		assert.Equal(t, int64(0), initialSavings)
		assert.Equal(t, int64(0), initialPerblockBucket)
		assert.Equal(t, int64(0), initialPervoteBucket)
		assert.Equal(t, int64(float64(initialSupply.Amount * int64(secsBetweenFills)) * continueRate / secsPerYear), supply.Amount - initialSupply.Amount)
		assert.Equal(t, int64(float64(initialSupply.Amount * int64(secsBetweenFills)) * (4 * continueRate / 5 ) / secsPerYear), savings - initialSavings)
		assert.Equal(t, int64(float64(initialSupply.Amount * int64(secsBetweenFills)) * (0.25 * continueRate / 5) / secsPerYear), balance.Amount - initialBalance.Amount)

		fromPerblockBucket := int64(float64(initialSupply.Amount * int64(secsBetweenFills)) * (0.25 * continueRate / 5) / secsPerYear)
		fromPervoteBucket := int64(float64(initialSupply.Amount * int64(secsBetweenFills)) * (0.75 * continueRate / 5) / secsPerYear)

		if fromPervoteBucket >= 100 * 1000 {
			assert.Equal(t, fromPerblockBucket + fromPervoteBucket, balance.Amount - initialBalance.Amount)
			assert.Equal(t, int64(0), pervoteBucket)
		} else {
			assert.Equal(t, fromPerblockBucket, balance.Amount - initialBalance.Amount)
			assert.Equal(t, fromPervoteBucket, pervoteBucket)
		}
	}

	{
		var ex string
		Try(func() {
			e.ClaimRewards(producer1)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	// defproducera waits for 23 hours and 55 minutes, can't claim rewards yet
	{
		e.ProduceBlock(common.Seconds(23 * 3600 + 55 * 60), 0)
		var ex string
		Try(func() {
			e.ClaimRewards(producer1)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	// wait 5 more minutes, defproducera can now claim rewards again
	{
		e.ProduceBlock(common.Seconds(5 * 60), 0)
		initialGlobalState := e.GetGlobalState()
		initialClaimTime := initialGlobalState["last_pervote_bucket_fill"].(uint64)
		initialPervoteBucket := initialGlobalState["pervote_bucket"].(int64)
		//initialPerblockBucket := initialGlobalState["perblock_bucket"].(int64)
		initialSavings := e.GetBalance(eosioSaving).Amount
		initialTotUnpaidBlocks := initialGlobalState["total_unpaid_blocks"].(uint32)
		initialTotVoteWeight := initialGlobalState["total_producer_vote_weight"].(float64)
		prod = e.GetProducerInfo(producer1)
		unpaidBlocks := prod["unpaid_blocks"].(uint32)
		assert.True(t, 1 < unpaidBlocks)
		assert.Equal(t, initialTotUnpaidBlocks, unpaidBlocks)
		assert.True(t, float64(0) < prod["total_votes"].(float64))
		assert.Equal(t, initialTotVoteWeight, prod["total_votes"].(float64))
		assert.True(t, uint64(0) < prod["last_claim_time"].(uint64))

		initialSupply := e.GetTokenSupply()
		initialBalance := e.GetBalance(producer1)
		assert.Equal(t, e.Success(), e.ClaimRewards(producer1))
		globalState := e.GetGlobalState()
		claimTime := globalState["last_pervote_bucket_fill"].(uint64)
		pervoteBucket := globalState["pervote_bucket"].(int64)
		//perblockBucket := globalState["perblock_bucket"].(int64)
		savings := e.GetBalance(eosioSaving).Amount
		totUnpaidBlocks := globalState["total_unpaid_blocks"].(uint32)

		prod = e.GetProducerInfo(producer1)
		assert.Equal(t, uint32(1), prod["unpaid_blocks"].(uint32))
		assert.Equal(t, uint32(1), totUnpaidBlocks)
		supply := e.GetTokenSupply()
		balance := e.GetBalance(producer1)

		assert.Equal(t, claimTime, prod["last_claim_time"].(uint64))
		usecsBetweenFills := claimTime - initialClaimTime

		assert.Equal(t, int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * continueRate / usecsPerYear), supply.Amount - initialSupply.Amount)
		assert.Equal(t, (supply.Amount - initialSupply.Amount) - (supply.Amount - initialSupply.Amount) / 5 , savings - initialSavings)

		toProducer := int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * continueRate / usecsPerYear) / 5
		toPerblockBucket := toProducer / 4
		toPervoteBucket := toProducer - toPerblockBucket

		if toPervoteBucket + initialPervoteBucket >= 100 * 10000 {
			assert.Equal(t, toPerblockBucket + toPervoteBucket +initialPervoteBucket, balance.Amount - initialBalance.Amount)
			assert.Equal(t, int64(0), pervoteBucket)
		} else {
			assert.Equal(t, toPerblockBucket, balance.Amount - initialBalance.Amount)
			assert.Equal(t, toPervoteBucket + initialPervoteBucket, pervoteBucket)
		}
	}
	// defproducerb tries to claim rewards but he's not on the list
	{
		e.RegProducer(producer2)
		e.RegProducer(producer3)
		initialSupply := e.GetTokenSupply()
		initialSavings := e.GetBalance(eosioSaving).Amount
		for i := 0; i < 7 * 52; i++ {
			e.ProduceBlock(common.Seconds(8 * 3600), 0)
			assert.Equal(t, e.Success(), e.ClaimRewards(producer3))
			e.ProduceBlock(common.Seconds(8 * 3600), 0)
			assert.Equal(t, e.Success(), e.ClaimRewards(producer2))
			e.ProduceBlock(common.Seconds(8 * 3600), 0)
			assert.Equal(t, e.Success(), e.ClaimRewards(producer1))
		}
		supply := e.GetTokenSupply()
		savings := e.GetBalance(eosioSaving).Amount

		// Amount issued per year is very close to the 5% inflation target. Small difference (500 tokens out of 50'000'000 issued)
		// is due to compounding every 8 hours in this test as opposed to theoretical continuous compounding
		assert.True(t, 500 * 10000 > int64(float64(initialSupply.Amount) * float64(0.05) - float64(supply.Amount - initialSupply.Amount)))
		assert.True(t, 500 * 10000 > int64(float64(initialSupply.Amount) * float64(0.04) - float64(savings - initialSavings)))
	}

	e.close()
}

func TestMultipleProducerPay(t *testing.T) {
	e := initEosioSystemTester()
	withinOne := func(a int64, b int64) bool {
		return math.Abs(float64(a-b)) <= 1
	}

	secsPerYear := 52 * 7 * 24 * 3600
	usecsPerYear := secsPerYear * 1000000
	contRate := 4.879 / 100
	net := CoreFromString("80.0000")
	cpu := CoreFromString("80.0000")
	voters := []common.AccountName{voter1, voter2, voter3, voter4}
	for _, v := range voters {
		e.CreateAccountWithResources(v, eosio, CoreFromString("1.0000"), false, net, cpu)
		e.Transfer(eosio, v, CoreFromString("100000000.0000"), eosio)
		assert.Equal(t, e.Success(), e.Stake(v, v, CoreFromString("30000000.0000"), CoreFromString("30000000.0000")))
	}

	// create accounts {defproducera, defproducerb, ..., defproducerz, abcproducera, ..., abcproducern} and register as producers
	var producersNames []common.AccountName
	{
		root := "defproducer"
		for c := 'a'; c <= 'z'; c++ {
			acc := common.N(root + string(c))
			producersNames = append(producersNames, acc)
		}
		root = "abcproducer"
		for c := 'a'; c <= 'n'; c++ {
			acc := common.N(root + string(c))
			producersNames = append(producersNames, acc)
		}
		e.SetupProducerAccounts(producersNames)
		for _, prod := range producersNames {
			assert.Equal(t, e.Success(), e.RegProducer(prod))
			e.ProduceBlocks(1, false)
			assert.True(t, math.Abs(e.GetProducerInfo(prod)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
		}
	}

	// producvotera votes for defproducera ... defproducerj
	// producvoterb votes for defproducera ... defproduceru
	// producvoterc votes for defproducera ... defproducerz
	// producvoterd votes for abcproducera ... abcproducern
	assert.Equal(t, e.Success(), e.Vote(voter1, producersNames[0:10], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter2, producersNames[0:21], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter3, producersNames[0:26], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter4, producersNames[26:39], common.AccountName(0)))

	{
		proda := e.GetProducerInfo(common.N("defproducera"))
		prodj := e.GetProducerInfo(common.N("defproducerj"))
		prodk := e.GetProducerInfo(common.N("defproducerk"))
		produ := e.GetProducerInfo(common.N("defproduceru"))
		prodv := e.GetProducerInfo(common.N("defproducerv"))
		prodz := e.GetProducerInfo(common.N("defproducerz"))
		assert.True(t, proda["unpaid_blocks"].(uint32) == 0 && prodz["unpaid_blocks"].(uint32) == 0)
		assert.True(t, proda["last_claim_time"].(uint64) == 0 && prodz["last_claim_time"].(uint64) == 0)

		// check vote ratios
		assert.True(t, proda["total_votes"].(float64) > 0 && prodz["total_votes"].(float64) > 0)
		assert.True(t, proda["total_votes"].(float64) == prodj["total_votes"].(float64))
		assert.True(t, prodk["total_votes"].(float64) == produ["total_votes"].(float64))
		assert.True(t, prodv["total_votes"].(float64) == prodz["total_votes"].(float64))
		assert.True(t, 2 * proda["total_votes"].(float64) == 3 * produ["total_votes"].(float64))
		assert.True(t, proda["total_votes"].(float64) == 3 * prodz["total_votes"].(float64))
	}

	// give a chance for everyone to produce blocks
	{
		e.ProduceBlocks(23*12+20, false)
		all21Produced := true
		for i := 0; i < 21; i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) == uint32(0) {
				all21Produced = false
			}
		}
		restDidntProduce := true
		for i := 21; i < len(producersNames); i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > uint32(0) {
				restDidntProduce = false
			}
		}
		assert.True(t, all21Produced && restDidntProduce)
	}

	var voteShares []float64
	{
		totalVotes := float64(0)
		thisVotes := float64(0)
		for i := 0; i < len(producersNames); i++ {
			thisVotes = e.GetProducerInfo(producersNames[i])["total_votes"].(float64)
			voteShares = append(voteShares, thisVotes)
			totalVotes += thisVotes
		}
		assert.True(t, math.Abs(e.GetGlobalState()["total_producer_vote_weight"].(float64) - totalVotes)/totalVotes < EPSINON)
		accum := float64(0)
		for i, v := range voteShares {
			voteShares[i] = v / totalVotes
			accum += voteShares[i]
		}
		assert.True(t, math.Abs(float64(1) - accum) < EPSINON)
		assert.True(t, math.Abs(float64(3./71.) - voteShares[0]) < EPSINON)
		assert.True(t, math.Abs(float64(1./71.) - voteShares[38]) < EPSINON)
	}

	{
		prodIndex := 2
		prodName := producersNames[prodIndex]
		initialGlobalState := e.GetGlobalState()
		initialClaimTime := initialGlobalState["last_pervote_bucket_fill"].(uint64)
		//initialPervoteBucket := initialGlobalState["pervote_bucket"].(int64)
		//initialPerblockBucket := initialGlobalState["perblock_bucket"].(int64)
		initialSavings := e.GetBalance(eosioSaving).Amount
		initialTotUnpaidBlocks := initialGlobalState["total_unpaid_blocks"].(uint32)
		initialSupply := e.GetTokenSupply()
		//initialBpayBalance := e.GetBalance(eosioBpay)
		//initialVpayBalance := e.GetBalance(eosioVpay)
		initialBalance := e.GetBalance(prodName)
		initialUnpaidBlocks := e.GetProducerInfo(prodName)["unpaid_blocks"].(uint32)

		assert.Equal(t, e.Success(), e.ClaimRewards(prodName))

		globalState := e.GetGlobalState()
		claimTime := globalState["last_pervote_bucket_fill"].(uint64)
		pervoteBucket := globalState["pervote_bucket"].(int64)
		perblockBucket := globalState["perblock_bucket"].(int64)
		savings := e.GetBalance(eosioSaving).Amount
		//totUnpaidBlocks := globalState["total_unpaid_blocks"].(uint32)
		supply := e.GetTokenSupply()
		bpayBalance := e.GetBalance(eosioBpay)
		vpayBalance := e.GetBalance(eosioVpay)
		balance := e.GetBalance(prodName)
		//unpaidBlocks := e.GetProducerInfo(prodName)["unpaid_blocks"].(uint32)

		usecsBetweenFills := claimTime - initialClaimTime
		//secsBetweenFills := usecsBetweenFills / 1000000

		expectSupplyGrowth := float64(initialSupply.Amount) * float64(usecsBetweenFills) * contRate / float64(usecsPerYear)
		assert.Equal(t, int64(expectSupplyGrowth), supply.Amount - initialSupply.Amount)

		assert.Equal(t, int64(expectSupplyGrowth) - int64(expectSupplyGrowth)/5, savings - initialSavings)

		expectedPerblockBucket := int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * (0.25 * contRate / 5.) / float64(usecsPerYear))
		expectedPervoteBucket := int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * (0.75 * contRate / 5.) / float64(usecsPerYear))

		fromPerblockBucket := int64(initialUnpaidBlocks) * expectedPerblockBucket / int64(initialTotUnpaidBlocks)
		fromPervoteBucket := int64(voteShares[prodIndex] * float64(expectedPervoteBucket))

		if fromPervoteBucket >= 100 * 10000 {
			assert.True(t, withinOne(fromPerblockBucket + fromPervoteBucket, balance.Amount - initialBalance.Amount))
			assert.True(t, withinOne(expectedPervoteBucket - fromPervoteBucket, pervoteBucket))
		} else {
			assert.True(t, withinOne(fromPerblockBucket, balance.Amount - initialBalance.Amount))
			assert.True(t, withinOne(expectedPervoteBucket, pervoteBucket))
			assert.True(t, withinOne(expectedPervoteBucket, vpayBalance.Amount))
			assert.True(t, withinOne(perblockBucket, bpayBalance.Amount))
		}

		e.ProduceBlocks(5, false)
		var ex string
		Try(func() {
			e.ClaimRewards(prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	{
		prodIndex := 23
		prodName := producersNames[prodIndex]
		assert.Equal(t, e.Success(), e.ClaimRewards(prodName))
		assert.Equal(t, int64(0), e.GetBalance(prodName).Amount)
		var ex string
		Try(func() {
			e.ClaimRewards(prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	// Wait for 23 hours. By now, pervote_bucket has grown enough
	// that a producer's share is more than 100 tokens.
	e.ProduceBlock(common.Seconds(23 * 3600), 0)
	{
		prodIndex := 15
		prodName := producersNames[prodIndex]
		initialGlobalState := e.GetGlobalState()
		initialClaimTime := initialGlobalState["last_pervote_bucket_fill"].(uint64)
		initialPervoteBucket := initialGlobalState["pervote_bucket"].(int64)
		initialPerblockBucket := initialGlobalState["perblock_bucket"].(int64)
		initialSavings := e.GetBalance(eosioSaving).Amount
		initialTotUnpaidBlocks := initialGlobalState["total_unpaid_blocks"].(uint32)
		initialSupply := e.GetTokenSupply()
		//initialBpayBalance := e.GetBalance(eosioBpay)
		//initialVpayBalance := e.GetBalance(eosioVpay)
		initialBalance := e.GetBalance(prodName)
		initialUnpaidBlocks := e.GetProducerInfo(prodName)["unpaid_blocks"].(uint32)

		assert.Equal(t, e.Success(), e.ClaimRewards(prodName))

		globalState := e.GetGlobalState()
		claimTime := globalState["last_pervote_bucket_fill"].(uint64)
		pervoteBucket := globalState["pervote_bucket"].(int64)
		perblockBucket := globalState["perblock_bucket"].(int64)
		savings := e.GetBalance(eosioSaving).Amount
		totUnpaidBlocks := globalState["total_unpaid_blocks"].(uint32)
		supply := e.GetTokenSupply()
		bpayBalance := e.GetBalance(eosioBpay)
		vpayBalance := e.GetBalance(eosioVpay)
		balance := e.GetBalance(prodName)
		unpaidBlocks := e.GetProducerInfo(prodName)["unpaid_blocks"].(uint32)

		usecsBetweenFills := claimTime - initialClaimTime

		expectSupplyGrowth := float64(initialSupply.Amount) * float64(usecsBetweenFills) * contRate / float64(usecsPerYear)
		assert.Equal(t, int64(expectSupplyGrowth), supply.Amount - initialSupply.Amount)
		assert.Equal(t, int64(expectSupplyGrowth) - int64(expectSupplyGrowth)/5, savings - initialSavings)

		expectedPerblockBucket := int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * (0.25 * contRate / 5.) / float64(usecsPerYear)) + initialPerblockBucket
		expectedPervoteBucket := int64(float64(initialSupply.Amount) * float64(usecsBetweenFills) * (0.75 * contRate / 5.) / float64(usecsPerYear)) + initialPervoteBucket

		fromPerblockBucket := int64(initialUnpaidBlocks) * expectedPerblockBucket / int64(initialTotUnpaidBlocks)
		fromPervoteBucket := int64(voteShares[prodIndex] * float64(expectedPervoteBucket))

		assert.True(t, withinOne(int64(initialTotUnpaidBlocks - totUnpaidBlocks), int64(initialUnpaidBlocks - unpaidBlocks)))
		if fromPervoteBucket >= 100 * 10000 {
			assert.True(t, withinOne(fromPerblockBucket + fromPervoteBucket, balance.Amount - initialBalance.Amount))
			assert.True(t, withinOne(expectedPervoteBucket - fromPervoteBucket, pervoteBucket))
			assert.True(t, withinOne(expectedPervoteBucket - fromPervoteBucket, vpayBalance.Amount))
			assert.True(t, withinOne(expectedPerblockBucket - fromPerblockBucket, perblockBucket))
			assert.True(t, withinOne(expectedPerblockBucket - fromPerblockBucket, bpayBalance.Amount))
		} else {
			assert.True(t, withinOne(fromPerblockBucket, balance.Amount - initialBalance.Amount))
			assert.True(t, withinOne(expectedPervoteBucket, pervoteBucket))
		}

		e.ProduceBlocks(5, false)
		var ex string
		Try(func() {
			e.ClaimRewards(prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	{
		prodIndex := 24
		prodName := producersNames[prodIndex]
		assert.Equal(t, e.Success(), e.ClaimRewards(prodName))
		assert.True(t, 100 * 10000 <= e.GetBalance(prodName).Amount)
		var ex string
		Try(func() {
			e.ClaimRewards(prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "already claimed rewards within past day"))
	}

	{
		rmvIndex := 5
		prodName := producersNames[rmvIndex]

		info := e.GetProducerInfo(prodName)
		assert.True(t, info["is_active"].(bool))
		assert.NotEqual(t, *ecc.NewPublicKeyNil(), info["producer_key"].(ecc.PublicKey))
		var ex string
		Try(func() {
			e.RmvProducer(prodName, prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "missing authority of eosio"))

		Try(func() {
			e.RmvProducer(producersNames[rmvIndex + 2], prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "missing authority of eosio"))

		assert.Equal(t, e.Success(), e.RmvProducer(eosio, prodName))
		{
			restDidntProduce := true
			for i := 21; i < len(producersNames); i++ {
				if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > 0{
					restDidntProduce = false
				}
			}
			assert.True(t, restDidntProduce)
		}

		e.ProduceBlocks(3*21*12, false)
		info = e.GetProducerInfo(prodName)
		initUnpaidBlocks := info["unpaid_blocks"].(uint32)
		assert.True(t, !info["is_active"].(bool))
		assert.Equal(t, *ecc.NewPublicKeyNil(), info["producer_key"].(ecc.PublicKey))
		Try(func() {
			e.ClaimRewards(prodName)
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer does not have an active key"))

		e.ProduceBlocks(3*21*12, false)
		assert.Equal(t, initUnpaidBlocks, e.GetProducerInfo(prodName)["unpaid_blocks"].(uint32))
		{
			prodWasReplaced := false
			for i := 21; i < len(producersNames); i++ {
				if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > 0{
					prodWasReplaced = true
				}
			}
			assert.True(t, prodWasReplaced)
		}
	}

	{
		var ex string
		Try(func() {
			e.RmvProducer(eosio, common.N("nonexistingp"))
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "producer not found"))
	}
	e.close()
}

func TestProducersUpgradeSystemContract(t *testing.T) {
	e := initEosioSystemTester()

	//install multisig contract
	msigAbiSer := e.InitializeMultisig()
	producersNames := e.ActiveAndVoteProducers()

	//helper function
	pushActionMsig := func(signer common.AccountName, name common.ActionName, data common.Variants, auth bool) ActionResult {
		actionTypeName := msigAbiSer.GetActionType(name)
		act := types.Action{}
		act.Account = eosioMsig
		act.Name = name
		act.Data = msigAbiSer.VariantToBinary(actionTypeName, &data, e.AbiSerializerMaxTime)
		var signerAuth common.AccountName
		if auth {
			signerAuth = signer
		} else {
			if signer == bob {
				signerAuth = alice
			} else {
				signerAuth = bob
			}
		}
		return e.PushAction(&act, signerAuth)
	}

	// test begins
	var prodPerms []types.PermissionLevel
	for _, x := range producersNames {
		prodPerms = append(prodPerms, types.PermissionLevel{Actor:x, Permission:common.DefaultConfig.ActiveName})
	}
	//prepare system contract with different hash (contract differs in one byte)
	wast, _ := ioutil.ReadFile("test_contracts/eosio.system.wast")
	eosioSystemWast := string(wast)
	eosioSystemWast = strings.Replace(eosioSystemWast, "producer votes must be unique and sorted", "Producer votes must be unique and sorted", 1)

	trx := types.SignedTransaction{}
	{
		code := wast2wasm([]byte(eosioSystemWast))
		setCode := SetCode{Account: eosio, VmType: 0, VmVersion: 0, Code: code}
		data, _ := rlp.EncodeToBytes(setCode)
		act := types.Action{
			Account:       setCode.GetAccount(),
			Name:          setCode.GetName(),
			Authorization: []types.PermissionLevel{{eosio, common.DefaultConfig.ActiveName}},
			Data:          data,
		}
		trx.Actions = append(trx.Actions, &act)
		e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta + 9, 0)
		fmt.Println(trx.Expiration.SecSinceEpoch())
		trx.Transaction.RefBlockNum = 2
		trx.Transaction.RefBlockPrefix = 3
	}
	data := common.Variants{
		"proposer":      alice,
		"proposal_name": common.N("upgrade1"),
		"requested":     prodPerms,
		"trx":           trx.Transaction,
	}

	assert.Equal(t, e.Success(), pushActionMsig(alice, common.N("propose"), data, true))

	// get 15 approvals
	for i := 0; i < 14; i++ {
		data = common.Variants{
			"proposer":      alice,
			"proposal_name": common.N("upgrade1"),
			"level":         types.PermissionLevel{Actor:producersNames[i], Permission:common.DefaultConfig.ActiveName},
		}
		assert.Equal(t, e.Success(), pushActionMsig(producersNames[i], common.N("approve"), data, true))
	}

	//should fail
	var ex string
	Try(func() {
		data = common.Variants{
			"proposer":      alice,
			"proposal_name": common.N("upgrade1"),
			"executer":      alice,
		}
		pushActionMsig(alice, common.N("exec"), data, true)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "transaction authorization failed"))

	// one more approval
	data = common.Variants{
		"proposer":      alice,
		"proposal_name": common.N("upgrade1"),
		"level":         types.PermissionLevel{Actor:producersNames[14], Permission:common.DefaultConfig.ActiveName},
	}
	assert.Equal(t, e.Success(), pushActionMsig(producersNames[14], common.N("approve"), data, true))

	data = common.Variants{
		"proposer":      alice,
		"proposal_name": common.N("upgrade1"),
		"executer":      alice,
	}
	pushActionMsig(alice, common.N("exec"), data, true)
	e.ProduceBlocks(250, false)
}

func TestProducerOnblockCheck(t *testing.T) {
	e := initEosioSystemTester()
	largeAsset := CoreFromString("80.0000")
	e.CreateAccountWithResources(voter1, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(voter2, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(voter3, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)

	// create accounts {defproducera, defproducerb, ..., defproducerz} and register as producers
	var producersNames []common.AccountName
	root := "defproducer"
	for c := 'a'; c <= 'z'; c++ {
		acc := common.N(root + string(c))
		producersNames = append(producersNames, acc)
	}
	e.SetupProducerAccounts(producersNames)

	for _, a := range producersNames {
		e.RegProducer(a)
	}

	assert.True(t, math.Abs(e.GetProducerInfo(producersNames[0])["total_votes"].(float64)-e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producersNames[len(producersNames)-1])["total_votes"].(float64)-e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	e.Transfer(eosio, voter1, CoreFromString("200000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(voter1, voter1, CoreFromString("70000000.0000"), CoreFromString("70000000.0000")))
	assert.Equal(t, e.Success(), e.Vote(voter1, producersNames[0:10], common.AccountName(0)))
	var ex string
	Try(func() {
		e.UnStake(voter1, voter1, CoreFromString("50.0000"), CoreFromString("50.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "cannot undelegate bandwidth until the chain is activated (at least 15% of all tokens participate in voting)"))

	// give a chance for everyone to produce blocks
	{
		e.ProduceBlocks(21*12, false)
		all21Produced := true
		for i := 0; i < 21; i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) == uint32(0) {
				all21Produced = false
			}
		}
		restDidntProduce := true
		for i := 21; i < len(producersNames); i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > uint32(0) {
				restDidntProduce = false
			}
		}
		assert.True(t, !all21Produced && restDidntProduce)
	}

	// stake across 15% boundary
	e.Transfer(eosio, voter2, CoreFromString("100000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(voter2, voter2, CoreFromString("4000000.0000"), CoreFromString("4000000.0000")))
	e.Transfer(eosio, voter3, CoreFromString("100000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(voter3, voter3, CoreFromString("2000000.0000"), CoreFromString("2000000.0000")))

	assert.Equal(t, e.Success(), e.Vote(voter2, producersNames[0:21], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter3, producersNames, common.AccountName(0)))

	// give a chance for everyone to produce blocks
	{
		e.ProduceBlocks(21*12, false)
		all21Produced := true
		for i := 0; i < 21; i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) == 0 {
				all21Produced = false
			}
		}
		restDidntProduce := true
		for i := 21; i < len(producersNames); i++ {
			if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > 0 {
				restDidntProduce = false
			}
		}
		assert.True(t, all21Produced && restDidntProduce)
		assert.Equal(t, e.Success(), e.ClaimRewards(producersNames[0]))
		assert.True(t, e.GetBalance(producersNames[0]).Amount > 0)
	}

	assert.Equal(t, e.Success(), e.UnStake(voter1, voter1, CoreFromString("50.0000"), CoreFromString("50.0000")))
	e.close()
}

func TestVotersActionsAffectProxyAndProducers(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	donald := common.N("donald111111")
	e.CreateAccountsWithResources([]common.AccountName{donald, producer1, producer2, producer3}, eosio)
	assert.Equal(t, e.Success(), e.RegProducer(producer1))
	assert.Equal(t, e.Success(), e.RegProducer(producer2))
	assert.Equal(t, e.Success(), e.RegProducer(producer3))

	//alice1111111 becomes a proxy
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.Equal(t, e.Proxy(alice), e.GetVoterInfo(alice))

	//alice1111111 makes stake and votes
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("30.0001"), CoreFromString("20.0001")))
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1, producer2}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("50.0002"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("50.0002"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//donald111111 becomes a proxy
	assert.Equal(t, e.Success(), e.RegProxy(donald))
	assert.Equal(t, e.Proxy(donald), e.GetVoterInfo(donald))

	//bob111111111 chooses alice1111111 as a proxy
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{}, alice))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//carol1111111 chooses alice1111111 as a proxy
	e.Issue(carol, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("30.0001"), CoreFromString("20.0001")))
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{}, alice))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("250.0007"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("250.0007"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//proxied voter carol1111111 increases stake
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("50.0000"), CoreFromString("70.0000")))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("320.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("370.0007"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("370.0007"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//proxied voter bob111111111 decreases stake
	assert.Equal(t, e.Success(), e.UnStake(bob, bob, CoreFromString("50.0001"), CoreFromString("50.0001")))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("220.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("270.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("270.0005"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//proxied voter carol1111111 chooses another proxy
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{}, donald))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("50.0001"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetVoterInfo(donald)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("170.0002"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("100.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("100.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	//bob111111111 switches to direct voting and votes for one of the same producers, but not for another one
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{producer2}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("50.0002"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("100.0003"))) <= EPSINON)
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)

	e.close()
}

func TestVoteBothProxyAndProducers(t *testing.T) {
	e := initEosioSystemTester()
	//alice1111111 becomes a proxy
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.Equal(t, e.Proxy(alice), e.GetVoterInfo(alice))

	//carol1111111 becomes a producer
	assert.Equal(t, e.Success(), e.RegProducer(carol))

	//bob111111111 chooses alice1111111 as a proxy
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))
	var ex string
	Try(func() {
		e.Vote(bob, []common.AccountName{carol}, alice)
	}).Catch(func(e Exception){
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "cannot vote for producers and proxy at same time"))

	e.close()
}

func TestSelectInvalidProxy(t *testing.T) {
	e := initEosioSystemTester()
	//accumulate proxied votes
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))

	//selecting account not registered as a proxy
	var ex string
	Try(func() {
		e.Vote(bob, []common.AccountName{}, alice)
	}).Catch(func(e Exception){
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "invalid proxy specified"))

	//selecting not existing account as a proxy
	Try(func() {
		e.Vote(bob, []common.AccountName{}, common.N("notexist"))
	}).Catch(func(e Exception){
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "invalid proxy specified"))

	e.close()
}

func TestDoubleRegisterUnregisterProxyKeepsVotes(t *testing.T) {
	e := initEosioSystemTester()
	//alice1111111 becomes a proxy
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("5.0000"), CoreFromString("5.0000")))
	assert.Equal(t, e.ProxyStake(alice, CoreFromString("10.0000")), e.GetVoterInfo(alice))

	//bob111111111 stakes and selects alice1111111 as a proxy
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{}, alice))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//double registering should fail without affecting total votes and stake
	{
		var ex string
		Try(func() {
			e.RegProxy(alice)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "action has no effect"))
	}
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//uregister
	assert.Equal(t, e.Success(), e.UnRegProxy(alice))
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)

	//double unregistering should not affect proxied_votes and stake
	{
		var ex string
		Try(func() {
			e.UnRegProxy(alice)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "action has no effect"))
	}
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= EPSINON)
}

func TestProxyCannotUseAnotherProxy(t *testing.T) {
	e := initEosioSystemTester()
	//alice1111111 and bob111111111 become proxies
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.Equal(t, e.Success(), e.RegProxy(bob))

	//proxy should not be able to use a proxy
	e.Issue(bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("100.0002"), CoreFromString("50.0001")))
	{
		var ex string
		Try(func() {
			e.Vote(bob, []common.AccountName{}, alice)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "account registered as a proxy is not allowed to use a proxy"))
	}

	//voter that uses a proxy should not be allowed to become a proxy
	e.Issue(carol, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(carol, carol, CoreFromString("100.0002"), CoreFromString("50.0001")))
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{}, alice))
	{
		var ex string
		Try(func() {
			e.RegProxy(carol)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "account that uses a proxy is not allowed to become a proxy"))
	}

	//proxy should not be able to use itself as a proxy
	{
		var ex string
		Try(func() {
			e.Vote(bob, []common.AccountName{}, bob)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "cannot proxy to self"))
	}
}

func configToVariant(config types.ChainConfig) common.Variants{
	return common.Variants{
		"max_block_net_usage": config.MaxBlockNetUsage ,
		"target_block_net_usage_pct": config.TargetBlockNetUsagePct ,
		"max_transaction_net_usage": config.MaxTransactionNetUsage ,
		"base_per_transaction_net_usage": config.BasePerTransactionNetUsage ,
		"context_free_discount_net_usage_num": config.ContextFreeDiscountNetUsageNum ,
		"context_free_discount_net_usage_den": config.ContextFreeDiscountNetUsageDen ,
		"max_block_cpu_usage": config.MaxBlockCpuUsage ,
		"target_block_cpu_usage_pct": config.TargetBlockCpuUsagePct ,
		"max_transaction_cpu_usage": config.MaxTransactionCpuUsage ,
		"min_transaction_cpu_usage": config.MinTransactionCpuUsage ,
		"max_transaction_lifetime": config.MaxTrxLifetime ,
		"deferred_trx_expiration_window": config.DeferredTrxExpirationWindow ,
		"max_transaction_delay": config.MaxTrxDelay ,
		"max_inline_action_size": config.MaxInlineActionSize ,
		"max_inline_action_depth": config.MaxInlineActionDepth ,
		"max_authority_depth": config.MaxAuthorityDepth ,
	}
}

func TestElectProducers(t *testing.T) {
	e := initEosioSystemTester()
	e.CreateAccountsWithResources([]common.AccountName{producer1, producer2, producer3}, eosio)
	assert.Equal(t, e.Success(), e.RegProducer(producer1))
	assert.Equal(t, e.Success(), e.RegProducer(producer2))
	assert.Equal(t, e.Success(), e.RegProducer(producer3))

	//stake more than 15% of total EOS supply to activate chain
	e.Transfer(eosio, alice, CoreFromString("600000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("300000000.0000"), CoreFromString("300000000.0000")))

	//vote for producers
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1}, common.AccountName(0)))
	e.ProduceBlocks(250, false)
	producerKeys := e.Control.HeadBlockState().ActiveSchedule.Producers
	assert.Equal(t, int(1), len(producerKeys))
	assert.Equal(t, producer1, producerKeys[0].ProducerName)

	// elect 2 producers
	e.Issue(bob, CoreFromString("80000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(bob, bob, CoreFromString("40000.0000"), CoreFromString("40000.0000")))
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{producer2}, common.AccountName(0)))
	e.ProduceBlocks(250, false)
	producerKeys = e.Control.HeadBlockState().ActiveSchedule.Producers
	assert.Equal(t, int(2), len(producerKeys))
	assert.Equal(t, producer1, producerKeys[0].ProducerName)
	assert.Equal(t, producer2, producerKeys[1].ProducerName)

	// elect 3 producers
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{producer2, producer3}, common.AccountName(0)))
	e.ProduceBlocks(250, false)
	producerKeys = e.Control.HeadBlockState().ActiveSchedule.Producers
	assert.Equal(t, int(3), len(producerKeys))
	assert.Equal(t, producer1, producerKeys[0].ProducerName)
	assert.Equal(t, producer2, producerKeys[1].ProducerName)
	assert.Equal(t, producer3, producerKeys[2].ProducerName)

	// try to go back to 2 producers and fail
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{producer3}, common.AccountName(0)))
	e.ProduceBlocks(250, false)
	producerKeys = e.Control.HeadBlockState().ActiveSchedule.Producers
	assert.Equal(t, int(3), len(producerKeys))

	e.close()
}

func TestBuyName(t *testing.T) {
	e := initEosioSystemTester()
	dan := common.N("dan")
	sam := common.N("sam")
	e.CreateAccountsWithResources([]common.AccountName{dan, sam}, eosio)
	e.Transfer(eosio, dan, CoreFromString("10000.0000"), eosio)
	e.Transfer(eosio, sam, CoreFromString("10000.0000"), eosio)
	e.StakeWithTransfer(eosio, sam, CoreFromString("80000000.0000"), CoreFromString("80000000.0000"))
	e.StakeWithTransfer(eosio, dan, CoreFromString("80000000.0000"), CoreFromString("80000000.0000"))
	e.RegProducer(eosio)
	assert.Equal(t, e.Success(), e.Vote(sam, []common.AccountName{eosio}, common.AccountName(0)))
	// wait 14 days after min required amount has been staked
	e.ProduceBlock(common.Days(7), 0)
	assert.Equal(t, e.Success(), e.Vote(dan, []common.AccountName{eosio}, common.AccountName(0)))
	e.ProduceBlock(common.Days(7), 0)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// dan shouldn't be able to create fail
	{
		var ex string
		Try(func() {
			e.CreateAccountsWithResources([]common.AccountName{common.N("fail")}, dan)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "no active bid for name"))
	}
	e.BidName(dan, common.N("nofail"), CoreFromString("1.0000"))

	// didn't increase bid by 10%
	{
		var ex string
		Try(func() {
			e.BidName(sam, common.N("nofail"), CoreFromString("1.0000"))
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "assertion failure with message: must increase bid by 10%"))
	}
	e.BidName(sam, common.N("nofail"), CoreFromString("2.0000"))
	e.ProduceBlock(common.Days(1), 0)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// dan shouldn't be able to do this, sam won
	{
		var ex string
		Try(func() {
			e.CreateAccountsWithResources([]common.AccountName{common.N("nofail")}, dan)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "only highest bidder can claim"))
	}
	e.CreateAccountsWithResources([]common.AccountName{common.N("nofail")}, sam)
	e.Transfer(eosio, common.N("nofail"), CoreFromString("1000.0000"), eosio)

	// only nofail can create test.nofail
	e.CreateAccountsWithResources([]common.AccountName{common.N("test.nofail")}, common.N("nofail"))
	{
		var ex string
		Try(func() {
			e.CreateAccountsWithResources([]common.AccountName{common.N("test.fail")}, dan)
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "only suffix may create this account"))
	}

	e.close()
}

func TestInvalidNames(t *testing.T) {
	e := initEosioSystemTester()
	dan := common.N("dan")
	e.CreateAccountsWithResources([]common.AccountName{dan}, eosio)
	{
		var ex string
		Try(func() {
			e.BidName(dan, common.N("abcdefg.12345"), CoreFromString("1.0000"))
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		fmt.Println(ex)
		assert.True(t, inString(ex, "you can only bid on top-level suffix"))
	}

	{
		var ex string
		Try(func() {
			e.BidName(dan, common.N(""), CoreFromString("1.0000"))
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "the empty name is not a valid account name to bid on"))
	}

	{
		var ex string
		Try(func() {
			e.BidName(dan, common.N("abcdefgh12345"), CoreFromString("1.0000"))
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		fmt.Println(ex)
		assert.True(t, inString(ex, "13 character names are not valid account names to bid on"))
	}

	{
		var ex string
		Try(func() {
			e.BidName(dan, common.N("abcdefg12345"), CoreFromString("1.0000"))
		}).Catch(func(e Exception){
			ex = e.DetailMessage()
		})
		assert.True(t, inString(ex, "accounts with 12 character names and no dots can be created without bidding required"))
	}

	e.close()
}

func TestMultipleNameBids(t *testing.T) {
	e := initEosioSystemTester()
	alice := common.N("alice")
	bob := common.N("bob")
	carl := common.N("carl")
	david := common.N("david")
	eve := common.N("eve")
	producer := common.N("producer")
	accounts := []common.AccountName{alice, bob, carl, david, eve}
	e.CreateAccountsWithResources(accounts, eosio)
	for _, a := range accounts {
		e.Transfer(eosio, a, CoreFromString("10000.0000"), eosio)
		assert.Equal(t, CoreFromString("10000.0000"), e.GetBalance(a))
	}
	e.CreateAccountsWithResources([]common.AccountName{producer}, eosio)
	assert.Equal(t, e.Success(), e.RegProducer(producer))

	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	// stake but but not enough to go live
	e.StakeWithTransfer(eosio, bob, CoreFromString("35000000.0000"), CoreFromString("35000000.0000"))
	e.StakeWithTransfer(eosio, carl, CoreFromString("35000000.0000"), CoreFromString("35000000.0000"))
	assert.Equal(t, e.Success(), e.Vote(bob, []common.AccountName{producer}, common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(carl, []common.AccountName{producer}, common.AccountName(0)))

	// start bids
	e.BidName(bob, common.N("prefa"), CoreFromString("1.0003"))
	assert.Equal(t, CoreFromString("9998.9997"), e.GetBalance(bob))
	e.BidName(bob, common.N("prefb"), CoreFromString("1.0000"))
	e.BidName(bob, common.N("prefc"), CoreFromString("1.0000"))
	assert.Equal(t, CoreFromString("9996.9997"), e.GetBalance(bob))

	e.BidName(carl, common.N("prefd"), CoreFromString("1.0000"))
	e.BidName(carl, common.N("prefe"), CoreFromString("1.0000"))
	assert.Equal(t, CoreFromString("9998.0000"), e.GetBalance(carl))

	var ex string
	Try(func() {
		e.BidName(bob, common.N("prefb"), CoreFromString("1.1001"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "assertion failure with message: account is already highest bidder"))
	Try(func() {
		e.BidName(alice, common.N("prefb"), CoreFromString("1.0999"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "assertion failure with message: must increase bid by 10%"))
	assert.Equal(t, CoreFromString("9996.9997"), e.GetBalance(bob))
	assert.Equal(t, CoreFromString("10000.0000"), e.GetBalance(alice))

	// alice outbids bob on prefb
	initialNamesBalance := e.GetBalance(eosioName)
	assert.Equal(t, e.Success(), e.BidName(alice, common.N("prefb"), CoreFromString("1.1001")))
	assert.Equal(t, CoreFromString("9997.9997"), e.GetBalance(bob))
	assert.Equal(t, CoreFromString("9998.8999"), e.GetBalance(alice))
	assert.Equal(t, initialNamesBalance.Amount + CoreFromString("0.1001").Amount, e.GetBalance(eosioName).Amount)

	// david outbids carl on prefd
	assert.Equal(t, CoreFromString("9998.0000"), e.GetBalance(carl))
	assert.Equal(t, CoreFromString("10000.0000"), e.GetBalance(david))
	assert.Equal(t, e.Success(), e.BidName(david, common.N("prefd"), CoreFromString("1.9900")))
	assert.Equal(t, CoreFromString("9999.0000"), e.GetBalance(carl))
	assert.Equal(t, CoreFromString("9998.0100"), e.GetBalance(david))

	// eve outbids carl on prefe
	assert.Equal(t, e.Success(), e.BidName(eve, common.N("prefe"), CoreFromString("1.7200")))

	e.ProduceBlock(common.Days(14), 0)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// highest bid is from david for prefd but no bids can be closed yet
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefd"), david, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))

	// stake enough to go above the 15% threshold
	e.StakeWithTransfer(eosio, alice, CoreFromString("10000000.0000"), CoreFromString("10000000.0000"))
	assert.Equal(t, uint32(0), e.GetProducerInfo(producer)["unpaid_blocks"].(uint32))
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer}, common.AccountName(0)))

	// need to wait for 14 days after going live
	e.ProduceBlocks(10, false)
	e.ProduceBlock(common.Days(2), 0)
	e.ProduceBlocks(10, false)
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefd"), david, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))

	// it's been 14 days, auction for prefd has been closed
	e.ProduceBlock(common.Days(12), 0)
	e.CreateAccountWithResources2(common.N("prefd"), david, 8000)
	e.ProduceBlocks(2, false)
	e.ProduceBlock(common.Hours(23), 0)

	// auctions for prefa, prefb, prefc, prefe haven't been closed
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefa"), bob, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefb"), alice, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefc"), bob, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefe"), eve, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))

	// attempt to create account with no bid
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefg"), alice, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "no active bid for name"))

	// changing highest bid pushes auction closing time by 24 hours
	assert.Equal(t, e.Success(), e.BidName(eve, common.N("prefb"), CoreFromString("2.1880")))
	e.ProduceBlock(common.Hours(22), 0)
	e.ProduceBlocks(2, false)
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefb"), eve, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))

	// but changing a bid that is not the highest does not push closing time
	assert.Equal(t, e.Success(), e.BidName(carl, common.N("prefe"), CoreFromString("2.0980")))
	e.ProduceBlock(common.Hours(2), 0)
	e.ProduceBlocks(2, false)
	// bid for prefb has closed, only highest bidder can claim
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefb"), alice, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "only highest bidder can claim"))
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefb"), carl, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "only highest bidder can claim"))
	e.CreateAccountWithResources2(common.N("prefb"), eve, 8000)

	Try(func() {
		e.CreateAccountWithResources2(common.N("prefe"), carl, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	e.ProduceBlock(common.Hours(24), 0)
	// by now bid for prefe has closed
	e.CreateAccountWithResources2(common.N("prefe"), carl, 8000)

	// prefe can now create *.prefe
	Try(func() {
		e.CreateAccountWithResources2(common.N("xyz.prefe"), carl, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "only suffix may create this account"))
	e.Transfer(eosio, common.N("prefe"), CoreFromString("10000.0000"), eosio)
	e.CreateAccountWithResources2(common.N("xyz.prefe"), common.N("prefe"), 8000)

	// other auctions haven't closed
	Try(func() {
		e.CreateAccountWithResources2(common.N("prefa"), bob, 8000)
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "auction for name is not closed yet"))
}

func TestVoteProducersInAndOut(t *testing.T) {
	e := initEosioSystemTester()
	net := CoreFromString("80.0000")
	cpu := CoreFromString("80.0000")
	voters := []common.AccountName{voter1, voter2, voter3, voter4}
	for _, vote := range voters {
		e.CreateAccountWithResources(vote, eosio, CoreFromString("1.0000"), false, net, cpu)
	}

	// create accounts {defproducera, defproducerb, ..., defproducerz} and register as producers
	var producersNames []common.AccountName
	{
		root := "defproducer"
		for c := 'a'; c <= 'z'; c++ {
			acc := common.N(root + string(c))
			producersNames = append(producersNames, acc)
		}
		e.SetupProducerAccounts(producersNames)
		for _, prod := range producersNames {
			assert.Equal(t, e.Success(), e.RegProducer(prod))
			e.ProduceBlocks(1, false)
			assert.True(t, math.Abs(e.GetProducerInfo(prod)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= EPSINON)
		}
	}

	for _, vote := range voters {
		e.Transfer(eosio, vote, CoreFromString("200000000.0000"), eosio)
		assert.Equal(t, e.Success(), e.Stake(vote, vote, CoreFromString("30000000.0000"), CoreFromString("30000000.0000")))
	}
	fmt.Println(producersNames[0:20])
	assert.Equal(t, e.Success(), e.Vote(voter1, producersNames[0:20], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter2, producersNames[0:21], common.AccountName(0)))
	assert.Equal(t, e.Success(), e.Vote(voter3, producersNames, common.AccountName(0)))

	// give a chance for everyone to produce blocks
	e.ProduceBlocks(23 * 12 + 20, false)
	all21Produced := true
	for i := 0; i < 21; i++ {
		if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) == uint32(0) {
			all21Produced = false
		}
	}
	restDidntProduce := true
	for i := 21; i < len(producersNames); i++ {
		if e.GetProducerInfo(producersNames[i])["unpaid_blocks"].(uint32) > uint32(0) {
			restDidntProduce = false
		}
	}
	assert.True(t, all21Produced && restDidntProduce)

	{
		e.ProduceBlock(common.Hours(7), 0)
		votedOutIndex := 20
		newProdIndex := 23
		assert.Equal(t, e.Success(), e.Stake(voter4, voter4, CoreFromString("40000000.0000"), CoreFromString("40000000.0000")))
		assert.Equal(t, e.Success(), e.Vote(voter4, []common.AccountName{producersNames[newProdIndex]}, common.AccountName(0)))
		assert.Equal(t, uint32(0), e.GetProducerInfo(producersNames[newProdIndex])["unpaid_blocks"].(uint32))
		e.ProduceBlocks(4 * 12 * 21, false)
		assert.True(t, e.GetProducerInfo(producersNames[newProdIndex])["unpaid_blocks"].(uint32) > uint32(0))
		initialUnpaidBlocks := e.GetProducerInfo(producersNames[votedOutIndex])["unpaid_blocks"].(uint32)
		e.ProduceBlocks(2 * 12 * 21, false)
		assert.Equal(t, initialUnpaidBlocks, e.GetProducerInfo(producersNames[votedOutIndex])["unpaid_blocks"].(uint32))
		e.ProduceBlock(common.Hours(24), 0)
		assert.Equal(t, e.Success(), e.Vote(voter4, []common.AccountName{producersNames[votedOutIndex]}, common.AccountName(0)))
		e.ProduceBlocks(2 * 12 * 21, false)
		assert.True(t, *ecc.NewPublicKeyNil() != e.GetProducerInfo(producersNames[votedOutIndex])["producer_key"].(ecc.PublicKey))
		assert.Equal(t, e.Success(), e.ClaimRewards(producersNames[votedOutIndex]))
	}
	e.close()
}

func TestSetParams(t *testing.T) {
	e := initEosioSystemTester()

	//install multisig contract
	msigAbiSer := e.InitializeMultisig()
	producersNames := e.ActiveAndVoteProducers()

	//helper function
	pushActionMsig := func(signer common.AccountName, name common.ActionName, data common.Variants, auth bool) ActionResult {
		actionTypeName := msigAbiSer.GetActionType(name)
		act := types.Action{}
		act.Account = eosioMsig
		act.Name = name
		act.Data = msigAbiSer.VariantToBinary(actionTypeName, &data, e.AbiSerializerMaxTime)
		var signerAuth common.AccountName
		if auth {
			signerAuth = signer
		} else {
			if signer == bob {
				signerAuth = alice
			} else {
				signerAuth = bob
			}
		}
		return e.PushAction(&act, signerAuth)
	}

	//test begins
	var prodPerms []types.PermissionLevel
	for _, x := range producersNames {
		prodPerms = append(prodPerms, types.PermissionLevel{Actor:x, Permission:common.DefaultConfig.ActiveName})
	}
	params := e.Control.GetGlobalProperties().Configuration

	//change some values
	params.MaxBlockNetUsage += 10
	params.MaxTrxLifetime += 1

	trx := types.SignedTransaction{}
	{
		data := common.Variants{"params": params}
		act := e.GetAction(eosio, common.N("setparams"), []types.PermissionLevel{{eosio, common.DefaultConfig.ActiveName}}, &data)
		trx.Actions = append(trx.Actions, act)
		e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
		fmt.Println(trx.Expiration.SecSinceEpoch())
		trx.Transaction.RefBlockNum = 2
		trx.Transaction.RefBlockPrefix = 3
	}
	data := common.Variants{
		"proposer":      alice,
		"proposal_name": common.N("setparams1"),
		"trx":           trx.Transaction,
		"requested":     prodPerms,
	}
	assert.Equal(t, e.Success(), pushActionMsig(alice, common.N("propose"), data, true))
}

func TestSetRamEffect(t *testing.T) {
	e := initEosioSystemTester()
	net := CoreFromString("8.0000")
	cpu := CoreFromString("8.0000")
	alice := common.N("aliceaccount")
	bobby := common.N("bobbyaccount")
	accounts := []common.AccountName{alice, bobby}
	for _, a := range accounts {
		e.CreateAccountWithResources(a, eosio, CoreFromString("1.0000"), false, net, cpu)
	}

	{
		e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
		assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))
		initByte := e.GetTotalStake(alice)["ram_bytes"].(uint64)
		assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("300.0000")))
		assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
		boughtBytes := e.GetTotalStake(alice)["ram_bytes"].(uint64) - initByte

		// after buying and selling balance should be 700 + 300 * 0.995 * 0.995 = 997.0075 (actually 997.0074 due to rounding fees up)
		assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
		assert.Equal(t, CoreFromString("997.0074"), e.GetBalance(alice))
	}

	{
		e.Transfer(eosio, bobby, CoreFromString("1000.0000"), eosio)
		assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(bobby))
		initByte := e.GetTotalStake(bobby)["ram_bytes"].(uint64)

		// bobby buys ram at current price
		assert.Equal(t, e.Success(), e.BuyRam(bobby, bobby, CoreFromString("300.0000")))
		assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(bobby))
		boughtBytes := e.GetTotalStake(bobby)["ram_bytes"].(uint64) - initByte

		// increase max_ram_size, ram bought by bobby loses part of its value
		var ex string
		Try(func() {
			e.SetRam(eosio, uint64(uint64(64) * 1024 * 1024 * 1024))
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "ram may only be increased"))
		Try(func() {
			e.SetRam(bobby, uint64(uint64(80) * 1024 * 1024 * 1024))
		}).Catch(func(e Exception) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "missing authority of eosio"))
		assert.Equal(t, e.Success(), e.SetRam(eosio, uint64(uint64(80) * 1024 * 1024 * 1024)))

		assert.Equal(t, e.Success(), e.SellRam(bobby, boughtBytes))
		assert.True(t, e.GetBalance(bobby).Amount > 9000000 && e.GetBalance(bobby).Amount < 9500000)
	}
}