package wasmgo

import (
	"bytes"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/wasmgo/compiler"

	//"fmt"
	"github.com/eosspark/eos-go/wasmgo/wagon/disasm"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm/leb128"

	ops "github.com/eosspark/eos-go/wasmgo/wagon/wasm/operators"
)

var (
	disabled bool
	depth    uint16
)

func Init(d bool, de uint16) {
	disabled = d
	depth = de
}

func Validate(d bool, code []byte) error {

	Init(d, 0)

	m, err := compiler.LoadModule(code)
	if err != nil {
		return err
	}

	err = ValidateModule(m.Base)
	return err
}

func ValidateModule(m *wasm.Module) error {

	for _, fn := range []func(m *wasm.Module) error{
		memoriesValidation,
		dataSegementsValidation,
		tablesValidation,
		globalsValidation,
		maximumFunctionStack,
		ensureApplyExported,
	} {
		if err := fn(m); err != nil {
			return err
		}
	}

	for _, f := range m.FunctionIndexSpace {
		d, err := disasm.Disassemble(f, m)
		if err != nil {
			return err
		}

		err = ValidateBody(d)
		if err != nil {
			return err
		}

	}

	return nil

}

func memoriesValidation(m *wasm.Module) error {
	if m.Memory == nil {
		return nil
	}
	if len(m.Memory.Entries) > 0 && m.Memory.Entries[0].Limits.Initial > MaximumLinearMemory/WasmPageSize {
		EosThrow(&WasmExecutionError{}, "Smart contract initial memory size must be less than or equal to %d KiB", MaximumLinearMemory/1024)
	}
	return nil
}

func dataSegementsValidation(m *wasm.Module) error {
	if m.Data == nil {
		return nil
	}
	for _, ds := range m.Data.Entries {
		if ds.Offset[0] != ops.I32Const {
			EosThrow(&WasmExecutionError{}, "Smart contract has unexpected memory base offset type")
		}

		r := bytes.NewReader(ds.Offset)
		r.ReadByte() //ops.I32Const
		offset, _ := leb128.ReadVarint32(r)

		if offset < 0 {
			EosThrow(&WasmExecutionError{}, "offset must positive")
		}

		if int(offset)+len(ds.Data) > MaximumLinearMemoryInit {
			EosThrow(&WasmExecutionError{}, "Smart contract data segments must lie in first %d KiB", MaximumLinearMemoryInit/1024)
		}
	}

	return nil
}

func tablesValidation(m *wasm.Module) error {
	if m.Table == nil {
		return nil
	}
	if len(m.Table.Entries) > 0 && m.Table.Entries[0].Limits.Initial > MaximumTableElements {
		EosThrow(&WasmExecutionError{}, "Smart contract table limited to %d elements", MaximumTableElements)
	}

	return nil
}

func globalsValidation(m *wasm.Module) error {

	if m.Global == nil {
		return nil
	}
	mutableGlobalsTotalSize := 0
	for _, global := range m.Global.Globals {
		if !global.Type.Mutable {
			continue
		}

		switch global.Type.Type {
		case wasm.ValueTypeI32, wasm.ValueTypeF32:
			mutableGlobalsTotalSize += 4
		case wasm.ValueTypeI64, wasm.ValueTypeF64:
			mutableGlobalsTotalSize += 8
		default:
			EosThrow(&WasmExecutionError{}, "Smart contract has unexpected global definition value type")
		}
	}

	if mutableGlobalsTotalSize > MaximumMutableGlobals {
		EosThrow(&WasmExecutionError{}, "Smart contract has more than %d bytes of mutable globals", MaximumMutableGlobals)
	}

	return nil
}

func maximumFunctionStack(m *wasm.Module) error {

	for i, functypeIndex := range m.Function.Types {

		functionStackUsage := 0

		for _, local := range m.Code.Bodies[i].Locals {
			functionStackUsage += int(local.Count * uint32(getTypeBitWidth(local.Type)/8))
		}

		for _, params := range m.Types.Entries[functypeIndex].ParamTypes {
			functionStackUsage += int(getTypeBitWidth(params) / 8)
		}

		if functionStackUsage > MaximumFuncLocalBytes {
			EosThrow(&WasmExecutionError{}, "Smart contract function has more than %d bytes of stack usage", MaximumFuncLocalBytes)
		}
	}

	return nil
}

func getTypeBitWidth(typ wasm.ValueType) int8 {

	switch typ {
	case wasm.ValueTypeI32:
		return 32
	case wasm.ValueTypeI64:
		return 64
	case wasm.ValueTypeF32:
		return 32
	case wasm.ValueTypeF64:
		return 64
	}

	EosThrow(&WasmExecutionError{}, "wasm valuetype unreachable")

	return 0

}

func ensureApplyExported(m *wasm.Module) error {

	if apply, ok := m.Export.Entries["apply"]; ok {
		if apply.Kind == wasm.ExternalFunction {

			index := int(apply.Index) - getImportFuctionNumber(m) //len(m.ImportEntries)

			functionSig := m.Types.Entries[m.Function.Types[index]]
			if len(functionSig.ParamTypes) == 3 &&
				functionSig.ParamTypes[0] == wasm.ValueTypeI64 &&
				functionSig.ParamTypes[0] == wasm.ValueTypeI64 &&
				functionSig.ParamTypes[0] == wasm.ValueTypeI64 &&
				len(functionSig.ReturnTypes) == 0 {
				return nil
			}
		}
	}

	EosThrow(&WasmExecutionError{}, "Smart contract's apply function not exported; non-existent; or wrong type")
	return nil
}

func getImportFuctionNumber(m *wasm.Module) int {

	if m.Import == nil {
		return 0
	}

	count := 0
	for _, entry := range m.Import.Entries {
		if entry.Type.Kind() == wasm.ExternalFunction {
			count++
		}
	}
	return count

}

func ValidateBody(d *disasm.Disassembly) error {

	//uint16 depth
	for _, instr := range d.Code {
		switch instr.Op.Code {
		case ops.Block, ops.Loop, ops.If, ops.Else, ops.End:
			nestedValidator(&instr, &depth)
		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			largeOffsetValidator(&instr)
		}
	}
	return nil
}

func nestedValidator(instr *disasm.Instr, depth *uint16) {

	if !disabled {
		if instr.Op.Code == ops.End {
			*depth--
		} else {
			*depth++
		}

		EosAssert(*depth < 1024, &WasmExecutionError{}, "Nested depth exceeded")
	}

}

func largeOffsetValidator(instr *disasm.Instr) {

	if instr.Immediates[1].(uint32) >= MaximumLinearMemory {
		EosThrow(&WasmExecutionError{}, "Smart contract used an invalid large memory store/load offset")
	}

}
