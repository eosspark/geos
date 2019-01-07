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

var BlockLogBackupDirExistName = reflect.TypeOf(BlockLogBackupDirExist{}).Name()

type BlockLogBackupDirExist struct {
	_BlockLogException
	Elog log.Messages
}

func NewBlockLogBackupDirExist(parent _BlockLogException, message log.Message) *BlockLogBackupDirExist {
	return &BlockLogBackupDirExist{parent, log.Messages{message}}
}

func (e BlockLogBackupDirExist) Code() int64 {
	return 3190004
}

func (e BlockLogBackupDirExist) Name() string {
	return BlockLogBackupDirExistName
}

func (e BlockLogBackupDirExist) What() string {
	return "block log backup dir already exists"
}

func (e *BlockLogBackupDirExist) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockLogBackupDirExist) GetLog() log.Messages {
	return e.Elog
}

func (e BlockLogBackupDirExist) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e BlockLogBackupDirExist) DetailMessage() string {
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

func (e BlockLogBackupDirExist) String() string {
	return e.DetailMessage()
}

func (e BlockLogBackupDirExist) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3190004,
		Name: BlockLogBackupDirExistName,
		What: "block log backup dir already exists",
	}

	return json.Marshal(except)
}

func (e BlockLogBackupDirExist) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockLogBackupDirExist):
		callback(&e)
		return true
	case func(BlockLogBackupDirExist):
		callback(e)
		return true
	default:
		return false
	}
}