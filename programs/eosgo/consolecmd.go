package main

// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

import (
	"github.com/eosspark/eos-go/console"
	"github.com/eosspark/eos-go/console/js/eosapi"
	"github.com/eosspark/eos-go/console/rpc"
	"github.com/eosspark/eos-go/log"
	//"github.com/eosspark/eos-go/plugins/net_plugin"
	//"github.com/eosspark/eos-go/plugins/wallet_plugin"
	"github.com/eosspark/eos-go/programs/cleos/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	consoleFlags = []cli.Flag{utils.JSpathFlag, utils.ExecFlag, utils.PreloadJSFlag}
	rpcFlags     = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
	}
	consoleCommand = cli.Command{
		Action:   localConsole,
		Name:     "console",
		Usage:    "Start an interactive JavaScript environment",
		Flags:    append(rpcFlags, consoleFlags...),
		Category: "CONSOLE COMMANDS",
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console.`,
	}

	attachCommand = cli.Command{
		Action:    remoteConsole,
		Name:      "attach",
		Usage:     "Start an interactive JavaScript environment (connect to node)",
		ArgsUsage: "[endpoint]",
		Flags:     append(consoleFlags, utils.DataDirFlag),
		Category:  "CONSOLE COMMANDS",
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console.
This command allows to open a console on a running geth node.`,
	}

	javascriptCommand = cli.Command{
		Action:    ephemeralConsole,
		Name:      "js",
		Usage:     "Execute the specified JavaScript files",
		ArgsUsage: "<jsfile> [jsfile...]",
		Flags:     consoleFlags,
		Category:  "CONSOLE COMMANDS",
		Description: `
The JavaScript VM exposes a node admin interface as well as the Ðapp
JavaScript API. See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console`,
	}
)

// localConsole starts a new geth node, attaching a JavaScript console to it at the
// same time.
func localConsole(ctx *cli.Context) error {

	// Register all the APIs exposed by the services
	handler := rpc.NewServer()
	apis := apis()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			log.Error(err.Error())
			panic(err)
		}
		//log.Debug("InProc registered :  namespace =%s, service = %s", api.Namespace, api.Service)
		log.Debug("InProc registered :  namespace =%s", api.Namespace)
	}

	client := rpc.DialInProc(handler)

	config := console.Config{
		DataDir: "./history_console_data",
		DocRoot: "testdata", //ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		//Printer:printer,
		//Preload:utils.MakeConsolePreloads(ctx),
		Preload: []string{"preload.js"},
	}

	console, err := console.New(config)
	if err != nil {
		log.Error("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	// If only a short execution was requested, evaluate and return
	if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

// remoteConsole will connect to a remote geth instance, attaching a JavaScript
// console to it.
func remoteConsole(ctx *cli.Context) error {
	//// Attach to a remotely running geth instance and start the JavaScript console
	//endpoint := ctx.Args().First()
	//if endpoint == "" {
	//	path := node.DefaultDataDir()
	//	if ctx.GlobalIsSet(utils.DataDirFlag.Name) {
	//		path = ctx.GlobalString(utils.DataDirFlag.Name)
	//	}
	//	if path != "" {
	//		if ctx.GlobalBool(utils.TestnetFlag.Name) {
	//			path = filepath.Join(path, "testnet")
	//		} else if ctx.GlobalBool(utils.RinkebyFlag.Name) {
	//			path = filepath.Join(path, "rinkeby")
	//		}
	//	}
	//	endpoint = fmt.Sprintf("%s/geth.ipc", path)
	//}
	//client, err := dialRPC(endpoint)
	//if err != nil {
	//	utils.Fatalf("Unable to attach to remote geth: %v", err)
	//}
	//config := console.Config{
	//	DataDir: utils.MakeDataDir(ctx),
	//	DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
	//	Client:  client,
	//	Preload: utils.MakeConsolePreloads(ctx),
	//}
	//
	//console, err := console.New(config)
	//if err != nil {
	//	utils.Fatalf("Failed to start the JavaScript console: %v", err)
	//}
	//defer console.Stop(false)
	//
	//if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
	//	console.Evaluate(script)
	//	return nil
	//}
	//
	//// Otherwise print the welcome screen and enter interactive mode
	//console.Welcome()
	//console.Interactive()

	return nil
}

// dialRPC returns a RPC client which connects to the given endpoint.
// The check for empty endpoint implements the defaulting logic
// for "geth attach" and "geth monitor" with no argument.
//func dialRPC(endpoint string) (*rpc.Client, error) {
//	if endpoint == "" {
//		endpoint = node.DefaultIPCEndpoint(clientIdentifier)
//	} else if strings.HasPrefix(endpoint, "rpc:") || strings.HasPrefix(endpoint, "ipc:") {
//		// Backwards compatibility with geth < 1.5 which required
//		// these prefixes.
//		endpoint = endpoint[4:]
//	}
//	return rpc.Dial(endpoint)
//}

// ephemeralConsole starts a new geth node, attaches an ephemeral JavaScript
// console to it, executes each of the files specified as arguments and tears
// everything down.
func ephemeralConsole(ctx *cli.Context) error {
	//// Create and start the node based on the CLI flags
	//node := makeFullNode(ctx)
	//startNode(ctx, node)
	//defer node.Stop()
	//
	//// Attach to the newly started node and start the JavaScript console
	//client, err := node.Attach()
	//if err != nil {
	//	utils.Fatalf("Failed to attach to the inproc geth: %v", err)
	//}
	//config := console.Config{
	//	DataDir: utils.MakeDataDir(ctx),
	//	DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
	//	Client:  client,
	//	Preload: utils.MakeConsolePreloads(ctx),
	//}
	//
	//console, err := console.New(config)
	//if err != nil {
	//	utils.Fatalf("Failed to start the JavaScript console: %v", err)
	//}
	//defer console.Stop(false)
	//
	//// Evaluate each of the specified JavaScript files
	//for _, file := range ctx.Args() {
	//	if err = console.Execute(file); err != nil {
	//		utils.Fatalf("Failed to execute %s: %v", file, err)
	//	}
	//}
	//// Wait for pending callbacks, but stop for Ctrl-C.
	//abort := make(chan os.Signal, 1)
	//signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	//
	//go func() {
	//	<-abort
	//	os.Exit(0)
	//}()
	//console.Stop(true)

	return nil
}

// apis returns the collection of RPC descriptors this node offers.
func apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "api",
			Version:   "1.0",
			Service:   eosapi.NewEosApi(),
		},
		//{
		//	Namespace: "net",
		//	Version:   "1.0",
		//	Service:   net_plugin.NewNetPlugin(),
		//},
		//{
		//	Namespace: "wallet",
		//	Version:   "1.0",
		//	Service:   wallet_plugin.NewWalletPlugin(),
		//},
	}
}
