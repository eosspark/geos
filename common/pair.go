package common

import (
	"github.com/eosspark/eos-go/rlp"
)

type Pair struct {
	First  interface{}
	Second interface{}
}

func MakePair(a, b interface{}) *Pair {
	return &Pair{First: a, Second: b}
}

func (p *Pair) GetIndex() []byte { //TODO
	first, _ := rlp.EncodeToBytes(p.First)
	out := first
	second, _ := rlp.EncodeToBytes(p.Second)
	out = append(out, second...)
	return out
}

type Tuple struct {
	First  interface{}
	Second interface{}
	Third  interface{}
}

func MakeTuple(a, b, c interface{}) *Tuple {
	return &Tuple{First: a, Second: b, Third: c}
}

func (p *Tuple) GetIndex() []byte { //TODO
	first, _ := rlp.EncodeToBytes(p.First)
	out := first
	second, _ := rlp.EncodeToBytes(p.Second)
	out = append(out, second...)
	t, _ := rlp.EncodeToBytes(p.Third)
	out = append(out, t...)
	return out
}

type Tuple4 struct {
	First  interface{}
	Second interface{}
	Third  interface{}
	Fourth interface{}
}

func MakeTuple4(a, b, c, d interface{}) *Tuple4 {
	return &Tuple4{First: a, Second: b, Third: c, Fourth: d}
}

func (p *Tuple4) GetIndex() []byte { //TODO
	first, _ := rlp.EncodeToBytes(p.First)
	out := first
	second, _ := rlp.EncodeToBytes(p.Second)
	out = append(out, second...)
	third, _ := rlp.EncodeToBytes(p.Third)
	out = append(out, third...)
	fourth, _ := rlp.EncodeToBytes(p.Fourth)
	out = append(out, fourth...)
	return out
}
