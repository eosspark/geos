package rlp_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	. "github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoderName(t *testing.T) {
	a := common.Name(common.N("eosio"))
	name, _ := rlp.EncodeToBytes(a)
	check := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55}
	assert.Equal(t, check, name)

	b := common.AccountName(common.N("eosio"))
	name, _ = rlp.EncodeToBytes(b)
	assert.Equal(t, check, name)
	c := common.PermissionName(common.N("eosio"))
	name, _ = rlp.EncodeToBytes(c)

	assert.Equal(t, check, name)
	d := common.ActionName(common.N("eosio"))
	name, _ = rlp.EncodeToBytes(d)
	assert.Equal(t, check, name)
	e := common.TableName(common.N("eosio"))
	name, _ = rlp.EncodeToBytes(e)
	assert.Equal(t, check, name)

	f := common.ScopeName(common.N("eosio"))
	name, _ = rlp.EncodeToBytes(f)
	assert.Equal(t, check, name)

}

func TestAsset(t *testing.T) {
	sys := common.Symbol{
		Precision: 5,
		Symbol:    "BTC",
	}
	asset := common.Asset{
		Amount: 10000,
		Symbol: sys,
	}

	re, err := rlp.EncodeToBytes(asset)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(re)

	var btc common.Asset
	err = rlp.DecodeBytes(re, &btc)
	result := common.Asset{Amount: 10000, Symbol: common.Symbol{Precision: 0x5, Symbol: "BTC"}}
	assert.Equal(t, result, btc)
	//fmt.Printf("%#v\n",btc)
}

func TestDecoder_Checksum256(t *testing.T) {
	s := crypto.NewSha256Byte(bytes.Repeat([]byte{1}, 32))

	buf, err := rlp.EncodeToBytes(s)
	assert.NoError(t, err)
	var re crypto.Sha256
	err = rlp.DecodeBytes(buf, &re)
	assert.NoError(t, err)

	assert.Equal(t, s, &re)
	//assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Empty_Checksum256(t *testing.T) {
	s := crypto.NewSha256Nil()
	buf, err := rlp.EncodeToBytes(s)
	assert.NoError(t, err)

	var re crypto.Sha256
	err = rlp.DecodeBytes(buf, &re)
	assert.NoError(t, err)
	assert.Equal(t, s, &re)
}

func TestDecoder_PublicKey(t *testing.T) {
	pk := ecc.MustNewPublicKey("EOS1111111111111111111111111111111114T1Anm")
	buf, err := rlp.EncodeToBytes(pk)
	assert.NoError(t, err)
	pub := ecc.NewPublicKeyNil()
	err = rlp.DecodeBytes(buf, pub)
	assert.NoError(t, err)
	assert.Equal(t, pk, *pub)

	pk = ecc.MustNewPublicKey("EOS7AfdKvHbtaQzAiDZ54n8cJ6FBzMAdXU3TFVDXRJGWHU16DPoyp")
	buf, err = rlp.EncodeToBytes(pk)
	assert.NoError(t, err)
	pub = ecc.NewPublicKeyNil()
	err = rlp.DecodeBytes(buf, pub)
	assert.NoError(t, err)

	pubre := ecc.MustNewPublicKeyFromData(buf)

	assert.Equal(t, pk.String(), pubre.String())
	assert.Equal(t, pk, *pub)
}

func TestDecoder_Signature(t *testing.T) {
	sig := ecc.MustNewSignatureFromData(bytes.Repeat([]byte{0}, 66))

	buf, err := rlp.EncodeToBytes(sig)
	assert.NoError(t, err)
	rsig := ecc.NewSigNil()
	err = rlp.DecodeBytes(buf, rsig)
	assert.NoError(t, err)
	assert.Equal(t, sig, *rsig)

	sig, err = ecc.NewSignature("SIG_K1_JujD18YwovAZp3nUEsQod8x9HXGqdbU4gptPWqSjgnPjNMGqoP5aNa4aZPjLZKtwBPBkXpBfVwcGopCjGzEGMJxGTLhfQ3")
	buf, err = rlp.EncodeToBytes(sig)
	assert.NoError(t, err)

	var re ecc.Signature
	err = rlp.DecodeBytes(buf, &re)
	assert.NoError(t, err)
	assert.Equal(t, sig, re)
	//fmt.Println(re)

	var signil ecc.Signature
	buf, err = rlp.EncodeToBytes(signil)
	assert.NoError(t, err)

}

func TestDecoder_BlockTimestamp(t *testing.T) {
	ts := types.NewBlockTimeStamp(common.TimePoint(0))
	buf, err := rlp.EncodeToBytes(ts)
	assert.NoError(t, err)
	var rts types.BlockTimeStamp
	err = rlp.DecodeBytes(buf, &rts)
	assert.NoError(t, err)
	assert.Equal(t, ts, rts)
}

type EncodeTestStruct struct {
	F1 string
	F2 int16
	F3 uint16
	F4 uint32
	F5 *crypto.Sha256
	F6 []string
	F7 [2]string
	//	F8  map[string]string
	F9  ecc.PublicKey
	F10 ecc.Signature
	F11 byte
	F12 uint64
	F13 []byte
	//F14 common.TimePoint
	F15 types.BlockTimeStamp
	//F16 common.Varuint32
	F17 bool
	F18 common.Asset
}

