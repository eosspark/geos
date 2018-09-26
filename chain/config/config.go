package config

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/cvm/exec"
)

var SystemAccountName = common.StringToName("eosio")
var NullAccountName = common.StringToName("eosio.null")
var ProducersAccountName = common.StringToName("eosio.prods")

// Active permission of producers account requires greater than 2/3 of the producers to authorize
var MajorityProducersPermissionName = common.StringToName("prod.major")
var MinorityProducersPermissionName = common.StringToName("prod.minor")

var RateLimitingPrecision uint32 = 1000 * 1000

var ActiveName uint64 = common.StringToName("active")

var ForkDBName = "forkdb.dat"
var DBFileName = "shared_memory.bin"
var ReversibleFileName = "shared_memory.bin"
var BlockFileName = "blog.log"
var DefaultBlocksDirName = "blocks"
var DefaultStateDirName = "state"
var DefaultStateSize uint64 = 0
var DefaultStateGuardSize uint64 = 0
var DefaultReversibleCacheSize uint64 = 0
var DefaultReversibleGuardSize uint64 = 0

var DefaultWasmRuntime = exec.WasmInterface{}
