package common

import (
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
)

var DefaultConfig Config

type billableSize struct {
	overhead uint64
	value    uint64
}

func init() {
	DefaultConfig.SystemAccountName = AccountName(N("eosio"))
	DefaultConfig.NullAccountName = AccountName(N("eosio.null"))
	DefaultConfig.ProducersAccountName = AccountName(N("eosio.prods"))
	DefaultConfig.MajorityProducersPermissionName = PermissionName(N("prod.major"))
	DefaultConfig.MinorityProducersPermissionName = PermissionName(N("prod.minor"))

	DefaultConfig.EosioAuthScope = AccountName(N("eosio.auth"))
	DefaultConfig.EosioAllScope = AccountName(N("eosio.all"))
	DefaultConfig.ActiveName = PermissionName(N("active"))
	DefaultConfig.OwnerName = PermissionName(N("owner"))
	DefaultConfig.EosioAnyName = PermissionName(N("eosio.any"))
	DefaultConfig.EosioCodeName = PermissionName(N("eosio.code"))

	DefaultConfig.RateLimitingPrecision = 1000 * 1000

	DefaultConfig.BillableAlignment = 16
	DefaultConfig.BillableSize = map[string]billableSize{
		"permission_level_weight":      {value: 24},
		"key_weight":                   {value: 8},
		"wait_weight":                  {value: 16},
		"shared_authority":             {value: 3*1 + 4},
		"permission_link_object":       {overhead: 32 * 3, value: 40 + 32},
		"permission_object":            {overhead: 5 * 32, value: 3*1 + 4 + 64 + 5*32},
		"table_id_object":              {overhead: 32 * 2, value: 44 + 32*2},
		"key_value_object":             {overhead: 32 * 2, value: 32 + 8 + 4 + 32*2},
		"index64_object":               {overhead: 32 * 3, value: 24 + 8 + 32*3},
		"index128_object":              {overhead: 32 * 3, value: 24 + 16 + 32*3},
		"index256_object":              {overhead: 32 * 3, value: 24 + 32 + 32*3},
		"index_double_object":          {overhead: 32 * 3, value: 24 + 8 + 32*3},
		"index_long_double_object":     {overhead: 32 * 3, value: 24 + 16 + 32*3},
		"generated_transaction_object": {overhead: 32 * 5, value: 96 + 4 + 32*5},
	}

	DefaultConfig.FixedNetOverheadOfPackedTrx = 16

	DefaultConfig.BlockIntervalMs = 500
	DefaultConfig.BlockIntervalUs = 1000 * DefaultConfig.BlockIntervalMs
	DefaultConfig.BlockTimestampEpochMs = 946684800000 // epoch is year 2000.
	DefaultConfig.BlockTimestamoEpochNanos = 1e6 * DefaultConfig.BlockTimestampEpochMs

	DefaultConfig.ProducerRepetitions = 12
	DefaultConfig.MaxProducers = 125
	DefaultConfig.MaxTrackedDposConfirmations = 1024

	DefaultConfig.Percent_100 = 10000
	DefaultConfig.Percent_1 = 100
	DefaultConfig.AccountCpuUsageAverageWindowMs = 24 * 60 * 60 * 1000
	DefaultConfig.AccountNetUsageAverageWindowMs = 24 * 60 * 60 * 1000
	DefaultConfig.BlockCpuUsageAverageWindowMs = 60 * 1000
	DefaultConfig.BlockSizeAverageWindowMs = 60 * 1000

	DefaultConfig.MaxBlockNetUsage = 1024 * 1024                                /// at 500ms blocks and 200byte trx, this enables ~10,000 TPS burst
	DefaultConfig.TargetBlockNetUsagePct = uint32(10 * DefaultConfig.Percent_1) /// we target 1000 TPS
	DefaultConfig.MaxTransactionNetUsage = uint32(DefaultConfig.MaxBlockNetUsage) / 2
	DefaultConfig.BasePerTransactionNetUsage = 12     // 12 bytes (11 bytes for worst case of transaction_receipt_header + 1 byte for static_variant tag)
	DefaultConfig.NetUsageLeeway = 500                // TODO: is this reasonable?
	DefaultConfig.ContextFreeDiscountNetUsageNum = 20 // TODO: is this reasonable?
	DefaultConfig.ContextFreeDiscountNetUsageDen = 100
	DefaultConfig.TransactionIdNetUsage = 32 // 32 bytes for the size of a transaction id

	DefaultConfig.MaxBlockCpuUsage = 200000 /// max block cpu usage in microseconds
	DefaultConfig.TargetBlockCpuUsagePct = uint32(10 * DefaultConfig.Percent_1)
	DefaultConfig.MaxTransactionCpuUsage = 3 * DefaultConfig.MaxBlockCpuUsage / 4 /// max trx cpu usage in microseconds
	DefaultConfig.MinTransactionCpuUsage = 100                                    /// min trx cpu usage in microseconds (10000 TPS equiv)

	DefaultConfig.MaxTrxLifetime = 60 * 60              // 1 hour
	DefaultConfig.DeferredTrxExpirationWindow = 10 * 60 // 10 minutes
	DefaultConfig.MaxTrxDelay = 45 * 24 * 3600          // 45 days
	DefaultConfig.MaxInlineActionSize = 4 * 1024        // 4 KB
	DefaultConfig.MaxInlineActionDepth = 4
	DefaultConfig.MaxAuthorityDepth = 6
	DefaultConfig.FixedNetOverheadOfPackedTrx = 16
	DefaultConfig.FixedOverheadSharedVectorRamBytes = 16
	DefaultConfig.OverheadPerRowPerIndexRamBytes = 32
	DefaultConfig.OverheadPerAccountRamBytes = 2 * 1024
	DefaultConfig.SetcodeRamBytesMultiplier = 10
	DefaultConfig.HashingChecktimeBlockSize = 10 * 1024

	DefaultConfig.ForkDbName = "forkdb.dat"
	DefaultConfig.DBFileName = "shared_memory.bin"
	DefaultConfig.ReversibleFileName = "shared_memory_tmp.bin" //wait db modify
	DefaultConfig.BlockFileName = "blog.log"
	DefaultConfig.DefaultBlocksDirName = "/tmp/data/blocks"
	DefaultConfig.DefaultReversibleBlocksDirName = "/tmp/data/reversible"
	DefaultConfig.DefaultStateDirName = "/tmp/data/state"
	DefaultConfig.DefaultStateSize = 1 * 1024 * 1024 * 1024
	DefaultConfig.DefaultStateGuardSize = 128 * 1024 * 1024
	DefaultConfig.DefaultReversibleCacheSize = 340 * 1024 * 1024
	DefaultConfig.DefaultReversibleGuardSize = 2 * 1024 * 1024
	DefaultConfig.MinNetUsageDeltaBetweenBaseAndMaxForTrx = 10 * 1024
}

