package wasmgo

import (
	"bytes"
	"errors"
	"github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"os"

	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/wasmgo/wasm"
)

const (
	MaximumLinearMemory     = 33 * 1024 * 1024 //bytes
	MaximumMutableGlobals   = 1024             //bytes
	MaximumTableElements    = 1024             //elements
	MaximumSectionElements  = 1024             //elements
	MaximumLinearMemoryInit = 64 * 1024        //bytes
	MaximumFuncLocalBytes   = 8192             //bytes
	MaximumCallDepth        = 250              //nested calls
	MaximumCodeSize         = 20 * 1024 * 1024 //bytes
	WasmPageSize            = 64 * 1024        //bytes

	// Assert(MaximumLinearMemory%WasmPageSize == 0, "MaximumLinearMemory must be mulitple of wasm page size")
	// Assert(MaximumMutableGlobals%4 == 0, "MaximumMutableGlobals must be mulitple of 4")
	// Assert(MaximumTableElements*8%4096 == 0, "maximum_table_elements*8 must be mulitple of 4096")
	// Assert(MaximumLinearMemoryInit%WasmPageSize == 0, "MaximumLinearMemoryInit must be mulitple of wasm page size")
	// Assert(MaximumFuncLocalBytes%8 == 0, "MaximumFuncLocalBytes must be mulitple of 8")
	// Assert(MaximumFuncLocalBytes > 32, "MaximumFuncLocalBytes must be greater than 32")
)

var (
	wasmGo    *WasmGo
	envModule *wasm.Module
	ignore    bool = false
	debug     bool = false
	//
)

type size_t int

type WasmGo struct {
	context EnvContext
	handles map[string]*func(*VM)
	vmCache map[crypto.Sha256]*VM

	//vmCache = make(map[crypto.Sha256]*VM)
	ilog log.Logger
}

