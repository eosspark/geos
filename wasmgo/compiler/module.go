package compiler

import (
	"bytes"
	"encoding/binary"
	"errors"
	//"github.com/eosspark/eos-go/wasmgo/injection"

	//"fmt"
	"github.com/eosspark/eos-go/wasmgo/wagon/disasm"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm"
	//"github.com/eosspark/eos-go/wasmgo/wagon/validate"
	"github.com/eosspark/eos-go/wasmgo/compiler/opcodes"
	"github.com/eosspark/eos-go/wasmgo/utils"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm/leb128"
)

type Module struct {
	Base                 *wasm.Module
	FunctionNames        map[int]string
	DisableFloatingPoint bool
}

type InterpreterCode struct {
	NumRegs    int
	NumParams  int
	NumLocals  int
	NumReturns int
	Bytes      []byte
	JITInfo    interface{}
	JITDone    bool
}

func importer(name string) (*wasm.Module, error) {
	return nil, errors.New("env module will never be imported")
}

func LoadModule(raw []byte) (*Module, error) {
	reader := bytes.NewReader(raw)

	//m, err := wasm.ReadModule(reader, importer)
	m, err := wasm.ReadModule(reader, nil)
	if err != nil {
		return nil, err
	}

	/*err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}*/

	functionNames := make(map[int]string)

	for _, sec := range m.Customs {
		if sec.Name == "name" {
			r := bytes.NewReader(sec.RawSection.Bytes)
			for {
				ty, err := leb128.ReadVarUint32(r)
				if err != nil || ty != 1 {
					break
				}
				payloadLen, err := leb128.ReadVarUint32(r)
				if err != nil {
					panic(err)
				}
				data := make([]byte, int(payloadLen))
				n, err := r.Read(data)
				if err != nil {
					panic(err)
				}
				if n != len(data) {
					panic("len mismatch")
				}
				{
					r := bytes.NewReader(data)
					for {
						count, err := leb128.ReadVarUint32(r)
						if err != nil {
							break
						}
						for i := 0; i < int(count); i++ {
							index, err := leb128.ReadVarUint32(r)
							if err != nil {
								panic(err)
							}
							nameLen, err := leb128.ReadVarUint32(r)
							if err != nil {
								panic(err)
							}
							name := make([]byte, int(nameLen))
							n, err := r.Read(name)
							if err != nil {
								panic(err)
							}
							if n != len(name) {
								panic("len mismatch")
							}
							functionNames[int(index)] = string(name)
							//fmt.Printf("%d -> %s\n", int(index), string(name))
						}
					}
				}
			}
			//fmt.Printf("%d function names written\n", len(functionNames))
		}
	}

	return &Module{
		Base:          m,
		FunctionNames: functionNames,
	}, nil
}

func (m *Module) CompileForInterpreter(gp GasPolicy) (_retCode []InterpreterCode, retErr error) {
	defer utils.CatchPanic(&retErr)

	ret := make([]InterpreterCode, 0)
	importTypeIDs := make([]int, 0)

	if m.Base.Import != nil {

		//add checktime
		// {
		// 	buf := &bytes.Buffer{}

		// 	binary.Write(buf, binary.LittleEndian, uint32(1)) // value ID
		// 	binary.Write(buf, binary.LittleEndian, opcodes.InvokeImport)
		// 	binary.Write(buf, binary.LittleEndian, uint32(0))
		// 	binary.Write(buf, binary.LittleEndian, opcodes.ReturnVoid)

		// 	code := buf.Bytes()

		// 	ret = append(ret, InterpreterCode{
		// 		NumRegs:    2,
		// 		NumParams:  0,
		// 		NumLocals:  0,
		// 		NumReturns: 0,
		// 		Bytes:      code,
		// 	})
		// 	tyID := injection.GetOrCreateCheckTimeSig(m.Base)
		// 	importTypeIDs = append(importTypeIDs, int(tyID))
		// }

		for i := 0; i < len(m.Base.Import.Entries); i++ {
			e := &m.Base.Import.Entries[i]
			if e.Type.Kind() != wasm.ExternalFunction {
				continue
			}
			tyID := e.Type.(wasm.FuncImport).Type
			ty := &m.Base.Types.Entries[int(tyID)]

			buf := &bytes.Buffer{}

			binary.Write(buf, binary.LittleEndian, uint32(1)) // value ID
			binary.Write(buf, binary.LittleEndian, opcodes.InvokeImport)
			binary.Write(buf, binary.LittleEndian, uint32(i))

			binary.Write(buf, binary.LittleEndian, uint32(0))
			if len(ty.ReturnTypes) != 0 {
				binary.Write(buf, binary.LittleEndian, opcodes.ReturnValue)
				binary.Write(buf, binary.LittleEndian, uint32(1))
			} else {
				binary.Write(buf, binary.LittleEndian, opcodes.ReturnVoid)
			}

			code := buf.Bytes()

			ret = append(ret, InterpreterCode{
				NumRegs:    2,
				NumParams:  len(ty.ParamTypes),
				NumLocals:  0,
				NumReturns: len(ty.ReturnTypes),
				Bytes:      code,
			})
			importTypeIDs = append(importTypeIDs, int(tyID))
		}
	}

	numFuncImports := len(ret)
	ret = append(ret, make([]InterpreterCode, len(m.Base.FunctionIndexSpace))...)

	//var index int = 0
	for i, f := range m.Base.FunctionIndexSpace {
		//fmt.Printf("Compiling function %d (%+v) with %d locals\n", i, f.Sig, len(f.Body.Locals))
		d, err := disasm.Disassemble(f, m.Base)
		if err != nil {
			panic(err)
		}

		//fmt.Printf("Compiling function %d (%+v) with %d locals %d structions \n", i, f.Sig, len(f.Body.Locals), len(d.Code))
		//injectDisassemble := injection.Inject(m.Base, d)

		compiler := NewSSAFunctionCompiler(m.Base, d)
		compiler.CallIndexOffset = numFuncImports
		compiler.Compile(importTypeIDs)
		if m.DisableFloatingPoint {
			compiler.FilterFloatingPoint()
		}
		if gp != nil {
			compiler.InsertGasCounters(gp)
		}
		//fmt.Println(compiler.Code)
		//fmt.Printf("%+v\n", compiler.NewCFGraph())
		numRegs := compiler.RegAlloc()
		//fmt.Println(compiler.Code)
		numLocals := 0
		for _, v := range f.Body.Locals {
			numLocals += int(v.Count)
		}
		ret[numFuncImports+i] = InterpreterCode{
			NumRegs:    numRegs,
			NumParams:  len(f.Sig.ParamTypes),
			NumLocals:  numLocals,
			NumReturns: len(f.Sig.ReturnTypes),
			Bytes:      compiler.Serialize(),
		}

		//index++
	}

	return ret, nil
}
