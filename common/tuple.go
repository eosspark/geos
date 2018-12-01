package common

import (
	"bytes"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/log"
)

type Tuple []interface{}

func MakeTuple(in ...interface{}) Tuple {
	out := make([]interface{}, 0, len(in)) //alloc capacity
	out = append(out, in...)
	return out
}

// used for pair encode
type Pair struct {
	First  interface{}
	Second interface{}
}

func MakePair(a interface{}, b interface{}) Pair {
	return Pair{a, b}
}

func (p *Pair) GetKey() []byte {
	byt, err := database.EncodeToBytes(p)
	if err != nil {
		log.Error("Pair GetKey is error :%s", err.Error())
		return nil
	}
	return byt
}

func ComparePair(first interface{}, second interface{}) int {
	return bytes.Compare(first.(*Pair).GetKey(), second.(*Pair).GetKey())
}
