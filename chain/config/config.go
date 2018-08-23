package config

import (
	"github.com/eosspark/eos-go/common"
)

var SystemAccountName    = common.StringToName("eosio")
var NullAccountName      = common.StringToName("eosio.null")
var ProducersAccountName = common.StringToName("eosio.prods")

// Active permission of producers account requires greater than 2/3 of the producers to authorize
var MajorityProducersPermissionName = common.StringToName("prod.major")
var MinorityProducersPermissionName = common.StringToName("prod.minor")




var RateLimitingPrecision uint32 = 1000 * 1000