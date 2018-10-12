package exception

type BlockValidateException struct{ logMessage }

func (e *BlockValidateException) ChainExceptions()         {}
func (e *BlockValidateException) BlockValidateExceptions() {}
func (e *BlockValidateException) Code() ExcTypes           { return 3030000 }
func (e *BlockValidateException) What() string             { return "Action validate exception" }

type UnlinkableBlockException struct{ logMessage }

func (e *UnlinkableBlockException) ChainExceptions()         {}
func (e *UnlinkableBlockException) BlockValidateExceptions() {}
func (e *UnlinkableBlockException) Code() ExcTypes           { return 3030001 }
func (e *UnlinkableBlockException) What() string             { return "Unlinkable block" }

type BlockTxOutputException struct{ logMessage }

func (e *BlockTxOutputException) ChainExceptions()         {}
func (e *BlockTxOutputException) BlockValidateExceptions() {}
func (e *BlockTxOutputException) Code() ExcTypes           { return 3030002 }
func (e *BlockTxOutputException) What() string {
	return "Transaction outputs in block do not match transaction outputs from applying block"
}

type BlockConcurrencyException struct{ logMessage }

func (e *BlockConcurrencyException) ChainExceptions()         {}
func (e *BlockConcurrencyException) BlockValidateExceptions() {}
func (e *BlockConcurrencyException) Code() ExcTypes           { return 3030003 }
func (e *BlockConcurrencyException) What() string {
	return "Block does not guarantee concurrent execution without conflicts"
}

type BlockLockException struct{ logMessage }

func (e *BlockLockException) ChainExceptions()         {}
func (e *BlockLockException) BlockValidateExceptions() {}
func (e *BlockLockException) Code() ExcTypes           { return 3030004 }
func (e *BlockLockException) What() string {
	return "Shard locks in block are incorrect or mal-formed"
}

type BlockResourceExhausted struct{ logMessage }

func (e *BlockResourceExhausted) ChainExceptions()         {}
func (e *BlockResourceExhausted) BlockValidateExceptions() {}
func (e *BlockResourceExhausted) Code() ExcTypes           { return 3030005 }
func (e *BlockResourceExhausted) What() string             { return "Block exhausted allowed resources" }

type BlockTooOldException struct{ logMessage }

func (e *BlockTooOldException) ChainExceptions()         {}
func (e *BlockTooOldException) BlockValidateExceptions() {}
func (e *BlockTooOldException) Code() ExcTypes           { return 3030006 }
func (e *BlockTooOldException) What() string             { return "Block is too old to push" }

type BlockFromTheFuture struct{ logMessage }

func (e *BlockFromTheFuture) ChainExceptions()         {}
func (e *BlockFromTheFuture) BlockValidateExceptions() {}
func (e *BlockFromTheFuture) Code() ExcTypes           { return 3030007 }
func (e *BlockFromTheFuture) What() string             { return "Block is from the future" }

type WrongSigningKey struct{ logMessage }

func (e *WrongSigningKey) ChainExceptions()         {}
func (e *WrongSigningKey) BlockValidateExceptions() {}
func (e *WrongSigningKey) Code() ExcTypes           { return 3030008 }
func (e *WrongSigningKey) What() string             { return "Block is not signed with expected key" }

type WrongProducer struct{ logMessage }

func (e *WrongProducer) ChainExceptions()         {}
func (e *WrongProducer) BlockValidateExceptions() {}
func (e *WrongProducer) Code() ExcTypes           { return 3030009 }
func (e *WrongProducer) What() string             { return "Block is not signed by expected producer" }
