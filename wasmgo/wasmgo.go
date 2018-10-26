package wasmgo

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"log"
	"reflect"

	"github.com/eosspark/eos-go/wasmgo/wagon/exec"
	"github.com/eosspark/eos-go/wasmgo/wagon/wasm"
)

var (
	wasmGo *WasmGo
	ignore bool = false
)

type size_t int

type WasmGo struct {
	context EnvContext
	handles map[string]interface{}
	vm      *exec.VM
}

func NewWasmGo() *WasmGo {

	if wasmGo != nil {
		return wasmGo
	}

	w := WasmGo{handles: make(map[string]interface{})}

	w.Register("action_data_size", actionDataSize)
	w.Register("read_action_data", readActionData)
	w.Register("current_receiver", currentReceiver)

	w.Register("require_auth", requireAuthorization)
	w.Register("has_auth", hasAuthorization)
	w.Register("require_auth2", requireAuth2)
	w.Register("require_recipient", requireRecipient)
	w.Register("is_account", isAccount)

	w.Register("prints", prints)
	w.Register("prints_l", printsl)
	w.Register("printi", printi)
	w.Register("printui", printui)
	w.Register("printi128", printi128)
	w.Register("printui128", printui128)
	w.Register("printsf", printsf)
	w.Register("printdf", printdf)
	w.Register("printqf", printqf)
	w.Register("printn", printn)
	w.Register("printhex", printhex)

	w.Register("assert_recover_key", assertRecoverKey)
	w.Register("recover_key", recoverKey)
	w.Register("assert_sha256", assertSha256)
	w.Register("assert_sha1", assertSha1)
	w.Register("assert_sha256", assertSha256)
	w.Register("assert_sha512", assertSha512)
	w.Register("assert_ripemd160", assertRipemd160)
	w.Register("sha1", sha1)
	w.Register("sha256", sha256)
	w.Register("sha512", sha512)
	w.Register("ripemd160", ripemd160)

	w.Register("db_store_i64", dbStoreI64)
	w.Register("db_update_i64", dbUpdateI64)
	w.Register("db_remove_i64", dbRemoveI64)
	w.Register("db_get_i64", dbGetI64)
	w.Register("db_next_i64", dbNextI64)
	w.Register("db_previous_i64", dbPreviousI64)
	w.Register("db_find_i64", dbFindI64)
	w.Register("db_lowerbound_i64", dbLowerboundI64)
	w.Register("db_upperbound_i64", dbUpperboundI64)
	w.Register("db_end_i64", dbEndI64)

	w.Register("db_idx64_store", dbIdx64Store)
	w.Register("db_idx64_remove", dbIdx64Remove)
	w.Register("db_idx64_update", dbIdx64Update)
	w.Register("db_idx64_find_secondary", dbIdx64findSecondary)
	w.Register("db_idx64_lowerbound", dbIdx64Lowerbound)
	w.Register("db_idx64_upperbound", dbIdx64Upperbound)
	w.Register("db_idx64_end", dbIdx64End)
	w.Register("db_idx64_next", dbIdx64Next)
	w.Register("db_idx64_previous", dbIdx64Previous)
	w.Register("db_idx64_find_primary", dbIdx64FindPrimary)

	w.Register("db_idx_double_store", dbIdxDoubleStore)
	w.Register("db_idx_double_remove", dbIdxDoubleRemove)
	w.Register("db_idx_double_update", dbIdxDoubleUpdate)
	w.Register("db_idx_double_find_secondary", dbIdxDoublefindSecondary)
	w.Register("db_idx_double_lowerbound", dbIdxDoubleLowerbound)
	w.Register("db_idx_double_upperbound", dbIdxDoubleUpperbound)
	w.Register("db_idx_double_end", dbIdxDoubleEnd)
	w.Register("db_idx_double_next", dbIdxDoubleNext)
	w.Register("db_idx_double_previous", dbIdxDoublePrevious)
	w.Register("db_idx_double_find_primary", dbIdxDoubleFindPrimary)

	w.Register("memcpy", memcpy)
	w.Register("memmove", memmove)
	w.Register("memcmp", memcmp)
	w.Register("memset", memset)
	w.Register("free", free)

	w.Register("check_transaction_authorization", checkTransactionAuthorization)
	w.Register("check_permission_authorization", checkPermissionAuthorization)
	w.Register("get_permission_last_used", getPermissionLastUsed)
	w.Register("get_account_creation_time", getAccountCreationTime)

	w.Register("is_feature_active", isFeatureActive)
	w.Register("activate_feature", activateFeature)
	w.Register("set_resource_limits", setResourceLimits)
	w.Register("get_resource_limits", getResourceLimits)
	w.Register("get_blockchain_parameters_packed", getBlockchainParametersPacked)
	w.Register("set_blockchain_parameters_packed", setBlockchainParametersPacked)
	w.Register("is_privileged", isPrivileged)
	w.Register("set_privileged", setPrivileged)

	w.Register("set_proposed_producers", setProposedProducers)
	w.Register("get_active_producers", getActiveProducers)

	w.Register("checktime", checkTime)
	w.Register("current_time", currentTime)
	w.Register("publication_time", publicationTime)
	w.Register("abort", abort)
	w.Register("eosio_assert", eosioAssert)
	w.Register("eosio_assert_message", eosioAssertMessage)
	w.Register("eosio_assert_code", eosioAssertCode)
	w.Register("eosio_exit", eosioExit)

	w.Register("send_inline", sendInline)
	w.Register("send_context_free_inline", sendContextFreeInline)
	w.Register("send_deferred", sendDeferred)
	w.Register("cancel_deferred", cancelDeferred)
	w.Register("read_transaction", readTransaction)
	w.Register("transaction_size", transactionSize)
	w.Register("expiration", expiration)
	w.Register("tapos_block_num", taposBlockNum)
	w.Register("tapos_block_prefix", taposBlockPrefix)
	w.Register("get_action", getAction)
	w.Register("get_context_free_data", getContextFreeData)

	wasmGo = &w

	return wasmGo
}