func TestDecoder_Encode(t *testing.T) {

	//func Now() TimePoint          { return TimePoint(time.Now().UTC().UnixNano() / 1e3) }
	//now := time.Date(2018, time.September, 26, 1, 2, 3, 4, time.UTC)
	//tstamp := common.TimePoint(now.UTC().UnixNano() / 1e3)

	blockts := types.BlockTimeStamp(common.TimePoint(0))
	s := &EncodeTestStruct{
		F1: "abc",
		F2: -75,
		F3: 99,
		F4: 999,
		F5: crypto.NewSha256Byte(bytes.Repeat([]byte{0}, 32)),
		F6: []string{"def", "789"},
		F7: [2]string{"foo", "bar"},
		// maps don't serialize deterministically.. we no want that.
		//		F8:  map[string]string{"foo": "bar", "hello": "you"},
		F9:  ecc.MustNewPublicKey("EOS1111111111111111111111111111111114T1Anm"),
		F10: *ecc.NewSigNil(),
		F11: byte(1),
		F12: uint64(87),
		F13: []byte{1, 2, 3, 4, 5},
		//F14: tstamp,
		F15: blockts,
		//F16: common.Varuint32{999},
		F17: true,
		F18: common.NewEOSAsset(100000),
	}

	buf, err := rlp.EncodeToBytes(s)
	assert.NoError(t, err)
	assert.Equal(t, "03616263b5ff6300e7030000000000000000000000000000000000000000000000000000000000000000000002036465660337383903666f6f036261720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001570000000000000005010203040500000000e70701a08601000000000004454f5300000000", hex.EncodeToString(buf))
	var re EncodeTestStruct
	err = rlp.DecodeBytes(buf, &re)
	assert.NoError(t, err)

	assert.Equal(t, "abc", s.F1)
	assert.Equal(t, int16(-75), s.F2)
	assert.Equal(t, uint16(99), s.F3)
	assert.Equal(t, uint32(999), s.F4)
	assert.Equal(t, crypto.NewSha256Nil(), s.F5)
	assert.Equal(t, []string{"def", "789"}, s.F6)
	assert.Equal(t, [2]string{"foo", "bar"}, s.F7)
	//	assert.Equal(t, map[string]string{"foo": "bar", "hello": "you"}, s.F8)
	assert.Equal(t, ecc.MustNewPublicKeyFromData(bytes.Repeat([]byte{0}, 34)), s.F9)
	assert.Equal(t, ecc.MustNewSignatureFromData(bytes.Repeat([]byte{0}, 66)), s.F10)
	assert.Equal(t, byte(1), s.F11)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, s.F13)
	//assert.Equal(t, tstamp, s.F14)
	assert.Equal(t, blockts, s.F15)
	//assert.Equal(t, common.Varuint32{999}, s.F16)
	assert.Equal(t, true, s.F17)
	assert.Equal(t, int64(100000), s.F18.Amount)
	assert.Equal(t, uint8(4), s.F18.Precision)
	assert.Equal(t, "EOS", s.F18.Symbol.Symbol)

}

type TreeSetExp struct {
	ActorWhitelist AccountNameSet //common.AccountName
}

func TestTreeSet(t *testing.T) {
	vals := []common.AccountName{common.AccountName(common.N("eos")), common.AccountName(common.N("io"))}
	A := TreeSetExp{}

	A.ActorWhitelist = *NewAccountNameSet()
	for _, val := range vals {
		A.ActorWhitelist.AddItem(val)
	}
	bytes, err := rlp.EncodeToBytes(A)
	fmt.Println(bytes, err)

	B := TreeSetExp{}
	B.ActorWhitelist = *NewAccountNameSet()
	err = rlp.DecodeBytes(bytes, &B)

	fmt.Println(B.ActorWhitelist.Values(), err)

}

type TreeSetPub struct {
	WhiteList PublicKeySet
}

func TestPub(t *testing.T) {

	vals := []ecc.PublicKey{ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"), ecc.MustNewPublicKey("EOS5kpVjpFXiFHwhbrSLndAqCdpLLUctXhq583WjFH5tqy2VLYhLc")}
	A := TreeSetPub{}

	A.WhiteList = *NewPublicKeySet()
	for _, val := range vals {
		A.WhiteList.AddItem(val)
	}
	bytes, err := rlp.EncodeToBytes(A)
	assert.NoError(t, err)

	B := TreeSetPub{}
	B.WhiteList = *NewPublicKeySet()
	err = rlp.DecodeBytes(bytes, &B)
	assert.NoError(t, err)
	fmt.Println(B.WhiteList.Values(), err)
}

func BenchmarkTreeSetInsert(b *testing.B) {
	b.StopTimer()
	set := NewAccountNameSet()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		set.Clear()
		b.StartTimer()
		for i := 0; i < 100000; i++ {
			set.AddItem(common.Name(i))
		}
	}

}

