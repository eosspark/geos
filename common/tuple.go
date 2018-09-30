package common

type Tuple []interface{}

func MakeTuple(in ...interface{}) (out Tuple) {
	for _, content := range in {
		out = append(out, content)
	}
	return
}
