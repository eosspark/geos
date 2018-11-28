package rlp_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
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
	F16 common.Varuint32
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
		F16: common.Varuint32{999},
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
	assert.Equal(t, common.Varuint32{999}, s.F16)
	assert.Equal(t, true, s.F17)
	assert.Equal(t, int64(100000), s.F18.Amount)
	assert.Equal(t, uint8(4), s.F18.Precision)
	assert.Equal(t, "EOS", s.F18.Symbol.Symbol)

}
