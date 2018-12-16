package template

import (
	"bytes"
	"encoding/json"
	"github.com/eosspark/eos-go/log"
	"reflect"
	"strconv"
)

// template type Exception(PARENT,CODE,WHAT)
const CODE = 0
const WHAT = ""
var ExceptionName = reflect.TypeOf(Exception{}).Name()

type PARENT interface {

}
type Exception struct {
	PARENT
	Elog []log.LogMessage
}

func (e Exception) Code() int64 {
	return CODE
}

func (e Exception) Name() string {
	return ExceptionName
}

func (e Exception) What() string {
	return WHAT
}

func (e *Exception) AppendLog(l log.LogMessage) {
	e.Elog = append(e.Elog, l)
}

func (e Exception) TopMessage() string {
	for _, log := range e.Elog {
		if msg := log.Message(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e Exception) String() string {
	return e.DetailMessage()
}

func (e Exception) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, log := range e.Elog {
		buffer.WriteString(log.Message())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e Exception) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: CODE,
		Name: ExceptionName,
		What: WHAT,
	}

	return json.Marshal(except)
}