package injection

import (
	//"fmt"
	"github.com/eosspark/eos-go/wasmgo/wagon/disasm"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm"

	ops "github.com/eosspark/eos-go/wasmgo/wagon/wasm/operators"
)

func Inject(m *wasm.Module) {

	importChecktime := wasm.ImportEntry{ModuleName: "env", FieldName: "checktime", Type: wasm.FuncImport{uint32(GetOrCreateCheckTimeSig(m))}}
	m.Import.Entries = append([]wasm.ImportEntry{importChecktime}, m.Import.Entries[0:]...)

	for i, f := range m.FunctionIndexSpace {
		d, err := disasm.Disassemble(f, m)
		if err != nil {
			panic(err)
		}
		injectDisassemble := InjectDisassembly(d)
		code, _ := disasm.Assemble(injectDisassemble.Code)
		m.FunctionIndexSpace[i].Body.Code = code
	}

	//shift all exported functions by 1
	for k, export := range m.Export.Entries {
		export.Index = export.Index + 1
		m.Export.Entries[k] = export
	}

	//update the start function
	if m.Start != nil {
		m.Start.Index++
	}

	//shift all table entries for call indirect
	for i, elememts := range m.Elements.Entries {
		for j, _ := range elememts.Elems {
			m.Elements.Entries[i].Elems[j]++
		}
	}

}

func InjectDisassembly(d *disasm.Disassembly) *disasm.Disassembly {
	//inject checktime

	disas := &disasm.Disassembly{MaxDepth: d.MaxDepth}
	checkTime := disasm.Instr{
		Op: ops.Op{Code: 0x10, Name: "call", Polymorphic: true},
		//Immediates: make([]interface{}, 1),
	}

	checkTime.Immediates = append(checkTime.Immediates, uint32(0))
	for _, instr := range d.Code {

		switch instr.Op.Code {
		case ops.Loop:
			disas.Code = append(disas.Code, instr)
			disas.Code = append(disas.Code, checkTime)
		case ops.Call:
			i := instr
			i.Immediates = nil
			for _, imm := range instr.Immediates {
				i.Immediates = append(i.Immediates, imm.(uint32)+1)
			}
			disas.Code = append(disas.Code, i)
		default:
			disas.Code = append(disas.Code, instr)

		}
	}

	return disas
}

func GetOrCreateCheckTimeSig(m *wasm.Module) int {

	for i, typ := range m.Types.Entries {
		if len(typ.ParamTypes) == 0 && len(typ.ReturnTypes) == 0 {
			return i
		}
	}

	m.Types.Entries = append(m.Types.Entries, wasm.FunctionSig{0, make([]wasm.ValueType, 0), make([]wasm.ValueType, 0)})
	return len(m.Types.Entries) - 1

}
