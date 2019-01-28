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

type PARENT interface{}

type Exception struct {
	PARENT
	Elog log.Messages
}

func New(parent PARENT, message log.Message) *Exception {
	return &Exception{parent, log.Messages{message}}
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

func (e *Exception) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e Exception) GetLog() log.Messages {
	return e.Elog
}

func (e Exception) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e Exception) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e Exception) String() string {
	return e.DetailMessage()
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

func (e Exception) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*Exception):
		callback(&e)
		return true
	case func(Exception):
		callback(e)
		return true
	default:
		return false
	}
}
