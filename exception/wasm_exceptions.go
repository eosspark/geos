package exception

type WasmException struct{ logMessage }

func (e *WasmException) ChainExceptions() {}
func (e *WasmException) WasmExceptions()  {}
func (e *WasmException) Code() ExcTypes   { return 3070000 }
func (e *WasmException) What() string     { return "WASM Exception" }
