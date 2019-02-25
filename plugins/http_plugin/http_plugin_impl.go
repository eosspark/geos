package http_plugin

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"net"
	"net/http"
)

type NextFunction = func(interface{})
type UrlResponseCallback = func(int, []byte)
type UrlHandler = func(source string, body []byte, cb UrlResponseCallback)

type HttpPluginImpl struct {
	UrlHandlers map[string]UrlHandler

	AccessControlAllowOrigin      string
	AccessControlAllowHeaders     string
	AccessControlMaxAge           string
	AccessControlAllowCredentials bool //default false
	MaxBodySize                   common.SizeT
	httpsCeryChain                string
	httpsKey                      string

	listenStr           string
	ListenEndpoint      *http.ServeMux
	HttpsListenEndpoint *http.ServeMux

	Listener net.Listener
	//Handlers *rpc.Server

	self *HttpPlugin
}

func NewHttpPluginImpl(io *asio.IoContext) *HttpPluginImpl {
	impl := new(HttpPluginImpl)
	impl.UrlHandlers = make(map[string]UrlHandler)
	impl.AccessControlAllowCredentials = false
	return impl
}
