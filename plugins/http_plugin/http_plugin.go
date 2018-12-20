package http_plugin

import (
	"encoding/json"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/http_plugin/fasthttp"
	"github.com/urfave/cli"
)

const (
	httpListenendpoint string = "127.0.0.1:8888"
)

var verboseHttpErrors = false

const HttpPlug = PluginTypeName("HttpPlugin")

var httpPlugin Plugin = App().RegisterPlugin(HttpPlug, NewHttpPlugin(App().GetIoService()))

type HttpPlugin struct {
	AbstractPlugin
	//ConfirmedBlock Signal //TODO signal ConfirmedBlock
	my *HttpPluginImpl
}

var hlog log.Logger

func NewHttpPlugin(io *asio.IoContext) *HttpPlugin {
	h := &HttpPlugin{}
	h.my = NewHttpPluginImpl(io)
	h.my.self = h

	hlog = log.New("http")
	hlog.SetHandler(log.TerminalHandler)
	return h
}

func (h *HttpPlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.StringFlag{
			Name:  "http-server-address",
			Usage: "The local IP and port to listen for incoming http connections; set blank to disable.",
			Value: httpListenendpoint,
		},
		cli.StringFlag{
			Name:  "https-server-address",
			Usage: "The local IP and port to listen for incoming https connections; leave blank to disable.",
		},
		cli.StringFlag{
			Name:  "https-certificate-chain-file",
			Usage: "FilenName with the certificate chain to present on https connections,PEM format. Required for https.",
		},
		cli.StringFlag{
			Name:  "https-private-key-file",
			Usage: "FilenName with https private key in PEM format. Required for https.",
		},
		cli.StringFlag{
			Name:  "access-control-allow-origin",
			Usage: "Specify the Access-control-Allow-Origin to be returned on each request.",
		},
		cli.StringFlag{
			Name:  "access-control-allow-headers",
			Usage: "Specify the Access-Control-Allow-Headers to be returned on each request.",
		},
		cli.StringFlag{
			Name:  "access-control-max-age",
			Usage: "Specify the Access-Control-Max-Age to be returned on each request.",
		},
		cli.BoolFlag{
			Name:  "access-control-allow-credentials",
			Usage: "Specify if Access-Control-Allow-Credentials: true should be returned on each request.",
		},
		cli.UintFlag{
			Name:  "max-body-size",
			Usage: "The maximum body size in bytes allowed for incoming RPC requests",
			Value: 1024 * 1024,
		},
		cli.BoolFlag{
			Name:  "verbose-http-errors",
			Usage: "Append the error log to HTTP responses",
		},
		cli.BoolFlag{ // default true
			Name:  "http-validate-host",
			Usage: "If set to false,then any incoming \"Host\" header is considered valid",
		},
		cli.StringSliceFlag{
			Name:  "http-alias",
			Usage: "Additionaly acceptable values for the \"Host\" header of incoming HTTP requests,can be specified multiple times. Include http/s_server_address by default.",
		},
	)
}

