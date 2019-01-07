// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var ProducerPrivKeyNotFoundName = reflect.TypeOf(ProducerPrivKeyNotFound{}).Name()

type ProducerPrivKeyNotFound struct {
	_ProducerException
	Elog log.Messages
}

func NewProducerPrivKeyNotFound(parent _ProducerException, message log.Message) *ProducerPrivKeyNotFound {
	return &ProducerPrivKeyNotFound{parent, log.Messages{message}}
}

func (e ProducerPrivKeyNotFound) Code() int64 {
	return 3170001
}

func (e ProducerPrivKeyNotFound) Name() string {
	return ProducerPrivKeyNotFoundName
}

func (e ProducerPrivKeyNotFound) What() string {
	return "Producer private key is not available"
}

func (e *ProducerPrivKeyNotFound) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ProducerPrivKeyNotFound) GetLog() log.Messages {
	return e.Elog
}

func (e ProducerPrivKeyNotFound) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ProducerPrivKeyNotFound) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e ProducerPrivKeyNotFound) String() string {
	return e.DetailMessage()
}

func (e ProducerPrivKeyNotFound) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3170001,
		Name: ProducerPrivKeyNotFoundName,
		What: "Producer private key is not available",
	}

	return json.Marshal(except)
}

func (e ProducerPrivKeyNotFound) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ProducerPrivKeyNotFound):
		callback(&e)
		return true
	case func(ProducerPrivKeyNotFound):
		callback(e)
		return true
	default:
		return false
	}
}