func BenchmarkTreeSetFromDecode(b *testing.B) {
	b.StopTimer()
	set := NewAccountNameSet()
	for i := 0; i < 100000; i++ {
		set.AddItem(common.Name(i))
	}

	bytes, err := rlp.EncodeToBytes(*set)
	//fmt.Println(bytes)
	if err != nil {
		b.Fatal(err)
	}

	ss := NewAccountNameSet()
	if err = rlp.DecodeBytes(bytes, ss); err != nil {
		b.Fatal(err)
	}

	ss.Clear()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		ss.Clear()
		b.StartTimer()
		rlp.DecodeBytes(bytes, ss)
	}

}

const bench = 10000

var pks [bench]ecc.PublicKey

func init() {
	for i := 0; i < bench; i++ {
		pri, _ := ecc.NewRandomPrivateKey()
		pks[i] = pri.PublicKey()
	}

	fmt.Println("init success", pks[bench-1])
}

func BenchmarkTreeSetInsert2(b *testing.B) {
	b.StopTimer()

	set := NewPublicKeySet()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		set.Clear()
		b.StartTimer()
		for i := 0; i < bench; i++ {
			set.AddItem(pks[i])
		}
	}

}

func BenchmarkTreeSetFromDecode2(b *testing.B) {
	b.StopTimer()
	set := NewPublicKeySet()
	for i := 0; i < bench; i++ {
		set.AddItem(pks[i])
	}

	bytes, err := rlp.EncodeToBytes(*set)
	//fmt.Println(bytes)
	if err != nil {
		b.Fatal(err)
	}

	ss := NewPublicKeySet()
	if err = rlp.DecodeBytes(bytes, ss); err != nil {
		b.Fatal(err)
	}

	ss.Clear()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		ss.Clear()
		b.StartTimer()
		rlp.DecodeBytes(bytes, ss)
	}

}

type TreeMapExp struct {
	ProducerToLastProduced AccountNameUint32Map
}

func TestTreeMap(t *testing.T) {
	keyvalue := map[common.Name]uint32{
		common.N("eos"):   100,
		common.N("sys"):   101,
		common.N("hello"): 90,
	}
	tree := TreeMapExp{}
	tree.ProducerToLastProduced = *NewAccountNameUint32Map()
	for k, v := range keyvalue {
		tree.ProducerToLastProduced.Put(k, v)
	}
	tree.ProducerToLastProduced.Each(func(key common.AccountName, value uint32) {

	})
	bytes, err := rlp.EncodeToBytes(tree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bytes)

	tree2 := TreeMapExp{}
	tree2.ProducerToLastProduced = *NewAccountNameUint32Map()

	err = rlp.DecodeBytes(bytes, &tree2)
	assert.NoError(t, err, err)
	tree2.ProducerToLastProduced.Each(func(key common.AccountName, value uint32) {

	})
	json1, _ := tree.ProducerToLastProduced.ToJSON()
	json2, _ := tree2.ProducerToLastProduced.ToJSON()
	assert.Equal(t, json1, json2)
}

var AccountNameUint32MapBytes = []byte{90, 91, 123, 34, 107, 101, 121, 34, 58, 34, 101, 111, 115, 105, 111, 34, 44, 34, 118, 97, 108, 34, 58, 49, 48, 48, 125, 44, 123, 34, 107, 101, 121, 34, 58, 34, 101, 111, 115, 105, 111, 46, 108, 108, 108, 108, 108, 108, 34, 44, 34, 118, 97, 108, 34, 58, 57, 57, 125, 44, 123, 34, 107, 101, 121, 34, 58, 34, 101, 111, 115, 105, 111, 46, 116, 111, 107, 101, 110, 34, 44, 34, 118, 97, 108, 34, 58, 57, 57, 125, 93}

func TestAccountNameUint32Map(t *testing.T) {
	type BlockHead struct {
		NewProducerToLastProduced *AccountNameUint32Map
	}
	var a BlockHead

	a.NewProducerToLastProduced = NewAccountNameUint32Map()
	a.NewProducerToLastProduced.Put(common.N("eosio"), 100)
	a.NewProducerToLastProduced.Put(common.N("eosio.token"), 99)
	a.NewProducerToLastProduced.Put(common.N("eosio.llllll"), 99)
	jsonree, _ := a.NewProducerToLastProduced.ToJSON()
	fmt.Println("*****:  ", jsonree)
	bytes, err := rlp.EncodeToBytes(a)
	fmt.Println(bytes, err)
	//
	//fmt.Printf("%v", bytes)
	//var str string
	//for i := 0; i < len(bytes); i++ {
	//	str1 := fmt.Sprintf("%d,", bytes[i])
	//	str = str + str1
	//	fmt.Println(str)
	//}

	//fmt.Println(str)
	//var test types.AccountNameUint32Map
	//test = *types.NewAccountNameUint32Map()
	//err = rlp.DecodeBytes(bytes,&test)
	//fmt.Println("end:",err,test)

	var test BlockHead
	test.NewProducerToLastProduced = NewAccountNameUint32Map()
	err = rlp.DecodeBytes(bytes, &test)
	fmt.Println("end:", err, test)

}