func (h *HttpPlugin) PluginInitialize(c *cli.Context) {

	Try(func() {
		h.my.listenStr = c.String("http-server-address")

		//handle := rpc.NewServer()
		//apis := plugins()
		//for _, api := range apis {
		//	if err := handle.RegisterName(api.Namespace, api.Service); err != nil {
		//		h.my.log.Error(err.Error())
		//		panic(err)
		//	}
		//	h.my.log.Debug("InProc registered :  namespace =%s", api.Namespace)
		//}
		//
		//h.my.Handlers = handle
		//// All APIs registered, start the HTTP listener
		//listener, err := net.Listen("tcp", listenStr)
		//if err != nil {
		//	h.my.log.Error("%s", err)
		//}
		//listener = netutil.LimitListener(listener,1)
		//
		//h.my.Listener = listener
		//h.my.log.Info("configured http to listen on %s", listenStr)

		//err := http.ListenAndServe(listenStr, h.my.ListenEndpoint)
		//if err != nil {
		//	h.my.log.Error("failed to configure https to listen on %s , %s", listenStr, err.Error())
		//	panic(err)
		//}
		//h.my.log.Info("configured http to listen on %s", h.my.listenStr)
	}).End()
	//Try(func() {
	//	h.my.log.Info("http plugin initialize")
	//	h.my.AccessControlAllowOrigin = c.String("access-control-allow-origin")
	//
	//	if c.IsSet("access-control-allow-origin") {
	//		h.my.AccessControlAllowOrigin = c.String("access-control-allow-origin")
	//		h.my.log.Info("configured http with Access-Control-Allow-Origin : %s", h.my.AccessControlAllowOrigin)
	//	}
	//	if c.IsSet("access-control-allow-headers") {
	//		h.my.AccessControlAllowHeaders = c.String("access-control-allow-headers")
	//		h.my.log.Info("configured http with Access-Control-Allow-Headers : %s", h.my.AccessControlAllowHeaders)
	//	}
	//	if c.IsSet("access-control-max-age") {
	//		h.my.AccessControlMaxAge = c.String("access-control-max-age")
	//		h.my.log.Info("configured http with Access-Control-Max-Age : %s", h.my.AccessControlMaxAge)
	//	}
	//
	//	h.my.AccessControlAllowCredentials = c.Bool("access-control-allow-credentials") //TODO
	//	if h.my.AccessControlAllowCredentials {
	//		h.my.log.Info("configured http with Access-Control-Allow-Credentials: true")
	//	}
	//
	//	//if c.IsSet("http-server-address") {
	//		listenStr := c.String("http-server-address")
	//		h.my.ListenEndpoint = http.NewServeMux()
	//		// httpPlugin.Handle(walletSetTimeOutFunc, walletPlugin.SetTimeOut())
	//	h.my.log.Info("configured http to listen on %s", "")
	//		err := http.ListenAndServe(listenStr, h.my.ListenEndpoint)
	//		if err != nil {
	//			h.my.log.Error("failed to configure https to listen on %s , %s", listenStr, err.Error())
	//			panic(err)
	//		}
	//		h.my.log.Info("configured http to listen on %s", h.my.listenStr)
	//	//}
	//
	//
	//	//if c.IsSet("https-server-address") {
	//	//	//	if( !options.count( "https-certificate-chain-file" ) ||
	//	//	//		options.at( "https-certificate-chain-file" ).as<string>().empty()) {
	//	//	//		elog( "https-certificate-chain-file is required for HTTPS" );
	//	//	//		return;
	//	//	//	}
	//	//	//	if( !options.count( "https-private-key-file" ) ||
	//	//	//		options.at( "https-private-key-file" ).as<string>().empty()) {
	//	//	//		elog( "https-private-key-file is required for HTTPS" );
	//	//	//		return;
	//	//	//	}
	//	//	//
	//	//	//	string lipstr = options.at( my->https_server_address_option_name ).as<string>();
	//	//	//	string host = lipstr.substr( 0, lipstr.find( ':' ));
	//	//	//	string port = lipstr.substr( host.size() + 1, lipstr.size());
	//	//	//tcp::resolver::query query( tcp::v4(), host.c_str(), port.c_str());
	//	//	//	try {
	//	//	//		my->https_listen_endpoint = *resolver.resolve( query );
	//	//	//		ilog( "configured https to listen on ${h}:${p} (TLS configuration will be validated momentarily)",
	//	//	//	("h", host)( "p", port ));
	//	//	//		my->https_cert_chain = options.at( "https-certificate-chain-file" ).as<string>();
	//	//	//		my->https_key = options.at( "https-private-key-file" ).as<string>();
	//	//	//	} catch ( const boost::system::system_error& ec ) {
	//	//	//	elog( "failed to configure https to listen on ${h}:${p} (${m})",
	//	//	//	("h", host)( "p", port )( "m", ec.what()));
	//	//	//	}
	//	//	//
	//	//	//	// add in resolved hosts and ports as well
	//	//	//	if (my->https_listen_endpoint) {
	//	//	//		my->add_aliases_for_endpoint(*my->https_listen_endpoint, host, port);
	//	//	//	}\
	//	//
	//	//	if !c.IsSet("https-certificate-chain-file") || len(c.String("https-certificate-chain-file")) == 0 {
	//	//		h.my.log.Error("https-certificate-chain-file is required for HTTPS")
	//	//		return
	//	//	}
	//	//
	//	//	if !c.IsSet("https-private-key-file") || len(c.String("https-private-key-file")) == 0 {
	//	//		h.my.log.Error("https-private-key-file is required for HTTPS")
	//	//		return
	//	//	}
	//	//
	//	//	lipStr := c.String("https-server-address")
	//	//	h.my.HttpsListenEndpoint = http.NewServeMux() //TODO https need to emplace
	//	//	err := http.ListenAndServe(lipStr, h.my.HttpsListenEndpoint)
	//	//	if err != nil {
	//	//		h.my.log.Error("failed to configure https to listen on %s , %s", lipStr, err.Error())
	//	//		panic(err)
	//	//	}
	//	//	h.my.log.Info("configured https to listen on %s (TLS configuration will be validated momentarily)", lipStr)
	//	//	h.my.httpsCeryChain = c.String("https-certificate-chain-file")
	//	//	h.my.httpsKey = c.String("https-private-key-file")
	//	//
	//	//}
	//
	//	h.my.MaxBodySize = common.SizeT(c.Uint64("max-body-size"))
	//	verboseHttpErrors = c.Bool("verbose-http-errors")
	//
	//}).FcLogAndRethrow().End()
}

