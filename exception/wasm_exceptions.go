package exception

import . "github.com/eosspark/eos-go/log"

type WasmException struct{ LogMessage }

func (WasmException) ChainExceptions() {}
func (WasmException) WasmExceptions()  {}
func (WasmException) Code() ExcTypes   { return 3070000 }
func (WasmException) What() string     { return "WASM Exception" }

type PageMemoryError struct{ LogMessage }

func (PageMemoryError) ChainExceptions() {}
func (PageMemoryError) WasmExceptions()  {}
func (PageMemoryError) Code() ExcTypes   { return 3070001 }
func (PageMemoryError) What() string     { return "Error in WASM page memory" }

type WasmExecutionError struct{ LogMessage }

func (WasmExecutionError) ChainExceptions() {}
func (WasmExecutionError) WasmExceptions()  {}
func (WasmExecutionError) Code() ExcTypes   { return 3070002 }
func (WasmExecutionError) What() string     { return "Runtime Error Processing WASM" }

type WasmSerializationError struct{ LogMessage }

func (WasmSerializationError) ChainExceptions() {}
func (WasmSerializationError) WasmExceptions()  {}
func (WasmSerializationError) Code() ExcTypes   { return 3070003 }
func (WasmSerializationError) What() string     { return "Serialization Error Processing WASM" }

type OverlappingMemoryError struct{ LogMessage }

func (OverlappingMemoryError) ChainExceptions() {}
func (OverlappingMemoryError) WasmExceptions()  {}
func (OverlappingMemoryError) Code() ExcTypes   { return 3070004 }
func (OverlappingMemoryError) What() string     { return "memcpy with overlapping memory" }

type BinaryenException struct{ LogMessage }

func (BinaryenException) ChainExceptions() {}
func (BinaryenException) WasmExceptions()  {}
func (BinaryenException) Code() ExcTypes   { return 3070005 }
func (BinaryenException) What() string     { return "binaryen exception" }
