package console

import (
	"fmt"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/console_plugin/console/rpc"
	"strings"
)

// startHTTP initializes and starts the HTTP RPC endpoint.
func StartHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string) error {
	// Short circuit if the HTTP endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts)
	if err != nil {
		return err
	}
	log.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
	// All listeners booted successfully
	//n.httpEndpoint = endpoint
	//n.httpListener = listener
	//n.httpHandler = handler

	fmt.Println(listener, handler)
	return nil

}

//rpc.StartHTTP("127.0.0.1:8888",)
