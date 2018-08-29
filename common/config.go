package common

var DefaultConfig Config

func init() {
	DefaultConfig = Config{
		SystemAccountName:    AccountName(StringToName("eosio")),
		NullAccountName:      AccountName(StringToName("eosio.null")),
		ProducersAccountName: AccountName(StringToName("eosio.prods")),

		MajorityProducersPermissionName: AccountName(StringToName("prod.major")),
		MinorityProducersPermissionName: AccountName(StringToName("prod.minor")),

		RateLimitingPrecision: 1000 * 1000,
	}

	DefaultConfig.BlockIntervalMs = 500
	DefaultConfig.BlockIntervalUs = 1000 * DefaultConfig.BlockIntervalMs
	DefaultConfig.BlockTimestampEpochMs = 946684800000
	DefaultConfig.BlockTimestamoEpochNanos = 1e6 * DefaultConfig.BlockTimestampEpochMs

	DefaultConfig.ProducerRepetitions = 12
	DefaultConfig.MaxProducers = 125
	DefaultConfig.MaxTrackedDposConfirmations = 1024
}

type Config struct {
	SystemAccountName    AccountName
	NullAccountName      AccountName
	ProducersAccountName AccountName

	// Active permission of producers account requires greater than 2/3 of the producers to authorize
	MajorityProducersPermissionName AccountName
	MinorityProducersPermissionName AccountName

	RateLimitingPrecision uint32

	BlockIntervalMs          int64
	BlockIntervalUs          int64
	BlockTimestampEpochMs    int64
	BlockTimestamoEpochNanos int64

	/**
	 *  The number of sequential blocks produced by a single producer
	 */
	ProducerRepetitions int
	MaxProducers        int

	MaxTrackedDposConfirmations int ///<
}