func (w *WasmGo) Apply(code_id *crypto.Sha256, code []byte, context EnvContext) {
	w.context = context

	bf := bytes.NewReader(code)

	m, err := wasm.ReadModule(bf, w.importer)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	// if *verify {
	// 	err = validate.VerifyModule(m)
	// 	if err != nil {
	// 		log.Fatalf("could not verify module: %v", err)
	// 	}
	// }

	if m.Export == nil {
		log.Fatalf("module has no export section")
	}

	vm, err := exec.NewVM(m, w)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	e, _ := m.Export.Entries["apply"]
	i := int64(e.Index)
	//fidx := m.Function.Types[int(i)]
	//ftype := m.Types.Entries[int(fidx)]

	w.vm = vm

	args := make([]uint64, 3)
	args[0] = uint64(context.GetReceiver())
	args[1] = uint64(context.GetCode())
	args[2] = uint64(context.GetAct())

	o, err := vm.ExecCode(i, args[0], args[1], args[2])
	if err != nil {
		fmt.Printf("\n")
		log.Printf("err=%v", err)
	}
	//if len(ftype.ReturnTypes) == 0 {
	//	fmt.Printf("\n")
	//}
	if o != nil {
		fmt.Printf("%[1]v (%[1]T)\n", o)
	}
}

func (w *WasmGo) Register(name string, handler interface{}) bool {
	if _, ok := w.handles[name]; ok {
		return false
	}

	w.handles[name] = handler
	return true
}

func (w *WasmGo) Add(handles map[string]interface{}) bool {
	for k, v := range handles {
		if _, ok := w.handles[k]; !ok {
			w.handles[k] = v
		}
	}
	return true
}

func (w *WasmGo) GetHandles() map[string]interface{} {
	return w.handles
}

func (w *WasmGo) GetHandle(name string) interface{} {

	if _, ok := w.handles[name]; ok {
		return w.handles[name]
	}

	return nil
}

func (w *WasmGo) importer(name string) (*wasm.Module, error) {

	if name == "env" {

		count := len(w.handles)

		m := wasm.NewModule()
		m.Types.Entries = make([]wasm.FunctionSig, count)
		m.FunctionIndexSpace = make([]wasm.Function, count)
		m.Export.Entries = make(map[string]wasm.ExportEntry, count)

		i := 0
		for k, v := range w.handles {

			// 1st param is *wasm_interface should be deleted
			numIn := reflect.TypeOf(v).NumIn() - 1
			args := make([]wasm.ValueType, numIn)
			for j := int(0); j < numIn; j++ {
				args[j] = reflect2wasm(reflect.TypeOf(v).In(j + 1).Kind())
			}

			numOut := reflect.TypeOf(v).NumOut()
			rtrns := make([]wasm.ValueType, numOut)
			for m := int(0); m < numOut; m++ {
				rtrns[m] = reflect2wasm(reflect.TypeOf(v).Out(m).Kind())
			}

			m.Types.Entries[i] = wasm.FunctionSig{
				//Form:        0,
				ParamTypes:  args,
				ReturnTypes: rtrns,
			}

			m.FunctionIndexSpace[i] = wasm.Function{
				Sig:  &m.Types.Entries[i],
				Host: reflect.ValueOf(v),
				Body: &wasm.FunctionBody{},
				Name: k,
			}

			m.Export.Entries[k] = wasm.ExportEntry{
				FieldStr: k,
				Kind:     wasm.ExternalFunction,
				Index:    uint32(i),
			}

			i++

		}

		return m, nil

	}

	return nil, errors.New("Only env module availible")

}