func (h *HttpPlugin) PluginStartup() {
	h.my.log.Info("http plugin startup")

	handler := h.Handler

	fasthttp.ListenAndAsyncServe(App().GetIoService(), h.my.listenStr, handler)

	//if(my->listen_endpoint) {
	//	try {
	//		my->create_server_for_endpoint(*my->listen_endpoint, my->server);
	//
	//		ilog("start listening for http requests");
	//		my->server.listen(*my->listen_endpoint);
	//		my->server.start_accept();
	//	} catch ( const fc::exception& e ){
	//	elog( "http service failed to start: ${e}", ("e",e.to_detail_string()));
	//	throw;
	//	} catch ( const std::exception& e ){
	//	elog( "http service failed to start: ${e}", ("e",e.what()));
	//	throw;
	//	} catch (...) {
	//	elog("error thrown from http io service");
	//	throw;
	//	}
	//}

	//if h.my.ListenEndpoint != nil {
	//	Try(func() {
	//
	//	})
	//}

	//if(my->unix_endpoint) {
	//	try {
	//		my->unix_server.clear_access_channels(websocketpp::log::alevel::all);
	//		my->unix_server.init_asio(&app().get_io_service());
	//		my->unix_server.set_max_http_body_size(my->max_body_size);
	//		my->unix_server.listen(*my->unix_endpoint);
	//		my->unix_server.set_http_handler([&](connection_hdl hdl) {
	//		my->handle_http_request<detail::asio_local_with_stub_log>( my->unix_server.get_con_from_hdl(hdl));
	//	});
	//		my->unix_server.start_accept();
	//	} catch ( const fc::exception& e ){
	//	elog( "unix socket service failed to start: ${e}", ("e",e.to_detail_string()));
	//	throw;
	//	} catch ( const std::exception& e ){
	//	elog( "unix socket service failed to start: ${e}", ("e",e.what()));
	//	throw;
	//	} catch (...) {
	//	elog("error thrown from unix socket io service");
	//	throw;
	//	}
	//}
	//
	//if(my->https_listen_endpoint) {
	//	try {
	//		my->create_server_for_endpoint(*my->https_listen_endpoint, my->https_server);
	//		my->https_server.set_tls_init_handler([this](websocketpp::connection_hdl hdl) -> ssl_context_ptr{
	//		return my->on_tls_init(hdl);
	//	});
	//
	//		ilog("start listening for https requests");
	//		my->https_server.listen(*my->https_listen_endpoint);
	//		my->https_server.start_accept();
	//	} catch ( const fc::exception& e ){
	//	elog( "https service failed to start: ${e}", ("e",e.to_detail_string()));
	//	throw;
	//	} catch ( const std::exception& e ){
	//	elog( "https service failed to start: ${e}", ("e",e.what()));
	//	throw;
	//	} catch (...) {
	//	elog("error thrown from https io service");
	//	throw;
	//	}
	//}
}

func (h *HttpPlugin) PluginShutdown() {
}

func (h *HttpPlugin) VerboseErrors() bool {
	return verboseHttpErrors
}

