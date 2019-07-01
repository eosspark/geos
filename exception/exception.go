package exception

import "github.com/eosspark/eos-go/log"

type ExcTypes = int64

const (
	UnspecifiedExceptionCode = ExcTypes(iota)
	UnhandledExceptionCode   ///< for unhandled 3rd party exceptions
	TimeoutExceptionCode     ///< timeout exceptions
	FileNotFoundExceptionCode
	ParseErrorExceptionCode
	InvalidArgExceptionCode
	KeyNotFoundExceptionCode
	BadCastExceptionCode
	OutOfRangeExceptionCode
	CanceledExceptionCode
	AssertExceptionCode
	_
	EofExceptionCode
	StdExceptionCode
	InvalidOperationExceptionCode
	UnknownHostExceptionCode
	NullOptionalCode
	UdtErrorCode
	AesErrorCode
	OverflowCode
	UnderflowCode
	DivideByZeroCode
)

type Exception interface {
	Code() int64
	Name() string
	What() string

	AppendLog(l log.Message)
	GetLog() log.Messages

	TopMessage() string
	DetailMessage() string
	String() string
	MarshalJSON() ([]byte, error)

	Callback(f interface{}) bool
}

type ChainExceptions interface {
	Exception
	ChainExceptions()
}

//using empty method implements from interface in order to upcasting
type _ChainException struct{}

func (_ChainException) ChainExceptions() {}

/**
 * 	chain_type_exception
 */
type ChainTypeExceptions interface {
	ChainExceptions
	ChainTypeExceptions()
}

type _ChainTypeException struct{ _ChainException }

func (_ChainTypeException) ChainTypeExceptions() {}

/**
 * fork_database_exception
 */
type ForkDatabaseExceptions interface {
	ChainExceptions
	ForkDatabaseExceptions()
}

type _ForkDatabaseException struct{ _ChainException }

func (_ForkDatabaseException) ForkDatabaseExceptions() {}

/**
 * 	block_validate_exception
 */
type BlockValidateExceptions interface {
	ChainExceptions
	BlockValidateExceptions()
}

type _BlockValidateException struct{ _ChainException }

func (_BlockValidateException) BlockValidateExceptions() {}

/**
 * transaction_exception
 */
type TransactionExceptions interface {
	ChainExceptions
	TransactionExceptions()
}

type _TransactionException struct{ _ChainException }

func (_TransactionException) TransactionExceptions() {}

/**
 * action_validate_exception
 */
type ActionValidateExceptions interface {
	ChainExceptions
	ActionValidateExceptions()
}

type _ActionValidateException struct{ _ChainException }

func (_ActionValidateException) ActionValidateExceptions() {}

/**
 * database_exception
 */
type DatabaseExceptions interface {
	ChainExceptions
	DatabaseExceptions()
}

type _DatabaseException struct{ _ChainException }

func (_DatabaseException) DatabaseExceptions() {}

type GuardExceptions interface {
	DatabaseExceptions
	GuardExceptions()
}

type _GuardException struct {
	_ChainException
	_DatabaseException
}

func (_GuardException) GuardExceptions() {}

/**
 * wasm_exception
 */
type WasmExceptions interface {
	ChainExceptions
	WasmExceptions()
}

type _WasmException struct{ _ChainException }

func (_WasmException) WasmExceptions() {}

/**
 * resource_exhausted_exception
 */
type ResourceExhaustedExceptions interface {
	ChainExceptions
	ResourceExhaustedExceptions()
}

type _ResourceExhaustedException struct{ _ChainException }

func (_ResourceExhaustedException) ResourceExhaustedExceptions() {}

type DeadlineExceptions interface {
	ResourceExhaustedExceptions
	DeadlineExceptions()
}

type _DeadlineException struct {
	_ChainException
	_ResourceExhaustedException
}

func (_DeadlineException) DeadlineExceptions() {}

/**
 * authorization_exception
 */
type AuthorizationExceptions interface {
	ChainExceptions
	AuthorizationExceptions()
}

type _AuthorizationException struct{ _ChainException }

func (_AuthorizationException) AuthorizationExceptions() {}

/**
 * misc_exception
 */
type MiscExceptions interface {
	ChainExceptions
	MiscExceptions()
}

type _MiscException struct{ _ChainException }

func (_MiscException) MiscExceptions() {}

/**
 * plugin_exception
 */
type PluginExceptions interface {
	ChainExceptions
	PluginExceptions()
}

type _PluginException struct{ _ChainException }

func (_PluginException) PluginExceptions() {}

/**
 * wallet_exception
 */
type WalletExceptions interface {
	ChainExceptions
	WalletExceptions()
}

type _WalletException struct{ _ChainException }

func (_WalletException) WalletExceptions() {}

/**
 * whitelist_blacklist_exception
 */

type WhitelistBlacklistExceptions interface {
	ChainExceptions
	WhitelistBlacklistExceptions()
}

type _WhitelistBlacklistException struct{ _ChainException }

func (_WhitelistBlacklistException) WhitelistBlacklistExceptions() {}

/**
 * controller_emit_signal_exception
 */
type ControllerEmitSignalExceptions interface {
	ChainExceptions
	ControllerEmitSignalExceptions()
}

type _ControllerEmitSignalException struct{ _ChainException }

func (_ControllerEmitSignalException) ControllerEmitSignalExceptions() {}

/**
 * abi_exception
 */
type AbiExceptions interface {
	ChainExceptions
	AbiExceptions()
}

type _AbiException struct{ _ChainException }

func (_AbiException) AbiExceptions() {}

/**
 * contract_exception
 */
type ContractExceptions interface {
	ChainExceptions
	ContractExceptions()
}

type _ContractException struct{ _ChainException }

func (_ContractException) ContractExceptions() {}

/**
 * producer_exception
 */
type ProducerExceptions interface {
	ChainExceptions
	ProducerExceptions()
}

type _ProducerException struct{ _ChainException }

func (_ProducerException) ProducerExceptions() {}

/**
 * reversible_blocks_exception
 */
type ReversibleBlocksExceptions interface {
	ChainExceptions
	ReversibleBlocksExceptions()
}

type _ReversibleBlocksException struct{ _ChainException }

func (_ReversibleBlocksException) ReversibleBlocksExceptions() {}

/**
 * block_log_exception
 */
type BlockLogExceptions interface {
	ChainExceptions
	BlockLogExceptions()
}

type _BlockLogException struct{ _ChainException }

func (_BlockLogException) BlockLogExceptions() {}

/**
 * http_exception
 */
type HttpExceptions interface {
	ChainExceptions
	HttpExceptions()
}

type _HttpException struct{ _ChainException }

func (_HttpException) HttpExceptions() {}

/**
 * resource_limit_exception
 */
type ResourceLimitExceptions interface {
	ChainExceptions
	ResourceLimitExceptions()
}

type _ResourceLimitException struct{ _ChainException }

func (_ResourceLimitException) ResourceLimitExceptions() {}

/**
 * mongo_db_exception
 */
type MongoDbExceptions interface {
	ChainExceptions
	MongoDbExceptions()
}

type _MongoDbException struct{ _ChainException }

func (_MongoDbException) MongoDbExceptions() {}

/**
 * contract_api_exception
 */
type ContractApiExceptions interface {
	ChainExceptions
	ContractApiExceptions()
}

type _ContractApiException struct{ _ChainException }

func (_ContractApiException) ContractApiExceptions() {}
