package net_plugin

const (
	netFuncBase    string = "/v1/net"
	netConnect     string = netFuncBase + "/connect"
	netDisconnect  string = netFuncBase + "/disconnect"
	netStatus      string = netFuncBase + "/status"
	netConnections string = netFuncBase + "/connections"
)

func PluginStartUp() {

	//netApi := http.NewServeMux()
	//netApi.Handle(netConnect,connect())
	//netApi.Handle(netDisconnect,disconnect())
	//netApi.Handle(netStatus,status())
	//netApi.Handle(netConnections,connections())

}

func PluginInitialize() {
	//try {
	//	const auto& _http_plugin = app().get_plugin<http_plugin>();
	//	if( !_http_plugin.is_on_loopback()) {
	//	wlog( "\n"
	//	"**********SECURITY WARNING**********\n"
	//	"*                                  *\n"
	//	"* --         Net API            -- *\n"
	//	"* - EXPOSED to the LOCAL NETWORK - *\n"
	//	"* - USE ONLY ON SECURE NETWORKS! - *\n"
	//	"*                                  *\n"
	//	"************************************\n" );
	//}
	//} FC_LOG_AND_RETHROW()
}
