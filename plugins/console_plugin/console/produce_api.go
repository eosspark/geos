package console

import (
	"github.com/eosspark/eos-go/common"
	"github.com/robertkrimen/otto"
)

type ProduceAPI struct {
	c *Console
}

func newProduceAPI(c *Console) *ProduceAPI {
	return &ProduceAPI{c: c}
}

func (n *ProduceAPI) Pause(call otto.FunctionCall) (response otto.Value) {
	err := DoHttpCall(nil, common.ProducerPause, nil)
	if err != nil {
		clog.Error("SetWhitelistBlacklist is error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.UndefinedValue()
}

func (n *ProduceAPI) Resume(call otto.FunctionCall) (response otto.Value) {
	err := DoHttpCall(nil, common.ProducerResume, nil)
	if err != nil {
		clog.Error("SetWhitelistBlacklist is error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.UndefinedValue()
}

func (n *ProduceAPI) Paused(call otto.FunctionCall) (response otto.Value) {
	err := DoHttpCall(nil, common.ProducerPaused, nil)
	if err != nil {
		clog.Error("SetWhitelistBlacklist is error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.UndefinedValue()
}

type SetWhitelistBlacklistParams struct {
	ActorWhitelist    []common.AccountName `json:"actorWhitelist"`
	ActorBlacklist    []common.AccountName `json:"actorBlacklist"`
	ContractWhitelist []common.AccountName `json:"contractWhitelist"`
	ContractBlacklist []common.AccountName `json:"contractBlacklist"`
	ActionBlacklist   []common.NamePair    `json:"actionBlacklist"`
	KeyBlacklist      []string             `json:"keyBlacklist"`
}

func (n *ProduceAPI) SetWhitelistBlacklist(call otto.FunctionCall) (response otto.Value) {
	var params SetWhitelistBlacklistParams
	readParams(&params, call)
	clog.Info("params : %v", params)

	err := DoHttpCall(nil, common.ProducerSetWhitelistBlacklist, params)
	if err != nil {
		clog.Error("SetWhitelistBlacklist is error: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.UndefinedValue()
}
