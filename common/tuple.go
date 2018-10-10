package common

type Tuple []interface{}

func MakeTuple(in ...interface{}) (out Tuple) {
	for _, content := range in {
		out = append(out, content)
	}
	return
}

// used for pair encode
type Pair struct {
	First interface{}
	Second interface{}
}

func MakePair(a interface{}, b interface{}) Pair {
	return Pair{a, b}
}