func NewWasmGo() *WasmGo {

	if wasmGo != nil {
		return wasmGo
	}

	w := WasmGo{handles: make(map[string]*func(*VM)), vmCache: make(map[crypto.Sha256]*VM)}

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

	w.Register("db_idx128_store", dbIdx128Store)
	w.Register("db_idx128_remove", dbIdx128Remove)
	w.Register("db_idx128_update", dbIdx128Update)
	w.Register("db_idx128_find_secondary", dbIdx128findSecondary)
	w.Register("db_idx128_lowerbound", dbIdx128Lowerbound)
	w.Register("db_idx128_upperbound", dbIdx128Upperbound)
	w.Register("db_idx128_end", dbIdx128End)
	w.Register("db_idx128_next", dbIdx128Next)
	w.Register("db_idx128_previous", dbIdx128Previous)
	w.Register("db_idx128_find_primary", dbIdx128FindPrimary)

	w.Register("db_idx256_store", dbIdx256Store)
	w.Register("db_idx256_remove", dbIdx256Remove)
	w.Register("db_idx256_update", dbIdx256Update)
	w.Register("db_idx256_find_secondary", dbIdx256findSecondary)
	w.Register("db_idx256_lowerbound", dbIdx256Lowerbound)
	w.Register("db_idx256_upperbound", dbIdx256Upperbound)
	w.Register("db_idx256_end", dbIdx256End)
	w.Register("db_idx256_next", dbIdx256Next)
	w.Register("db_idx256_previous", dbIdx256Previous)
	w.Register("db_idx256_find_primary", dbIdx256FindPrimary)

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

	w.Register("db_idx_long_double_store", dbIdxLongDoubleStore)
	w.Register("db_idx_long_double_remove", dbIdxLongDoubleRemove)
	w.Register("db_idx_long_double_update", dbIdxLongDoubleUpdate)
	w.Register("db_idx_long_double_find_secondary", dbIdxLongDoublefindSecondary)
	w.Register("db_idx_long_double_lowerbound", dbIdxLongDoubleLowerbound)
	w.Register("db_idx_long_double_upperbound", dbIdxLongDoubleUpperbound)
	w.Register("db_idx_long_double_end", dbIdxLongDoubleEnd)
	w.Register("db_idx_long_double_next", dbIdxLongDoubleNext)
	w.Register("db_idx_long_double_previous", dbIdxLongDoublePrevious)
	w.Register("db_idx_long_double_find_primary", dbIdxLongDoubleFindPrimary)

	w.Register("memcpy", memcpy)
	w.Register("memmove", memmove)
	w.Register("memcmp", memcmp)
	w.Register("memset", memset)
	//w.Register("free", free)

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

	w.Register("__ashlti3", ashlti3)
	w.Register("__ashrti3", ashrti3)
	w.Register("__divti3", divti3)
	w.Register("__lshlti3", lshlti3)
	w.Register("__lshrti3", lshrti3)
	w.Register("__modti3", modti3)
	w.Register("__multi3", multi3)
	w.Register("__udivti3", udivti3)
	w.Register("__umodti3", umodti3)

	w.Register("__addtf3", addtf3)
	w.Register("__subtf3", subtf3)
	w.Register("__multf3", multf3)
	w.Register("__divtf3", divtf3)

	w.Register("__negtf2", negtf2)
	w.Register("__extendsftf2", extendsftf2)

	w.Register("__extenddftf2", extenddftf2)

	w.Register("__trunctfdf2", trunctfdf2)
	w.Register("__trunctfsf2", trunctfsf2)
	w.Register("__fixtfsi", fixtfsi)
	w.Register("__fixtfdi", fixtfdi)
	w.Register("__fixtfti", fixtfti)
	w.Register("__fixunstfsi", fixunstfsi)
	w.Register("__fixunstfdi", fixunstfdi)
	w.Register("__fixunstfti", fixunstfti)
	w.Register("__fixsfti", fixsfti)
	w.Register("__fixdfti", fixdfti)
	w.Register("__fixunssfti", fixunssfti)
	w.Register("__fixunsdfti", fixunsdfti)
	w.Register("__floatsidf", floatsidf)
	w.Register("__floatsitf", floatsitf)
	w.Register("__floatditf", floatditf)
	w.Register("__floatunsitf", floatunsitf)
	w.Register("__floatunditf", floatunditf)
	w.Register("__floattidf", floattidf)
	w.Register("__floatuntidf", floatuntidf)

	w.Register("___cmptf2", cmptf2)
	w.Register("__eqtf2", eqtf2)
	w.Register("__netf2", netf2)
	w.Register("__getf2", getf2)
	w.Register("__gttf2", gttf2)
	w.Register("__letf2", letf2)
	w.Register("__lttf2", lttf2)
	w.Register("__cmptf2", __cmptf2)
	w.Register("__unordtf2", unordtf2)

	wasmGo = &w

	wasmGo.ilog = log.New("wasmgo")
	logHandler := log.StreamHandler(os.Stdout, log.TerminalFormat(true))
	//wasmGo.ilog.SetHandler(log.LvlFilterHandler(log.LvlDebug, logHandler))
	wasmGo.ilog.SetHandler(log.LvlFilterHandler(log.LvlInfo, logHandler))
	return wasmGo
}

func (w *WasmGo) Apply(code_id *crypto.Sha256, code []byte, context EnvContext) {
	w.context = context

	var vm *VM = w.vmCache[*code_id]
	if vm != nil {
		vm.WasmGo = w
	} else {
		context.PauseBillingTimer()
		bf := bytes.NewReader(code)

		m, err := wasm.ReadModule(bf, w.importer)
		if err != nil {
			w.ilog.Error("could not read module: %v", err)
		}

		//if *verify {
		//if true {
		//	err = validate.VerifyModule(m)
		//	if err != nil {
		//		log.Fatalf("could not verify module: %v", err)
		//	}
		//}

		if m.Export == nil {
			w.ilog.Error("module has no export section")
		}

		vm, err = NewVM(m, w)
		if err != nil {
			w.ilog.Error("could not create VM: %v", err)
		}

		w.vmCache[*code_id] = vm

		//fidx := m.Function.Types[int(i)]
		//ftype := m.Types.Entries[int(fidx)]
		context.ResumeBillingTimer()
	}

	err := vm.ExecStart()
	if err != nil {
		w.ilog.Error("err=%v", err)
	}

	e, _ := vm.module.Export.Entries["apply"]
	i := int64(e.Index)

	args := make([]uint64, 3)
	args[0] = uint64(context.GetReceiver())
	args[1] = uint64(context.GetCode())
	args[2] = uint64(context.GetAct())

	o, err := vm.ExecCode(i, args[0], args[1], args[2])
	if err != nil {
		w.ilog.Error("err=%v", err)
	}
	//if len(ftype.ReturnTypes) == 0 {
	//	fmt.Printf("\n")
	//}
	if o != nil {
		w.ilog.Error("%[1]v (%[1]T)\n", o)
	}

	vm.ClearMemory()
	//w.vmCache[*code_id] = vm
}

