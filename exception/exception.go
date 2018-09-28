package exception

import "fmt"

type ExcTypes int

const (
	ChainException = ExcTypes(3000000)
)

type Exception struct {
	code    ExcTypes
	message string
}

func EosAssert(expr bool, excType ExcTypes, format string, args ...interface{}) {
	if !expr {
		throwException(excType, format, args)
	}
}

func throwException(excType ExcTypes, format string, args ...interface{}) {
	panic(Exception{excType, fmt.Sprintf(format, args)})
}
