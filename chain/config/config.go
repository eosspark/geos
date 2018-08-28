package config

import (
	"github.com/eosspark/eos-go/common"
)

var DefaultConfig = Config{
	SystemAccountName:    common.AccountName(common.StringToName("eosio")),
	NullAccountName:      common.AccountName(common.StringToName("eosio.null")),
	ProducersAccountName: common.AccountName(common.StringToName("eosio.prods")),

	MajorityProducersPermissionName: common.AccountName(common.StringToName("prod.major")),
	MinorityProducersPermissionName: common.AccountName(common.StringToName("prod.minor")),

	RateLimitingPrecision: 1000 * 1000,
}

type Config struct {
	SystemAccountName    common.AccountName
	NullAccountName      common.AccountName
	ProducersAccountName common.AccountName

	// Active permission of producers account requires greater than 2/3 of the producers to authorize
	MajorityProducersPermissionName common.AccountName
	MinorityProducersPermissionName common.AccountName

	RateLimitingPrecision uint32
}