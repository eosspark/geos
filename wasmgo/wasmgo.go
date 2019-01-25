package wasmgo

import (
	"fmt"
	"github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"os"
	//"time"
	//"github.com/eosspark/eos-go/wasmgo/wasm"
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

var wasmGo *WasmGo

//type size_t int

type WasmGo struct {
	context EnvContext
	vmCache map[crypto.Sha256]*VirtualMachine

	ilog log.Logger
}

func NewWasmGo() *WasmGo {

	if wasmGo != nil {
		return wasmGo
	}

	w := WasmGo{vmCache: make(map[crypto.Sha256]*VirtualMachine)}

	wasmGo = &w

	wasmGo.ilog = log.New("wasmgo")
	logHandler := log.StreamHandler(os.Stdout, log.TerminalFormat(true))
	wasmGo.ilog.SetHandler(log.LvlFilterHandler(log.LvlDebug, logHandler))
	//wasmGo.ilog.SetHandler(log.LvlFilterHandler(log.LvlInfo, logHandler))
	return wasmGo
}

func (w *WasmGo) Apply(codeId *crypto.Sha256, code []byte, context EnvContext) {
	w.context = context

	var vm *VirtualMachine = w.vmCache[*codeId]
	if vm != nil {
		vm.WasmGo = w
		vm.Reset()
	} else {
		context.PauseBillingTimer()
		//vm, err := NewVirtualMachine(w, code, exec.VMConfig{
		//	EnableJIT:false,
		//	MaxMemoryPages: MaximumLinearMemory / WasmPageSize,
		//	DefaultMemoryPages: 1,
		//	DefaultTableSize:   65536}, new(Resolver), nil)
		var err error
		vm, err = NewVirtualMachine(w, code, VMConfig{
			EnableJIT:          false,
			MaxMemoryPages:     MaximumLinearMemory / WasmPageSize,
			DefaultMemoryPages: 1,
			DefaultTableSize:   65536,
		}, new(Resolver), nil)

		if err != nil {
			w.ilog.Error("could not create VM: %v", err)
		}

		w.vmCache[*codeId] = vm
		context.ResumeBillingTimer()
	}

	//start := time.Now()
	entryID, ok := vm.GetFunctionExport("apply")
	if !ok {
		w.ilog.Info("Entry function %s not found", "apply")
	}

	if vm.Module.Base.Start != nil {
		startID := int(vm.Module.Base.Start.Index)
		_, err := vm.Run(startID)
		if err != nil {
			// vm.PrintStackTrace()
			// panic(err)
			w.ilog.Error("vm execute err: %v", err)
		}
	}

	args := make([]int64, 3)
	args[0] = int64(context.GetReceiver())
	args[1] = int64(context.GetCode())
	args[2] = int64(context.GetAct())

	// Run the WebAssembly module's entry function.
	_, err := vm.Run(entryID, args...)
	if err != nil {
		// vm.PrintStackTrace()
		// panic(err)
		w.ilog.Error("vm execute err: %v", err)
	}
	//end := time.Now()
	//w.ilog.Info("return value = %d, duration = %v", ret, end.Sub(start))

	//clear VM status
}

// Resolver defines imports for WebAssembly modules ran in Life.
type Resolver struct {
	tempRet0 int64
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) FunctionImport {
	//fmt.Printf("Resolve func: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "action_data_size":
			return actionDataSize
		case "read_action_data":
			return readActionData
		case "current_receiver":
			return currentReceiver

		case "require_auth":
			return requireAuthorization
		case "has_auth":
			return hasAuthorization
		case "require_auth2":
			return requireAuth2
		case "require_recipient":
			return requireRecipient
		case "is_account":
			return isAccount

		case "__ashlti3":
			return ashlti3
		case "__ashrti3":
			return ashrti3
		case "__lshlti3":
			return lshrti3
		case "__lshrti3":
			return ashrti3
		case "__divti3":
			return divti3
		case "__udivti3":
			return udivti3
		case "__multi3":
			return multi3
		case "__modti3":
			return modti3
		case "__umodti3":
			return umodti3
		case "__addtf3":
			return addtf3
		case "__subtf3":
			return subtf3
		case "__multf3":
			return multf3
		case "__divtf3":
			return divtf3
		case "__negtf2":
			return negtf2
		case "__extendsftf2":
			return extendsftf2
		case "__extenddftf2":
			return extenddftf2
		case "__trunctfdf2":
			return trunctfdf2
		case "__trunctfsf2":
			return trunctfsf2
		case "__fixtfsi":
			return fixtfsi
		case "__fixtfdi":
			return fixtfdi
		case "__fixtfti":
			return fixtfti
		case "__fixunstfsi":
			return fixunstfsi
		case "__fixunstfdi":
			return fixunstfdi
		case "__fixunstfti":
			return fixunstfti
		case "__fixsfti":
			return fixsfti
		case "__fixdfti":
			return fixdfti
		case "__fixunssfti":
			return fixunssfti
		case "__fixunsdfti":
			return fixunsdfti
		case "__floatsidf":
			return floatsidf
		case "__floatsitf":
			return floatsitf
		case "__floatditf":
			return floatditf
		case "__floatunsitf":
			return floatunsitf
		case "__floatunditf":
			return floatunditf
		case "__floattidf":
			return floattidf
		case "__floatuntidf":
			return floatuntidf
		case "___cmptf2":
			return cmptf2
		case "__eqtf2":
			return eqtf2
		case "__netf2":
			return netf2
		case "__getf2":
			return getf2
		case "__gttf2":
			return gttf2
		case "__letf2":
			return letf2
		case "__lttf2":
			return lttf2
		case "__cmptf2":
			return __cmptf2
		case "__unordtf2":
			return unordtf2

		case "assert_recover_key":
			return assertRecoverKey
		case "recover_key":
			return recoverKey
		case "assert_sha256":
			return assertSha256
		case "assert_sha1":
			return assertSha1
		case "assert_sha512":
			return assertSha512
		case "assert_ripemd160":
			return assertRipemd160
		case "sha1":
			return sha1
		case "sha256":
			return sha256
		case "sha512":
			return sha512
		case "ripemd160":
			return ripemd160

		case "db_store_i64":
			return dbStoreI64
		case "db_update_i64":
			return dbUpdateI64
		case "db_remove_i64":
			return dbRemoveI64
		case "db_get_i64":
			return dbGetI64
		case "db_next_i64":
			return dbNextI64
		case "db_previous_i64":
			return dbPreviousI64
		case "db_find_i64":
			return dbFindI64
		case "db_lowerbound_i64":
			return dbLowerboundI64
		case "db_upperbound_i64":
			return dbUpperboundI64
		case "db_end_i64":
			return dbEndI64

		case "db_idx64_store":
			return dbIdx64Store
		case "db_idx64_remove":
			return dbIdx64Remove
		case "db_idx64_update":
			return dbIdx64Update
		case "db_idx64_find_secondary":
			return dbIdx64findSecondary
		case "db_idx64_lowerbound":
			return dbIdx64Lowerbound
		case "db_idx64_upperbound":
			return dbIdx64Upperbound
		case "db_idx64_end":
			return dbIdx64End
		case "db_idx64_next":
			return dbIdx64Next
		case "db_idx64_previous":
			return dbIdx64Previous
		case "db_idx64_find_primary":
			return dbIdx64FindPrimary

		case "db_idx128_store":
			return dbIdx128Store
		case "db_idx128_remove":
			return dbIdx128Remove
		case "db_idx128_update":
			return dbIdx128Update
		case "db_idx128_find_secondary":
			return dbIdx128findSecondary
		case "db_idx128_lowerbound":
			return dbIdx128Lowerbound
		case "db_idx128_upperbound":
			return dbIdx128Upperbound
		case "db_idx128_end":
			return dbIdx128End
		case "db_idx128_next":
			return dbIdx128Next
		case "db_idx128_previous":
			return dbIdx128Previous
		case "db_idx128_find_primary":
			return dbIdx128FindPrimary

		case "db_idx256_store":
			return dbIdx256Store
		case "db_idx256_remove":
			return dbIdx256Remove
		case "db_idx256_update":
			return dbIdx256Update
		case "db_idx256_find_secondary":
			return dbIdx256findSecondary
		case "db_idx256_lowerbound":
			return dbIdx256Lowerbound
		case "db_idx256_upperbound":
			return dbIdx256Upperbound
		case "db_idx256_end":
			return dbIdx256End
		case "db_idx256_next":
			return dbIdx256Next
		case "db_idx256_previous":
			return dbIdx256Previous
		case "db_idx256_find_primary":
			return dbIdx256FindPrimary

		case "db_idx_double_store":
			return dbIdxDoubleStore
		case "db_idx_double_remove":
			return dbIdxDoubleRemove
		case "db_idx_double_update":
			return dbIdxDoubleUpdate
		case "db_idx_double_find_secondary":
			return dbIdxDoublefindSecondary
		case "db_idx_double_lowerbound":
			return dbIdxDoubleLowerbound
		case "db_idx_double_upperbound":
			return dbIdxDoubleUpperbound
		case "db_idx_double_end":
			return dbIdxDoubleEnd
		case "db_idx_double_next":
			return dbIdxDoubleNext
		case "db_idx_double_previous":
			return dbIdxDoublePrevious
		case "db_idx_double_find_primary":
			return dbIdxDoubleFindPrimary

		case "db_idx_long_double_store":
			return dbIdxLongDoubleStore
		case "db_idx_long_double_remove":
			return dbIdxLongDoubleRemove
		case "db_idx_long_double_update":
			return dbIdxLongDoubleUpdate
		case "db_idx_long_double_find_secondary":
			return dbIdxLongDoublefindSecondary
		case "db_idx_long_double_lowerbound":
			return dbIdxLongDoubleLowerbound
		case "db_idx_long_double_upperbound":
			return dbIdxLongDoubleUpperbound
		case "db_idx_long_double_end":
			return dbIdxLongDoubleEnd
		case "db_idx_long_double_next":
			return dbIdxLongDoubleNext
		case "db_idx_long_double_previous":
			return dbIdxLongDoublePrevious
		case "db_idx_long_double_find_primary":
			return dbIdxLongDoubleFindPrimary

		case "memcpy":
			return memcpy
		case "memmove":
			return memmove
		case "memcmp":
			return memcmp
		case "memset":
			return memset

		case "prints":
			return prints
		case "prints_l":
			return printsl
		case "printi":
			return printi
		case "printui":
			return printui
		case "printi128":
			return printi128
		case "printui128":
			return printui128
		case "printsf":
			return printsf
		case "printdf":
			return printdf
		case "printqf":
			return printqf
		case "printn":
			return printn
		case "printhex":
			return printhex

		//case "assert_recover_key":
		//	return assertRecoverKey
		//case "recover_key":
		//	return recoverKey
		//case "assert_sha256":
		//	return assertSha256
		//case "assert_sha1":
		//	return assertSha1
		//case "assert_sha512":
		//	return assertSha512
		//case "assert_ripemd160":
		//	return assertRipemd160
		//case "sha1":
		//	return sha1
		//case "sha256":
		//	return sha256
		//case "sha512":
		//	return sha512
		//case "ripemd160":
		//	return ripemd160

		case "check_transaction_authorization":
			return checkTransactionAuthorization
		case "check_permission_authorization":
			return checkPermissionAuthorization
		case "get_permission_last_used":
			return getPermissionLastUsed
		case "get_account_creation_time":
			return getAccountCreationTime

		case "is_feature_active":
			return isFeatureActive
		case "activate_feature":
			return activateFeature
		case "set_resource_limits":
			return setResourceLimits
		case "get_resource_limits":
			return getResourceLimits
		case "get_blockchain_parameters_packed":
			return getBlockchainParametersPacked
		case "set_blockchain_parameters_packed":
			return setBlockchainParametersPacked
		case "is_privileged":
			return isPrivileged
		case "set_privileged":
			return setPrivileged

		case "set_proposed_producers":
			return setProposedProducers
		case "get_active_producers":
			return getActiveProducers

		case "checktime":
			return checkTime
		case "current_time":
			return currentTime
		case "publication_time":
			return publicationTime
		case "abort":
			return abort
		case "eosio_assert":
			return eosioAssert
		case "eosio_assert_message":
			return eosioAssertMessage
		case "eosio_assert_code":
			return eosioAssertCode
		case "eosio_exit":
			return eosioExit

		case "send_inline":
			return sendInline
		case "send_context_free_inline":
			return sendContextFreeInline
		case "send_deferred":
			return sendDeferred
		case "cancel_deferred":
			return cancelDeferred
		case "read_transaction":
			return readTransaction
		case "transaction_size":
			return transactionSize
		case "expiration":
			return expiration
		case "tapos_block_num":
			return taposBlockNum
		case "tapos_block_prefix":
			return taposBlockPrefix
		case "get_action":
			return getAction
		case "get_context_free_data":
			return getContextFreeData
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	fmt.Printf("Resolve global: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

// func (w *WasmGo) Register(name string, handler func(*VM)) bool {
// 	if _, ok := w.handles[name]; ok {
// 		return false
// 	}

// 	w.handles[name] = &handler
// 	return true
// }

// func (w *WasmGo) GetHandle(name string) *func(*VM) {

// 	if _, ok := w.handles[name]; ok {
// 		return w.handles[name]
// 	}

// 	return nil
// }

// func (w *WasmGo) importer(name string) (*wasm.Module, error) {

// 	return nil, errors.New("env module will never be imported")

// }

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

func setMemory(vm *VirtualMachine, mIndex int, data []byte, dIndex int, bufferSize int) {
	//try.EosAssert(!(bufferSize > (1<<16) || mIndex+bufferSize > (1<<16)), &exception.OverlappingMemoryError{}, "access violoation")
	try.EosAssert(!(bufferSize > cap(vm.Memory) || mIndex+bufferSize > cap(vm.Memory)), &exception.OverflowException{}, "memory overflow")
	copy(vm.Memory[mIndex:mIndex+bufferSize], data[dIndex:dIndex+bufferSize])
}

func getMemory(vm *VirtualMachine, mIndex int, bufferSize int) []byte {
	// if debug {
	// 	fmt.Println("getMemory")
	// }

	cap := cap(vm.Memory)
	if cap < mIndex || cap < mIndex+bufferSize {
		try.EosAssert(false, &exception.OverflowException{}, "memory overflow")
		//fmt.Println("getMemory heap Memory out of bound")
		return nil
	}

	bytes := make([]byte, bufferSize)
	copy(bytes[:], vm.Memory[mIndex:mIndex+bufferSize])
	return bytes
}

func setUint128(vm *VirtualMachine, index int, val *eos_math.Uint128) {

	//fmt.Println("setUint128")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint128(vm *VirtualMachine, index int) *eos_math.Uint128 {

	//fmt.Println("getUint128")
	var ret eos_math.Uint128
	c := getMemory(vm, index, 16)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setUint256(vm *VirtualMachine, index int, val *eos_math.Uint256) {

	//fmt.Println("setUint256")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint256(vm *VirtualMachine, index int) *eos_math.Uint256 {

	//fmt.Println("getUint256")
	var ret eos_math.Uint256
	c := getMemory(vm, index, 32)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setDouble(vm *VirtualMachine, index int, val *eos_math.Float64) {

	//fmt.Println("setDouble")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getDouble(vm *VirtualMachine, index int) *eos_math.Float64 {

	//fmt.Println("getDouble")
	var ret eos_math.Float64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setFloat128(vm *VirtualMachine, index int, val *eos_math.Float128) {

	//fmt.Println("setUint128")
	c, _ := rlp.EncodeToBytes(*val)
	setMemory(vm, index, c, 0, len(c))
}

func getFloat128(vm *VirtualMachine, index int) *eos_math.Float128 {

	//fmt.Println("Float128")
	var ret eos_math.Float128
	c := getMemory(vm, index, 16)
	rlp.DecodeBytes(c, &ret)
	return &ret
}

func setUint64(vm *VirtualMachine, index int, val uint64) {

	//fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(vm, index, c, 0, len(c))
}

func getUint64(vm *VirtualMachine, index int) uint64 {

	//fmt.Println("getUint64")
	var ret uint64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func setFloat64(vm *VirtualMachine, index int, val float64) {

	//fmt.Println("setUint64")
	c, _ := rlp.EncodeToBytes(val)
	setMemory(vm, index, c, 0, len(c))
}

func getFloat64(vm *VirtualMachine, index int) float64 {

	//fmt.Println("getUint64")
	var ret float64
	c := getMemory(vm, index, 8)
	rlp.DecodeBytes(c, &ret)
	return ret
}

func getStringLength(vm *VirtualMachine, index int) int {
	var size int
	var i int
	memory := vm.Memory
	for i = 0; i < 512; i++ {
		if memory[index+i] == 0 {
			break
		}
		size++
	}

	return size
}

func getBytes(vm *VirtualMachine, index int, datalen int) []byte {
	return vm.Memory[index : index+datalen]
}

// func setSha256(vm *VM, index int, s []byte)    { setMemory(vm, index, s, 0, 32) }
// func getSha256(vm *VM, index int) []byte       { return getMemory(vm, index, 32) }
// func setSha512(vm *VM, index int, s []byte)    { setMemory(vm, index, s, 0, 64) }
// func getSha512(vm *VM, index int) []byte       { return getMemory(vm, index, 64) }
// func setSha1(vm *VM, index int, s []byte)      { setMemory(vm, index, s, 0, 20) }
// func getSha1(vm *VM, index int) []byte         { return getMemory(vm, index, 20) }
// func setRipemd160(vm *VM, index int, r []byte) { setMemory(vm, index, r, 0, 20) }
// func getRipemd160(vm *VM, index int) []byte    { return getMemory(vm, index, 20) }