//func (h *HttpPlugin) httpHandler(ctx *fasthttp.RequestCtx) {
//	ctx.SetContentType("text/plain; charset=utf8")
//	// Set arbitrary headers
//	ctx.Response.Header.Set("X-My-Header", "my-header-value")
//	// Set cookies
//	var c fasthttp.Cookie
//	c.SetKey("cookie-name")
//	c.SetValue("cookie-value")
//	ctx.Response.Header.SetCookie(&c)
//
//	h.my.log.Info("%s", ctx.Path())
//	h.my.log.Info("%s", ctx.Request.Body())
//
//	var Inputs map[string]string
//	path := string(ctx.Path())
//	body := bytes.NewBuffer(ctx.Request.Body())
//
//	chainROApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadOnlyApi()
//	chainRWApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadWriteApi()
//	walletMgr := App().GetPlugin(wallet_plugin.WalletPlug).(*wallet_plugin.WalletPlugin).GetWalletManager()
//
//	var reJson []byte
//	var err error
//
//	switch path {
//	case getInfoFunc:
//		resp := chainROApi.GetInfo()
//		reJson, err = json.Marshal(resp)
//
//	case getBlockFunc:
//		params := chain_plugin.GetBlockParams{}
//		if err := json.NewDecoder(body).Decode(&params); err != nil {
//			fmt.Println("error:", err)
//		}
//
//		resp := chainROApi.GetBlock(params)
//		reJson, err = json.Marshal(resp)
//
//	case getRequiredKeys:
//
//		params := &chain_plugin.GetRequiredKeysParams{}
//		if err := json.NewDecoder(body).Decode(params); err != nil {
//			fmt.Println(" get balance error:", err)
//		}
//		resp := chainROApi.GetRequiredKeys(params)
//		reJson, err = json.Marshal(resp)
//	case getCurrencyBalanceFunc:
//
//		if err := json.NewDecoder(body).Decode(&Inputs); err != nil {
//			fmt.Println(" get balance error:", err)
//			//return nil
//		}
//		log.Error("%#v", Inputs)
//
//		fmt.Println(Inputs)
//
//	case pushTxnFunc:
//		var packedTrx types.PackedTransaction
//		if err := json.NewDecoder(body).Decode(&packedTrx); err != nil {
//			fmt.Println(" decode packedTransaction:", err)
//		}
//		log.Error("%#v", Inputs)
//		chainRWApi.PushTransaction(&packedTrx)
//		reJson, err = json.Marshal("push transaction test")
//	case walletCreate:
//		var name string
//		var result string
//		if err := json.NewDecoder(body).Decode(&name); err != nil {
//			fmt.Println(" decode wallet name:", err)
//		}
//		if len(name) == 0 {
//			name = "default"
//		}
//		password, err := walletMgr.Create(name)
//		if err != nil {
//			result = err.Error()
//		} else {
//			result = password
//		}
//		reJson, err = json.Marshal(result)
//
//	case walletImportKey:
//		type Params struct {
//			Name string
//			Key  string
//		}
//		var params Params
//		if err := json.NewDecoder(body).Decode(&params); err != nil {
//			fmt.Println(" get balance error:", err)
//		}
//		err := walletMgr.ImportKey(params.Name, params.Key)
//		if err != nil {
//			reJson, err = json.Marshal(err.Error())
//		} else {
//			walletKey, err := ecc.NewPrivateKey(params.Key)
//			if err != nil {
//				fmt.Println("invalid private key %s", params.Key)
//				//try.EosThrow(&exception.PrivateKeyTypeException{}, "Invalid private key: %s", walletKeyStr)
//			}
//
//			re := fmt.Sprintf("imported private key for: %s", walletKey.PublicKey().String())
//			reJson, err = json.Marshal(re)
//		}
//
//	case walletCreateKey:
//		priKey, _ := ecc.NewRandomPrivateKey()
//		type Keys struct {
//			Pri ecc.PrivateKey `json:"Private Key"`
//			Pub ecc.PublicKey  `json:"Public Key"`
//		}
//		resp := &Keys{Pri: *priKey, Pub: priKey.PublicKey()}
//		reJson, err = json.Marshal(resp)
//
//	case walletSignTrx:
//		type WalletSignedTrx struct {
//			Trx  types.SignedTransaction `json:"signed_transaction"`
//			Keys []ecc.PublicKey         `json:"keys"`
//			ID   common.ChainIdType      `id`
//		}
//		var params WalletSignedTrx
//		if err := json.NewDecoder(body).Decode(&params); err != nil {
//			fmt.Println(" get balance error:", err)
//
//		}
//
//		h.my.log.Info("%s", err)
//		h.my.log.Info("%#v", params)
//		fmt.Println()
//
//		trx, err := walletMgr.SignTransaction(&params.Trx, params.Keys, params.ID)
//		if err != nil {
//			fmt.Println(err)
//		}
//		h.my.log.Debug("%#v", trx)
//		if err != nil {
//			reJson, err = json.Marshal(err.Error())
//		} else {
//			h.my.log.Debug("%#v", trx)
//			reJson, err = json.Marshal(trx)
//		}
//	default:
//		//return nil
//	}
//
//	fmt.Println(string(reJson), err)
//	fmt.Fprintf(ctx, "%s", reJson)
//
//}

