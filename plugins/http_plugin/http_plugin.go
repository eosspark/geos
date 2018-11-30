package http_plugin

import (
	"github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/console_plugin/console/js/eosapi"
	"github.com/eosspark/eos-go/plugins/http_plugin/rpc"
	"github.com/eosspark/eos-go/plugins/net_plugin"
	"github.com/eosspark/eos-go/plugins/wallet_plugin"
	"github.com/urfave/cli"
	"net"
	"net/http"
	"time"
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

func NewHttpPlugin(io *asio.IoContext) *HttpPlugin {
	h := &HttpPlugin{}
	h.my = NewHttpPluginImpl(io)
	h.my.self = h

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

	try.Try(func() {
		listenStr := c.String("http-server-address")
		handle := rpc.NewServer()
		apis := plugins()
		for _, api := range apis {
			if err := handle.RegisterName(api.Namespace, api.Service); err != nil {
				h.my.log.Error(err.Error())
				panic(err)
			}
			h.my.log.Debug("InProc registered :  namespace =%s", api.Namespace)
		}

		h.my.Handlers = handle
		// All APIs registered, start the HTTP listener
		listener, err := net.Listen("tcp", listenStr)
		if err != nil {
			h.my.log.Error("%s", err)
		}
		h.my.Listener = listener
		h.my.log.Info("configured http to listen on %s", listenStr)

		//err := http.ListenAndServe(listenStr, h.my.ListenEndpoint)
		//if err != nil {
		//	h.my.log.Error("failed to configure https to listen on %s , %s", listenStr, err.Error())
		//	panic(err)
		//}
		//h.my.log.Info("configured http to listen on %s", h.my.listenStr)
	})
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
	if h.my.Listener != nil {
		server := http.Server{
			Handler:      h.my.Handlers,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		go server.Serve(h.my.Listener)
	}

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

// void http_plugin::add_handler(const string& url, const url_handler& handler) {
//   ilog( "add api url: ${c}", ("c",url) );
//   app().get_io_service().post([=](){
//     my->url_handlers.insert(std::make_pair(url,handler));
//   });
// }

/*func (h *HttpPlugin) AddHandler(url string, handler *http.Handler) {
	h.my.log.Info("add api url: %s", url)

	h.my.UrlHandlers[url] = handler
}*/

//{
//std::string("/v1/""net""/""connect"),
//[ & net_mgr](string, string body, url_response_callback cb) mutable {
//try {
//if (body.empty())
//  body = "{}";
//auto result = net_mgr.connect(fc::json::from_string(body).as < std::string > ());
//cb(201, fc::json::to_string(result));
//} catch (...) {
//http_plugin::handle_exception("net", "connect", body, cb);
//}
//}
//},

// {
// 	std::string("/v1/""chain""/""get_info"), [this, ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_info(fc::json::from_string(body).as < chain_apis::read_only::get_info_params > ());
// 			cb(200 l, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_info", body, cb);
// 		}
// 	}
// }

// {
// 	std::string("/v1/""chain""/""get_block"),
// 	  [this, ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_block(fc::json::from_string(body).as < chain_apis::read_only::get_block_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_block", body, cb);
// 		}
// 	}
// },
// apis returns the collection of RPC descriptors this node offers.
func plugins() []rpc.API {
	return []rpc.API{
		{
			Namespace: "api",
			Version:   "1.0",
			Service:   eosapi.NewEosApi(),
		},
		{
			Namespace: "chain",
			Version:   "1.0",
			Service:   App().GetPlugin("ChainPlugin").(*chain_plugin.ChainPlugin).GetReadOnlyApi(),
		},
		{
			Namespace: "chain",
			Version:   "1.0",
			Service:   App().GetPlugin("ChainPlugin").(*chain_plugin.ChainPlugin).GetReadWriteApi(),
		},
		{
			Namespace: "wallet",
			Version:   "1.0",
			Service:   App().GetPlugin("WalletPlugin").(*wallet_plugin.WalletPlugin).My,
		},
		{
			Namespace: "net",
			Version:   "1.0",
			Service:   App().GetPlugin("NetPlugin").(*net_plugin.NetPlugin),
		},
	}
}
