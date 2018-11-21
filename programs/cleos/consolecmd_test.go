package main

import (
	"github.com/eosspark/eos-go/console"
	"github.com/eosspark/eos-go/console/rpc"
	"github.com/eosspark/eos-go/log"
	"testing"
)

func TestConsoleWelcome(t *testing.T) {

	// Register all the APIs exposed by the services
	handler := rpc.NewServer()
	apis := apis()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			log.Error(err.Error())
			panic(err)
		}
		log.Debug("InProc registered :  namespace =%s, service = %s", api.Namespace, api.Service)
	}

	client := rpc.DialInProc(handler)

	config := console.Config{
		DataDir: "./history_console_data",
		DocRoot: "testdata",
		Client:  client,
		//Printer:printer,
		//Preload:utils.MakeConsolePreloads(ctx),
		Preload: []string{"preload.js"}, //TODO
	}

	console, err := console.New(config)
	if err != nil {
		log.Error("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	// If only a short execution was requested, evaluate and return
	//if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
	//	console.Evaluate(script)
	//	return nil
	//}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()

	log.Warn("preloaded:  ")

	err = console.Evaluate("preloaded")

	if err != nil {
		log.Error(err.Error())
	}

	//var resp string
	//if err := client.Call(&resp, "api_createKey"); err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(resp)

	//var info eosapi.InfoResp
	//if err := client.Call(&info, "api_getInfo" ); err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(info)
	//
	//var status net_plugin.PeerStatus
	//if err := client.Call(&status, "net_plugin_status", "127.0.0.1:9876"); err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(status)

	//var slicepeerstate  []net_plugin.PeerStatus
	//if err := client.Call(&slicepeerstate, "net_plugin_connections" ); err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(slicepeerstate)
	console.Interactive()

}
