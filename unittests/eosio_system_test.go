package unittests

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuySell(t *testing.T) {
	e := initEosioSystemTester()
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.Stake(eosio, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(uint64(alice))
	initBytes := uint64(total["ram_bytes"].(int64))
	initialRamBalance := e.GetBalance(eosioRam)
	initialRamFeeBalance := e.GetBalance(eosioRamFee)
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("200.0000")))
	assert.Equal(t, CoreFromString("800.0000"), e.GetBalance(alice))
	assert.Equal(t, initialRamBalance.Add(CoreFromString("199.0000")), e.GetBalance(eosioRam))
	assert.Equal(t, initialRamFeeBalance.Add(CoreFromString("1.0000")), e.GetBalance(eosioRamFee))

	total = e.GetTotalStake(uint64(alice))
	bytes := uint64(total["ram_bytes"].(int64))
	boughtBytes := bytes - initBytes

	assert.Equal(t, true, 0 < boughtBytes)
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("998.0049"), e.GetBalance(alice))
	total = e.GetTotalStake(uint64(alice))
	assert.Equal(t, initBytes, uint64(total["ram_bytes"].(int64)))

	e.Transfer(eosio, alice, CoreFromString("100000000.0000"), eosio)
	assert.Equal(t, CoreFromString("100000998.0049"), e.GetBalance(alice))
	assert.Equal(t, e.Success(), e.BuyRam(alice, alice, CoreFromString("10000000.0000")))
	assert.Equal(t, CoreFromString("90000998.0049"), e.GetBalance(alice))

	total = e.GetTotalStake(uint64(alice))
	bytes = uint64(total["ram_bytes"].(int64))
	boughtBytes = bytes - initBytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))

	total = e.GetTotalStake(uint64(alice))
	bytes = uint64(total["ram_bytes"].(int64))
	boughtBytes = bytes - initBytes
	assert.Equal(t, initBytes, uint64(total["ram_bytes"].(int64)))
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

	newTotal := e.GetTotalStake(uint64(alice))
	newBytes := uint64(newTotal["ram_bytes"].(int64))
	boughtBytes = newBytes - bytes
	assert.Equal(t, e.Success(), e.SellRam(alice, boughtBytes))
	assert.Equal(t, CoreFromString("99901242.4179"), e.GetBalance(alice))

	newTotal = e.GetTotalStake(uint64(alice))
	startBytes := uint64(total["ram_bytes"].(int64))

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

	finalTotal := e.GetTotalStake(uint64(alice))
	endBytes := uint64(finalTotal["ram_bytes"].(int64))
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

	total := e.GetTotalStake(uint64(alice))
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

	e.ProduceBlock(common.Hours(1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("1000.0000"), e.GetBalance(alice))
	assert.Equal(t, initEosioStakeBalance, e.GetBalance(eosioStake))

	assert.Equal(t, e.Success(), e.Stake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))
	total = e.GetTotalStake(uint64(bob))
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	//   REQUIRE_MATCHING_OBJECT( voter( "alice1111111", core_from_string("300.0000")), get_voter_info( "alice1111111" ) );
	bytes := total["ram_bytes"].(int64)
	assert.True(t, 0 < bytes)

	assert.Equal(t, e.Success(), e.UnStake(alice, bob, CoreFromString("200.0000"), CoreFromString("100.0000")))
	total = e.GetTotalStake(uint64(bob))
	assert.Equal(t, CoreFromString("10.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("10.0000"), total["cpu_weight"].(common.Asset))

	e.ProduceBlock(common.Hours(3*24-1), 0)
	e.ProduceBlocks(1, false)
	assert.Equal(t, CoreFromString("700.0000"), e.GetBalance(alice))

	e.ProduceBlock(common.Hours(1),0)
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

	e.Transfer(eosio, bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.StakeWithTransfer(bob, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

	total := e.GetTotalStake(uint64(alice))
	assert.Equal(t, CoreFromString("210.0000"), total["net_weight"].(common.Asset))
	assert.Equal(t, CoreFromString("110.0000"), total["cpu_weight"].(common.Asset))

	e.close()
}

func TestStakeToSelfWithTransfer(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))
	e.Transfer(eosio, alice, CoreFromString("1000.0000"), eosio)
	Try(func() {
		e.StakeWithTransfer(alice, alice, CoreFromString("200.0000"), CoreFromString("100.0000"))
	}).Catch(func(e Exception) {
		assert.True(t, inString(e.DetailMessage(), "cannot use transfer flag if delegating to self"))
	}).End()
	e.close()
}

func TestStakeWhilePendingRefund(t *testing.T) {
	e := initEosioSystemTester()
	e.Cross15PercentThreshold()
	e.Issue(eosio, CoreFromString("1000.0000"), eosio)
	e.Issue(eosioStake, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, CoreFromString("0.0000"), e.GetBalance(alice))

	e.Transfer(eosio, bob, CoreFromString("1000.0000"), eosio)
	assert.Equal(t, e.Success(), e.StakeWithTransfer(bob, alice, CoreFromString("200.0000"), CoreFromString("100.0000")))

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