func (w *WasmGo) Register(name string, handler func(*VM)) bool {
	if _, ok := w.handles[name]; ok {
		return false
	}

	w.handles[name] = &handler
	return true
}

func (w *WasmGo) GetHandle(name string) *func(*VM) {

	if _, ok := w.handles[name]; ok {
		return w.handles[name]
	}

	return nil
}

func (w *WasmGo) importer(name string) (*wasm.Module, error) {

	return nil, errors.New("env module will never be imported")

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

func setMemory(vm *VM, mIndex int, data []byte, dIndex int, bufferSize int) {
	//try.EosAssert(!(bufferSize > (1<<16) || mIndex+bufferSize > (1<<16)), &exception.OverlappingMemoryError{}, "access violoation")
	try.EosAssert(!(bufferSize > cap(vm.Memory()) || mIndex+bufferSize > cap(vm.Memory())), &exception.OverflowException{}, "memory overflow")
	copy(vm.Memory()[mIndex:mIndex+bufferSize], data[dIndex:dIndex+bufferSize])
}

func getMemory(vm *VM, mIndex int, bufferSize int) []byte {
	// if debug {
	// 	fmt.Println("getMemory")
	// }

	cap := cap(vm.Memory())
	if cap < mIndex || cap < mIndex+bufferSize {
		try.EosAssert(false, &exception.OverflowException{}, "memory overflow")
		//fmt.Println("getMemory heap Memory out of bound")
		return nil
	}

	bytes := make([]byte, bufferSize)
	copy(bytes[:], vm.Memory()[mIndex:mIndex+bufferSize])
	return bytes
}

func setUint128(vm *VM, index int, val *eos_math.Uint128) {

	//fmt.Println("setUint128")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint128(vm *VM, index int) *eos_math.Uint128 {

	//fmt.Println("getUint128")
	var ret eos_math.Uint128
	c := getMemory(vm, index, 16)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setUint256(vm *VM, index int, val *eos_math.Uint256) {

	//fmt.Println("setUint256")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint256(vm *VM, index int) *eos_math.Uint256 {

	//fmt.Println("getUint256")
	var ret eos_math.Uint256
	c := getMemory(vm, index, 32)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setDouble(vm *VM, index int, val *eos_math.Float64) {

	//fmt.Println("setDouble")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getDouble(vm *VM, index int) *eos_math.Float64 {

	//fmt.Println("getDouble")
	var ret eos_math.Float64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setFloat128(vm *VM, index int, val *eos_math.Float128) {

	//fmt.Println("setUint128")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getFloat128(vm *VM, index int) *eos_math.Float128 {

	//fmt.Println("Float128")
	var ret eos_math.Float128
	c := getMemory(vm, index, 16)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setUint64(vm *VM, index int, val uint64) {

	//fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint64(vm *VM, index int) uint64 {

	//fmt.Println("getUint64")
	var ret uint64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func setFloat64(vm *VM, index int, val float64) {

	//fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(vm, index, c, 0, len(c))
}

func getFloat64(vm *VM, index int) float64 {

	//fmt.Println("getUint64")
	var ret float64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func getStringLength(vm *VM, index int) int {
	var size int
	var i int
	memory := vm.Memory()
	for i = 0; i < 512; i++ {
		if memory[index+i] == 0 {
			break
		}
		size++
	}

	return size
}

func getBytes(vm *VM, index int, datalen int) []byte {
	return vm.Memory()[index : index+datalen]
}
func setSha256(vm *VM, index int, s []byte)    { setMemory(vm, index, s, 0, 32) }
func getSha256(vm *VM, index int) []byte       { return getMemory(vm, index, 32) }
func setSha512(vm *VM, index int, s []byte)    { setMemory(vm, index, s, 0, 64) }
func getSha512(vm *VM, index int) []byte       { return getMemory(vm, index, 64) }
func setSha1(vm *VM, index int, s []byte)      { setMemory(vm, index, s, 0, 20) }
func getSha1(vm *VM, index int) []byte         { return getMemory(vm, index, 20) }
func setRipemd160(vm *VM, index int, r []byte) { setMemory(vm, index, r, 0, 20) }
func getRipemd160(vm *VM, index int) []byte    { return getMemory(vm, index, 20) }
