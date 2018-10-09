package common

var DefaultConfig Config

type billableSize struct {
	overhead uint64
	value    uint64
}

func init() {
	DefaultConfig = Config{
		SystemAccountName:    AccountName(N("eosio")),
		NullAccountName:      AccountName(N("eosio.null")),
		ProducersAccountName: AccountName(N("eosio.prods")),

		MajorityProducersPermissionName: AccountName(N("prod.major")),
		MinorityProducersPermissionName: AccountName(N("prod.minor")),

		EosioAuthScope: AccountName(N("eosio.auth")),
		EosioAllScope:  AccountName(N("eosio.all")),
		ActiveName:     AccountName(N("active")),
		OwnerName:      AccountName(N("owner")),
		EosioAnyName:   AccountName(N("eosio.any")),
		EosioCodeName:  AccountName(N("eosio.code")),

		RateLimitingPrecision: 1000 * 1000,

		BillableAlignment: 16,
		BillableSize: map[string]billableSize{
			"permission_level_weight": {value: 24},
			"key_weight":              {value: 8},
			"wait_weight":             {value: 16},
			"shared_authority":        {value: 3*1 + 4},
			"permission_link_object":  {overhead: 32 * 3, value: 40 + 32},
			"permission_object":       {overhead: 5 * 32, value: 3*1 + 4 + 64 + 5*32},
		},
		FixedNetOverheadOfPackedTrx: 16,
	}

	DefaultConfig.BlockIntervalMs = 500
	DefaultConfig.BlockIntervalUs = 1000 * DefaultConfig.BlockIntervalMs
	DefaultConfig.BlockTimestampEpochMs = 946684800000
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

	DefaultConfig.ForkDBName = "forkdb.dat"
	DefaultConfig.DBFileName = "shared_memory.bin"
	DefaultConfig.ReversibleFileName = "shared_memory_tmp.bin"			//wait db modify
	DefaultConfig.BlockFileName = "blog.log"
	DefaultConfig.DefaultBlocksDirName = "/tmp/data/blocks"
	DefaultConfig.DefaultReversibleBlocksDirName="reversible"
	DefaultConfig.DefaultStateDirName = "/tmp/data/state"
	DefaultConfig.DefaultStateSize = 0
	DefaultConfig.DefaultStateGuardSize = 0
	DefaultConfig.DefaultReversibleCacheSize = 0
	DefaultConfig.DefaultReversibleGuardSize = 0
}

type Config struct {
	SystemAccountName    AccountName
	NullAccountName      AccountName
	ProducersAccountName AccountName

	// Active permission of producers account requires greater than 2/3 of the producers to authorize
	MajorityProducersPermissionName AccountName
	MinorityProducersPermissionName AccountName

	EosioAuthScope AccountName
	EosioAllScope  AccountName
	ActiveName     AccountName
	OwnerName      AccountName
	EosioAnyName   AccountName
	EosioCodeName  AccountName

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
	DeferredTrxExpirationWindow uint32 ///< the number of seconds after the time a deferred transaction can first execute until it expires
	MaxTrxDelay                 uint32 ///< the maximum number of seconds that can be imposed as a delay requirement by authorization checks
	MaxInlineActionSize         uint32 ///< maximum allowed size (in bytes) of an inline action
	MaxInlineActionDepth        uint16 ///< recursion depth limit on sending inline actions
	MaxAuthorityDepth           uint16 ///< recursion depth limit for checking if an authority is satisfied
	/**************************chain_config end****************************/

	ForkDBName	string
	DBFileName string
	ReversibleFileName string
	BlockFileName string
	DefaultBlocksDirName string
	DefaultReversibleBlocksDirName string
	DefaultStateDirName string
	DefaultStateSize uint64
	DefaultStateGuardSize uint64
	DefaultReversibleCacheSize uint64
	DefaultReversibleGuardSize uint64
	//FixedNetOverheadOfPackedTrx uint32 // TODO: C++ default value 16 and is this reasonable?
}

func (c *Config) Validate() {
	/*EOS_ASSERT( target_block_net_usage_pct <= config::percent_100, action_validate_exception,
		"target block net usage percentage cannot exceed 100%" );
	EOS_ASSERT( target_block_net_usage_pct >= config::percent_1/10, action_validate_exception,
		"target block net usage percentage must be at least 0.1%" );
	EOS_ASSERT( target_block_cpu_usage_pct <= config::percent_100, action_validate_exception,
		"target block cpu usage percentage cannot exceed 100%" );
	EOS_ASSERT( target_block_cpu_usage_pct >= config::percent_1/10, action_validate_exception,
		"target block cpu usage percentage must be at least 0.1%" );

	EOS_ASSERT( max_transaction_net_usage < max_block_net_usage, action_validate_exception,
		"max transaction net usage must be less than max block net usage" );
	EOS_ASSERT( max_transaction_cpu_usage < max_block_cpu_usage, action_validate_exception,
		"max transaction cpu usage must be less than max block cpu usage" );

	EOS_ASSERT( base_per_transaction_net_usage < max_transaction_net_usage, action_validate_exception,
		"base net usage per transaction must be less than the max transaction net usage" );
	EOS_ASSERT( (max_transaction_net_usage - base_per_transaction_net_usage) >= config::min_net_usage_delta_between_base_and_max_for_trx,
		action_validate_exception,
		"max transaction net usage must be at least ${delta} bytes larger than base net usage per transaction",
		("delta", config::min_net_usage_delta_between_base_and_max_for_trx) );
	EOS_ASSERT( context_free_discount_net_usage_den > 0, action_validate_exception,
		"net usage discount ratio for context free data cannot have a 0 denominator" );
	EOS_ASSERT( context_free_discount_net_usage_num <= context_free_discount_net_usage_den, action_validate_exception,
		"net usage discount ratio for context free data cannot exceed 1" );

	EOS_ASSERT( min_transaction_cpu_usage <= max_transaction_cpu_usage, action_validate_exception,
		"min transaction cpu usage cannot exceed max transaction cpu usage" );
	EOS_ASSERT( max_transaction_cpu_usage < (max_block_cpu_usage - min_transaction_cpu_usage), action_validate_exception,
		"max transaction cpu usage must be at less than the difference between the max block cpu usage and the min transaction cpu usage" );

	EOS_ASSERT( 1 <= max_authority_depth, action_validate_exception,
		"max authority depth should be at least 1" );*/
}

func BillableSizeV(kind string) uint64 {
	return (DefaultConfig.BillableSize[kind].value + DefaultConfig.BillableAlignment - 1) / DefaultConfig.BillableAlignment * DefaultConfig.BillableAlignment
}

func EosPercent(value uint64, percentage uint32) uint64 {
	return (value * uint64(percentage)) / uint64(DefaultConfig.Percent_100)
}

const ()
