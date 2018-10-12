package exception

type ChainExceptions interface {
	Exception
	ChainExceptions()
}

type ChainException struct{ logMessage }

func (ChainException) ChainExceptions() {}
func (ChainException) Code() ExcTypes   { return 3000000 }
func (ChainException) What() string     { return "blockchain exception" }

/**
 *  chain_exception
 *   |- chain_type_exception 			 >3010000
 *   |- fork_database_exception 		 >3020000
 *   |- block_validate_exception    	 >3030000
 *   |- transaction_exception	    	 >3040000
 *   |- action_validate_exception   	 >3050000
 *   |- database_exception				 >3060000
	     |- guard_exception	   >3060100
 *   |- wasm_exception					 >3070000
 *   |- resource_exhausted_exception 	 >3080000
		 |- deadline_exception >3081000
 *   |- authorization_exception		 	 >3090000
 *   |- misc_exception				 	 >3100000
 *   |- plugin_exception			 	 >3110000
 *   |- wallet_exception				 >3120000
 *   |- whitelist_blacklist_exception	 >3130000
 *   |- controller_emit_signal_exception >3140000
 *   |- abi_exception					 >3150000
 *   |- contract_exception				 >3160000
 *   |- producer_exception				 >3170000
 *   |- reversible_blocks_exception 	 >3180000
 *   |- block_log_exception				 >3190000
 *   |- http_exception					 >3200000
 *   |- resource_limit_exception		 >3210000
 *   |- mongo_db_exception 				 >3220000
 *   |- contract_api_exception  		 >3230000
 */

/**
 * 	chain_type_exception
 */
type ChainTypeExceptions interface {
	ChainExceptions
	ChainTypeExceptions()
}

/**
 * fork_database_exception
 */
type ForkDatabaseExceptions interface {
	ChainExceptions
	ForkDatabaseExceptions()
}

/**
 * 	block_validate_exception
 */
type BlockValidateExceptions interface {
	ChainExceptions
	BlockValidateExceptions()
}

/**
 * transaction_exception
 */
type TransactionExceptions interface {
	ChainExceptions
	TransactionExceptions()
}

/**
 * action_validate_exception
 */
type ActionValidateExceptions interface {
	ChainExceptions
	ActionValidateExceptions()
}

/**
 * database_exception
 */
type DatabaseExceptions interface {
	ChainExceptions
	DatabaseExceptions()
}

type GuardExceptions interface {
	DatabaseExceptions
	GuardExceptions()
}

/**
 * wasm_exception
 */
type WasmExceptions interface {
	ChainExceptions
	WasmExceptions()
}

/**
 * resource_exhausted_exception
 */
type ResourceExhaustedExceptions interface {
	ChainExceptions
	ResourceExhaustedExceptions()
}

type DeadlineExceptions interface {
	ResourceExhaustedExceptions
	DeadlineExceptions()
}

/**
 * authorization_exception
 */
type AuthorizationExceptions interface {
	ChainExceptions
	AuthorizationExceptions()
}

/**
 * misc_exception
 */
type MiscExceptions interface {
	ChainExceptions
	MiscExceptions()
}

/**
 * plugin_exception
 */
type PluginExceptions interface {
	ChainExceptions
	PluginExceptions()
}

/**
 * wallet_exception
 */
type WalletExceptions interface {
	ChainExceptions
	WalletExceptions()
}

/**
 * whitelist_blacklist_exception
 */

type WhitelistBlacklistExceptions interface {
	ChainExceptions
	WhitelistBlacklistExceptions()
}

/**
 * controller_emit_signal_exception
 */
type ControllerEmitSignalExceptions interface {
	ChainExceptions
	ControllerEmitSignalExceptions()
}

/**
 * abi_exception
 */
type AbiExceptions interface {
	ChainExceptions
	AbiExceptions()
}

/**
 * contract_exception
 */
type ContractExceptions interface {
	ChainExceptions
	ContractExceptions()
}

/**
 * producer_exception
 */
type ProducerExceptions interface {
	ChainExceptions
	ProducerExceptions()
}

/**
 * reversible_blocks_exception
 */
type ReversibleBlocksExceptions interface {
	ChainExceptions
	ReversibleBlocksExceptions()
}

/**
 * block_log_exception
 */
type BlockLogExceptions interface {
	ChainExceptions
	BlockLogExceptions()
}

/**
 * http_exception
 */
type HttpExceptions interface {
	ChainExceptions
	HttpExceptions()
}

/**
 * resource_limit_exception
 */
type ResourceLimitExceptions interface {
	ChainExceptions
	ResourceLimitExceptions()
}

/**
 * mongo_db_exception
 */
type MongoDbExceptions interface {
	ChainExceptions
	MongoDbExceptions()
}

/**
 * contract_api_exception
 */
type ContractApiExceptions interface {
	ChainExceptions
	ContractApiExceptions()
}
