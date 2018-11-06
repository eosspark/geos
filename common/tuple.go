package common

type Tuple []interface{}

func MakeTuple(in ...interface{}) Tuple {
	out := make([]interface{}, 0, len(in)) //alloc capacity
	out = append(out, in...)
	return out
}

// used for pair encode
type Pair struct {
	First interface{}
	Second interface{}
}

func MakePair(a interface{}, b interface{}) Pair {
	return Pair{a, b}
}
