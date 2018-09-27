package common

import (
	"github.com/eosspark/eos-go/rlp"
)

type Pair struct {
	First  interface{}
	Second interface{}
}

func (p *Pair) GetIndex() []byte {
	f, _ := rlp.EncodeToBytes(p.First)
	s, _ := rlp.EncodeToBytes(p.Second)
	f = append(f, s...)
	return f
}

type Tuple struct {
	First  interface{}
	Second interface{}
	Third  interface{}
}

func (p *Tuple) GetIndex() []byte {
	f, _ := rlp.EncodeToBytes(p.First)
	s, _ := rlp.EncodeToBytes(p.Second)
	t, _ := rlp.EncodeToBytes(p.Third)
	f = append(f, s...)
	f = append(f, t...)
	return f
}

//func (p *Tuple)Get(i int)interface{}{
//
//	switch i{
//	case 0:return p.First
//	case 1: return p.Second
//	case 2: return p.Third
//	default:
//		return nil
//	}
//}