type Config struct {
	SystemAccountName    AccountName
	NullAccountName      AccountName
	ProducersAccountName AccountName

	// Active permission of producers account requires greater than 2/3 of the producers to authorize
	MajorityProducersPermissionName PermissionName
	MinorityProducersPermissionName PermissionName

	EosioAuthScope AccountName
	EosioAllScope  AccountName
	ActiveName     PermissionName
	OwnerName      PermissionName
	EosioAnyName   PermissionName
	EosioCodeName  PermissionName

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

	FixedNetOverheadOfPackedTrx       uint32 //TODO: C++ default value 16 and is this reasonable?
	FixedOverheadSharedVectorRamBytes uint32
	OverheadPerRowPerIndexRamBytes    uint32 ///< overhead accounts for basic tracking structures in a row per index
	OverheadPerAccountRamBytes        uint32 //= 2*1024; ///< overhead accounts for basic account storage and pre-pays features like account recovery
	SetcodeRamBytesMultiplier         uint32 //= 10;     ///< multiplier on contract size to account for multiple copies and cached compilation

	HashingChecktimeBlockSize uint32 //= 10*1024;

	BillableAlignment uint64
	BillableSize      map[string]billableSize

	MaxTrackedDposConfirmations int ///<
	TransactionIdNetUsage       uint32

	Percent_100 int
	Percent_1   int

	AccountCpuUsageAverageWindowMs uint32
	AccountNetUsageAverageWindowMs uint32
	BlockCpuUsageAverageWindowMs   uint32
	BlockSizeAverageWindowMs       uint32

	/**************************chain_config start****************************/
	MaxBlockNetUsage               uint64 ///< the maxiumum net usage in instructions for a block
	TargetBlockNetUsagePct         uint32 ///< the target percent (1% == 100, 100%= 10,000) of maximum net usage; exceeding this triggers congestion handling
	MaxTransactionNetUsage         uint32 ///< the maximum objectively measured net usage that the chain will allow regardless of account limits
	BasePerTransactionNetUsage     uint32 ///< the base amount of net usage billed for a transaction to cover incidentals
	NetUsageLeeway                 uint32
	ContextFreeDiscountNetUsageNum uint32 ///< the numerator for the discount on net usage of context-free data
	ContextFreeDiscountNetUsageDen uint32 ///< the denominator for the discount on net usage of context-free data

	MaxBlockCpuUsage       uint32 ///< the maxiumum billable cpu usage (in microseconds) for a block
	TargetBlockCpuUsagePct uint32 ///< the target percent (1% == 100, 100%= 10,000) of maximum cpu usage; exceeding this triggers congestion handling
	MaxTransactionCpuUsage uint32 ///< the maximum billable cpu usage (in microseconds) that the chain will allow regardless of account limits
	MinTransactionCpuUsage uint32 ///< the minimum billable cpu usage (in microseconds) that the chain requires

	MaxTrxLifetime uint32
	//MaxTransactionLifetime      uint32 ///< the maximum number of seconds that an input transaction's expiration can be ahead of the time of the block in which it is first included
	DeferredTrxExpirationWindow             uint32 ///< the number of seconds after the time a deferred transaction can first execute until it expires
	MaxTrxDelay                             uint32 ///< the maximum number of seconds that can be imposed as a delay requirement by authorization checks
	MaxInlineActionSize                     uint32 ///< maximum allowed size (in bytes) of an inline action
	MaxInlineActionDepth                    uint16 ///< recursion depth limit on sending inline actions
	MaxAuthorityDepth                       uint16 ///< recursion depth limit for checking if an authority is satisfied
	MinNetUsageDeltaBetweenBaseAndMaxForTrx uint32
	/**************************chain_config end****************************/

	ForkDbName                     string
	DBFileName                     string
	ReversibleFileName             string
	BlockFileName                  string
	DefaultBlocksDirName           string
	DefaultReversibleBlocksDirName string
	DefaultStateDirName            string
	DefaultStateSize               uint64
	DefaultStateGuardSize          uint64
	DefaultReversibleCacheSize     uint64
	DefaultReversibleGuardSize     uint64
	//FixedNetOverheadOfPackedTrx uint32 // TODO: C++ default value 16 and is this reasonable?
}