func (h *HttpPlugin) AddHandler(url string, handler UrlHandler) {
	h.my.log.Info("add api url: %s", url)
	App().GetIoService().Post(func(err error) {
		h.my.UrlHandlers[url] = handler
	})
}

func (h *HttpPlugin) Handler(ctx *fasthttp.RequestCtx) {
	h.my.log.Info("%s", ctx.Path())
	h.my.log.Info("%s", ctx.Request.Body())

	//fmt.Fprintf(ctx, "%s", re)
	ctx.SetContentType("text/plain; charset=utf8")
	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")
	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)

	resource := string(ctx.Path())
	body := ctx.Request.Body()

	handler, ok := h.my.UrlHandlers[resource]
	if !ok {
		h.my.log.Debug("404 - not found: %s", resource)
		ctx.NotFound()
	} else {
		//con->defer_http_response();
		h.my.log.Debug("handle getBlock")
		handler(resource, body, func(code int, body []byte) {
			ctx.SetBody([]byte(body))
			ctx.SetStatusCode(code)
		})
	}
}

func (h *HttpPlugin) IsOnLoopBack() bool { //TODO
	//return (!my->listen_endpoint || my->listen_endpoint->address().is_loopback()) && (!my->https_listen_endpoint || my->https_listen_endpoint->address().is_loopback());

	return false
}
func (h *HttpPlugin) IsSecure() bool { //TODO
	//return (!my->listen_endpoint || my->listen_endpoint->address().is_loopback());
	return false
}

/**
 * @brief Structure used to create JSON error responses
 */

const detailsLimit int = 10

type errorDetail struct {
	message    string
	file       string
	lineNumber uint64
	method     string
}
type errorInfo struct {
	code    int64
	name    string
	what    string
	details []errorDetail
}

func newErrorInfo(exc Exception, includeLog bool) errorInfo {
	e := errorInfo{}
	e.code = int64(exc.Code())
	e.name = exc.String()
	e.what = exc.What()
	if includeLog {
		for _, itr := range exc.GetLog() {
			// Prevent sending trace that are too big
			if len(e.details) >= detailsLimit {
				break
			}
			// Append error
			detail := errorDetail{
				message: itr.GetMessage(),
				//file:,
				//lineNumber:,
				//method:,
			}
			e.details = append(e.details, detail)
		}
	}
	return e
}

type errorResults struct {
	code    uint16
	message string
	error   errorInfo
}

func HandleException(e interface{}, apiName, callName, body string, cb UrlResponseCallback) {
	Try(func() {
		Try(func() {
			Throw(e)
		}).Catch(func(e *UnsatisfiedAuthorization) {
			results := errorResults{code: 401, message: "UnAuthorized", error: newErrorInfo(e, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(401, re)
		}).Catch(func(e *TxDuplicate) {
			results := errorResults{409, "Conflict", newErrorInfo(e, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(409, re)
		}).Catch(func(e *EofException) {
			results := errorResults{422, "Unprocessable Entity", newErrorInfo(e, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(422, re)
			hlog.Error("Unable to parse arguments to %s.%s", apiName, callName)
			hlog.Debug("Bad arguments: %s", body)
		}).Catch(func(e Exception) {
			results := errorResults{500, "Internal Service Error", newErrorInfo(e, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(500, re)
			if e.Code() != (GreylistNetUsageExceeded{}).Code() && e.Code() != (GreylistCpuUsageExceeded{}).Code() {
				hlog.Error("FC Exception encountered while processing %s.%s", apiName, callName)
				hlog.Debug("Exception Details: %s", e.DetailMessage())
			}
		}).Catch(func(e error) {
			results := errorResults{500, "Internal Service Error",
				newErrorInfo(&FcException{Elog: log.Messages{log.FcLogMessage(log.LvlError, e.Error())}}, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(500, re)
			hlog.Error("STD Exception encountered while processing %s.%s", apiName, callName)
			hlog.Debug("Exception Details: %s", e.Error())
		}).Catch(func(interface{}) {
			results := errorResults{500, "Internal Service Error",
				newErrorInfo(&FcException{Elog: log.Messages{log.FcLogMessage(log.LvlError, "Unknown Exception")}}, verboseHttpErrors)}
			re, _ := json.Marshal(results)
			cb(500, re)
			hlog.Error("Unknown Exception encountered while processing %s.%s", apiName, callName)
		})
	}).Catch(func(interface{}) {
		hlog.Error("Exception attempting to handle exception for %s.%s", apiName, callName)
	}).End()

}
