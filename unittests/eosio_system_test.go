package unittests

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"math"
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
	var ex string
	Try(func() {
		e.StakeWithTransfer(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000"))
	}).Catch(func(e Exception) {
		ex = e.DetailMessage()
	}).End()
	assert.True(t, inString(ex, "cannot use transfer flag if delegating to self"))
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
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("33.3333"))) <= math.Pow10(-4))

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
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.2220"))) <= math.Pow10(-4))

	//but eos should still be at stake
	assert.Equal(t, CoreFromString("1955.5556"), e.GetBalance(bob))

	//carol1111111 unstakes rest of eos
	assert.Equal(t, e.Success(), e.UnStake(carol, carol, CoreFromString("20.0000"), CoreFromString("0.2220")))

	//should decrease alice1111111's total_votes to zero
	prod = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(prod["total_votes"].(float64) - float64(0)) <= math.Pow10(-4))

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
	assert.True(t, math.Abs(prod["total_votes"].(float64) - float64(0)) <= math.Pow10(-4))

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
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= math.Pow10(-4))
	bobInfo := e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= math.Pow10(-4))

	//carol1111111 votes for alice1111111 (but revokes vote for bob111111111)
	assert.Equal(t, e.Success(), e.Vote(carol, []common.AccountName{alice}, common.AccountName(0)))
	aliceInfo = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("20.0005"))) <= math.Pow10(-4))
	bobInfo = e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))

	//alice1111111 votes for herself and bob111111111
	e.Issue(alice, CoreFromString("2.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("1.0000"), CoreFromString("1.0000")))
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{alice, bob}, common.AccountName(0)))
	aliceInfo = e.GetProducerInfo(alice)
	assert.True(t, math.Abs(aliceInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("22.0005"))) <= math.Pow10(-4))
	bobInfo = e.GetProducerInfo(bob)
	assert.True(t, math.Abs(bobInfo["total_votes"].(float64) - e.Stake2Votes(CoreFromString("2.0000"))) <= math.Pow10(-4))

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
	assert.True(t, math.Abs(e.GetVoterInfo(alice)["proxied_vote_weight"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))

	//vote for producers
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1, producer2}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))

	//vote for another producers
	assert.Equal(t, e.Success(), e.Vote(alice, []common.AccountName{producer1, producer3}, common.AccountName(0)))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))

	//unregister proxy
	assert.Equal(t, e.Success(), e.UnRegProxy(alice))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))

	//register proxy again
	assert.Equal(t, e.Success(), e.RegProxy(alice))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("150.0003"))) <= math.Pow10(-4))

	//stake increase by proxy itself affects producers
	e.Issue(alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(alice, alice, CoreFromString("30.0001"), CoreFromString("20.0001")))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("200.0005"))) <= math.Pow10(-4))

	//stake decrease by proxy itself affects producers
	assert.Equal(t, e.Success(), e.UnStake(alice, alice, CoreFromString("10.0001"), CoreFromString("10.0001")))
	assert.True(t, math.Abs(e.GetProducerInfo(producer1)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("180.0003"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer2)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))
	assert.True(t, math.Abs(e.GetProducerInfo(producer3)["total_votes"].(float64) - e.Stake2Votes(CoreFromString("180.0003"))) <= math.Pow10(-4))

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
	assert.True(t, math.Abs(prod["total_votes"].(float64) - e.Stake2Votes(CoreFromString("0.0000"))) <= math.Pow10(-4))

	e.Transfer(eosio, producer1, CoreFromString("400000000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(producer1, producer1, CoreFromString("100000000.0000"), CoreFromString("100000000.0000")))
	assert.Equal(t, e.Success(), e.Vote(producer1, []common.AccountName{producer1}, common.AccountName(0)))

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
		assert.True(t, 500 * 1000 > int64(float64(initialSupply.Amount) * float64(0.05) - float64(supply.Amount - initialSupply.Amount)))
		assert.True(t, 500 * 1000 > int64(float64(initialSupply.Amount) * float64(0.04) - float64(savings - initialSavings)))
	}

	e.close()
}

func TestMultipleProducerPay(t *testing.T) {

}

func TestProducersUpgradeSystemContract(t *testing.T) {

}

func TestProducerOnblockCheck(t *testing.T) {
	e := initEosioSystemTester()
	largeAsset := CoreFromString("80.0000")
	e.CreateAccountWithResources(producer1, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(producer2, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)
	e.CreateAccountWithResources(producer3, eosio, CoreFromString("1.0000"), false, largeAsset, largeAsset)

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

}