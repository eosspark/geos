package chain

import "fmt"

type Exception struct {
	Code        string
	Description string
	Types       string
}

func NewException(code string, desc string, types string) Exception {
	return Exception{
		Code:        code,
		Description: desc,
		Types:       types,
	}
}

func (ex *Exception) ToString() string {
	str := fmt.Sprintf("s%:d%:d%", ex.Types, ex.Code, ex.Description)
	return str
}
