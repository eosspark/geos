package unittests

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestBuySell(t *testing.T) {
	e := initEosioSystemTester()
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(eosio, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(alice)
	initBytes := uint64(total["ram_bytes"].(float64))
	initialRamBalance := e.GetBalance(eosioRam)
	initialRamFeeBalance := e.GetBalance(eosioRamFee)
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("200.0000")))
	assert.Equal(t, CoreFromString("800.0000"), e.GetBalance(alice))
	assert.Equal(t, initialRamBalance.Add(CoreFromString("199.0000")), e.GetBalance(eosioRam))
	assert.Equal(t, initialRamFeeBalance.Add(CoreFromString("1.0000")), e.GetBalance(eosioRamFee))

	total = e.GetTotalStake(alice)
	bytes := /*e.ToUint64(total["ram_bytes"])*/uint64(total["ram_bytes"].(float64))

	boughtBytes := bytes - initBytes

	assert.Equal(t, true, 0 < boughtBytes)
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("998.0049"), e.GetBalance(alice))
	total = e.GetTotalStake(alice)
	assert.Equal(t, initBytes, uint64(total["ram_bytes"].(float64)))

	e.Transfer(eosio, alice, CoreFromString("100000000.0000"), eosio)
	assert.Equal(t, CoreFromString("100000998.0049"), e.GetBalance(alice))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, CoreFromString("90000998.0049"), e.GetBalance(alice))

	total = e.GetTotalStake(alice)
	a, _ := strconv.Atoi(total["ram_bytes"].(string))
	bytes = uint64(a)
	boughtBytes = bytes - initBytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))

	total = e.GetTotalStake(alice)
	bytes = uint64(total["ram_bytes"].(float64))
	boughtBytes = bytes - initBytes
	assert.Equal(t, initBytes, uint64(total["ram_bytes"].(float64)))
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
	newBytes := uint64(newTotal["ram_bytes"].(float64))
	boughtBytes = newBytes - bytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("99901242.4179"), e.GetBalance(alice))

	newTotal = e.GetTotalStake(alice)
	startBytes := uint64(total["ram_bytes"].(float64))

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
	endBytes := uint64(finalTotal["ram_bytes"].(float64))
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
	bytes := total["ram_bytes"].(int64)
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
	assert.Equal(t, alice.String(), info["owner"].(string))
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
	assert.Equal(t, alice.String(), info["owner"].(string))
	assert.Equal(t, key.String(), info["producer_key"].(string))
	assert.Equal(t, "http://block.two", info["url"].(string))
	assert.Equal(t, float64(1), info["location"].(float64))

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
	assert.Equal(t, alice.String(), info["owner"].(string))
	assert.Equal(t, key2.String(), info["producer_key"].(string))
	assert.Equal(t, "http://block.two", info["url"].(string))
	assert.Equal(t, float64(2), info["location"].(float64))

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
	assert.Equal(t, ecc.NewPublicKeyNil().String(), info["producer_key"].(string))

	//everything else should stay the same
	assert.Equal(t, alice.String(), info["owner"].(string))
	assert.Equal(t, "http://block.two", info["url"].(string))
	assert.Equal(t, float64(2), info["location"].(float64))

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

func TestAccountName(t *testing.T) {
	//a :=common.AccountName(uint64(1))
	//fmt.Println(a)
	//fmt.Printf("%d\n",common.N("............1"))

	a := []byte{16, 66, 8, 87, 33, 157, 232, 173, 0, 0, 0, 0, 0, 0, 0, 0, 2, 16, 66, 8, 87, 33, 157, 232, 173}
	type VoteProducer struct {
		Voter     common.AccountName
		Proxy     common.AccountName
		Producers []common.AccountName
	}
	var Vote VoteProducer
	err := rlp.DecodeBytes(a, &Vote)
	fmt.Println(err)
	fmt.Printf("%v", Vote)
}
