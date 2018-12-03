package console

//
//import (
//	"fmt"
//	"github.com/eosspark/eos-go/log"
//	"github.com/eosspark/eos-go/plugins/console_plugin/console/rpc"
//	"strings"
//	"reflect"
//	"os"
//)
//
//// startHTTP initializes and starts the HTTP RPC endpoint.
//func StartHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string) error {
//	// Short circuit if the HTTP endpoint isn't being exposed
//	if endpoint == "" {
//		return nil
//	}
//	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts)
//	if err != nil {
//		return err
//	}
//	log.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
//	// All listeners booted successfully
//	//n.httpEndpoint = endpoint
//	//n.httpListener = listener
//	//n.httpHandler = handler
//
//	fmt.Println(listener, handler)
//	return nil
//
//}
//
////rpc.StartHTTP("127.0.0.1:8888",)
//
//
//// startRPC is a helper method to start all the various RPC endpoint during node
//// startup. It's not meant to be called at any time afterwards as it makes certain
//// assumptions about the state of the node.
//func startRPC(apis []rpc.API) error {
//	// Gather all the possible APIs to surface
//
//	// Start the various API endpoints, terminating all in case of errors
//	if err := n.startInProc(apis); err != nil {
//		return err
//	}
//	if err := n.startIPC(apis); err != nil {
//		n.stopInProc()
//		return err
//	}
//	if err := n.startHTTP(n.httpEndpoint, apis, n.config.HTTPModules, n.config.HTTPCors, n.config.HTTPVirtualHosts, n.config.HTTPTimeouts); err != nil {
//		n.stopIPC()
//		n.stopInProc()
//		return err
//	}
//	if err := n.startWS(n.wsEndpoint, apis, n.config.WSModules, n.config.WSOrigins, n.config.WSExposeAll); err != nil {
//		n.stopHTTP()
//		n.stopIPC()
//		n.stopInProc()
//		return err
//	}
//	// All API endpoints started successfully
//	n.rpcAPIs = apis
//	return nil
//}
//
//
//// startInProc initializes an in-process RPC endpoint.
//func startInProc(apis []rpc.API) error {
//	// Register all the APIs exposed by the services
//	handler := rpc.NewServer()
//	for _, api := range apis {
//		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
//			return err
//		}
//		n.log.Debug("InProc registered", "namespace", api.Namespace)
//	}
//	n.inprocHandler = handler
//	return nil
//}
//
//// stopInProc terminates the in-process RPC endpoint.
//func  stopInProc() {
//	if n.inprocHandler != nil {
//		n.inprocHandler.Stop()
//		n.inprocHandler = nil
//	}
//}
//
//// startIPC initializes and starts the IPC RPC endpoint.
//func  startIPC(apis []rpc.API) error {
//	if n.ipcEndpoint == "" {
//		return nil // IPC disabled.
//	}
//	listener, handler, err := rpc.StartIPCEndpoint(n.ipcEndpoint, apis)
//	if err != nil {
//		return err
//	}
//	n.ipcListener = listener
//	n.ipcHandler = handler
//	n.log.Info("IPC endpoint opened", "url", n.ipcEndpoint)
//	return nil
//}
//
//// stopIPC terminates the IPC RPC endpoint.
//func  stopIPC() {
//	if n.ipcListener != nil {
//		n.ipcListener.Close()
//		n.ipcListener = nil
//
//		n.log.Info("IPC endpoint closed", "endpoint", n.ipcEndpoint)
//	}
//	if n.ipcHandler != nil {
//		n.ipcHandler.Stop()
//		n.ipcHandler = nil
//	}
//}
//
//// startHTTP initializes and starts the HTTP RPC endpoint.
//func  startHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string, timeouts rpc.HTTPTimeouts) error {
//	// Short circuit if the HTTP endpoint isn't being exposed
//	if endpoint == "" {
//		return nil
//	}
//	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, timeouts)
//	if err != nil {
//		return err
//	}
//	n.log.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
//	// All listeners booted successfully
//	n.httpEndpoint = endpoint
//	n.httpListener = listener
//	n.httpHandler = handler
//
//	return nil
//}
//
//// stopHTTP terminates the HTTP RPC endpoint.
//func  stopHTTP() {
//	if n.httpListener != nil {
//		n.httpListener.Close()
//		n.httpListener = nil
//
//		n.log.Info("HTTP endpoint closed", "url", fmt.Sprintf("http://%s", n.httpEndpoint))
//	}
//	if n.httpHandler != nil {
//		n.httpHandler.Stop()
//		n.httpHandler = nil
//	}
//}
//
//// startWS initializes and starts the websocket RPC endpoint.
//func  startWS(endpoint string, apis []rpc.API, modules []string, wsOrigins []string, exposeAll bool) error {
//	// Short circuit if the WS endpoint isn't being exposed
//	if endpoint == "" {
//		return nil
//	}
//	listener, handler, err := rpc.StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
//	if err != nil {
//		return err
//	}
//	n.log.Info("WebSocket endpoint opened", "url", fmt.Sprintf("ws://%s", listener.Addr()))
//	// All listeners booted successfully
//	n.wsEndpoint = endpoint
//	n.wsListener = listener
//	n.wsHandler = handler
//
//	return nil
//}
//
//// stopWS terminates the websocket RPC endpoint.
//func  stopWS() {
//	if n.wsListener != nil {
//		n.wsListener.Close()
//		n.wsListener = nil
//
//		n.log.Info("WebSocket endpoint closed", "url", fmt.Sprintf("ws://%s", n.wsEndpoint))
//	}
//	if n.wsHandler != nil {
//		n.wsHandler.Stop()
//		n.wsHandler = nil
//	}
//}
//
//// Stop terminates a running node along with all it's services. In the node was
//// not started, an error is returned.
//func  Stop() error {
//	n.lock.Lock()
//	defer n.lock.Unlock()
//
//	// Short circuit if the node's not running
//	if n.server == nil {
//		return ErrNodeStopped
//	}
//
//	// Terminate the API, services and the p2p server.
//	n.stopWS()
//	n.stopHTTP()
//	n.stopIPC()
//	n.rpcAPIs = nil
//	failure := &StopError{
//		Services: make(map[reflect.Type]error),
//	}
//	for kind, service := range n.services {
//		if err := service.Stop(); err != nil {
//			failure.Services[kind] = err
//		}
//	}
//	n.server.Stop()
//	n.services = nil
//	n.server = nil
//
//	// Release instance directory lock.
//	if n.instanceDirLock != nil {
//		if err := n.instanceDirLock.Release(); err != nil {
//			n.log.Error("Can't release datadir lock", "err", err)
//		}
//		n.instanceDirLock = nil
//	}
//
//	// unblock n.Wait
//	close(n.stop)
//
//	// Remove the keystore if it was created ephemerally.
//	var keystoreErr error
//	if n.ephemeralKeystore != "" {
//		keystoreErr = os.RemoveAll(n.ephemeralKeystore)
//	}
//
//	if len(failure.Services) > 0 {
//		return failure
//	}
//	if keystoreErr != nil {
//		return keystoreErr
//	}
//	return nil
//}
//
