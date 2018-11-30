package console_plugin

import (
	"github.com/eosspark/eos-go/plugins/http_plugin/rpc"
	"strings"
)

// startRPC is a helper method to start all the various RPC endpoint during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (impl *ConsolePluginImpl) startRPC(apis []rpc.API) error {
	// Gather all the possible APIs to surface

	//// Start the various API endpoints, terminating all in case of errors
	//if err := impl.startInProc(apis); err != nil {
	//	return err
	//}
	//if err := impl.startIPC(apis); err != nil {
	//	impl.stopInProc()
	//	return err
	//}
	if err := impl.startHTTP(impl.httpEndpoint, apis, impl.config.HTTPModules, impl.config.HTTPCors, impl.config.HTTPVirtualHosts, impl.config.HTTPTimeouts); err != nil {
		//impl.stopIPC()
		//impl.stopInProc()
		return err
	}
	//if err := impl.startWS(n.wsEndpoint, apis, n.config.WSModules, n.config.WSOrigins, n.config.WSExposeAll); err != nil {
	//	impl.stopHTTP()
	//	impl.stopIPC()
	//	impl.stopInProc()
	//	return err
	//}
	// All API endpoints started successfully
	impl.rpcAPIs = apis
	return nil
}

//// startInProc initializes an in-process RPC endpoint.
//func (impl *ConsolePluginImpl)startInProc(apis []rpc.API) error {
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
//func  (impl *ConsolePluginImpl)stopInProc() {
//	if n.inprocHandler != nil {
//		n.inprocHandler.Stop()
//		n.inprocHandler = nil
//	}
//}

//// startIPC initializes and starts the IPC RPC endpoint.
//func  (impl *ConsolePluginImpl)startIPC(apis []rpc.API) error {
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
//func  (impl *ConsolePluginImpl)stopIPC() {
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

// startHTTP initializes and starts the HTTP RPC endpoint.
func (impl *ConsolePluginImpl) startHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string, timeouts rpc.HTTPTimeouts) error {
	// Short circuit if the HTTP endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, timeouts)
	if err != nil {
		return err
	}
	impl.log.Info("HTTP endpoint opened,url : http://%s%s%s", endpoint, strings.Join(cors, ","), strings.Join(vhosts, ","))
	// All listeners booted successfully
	impl.httpEndpoint = endpoint
	impl.httpListener = listener
	impl.httpHandler = handler

	return nil
}

// stopHTTP terminates the HTTP RPC endpoint.
func (impl *ConsolePluginImpl) stopHTTP() {
	if impl.httpListener != nil {
		impl.httpListener.Close()
		impl.httpListener = nil

		impl.log.Info("HTTP endpoint closed, url： http://%s", impl.httpEndpoint)
	}
	if impl.httpHandler != nil {
		impl.httpHandler.Stop()
		impl.httpHandler = nil
	}
}

//// startWS initializes and starts the websocket RPC endpoint.
//func  (impl *ConsolePluginImpl)startWS(endpoint string, apis []rpc.API, modules []string, wsOrigins []string, exposeAll bool) error {
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
//func  (impl *ConsolePluginImpl)stopWS() {
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

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (impl *ConsolePluginImpl) Stop() error {
	//n.lock.Lock()
	//defer n.lock.Unlock()
	//
	//// Short circuit if the node's not running
	//if n.server == nil {
	//	return ErrNodeStopped
	//}
	//
	//// Terminate the API, services and the p2p server.
	//n.stopWS()
	//n.stopHTTP()
	//n.stopIPC()
	//n.rpcAPIs = nil
	//failure := &StopError{
	//	Services: make(map[reflect.Type]error),
	//}
	//for kind, service := range n.services {
	//	if err := service.Stop(); err != nil {
	//		failure.Services[kind] = err
	//	}
	//}
	//n.server.Stop()
	//n.services = nil
	//n.server = nil
	//
	//// Release instance directory lock.
	//if n.instanceDirLock != nil {
	//	if err := n.instanceDirLock.Release(); err != nil {
	//		n.log.Error("Can't release datadir lock", "err", err)
	//	}
	//	n.instanceDirLock = nil
	//}
	//
	//// unblock n.Wait
	//close(n.stop)
	//
	//// Remove the keystore if it was created ephemerally.
	//var keystoreErr error
	//if n.ephemeralKeystore != "" {
	//	keystoreErr = os.RemoveAll(n.ephemeralKeystore)
	//}
	//
	//if len(failure.Services) > 0 {
	//	return failure
	//}
	//if keystoreErr != nil {
	//	return keystoreErr
	//}
	return nil
}
