package config

import (
	"github.com/eosspark/eos-go/common"
)

var SystemAccountName = common.S("eosio")
var NullAccountName = common.S("eosio.null")
var ProducersAccountName = common.S("eosio.prods")

// Active permission of producers account requires greater than 2/3 of the producers to authorize
var MajorityProducersPermissionName = common.S("prod.major")
var MinorityProducersPermissionName = common.S("prod.minor")

var RateLimitingPrecision uint32 = 1000 * 1000

var ActiveName uint64 = common.S("active")

var ForkDBName = "forkdb.dat"
var DBFileName = "shared_memory.bin"
var ReversibleFileName = "shared_memory_tmp.bin" //wait db modify
var BlockFileName = "blog.log"
var DefaultBlocksDirName = "blocks"
var DefaultReversibleBlocksDirName = "reversible"
var DefaultStateDirName = "state"
var DefaultStateSize uint64 = 0
var DefaultStateGuardSize uint64 = 0
var DefaultReversibleCacheSize uint64 = 0
var DefaultReversibleGuardSize uint64 = 0

//var DefaultWasmRuntime = exec.WasmInterface{}
