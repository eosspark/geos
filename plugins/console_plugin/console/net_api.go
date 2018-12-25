package console

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/net_plugin"
	"github.com/robertkrimen/otto"
)

type NetAPI struct {
	c   *Console
	log log.Logger
}

func newNetAPI(c *Console) *NetAPI {
	n := &NetAPI{
		c: c,
	}
	n.log = log.New("netAPI")
	n.log.SetHandler(log.TerminalHandler)
	return n
}

func (n *NetAPI) Connect(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var connectInfo string
	err = DoHttpCall(&connectInfo, common.NetConnect, host)
	if err != nil {
		fmt.Println(err)
		return otto.FalseValue()
	}

	return getJsResult(call, connectInfo)
}

func (n *NetAPI) DisConnect(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var result string
	err = DoHttpCall(&result, common.NetDisconnect, host)
	if err != nil {
		fmt.Println(err)
		return otto.FalseValue()
	}
	return getJsResult(call, result)
}

func (n *NetAPI) Status(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var result net_plugin.PeerStatus
	err = DoHttpCall(&result, common.NetStatus, host)
	if err != nil {
		fmt.Println(err)
		return otto.FalseValue()
	}
	return getJsResult(call, result)
}

func (n *NetAPI) Connections(call otto.FunctionCall) (response otto.Value) {
	host, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var result []net_plugin.PeerStatus
	err = DoHttpCall(&result, common.NetConnections, host)
	if err != nil {
		fmt.Println(err)
		return otto.FalseValue()
	}
	return getJsResult(call, result)
}
