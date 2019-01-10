package exception

/*
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

//TODO: go get gotemplate
//go:generate go install github.com/eosspark/eos-go/log/...
//go:generate go install github.com/eosspark/eos-go/exception/...
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "StdException(Exception,StdExceptionCode,\"golang standard error\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "FcException(Exception,UnspecifiedExceptionCode,\"unspecified\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnHandledException(Exception,UnhandledExceptionCode,\"unhandled\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TimeoutException(Exception,TimeoutExceptionCode,\"Timeout\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "FileNotFoundException(Exception,FileNotFoundExceptionCode,\"File Not Found\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ParseErrorException(Exception,ParseErrorExceptionCode,\"Parse Error\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidArgException(Exception,InvalidArgExceptionCode,\"Invalid Argument\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "KeyNotFoundException(Exception,KeyNotFoundExceptionCode,\"Key Not Found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BadCastException(Exception,BadCastExceptionCode,\"Bad Cast\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "OutOfRangeException(Exception,OutOfRangeExceptionCode,\"Out of Range\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidOperationException(Exception,InvalidOperationExceptionCode,\"Invalid Operation\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnknownHostException(Exception,UnknownHostExceptionCode,\"Unknown Host\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "CanceledException(Exception,CanceledExceptionCode,\"Canceled\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AssertException(Exception,AssertExceptionCode,\"Assert Exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "EofException(Exception,EofExceptionCode,\"End Of File\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "NullOptional(Exception,NullOptionalCode,\"null optionale\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UdtException(Exception,UdtErrorCode,\"UDT error\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AesException(Exception,UdtErrorCode,\"AES error\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "OverflowException(Exception,OverflowCode,\"Integer Overflow\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnderflowException(Exception,UnderflowCode,\"Integer Underflow\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DivideByZeroException(Exception,DivideByZeroCode,\"Integer Divide By Zero\")"

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ChainException(_ChainException,3000000,\"blockchain exception\")"

//_ChainTypeException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ChainTypeException(_ChainException,3010000,\"chain type exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "NameTypeException(_ChainException,3010001,\"Invalid name\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PublicKeyTypeException(_ChainException,3010002,\"Invalid public key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PrivateKeyTypeException(_ChainException,3010003,\"Invalid private key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AuthorityTypeException(_ChainException,3010004,\"Invalid authority\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActionTypeException(_ChainException,3010005,\"Invalid action\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TransactionTypeException(_ChainException,3010006,\"Invalid transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiTypeException(_ChainException,3010007,\"Invalid ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockIdTypeException(_ChainException,3010008,\"Invalid block ID\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TransactionIdTypeException(_ChainException,3010009,\"Invalid transaction ID\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PackedTransactionTypeException(_ChainException,3010010,\"Invalid packed transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AssetTypeException(_ChainException,3010011,\"Invalid asset\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ChainIdTypeException(_ChainException,3010012,\"Invalid chain ID\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "FixedKeyTypeException(_ChainException,3010013,\"Invalid fixed key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "SymbolTypeException(_ChainException,3010014,\"Invalid symbol\")"

//_ForkDatabaseException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ForkDatabaseException(_ForkDatabaseException,3020000,\"Fork database exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ForkDbBlockNotFound  (_ForkDatabaseException,3020001,\"Block can not be found\")"

//_BlockValidateException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockValidateException  (_BlockValidateException,3030000,\"Block exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnlinkableBlockException(_BlockValidateException,3030001,\"Unlinkable block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockTxOutputException(_BlockValidateException,3030002,\"Transaction outputs in block do not match transaction outputs from applying block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockConcurrencyException(_BlockValidateException,3030003,\"Block does not guarantee concurrent execution without conflicts\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLockException(_BlockValidateException,3030004,\"Shard locks in block are incorrect or mal-formed\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockResourceExhausted(_BlockValidateException,3030005,\"Block exhausted allowed resources\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockTooOldException(_BlockValidateException,3030006,\"Block is too old to push\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockFromTheFuture(_BlockValidateException,3030007,\"Block is from the future\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WrongSigningKey(_BlockValidateException,3030008,\"Block is not signed with expected key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WrongProducer(_BlockValidateException,3030009,\"Block is not signed by expected producer\")"

//_TransactionException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TransactionException (_TransactionException,3040000,\"Transaction exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxDecompressionError (_TransactionException,3040001,\"Error decompressing transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxNoAction (_TransactionException,3040002,\"Transaction should have at least one normal action\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxNoAuths (_TransactionException,3040003,\"Transaction should have at least one required authority\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "CfaIrrelevantAuth (_TransactionException,3040004,\"Context-free action should have no required authority\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ExpiredTxException (_TransactionException,3040005,\"Expired Transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxExpTooFarException (_TransactionException,3040006,\"Transaction Expiration Too Far\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidRefBlockException (_TransactionException,3040007,\"Invalid Reference Block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxDuplicate (_TransactionException,3040008,\"Duplicate transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DeferredTxDuplicate (_TransactionException,3040009,\"Duplicate deferred transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "CfaInsideGeneratedTx (_TransactionException,3040010,\"Context free action is not allowed inside generated transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxNotFound (_TransactionException,3040011,\"The transaction can not be found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TooManyTxAtOnce (_TransactionException,3040012,\"Pushing too many transactions at once\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxTooBig (_TransactionException,3040013,\"Transaction is too big\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnknownTransactionCompression (_TransactionException,3040014,\"Unknown transaction compression\")"

//_ActionValidateException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActionValidateException (_ActionValidateException,3050000,\"Transaction exceeded the current CPU usage limit imposed on the transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AccountNameExistsException (_ActionValidateException,3050001,\"Account name already exists\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidActionArgsException (_ActionValidateException,3050002,\"Invalid Action Arguments\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "EosioAssertMessageException (_ActionValidateException,3050003,\"eosio_assert_message assertion failure\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "EosioAssertCodeException (_ActionValidateException,3050004,\"eosio_assert_code assertion failure\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActionNotFoundException (_ActionValidateException,3050005,\"Action can not be found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActionDataAndStructMismatch (_ActionValidateException,3050006,\"Mismatch between action data and its struct\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnaccessibleApi (_ActionValidateException,3050007,\"Attempt to use unaccessible API\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbortCalled (_ActionValidateException,3050008,\"Abort Calle\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InlineActionTooBig (_ActionValidateException,3050009,\"Inline Action exceeds maximum size limit\")"

//_DatabaseException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DatabaseException (_DatabaseException,3060000,\"Database exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PermissionQueryException (_DatabaseException,3060001,\"Permission Query Exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AccountQueryException (_DatabaseException,3060002,\"Account Query Exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractTableQueryException (_DatabaseException,3060003,\"Contract Table Query Exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractQueryException (_DatabaseException,3060004,\"Contract Query Exception\")"
////_GuardException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "GuardException (_GuardException,3060100,\"Database exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DatabaseGuardException (_GuardException,3060101,\"Database usage is at unsafe levels\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ReversibleGuardException (_GuardException,3060102,\"Reversible block log usage is at unsafe levels\")"

//_WasmException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WasmException (_WasmException,3070000,\"WASM Exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PageMemoryError (_WasmException,3070001,\"Error in WASM page memory\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WasmExecutionError (_WasmException,3070002,\"Runtime Error Processing WASM\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WasmSerializationError (_WasmException,3070003,\"Serialization Error Processing WASM\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "OverlappingMemoryError (_WasmException,3070004,\"memcpy with overlapping memory\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BinaryenException (_WasmException,3070005,\"binaryen exception\")"

//_ResourceExhaustedException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ResourceExhaustedException (_ResourceExhaustedException,3080000,\"Resource exhausted exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "RamUsageExceeded (_ResourceExhaustedException,3080001,\"Account using more than allotted RAM usage\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxNetUsageExceeded (_ResourceExhaustedException,3080002,\"Transaction exceeded the current network usage limit imposed on the transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockNetUsageExceeded (_ResourceExhaustedException,3080003,\"Transaction network usage is too much for the remaining allowable usage of the current block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxCpuUsageExceeded (_ResourceExhaustedException,3080004,\"Transaction exceeded the current CPU usage limit imposed on the transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockCpuUsageExceeded (_ResourceExhaustedException,3080005,\"Transaction CPU usage is too much for the remaining allowable usage of the current block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "GreylistNetUsageExceeded (_ResourceExhaustedException,3080007,\"Transaction exceeded the current greylisted account network usage limit\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "GreylistCpuUsageExceeded (_ResourceExhaustedException,3080008,\"Transaction exceeded the current greylisted account CPU usage limit\")"
////_DeadlineException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DeadlineException (_DeadlineException,3080006,\"Transaction exceeded the current greylisted account CPU usage limit\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "LeewayDeadlineException (_DeadlineException,3081001,\"Transaction reached the deadline set due to leeway on account CPU limits\")"

//_AuthorizationException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AuthorizationException (_AuthorizationException,3090000,\"Authorization exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxDuplicateSig (_AuthorizationException,3090001,\"Duplicate signature included\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TxIrrelevantSig (_AuthorizationException,3090002,\"Irrelevant signature included\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnsatisfiedAuthorization (_AuthorizationException,3090003,\"Provided keys, permissions, and delays do not satisfy declared authorizations\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingAuthException (_AuthorizationException,3090004,\"Missing required authority\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "IrrelevantAuthException (_AuthorizationException,3090005,\"Irrelevant authority included\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InsufficientDelayException (_AuthorizationException,3090006,\"Insufficient delay\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidPermission (_AuthorizationException,3090007,\"Invalid Permission\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnlinkableMinPermissionAction (_AuthorizationException,3090008,\"The action is not allowed to be linked with minimum permission\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidParentPermission (_AuthorizationException,3090009,\"The parent permission is invalid\")"

//_MiscException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MiscException (_MiscException,3100000,\"Miscellaneous exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "RateLimitingStateInconsistent (_MiscException,3100001,\"Internal state is no longer consistent\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnknownBlockException (_MiscException,3100002,\"Unknown block\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnknownTransactionException (_MiscException,3100003,\"Unknown transaction\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "FixedReversibleDbException (_MiscException,3100004,\"Corrupted reversible block database was fixed\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ExtractGenesisStateException (_MiscException,3100005,\"Extracted genesis state from blocks.log\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "SubjectiveBlockProductionException (_MiscException,3100006,\"Subjective exception thrown during block production\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MultipleVoterInfo (_MiscException,3100007,\"Multiple voter info detected\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnsupportedFeature (_MiscException,3100008,\"Feature is currently unsupported\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "NodeManagementSuccess (_MiscException,3100009,\"Node management operation successfully executed\")"

//_PluginException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PluginException (_PluginException,3110000,\"Plugin exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingChainApiPluginException (_PluginException,3110001,\"Missing Chain API Plugin\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingWalletApiPluginException (_PluginException,3110002,\"Missing Wallet API Plugin\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingHistoryApiPluginException (_PluginException,3110003,\"Missing History API Plugin\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingNetApiPluginException (_PluginException,3110004,\"Missing Net API Plugin\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingChainPluginException (_PluginException,3110005,\"Missing Chain Plugin\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PluginConfigException (_PluginException,3110006,\"Incorrect plugin configuration\")"

//_WalletException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletException (_WalletException,3120000,\"Invalid contract vm version\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletExistException (_WalletException,3120001,\"Wallet already exists\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletNonexistentException (_WalletException,3120002,\"Nonexistent wallet\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletLockedException (_WalletException,3120003,\"Locked wallet\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletMissingPubKeyException (_WalletException,3120004,\"Missing public key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletInvalidPasswordException (_WalletException,3120005,\"Invalid wallet password\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletNotAvailableException (_WalletException,3120006,\"No available wallet\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WalletUnlockedException (_WalletException,3120007,\"Already unlocked\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "KeyExistException (_WalletException,3120008,\"Key already exists\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "KeyNonexistentException (_WalletException,3120009,\"Nonexistent key\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnsupportedKeyTypeException (_WalletException,3120010,\"Unsupported key type\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidLockTimeoutException (_WalletException,3120011,\"Wallet lock timeout is invalid\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "SecureEnclaveException (_WalletException,3120012,\"Secure Enclave Exception\")"

//_WhitelistBlacklistException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WhitelistBlacklistException (_WhitelistBlacklistException,3130000,\"Actor or contract whitelist/blacklist exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActorWhitelistException (_WhitelistBlacklistException,3130001,\"Authorizing actor of transaction is not on the whitelist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActorBlacklistException (_WhitelistBlacklistException,3130002,\"Authorizing actor of transaction is on the blacklist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractWhitelistException (_WhitelistBlacklistException,3130003,\"Contract to execute is not on the whitelist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractBlacklistException (_WhitelistBlacklistException,3130004,\"Contract to execute is on the blacklist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ActionBlacklistException (_WhitelistBlacklistException,3130005,\"Action to execute is on the blacklist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "KeyBlacklistException (_WhitelistBlacklistException,3130006,\"Public key in authority is on the blacklist\")"

//_ControllerEmitSignalException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ControllerEmitSignalException (_ControllerEmitSignalException,3140000,\"Exceptions that are allowed to bubble out of emit calls in controller\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "CheckpointException (_ControllerEmitSignalException,3140001,\"Block does not match checkpoint\")"

//_AbiException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiException (_AbiException,3150000,\"ABI exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiNotFoundException (_AbiException,3150001,\"No ABI Found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidRicardianClauseException (_AbiException,3150002,\"Invalid Ricardian Clause\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidActionClauseException (_AbiException,3150003,\"Invalid Ricardian Action\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidTypeInsideAbi (_AbiException,3150004,\"The type defined in the ABI is invalid\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiTypeDefException (_AbiException,3150005,\"Duplicate type definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiStructDefException (_AbiException,3150006,\"Duplicate struct definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiActionDefException (_AbiException,3150007,\"Duplicate action definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiTableDefException (_AbiException,3150008,\"Duplicate table definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiErrMsgDefException (_AbiException,3150009,\"Duplicate error message definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiSerializationDeadlineException (_AbiException,3150010,\"ABI serialization time has exceeded the deadline\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiRecursionDepthException (_AbiException,3150011,\"ABI recursive definition has exceeded the max recursion depth\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiCircularDefException (_AbiException,3150012,\"Circular definition is detected in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnpackException (_AbiException,3150013,\"Unpack data exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "PackException (_AbiException,3150014,\"Pack data exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DuplicateAbiVariantDefException (_AbiException,3150015,\"Duplicate variant definition in the ABI\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "UnsupportedAbiVersionException (_AbiException,3150016,\"ABI has an unsupported version\")"

//_ContractException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractException (_ContractException,3160000,\"Contract exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidTablePayer (_ContractException,3160001,\"The payer of the table data is invalid\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TableAccessViolation (_ContractException,3160002,\"Table access violation\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidTableIterator (_ContractException,3160003,\"Invalid table iterator\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TableNotInCache (_ContractException,3160004,\"Table can not be found inside the cache\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "TableOperationNotPermitted (_ContractException,3160005,\"The table operation is not allowed\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidContractVmType (_ContractException,3160006,\"Invalid contract vm type\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidContractVmVersion (_ContractException,3160007,\"Invalid contract vm version\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "SetExactCode (_ContractException,3160008,\"Contract is already running this version of code\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "WastFileNotFound (_ContractException,3160009,\"No wast file found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "AbiFileNotFound (_ContractException,3160010,\"No abi file found\")"

//_ProducerException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ProducerException (_ProducerException,3170000,\"Producer exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ProducerPrivKeyNotFound (_ProducerException,3170001,\"Producer private key is not available\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "MissingPendingBlockState (_ProducerException,3170002,\"Pending block state is missing\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ProducerDoubleConfirm (_ProducerException,3170003,\"Producer is double confirming known range\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ProducerScheduleException (_ProducerException,3170004,\"Producer schedule exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ProducerNotInSchedule (_ProducerException,3170006,\"The producer is not part of current schedule\")"

//_ReversibleBlocksException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ReversibleBlocksException (_ReversibleBlocksException,3180000,\"Reversible Blocks exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidReversibleBlocksDir (_ReversibleBlocksException,3180001,\"Invalid reversible blocks directory\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ReversibleBlocksBackupDirExist (_ReversibleBlocksException,3180002,\"Backup directory for reversible blocks already exist\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "GapInReversibleBlocksDb (_ReversibleBlocksException,3180003,\"Gap in the reversible blocks database\")"

//_BlockLogException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLogException (_BlockLogException,3190000,\"Block log exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLogUnsupportedVersion (_BlockLogException,3190001,\"unsupported version of block log\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLogAppendFail (_BlockLogException,3190002,\"fail to append block to the block log\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLogNotFound (_BlockLogException,3190003,\"block log can not be found\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "BlockLogBackupDirExist (_BlockLogException,3190004,\"block log backup dir already exists\")"

//_HttpException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "HttpException (_HttpException,3200000,\"http exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidHttpClientRootCert (_HttpException,3200001,\"invalid http client root certificate\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidHttpResponse (_HttpException,3200002,\"invalid http response\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ResolvedToMultiplePorts (_HttpException,3200003,\"service resolved to multiple ports\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "FailToResolveHost (_HttpException,3200004,\"fail to resolve host\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "HttpRequestFail (_HttpException,3200005,\"http request fail\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "InvalidHttpRequest (_HttpException,3200006,\"invalid http request\")"

//_ResourceLimitException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ResourceLimitException (_ResourceLimitException,3210000,\"Resource limit exception\")"

//_MongoDbException

//_ContractApiException
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ContractApiException (_ContractApiException,3230000,\"Contract API exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "CryptoApiException (_ContractApiException,3230001,\"Crypto API exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "DbApiException (_ContractApiException,3230002,\"Database API exception\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ArithmeticException (_ContractApiException,3230003,\"Arithmetic exception\")"

//go:generate go build .

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "ExplainedException(Exception,9000000,\"explained exception,see error log\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/exception/template" "LocalizedException(Exception,10000000,\"an error occured\")"

func registerException(e Exception) {}
