package exception

type ChainExceptions interface {
	Exception
	ChainExceptions()
}

type ChainException struct{ logMessage }

func (e *ChainException) ChainExceptions() {}
func (e *ChainException) Code() ExcTypes   { return 3000000 }
func (e *ChainException) What() string     { return "blockchain exception" }

/**
 *  chain_exception
 *   |- chain_type_exception
 *   |- fork_database_exception
 *   |- block_validate_exception
 *   |- transaction_exception
 *   |- action_validate_exception
 *   |- database_exception
 *   |- wasm_exception
 *   |- resource_exhausted_exception
 *   |- misc_exception
 *   |- plugin_exception
 *   |- wallet_exception
 *   |- whitelist_blacklist_exception
 *   |- controller_emit_signal_exception
 *   |- abi_exception
 *   |- contract_exception
 *   |- producer_exception
 *   |- reversible_blocks_exception
 *   |- block_log_exception
 *   |- resource_limit_exception
 *   |- mongo_db_exception
 *   |- contract_api_exception
 */

/**
 * 	chain_type_exception
 */
type ChainTypeExceptions interface {
	ChainExceptions
	ChainTypeExceptions()
}

type ChainTypeException struct{ logMessage }

func (e *ChainTypeException) ChainExceptions()     {}
func (e *ChainTypeException) ChainTypeExceptions() {}
func (e *ChainTypeException) Code() ExcTypes       { return 3010000 }
func (e *ChainTypeException) What() string         { return "chain type exception" }

type NameTypeException struct{ logMessage }

func (e *NameTypeException) ChainExceptions()     {}
func (e *NameTypeException) ChainTypeExceptions() {}
func (e *NameTypeException) Code() ExcTypes       { return 3010001 }
func (e *NameTypeException) What() string         { return "Invalid name" }

/**
 * fork_database_exception
 */
type ForkDatabaseExceptions interface {
	ChainExceptions
	ForkDatabaseExceptions()
}

type ForkDatabaseException struct{ logMessage }

func (e *ForkDatabaseException) ChainExceptions()        {}
func (e *ForkDatabaseException) ForkDatabaseExceptions() {}
func (e *ForkDatabaseException) Code() ExcTypes          { return 3020000 }
func (e *ForkDatabaseException) What() string            { return "Fork database exception" }

type ForkDbBlockNotFound struct{ logMessage }

func (e *ForkDbBlockNotFound) ChainExceptions()        {}
func (e *ForkDbBlockNotFound) ForkDatabaseExceptions() {}
func (e *ForkDbBlockNotFound) Code() ExcTypes          { return 3020001 }
func (e *ForkDbBlockNotFound) What() string            { return "Block can not be found" }

/**
 * 	block_validate_exception
 */
type BlockValidateExceptions interface {
	ChainExceptions
	BlockValidateExceptions()
}

type BlockValidateException struct{ logMessage }

func (e *BlockValidateException) ChainExceptions()         {}
func (e *BlockValidateException) BlockValidateExceptions() {}
func (e *BlockValidateException) Code() ExcTypes           { return 3030000 }
func (e *BlockValidateException) What() string             { return "Block exception" }

type DatabaseExceptions interface {
	ChainExceptions
	DatabaseExceptions()
}

type DatabaseException struct{ logMessage }

func (e *DatabaseException) ChainExceptions()    {}
func (e *DatabaseException) DatabaseExceptions() {}
func (e *DatabaseException) Code() ExcTypes      { return 3060000 }
func (e *DatabaseException) What() string        { return "Database exception" }

type GuardExceptions interface {
	DatabaseExceptions
	GuardExceptions()
}