// const (
// 	Invalid Kind = iota
// 	Bool
// 	Int
// 	Int8
// 	Int16
// 	Int32
// 	Int64
// 	Uint
// 	Uint8
// 	Uint16
// 	Uint32
// 	Uint64
// 	Uintptr
// 	Float32
// 	Float64
// 	Complex64
// 	Complex128
// 	Array
// 	Chan
// 	Func
// 	Interface
// 	Map
// 	Ptr
// 	Slice
// 	String
// 	Struct
// 	UnsafePointer
// )

func reflect2wasm(kind reflect.Kind) wasm.ValueType {

	switch kind {
	case reflect.Float64:
		return wasm.ValueTypeF64
	case reflect.Float32:
		return wasm.ValueTypeF32
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return wasm.ValueTypeI32
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Struct:
		return wasm.ValueTypeI32
	case reflect.Ptr:
		return wasm.ValueTypeI64
	default:
		//panic(fmt.Sprintf("exec: return value %d invalid kind=%v", kind))
		return wasm.ValueTypeI64
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func i2b(i int) bool {
	if i > 0 {
		return true
	}
	return false
}
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func setMemory(w *WasmGo, mIndex int, data []byte, dIndex int, bufferSize int) {
	fmt.Println("setMemory")
	copy(w.vm.Memory()[mIndex:mIndex+bufferSize], data[dIndex:dIndex+bufferSize])
}

func getMemory(w *WasmGo, mIndex int, bufferSize int) []byte {
	fmt.Println("getMemory")

	cap := cap(w.vm.Memory())
	if cap < mIndex || cap < mIndex+bufferSize {
		exception.EosAssert(false, &exception.OverlappingMemoryError{}, "memcpy can only accept non-aliasing pointers")
		//fmt.Println("getMemory heap Memory out of bound")
		return nil
	}

	bytes := make([]byte, bufferSize)
	copy(bytes[:], w.vm.Memory()[mIndex:mIndex+bufferSize])
	return bytes
}

func setUint64(w *WasmGo, index int, val uint64) {

	fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(w, index, c, 0, len(c))
}

func getUint64(w *WasmGo, index int) uint64 {

	fmt.Println("getUint64")
	var ret uint64
	c := getMemory(w, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func setFloat64(w *WasmGo, index int, val float64) {

	fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(w, index, c, 0, len(c))
}

func getFloat64(w *WasmGo, index int) float64 {

	fmt.Println("getUint64")
	var ret float64
	c := getMemory(w, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func getStringLength(w *WasmGo, index int) int {
	var size int
	var i int
	for i = 0; i < 512; i++ {
		if w.vm.Memory()[index+i] == 0 {
			break
		}
		size++
	}

	return size
}

func getBytes(w *WasmGo, index int, datalen int) []byte {
	return w.vm.Memory()[index : index+datalen]
}
func setSha256(w *WasmGo, index int, s []byte)    { setMemory(w, index, s, 0, 32) }
func getSha256(w *WasmGo, index int) []byte       { return getMemory(w, index, 32) }
func setSha512(w *WasmGo, index int, s []byte)    { setMemory(w, index, s, 0, 64) }
func getSha512(w *WasmGo, index int) []byte       { return getMemory(w, index, 64) }
func setSha1(w *WasmGo, index int, s []byte)      { setMemory(w, index, s, 0, 20) }
func getSha1(w *WasmGo, index int) []byte         { return getMemory(w, index, 20) }
func setRipemd160(w *WasmGo, index int, r []byte) { setMemory(w, index, r, 0, 20) }
func getRipemd160(w *WasmGo, index int) []byte    { return getMemory(w, index, 20) }