func (c *Config) Validate() {
	try.EosAssert(c.TargetBlockNetUsagePct <= uint32(c.Percent_100), &exception.ActionValidateException{},
		"target block net usage percentage cannot exceed 100%")
	try.EosAssert(c.TargetBlockNetUsagePct >= uint32(c.Percent_1/10), &exception.ActionValidateException{},
		"target block net usage percentage must be at least 0.1%")
	try.EosAssert(c.TargetBlockCpuUsagePct <= uint32(c.Percent_100), &exception.ActionValidateException{},
		"target block cpu usage percentage cannot exceed 100%")
	try.EosAssert(c.TargetBlockCpuUsagePct >= uint32(c.Percent_1/10), &exception.ActionValidateException{},
		"target block cpu usage percentage must be at least 0.1%")

	try.EosAssert(uint64(c.MaxTransactionNetUsage) < c.MaxBlockNetUsage, &exception.ActionValidateException{},
		"max transaction net usage must be less than max block net usage")
	try.EosAssert(c.MaxTransactionCpuUsage < c.MaxBlockCpuUsage, &exception.ActionValidateException{},
		"max transaction cpu usage must be less than max block cpu usage")

	try.EosAssert(c.BasePerTransactionNetUsage < c.MaxTransactionNetUsage, &exception.ActionValidateException{},
		"base net usage per transaction must be less than the max transaction net usage")
	try.EosAssert((c.MaxTransactionNetUsage-c.BasePerTransactionNetUsage) >= c.MinNetUsageDeltaBetweenBaseAndMaxForTrx,
		&exception.ActionValidateException{},
		"max transaction net usage must be at least: %s bytes larger than base net usage per transaction",
		c.MinNetUsageDeltaBetweenBaseAndMaxForTrx)
	try.EosAssert(c.ContextFreeDiscountNetUsageDen > 0, &exception.ActionValidateException{},
		"net usage discount ratio for context free data cannot have a 0 denominator")
	try.EosAssert(c.ContextFreeDiscountNetUsageNum <= c.ContextFreeDiscountNetUsageDen, &exception.ActionValidateException{},
		"net usage discount ratio for context free data cannot exceed 1")

	try.EosAssert(c.MinTransactionCpuUsage <= c.MaxTransactionCpuUsage, &exception.ActionValidateException{},
		"min transaction cpu usage cannot exceed max transaction cpu usage")
	try.EosAssert(c.MaxTransactionCpuUsage < (c.MaxBlockCpuUsage-c.MinTransactionCpuUsage), &exception.ActionValidateException{},
		"max transaction cpu usage must be at less than the difference between the max block cpu usage and the min transaction cpu usage")

	try.EosAssert(1 <= c.MaxAuthorityDepth, &exception.ActionValidateException{},
		"max authority depth should be at least 1")
}

func BillableSizeV(kind string) uint64 {
	return (DefaultConfig.BillableSize[kind].value + DefaultConfig.BillableAlignment - 1) / DefaultConfig.BillableAlignment * DefaultConfig.BillableAlignment
}

func EosPercent(value uint64, percentage uint32) uint64 {
	return (value * uint64(percentage)) / uint64(DefaultConfig.Percent_100)
}

//const ()
