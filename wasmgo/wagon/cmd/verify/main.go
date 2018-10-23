package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"

	//"github.com/eosspark/eos-go/rlp"
	//"github.com/stretchr/testify/assert"
	"reflect"
)

type SecondaryKeyInterface interface {
	SetValue(value interface{})
	GetValue() interface{}
}

type Uint64_t struct {
	Value uint64
}

func (u *Uint64_t) SetValue(value interface{}) {
	u.Value = reflect.ValueOf(value).Uint()
	u.Value += 100
}
func (u *Uint64_t) GetValue() interface{} {
	return u.Value
}

func main() {

	// sk := Uint64_t{100}
	// sk.SetValue(uint64(500))

	// fmt.Println(reflect.ValueOf(sk.GetValue()).Uint())

	//p := &int64{0}
	//q := &int(p)

	// p := &Uint64_t{99}
	// q := unsafe.Pointer (p)
	// k := (*byte) (q)
	// fmt.Println(*k)

	// a := []byte{1, 2, 3, 4, 5}
	// b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	// fmt.Println(unsafe.Pointer (&a))
	// fmt.Println(unsafe.Pointer (&b))

	// b = a
	// fmt.Println(b)
	// fmt.Println(unsafe.Pointer (&a))
	// fmt.Println(unsafe.Pointer (&b))

	//wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	//privKey, err := ecc.NewPrivateKey(wif)
	//
	//chainID, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	//payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000")
	//digest := sigDigest(chainID, payload)
	//
	//sig, err := privKey.Sign(digest)
	//
	////fromEOSIOC := "SIG_K1_K2WBNtiTY8o4mqFSz7HPnjkiT9JhUYGFa81RrzaXr3aWRF1F8qwVfutJXroqiL35ZiHTcvn8gPWGYJDwnKZTCcbAGJy4i9"
	////assert.Equal(t, fromEOSIOC, sig.String())
	//
	//pubKey, err := sig.PublicKey(digest)
	//
	//fmt.Println(err,pubKey.String())
	////assert.Equal(t, "PUB_K1_5DguRMaGh72NvbVX5LKHTb5cvbRmAxgrm9i2NNPKv5TC7FadXs", pubKey.String())

	// b := make([]bool, 10)
	// if !b[0] {
	// 	fmt.Println("b0 is false")
	// }

	actionTrace := types.ActionTrace{}
	var actionTraces []*types.ActionTrace

	actionTraces = append(actionTraces, &actionTrace)
	updateTrace(&actionTrace)
	updateTrace(actionTraces[len(actionTraces)-1])

	fmt.Println(actionTrace.BlockNum)
	fmt.Println(actionTraces[len(actionTraces)-1].BlockNum)

}

func updateTrace(actionTrace *types.ActionTrace) {
	actionTrace.BlockNum = 100
	return
}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
}
