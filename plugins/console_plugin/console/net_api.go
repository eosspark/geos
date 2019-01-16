package console

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/net_plugin"
	"github.com/robertkrimen/otto"
)

//NetAPI interacts with local p2p network connections
type NetAPI struct {
	c *Console
}

func newNetAPI(c *Console) *NetAPI {
	n := &NetAPI{
		c: c,
	}
	return n
}

//Connect starts a new connection to a peer
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

//Disconnect closes an existing connection
func (n *NetAPI) Disconnect(call otto.FunctionCall) (response otto.Value) {
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

//Status status of existing connection
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

//Peers status of exiting connection
func (n *NetAPI) Peers(call otto.FunctionCall) (response otto.Value) {
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
