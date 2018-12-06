package exception

import _ "github.com/eosspark/eos-go/log"

type BlockValidateException struct{ ELog }

func (BlockValidateException) ChainExceptions()         {}
func (BlockValidateException) BlockValidateExceptions() {}
func (BlockValidateException) Code() ExcTypes           { return 3030000 }
func (BlockValidateException) What() string             { return "Action validate exception" }

type UnlinkableBlockException struct{ ELog }

func (UnlinkableBlockException) ChainExceptions()         {}
func (UnlinkableBlockException) BlockValidateExceptions() {}
func (UnlinkableBlockException) Code() ExcTypes           { return 3030001 }
func (UnlinkableBlockException) What() string             { return "Unlinkable block" }

type BlockTxOutputException struct{ ELog }

func (BlockTxOutputException) ChainExceptions()         {}
func (BlockTxOutputException) BlockValidateExceptions() {}
func (BlockTxOutputException) Code() ExcTypes           { return 3030002 }
func (BlockTxOutputException) What() string {
	return "Transaction outputs in block do not match transaction outputs from applying block"
}

type BlockConcurrencyException struct{ ELog }

func (BlockConcurrencyException) ChainExceptions()         {}
func (BlockConcurrencyException) BlockValidateExceptions() {}
func (BlockConcurrencyException) Code() ExcTypes           { return 3030003 }
func (BlockConcurrencyException) What() string {
	return "Block does not guarantee concurrent execution without conflicts"
}

type BlockLockException struct{ ELog }

func (BlockLockException) ChainExceptions()         {}
func (BlockLockException) BlockValidateExceptions() {}
func (BlockLockException) Code() ExcTypes           { return 3030004 }
func (BlockLockException) What() string {
	return "Shard locks in block are incorrect or mal-formed"
}

type BlockResourceExhausted struct{ ELog }

func (BlockResourceExhausted) ChainExceptions()         {}
func (BlockResourceExhausted) BlockValidateExceptions() {}
func (BlockResourceExhausted) Code() ExcTypes           { return 3030005 }
func (BlockResourceExhausted) What() string             { return "Block exhausted allowed resources" }

type BlockTooOldException struct{ ELog }

func (BlockTooOldException) ChainExceptions()         {}
func (BlockTooOldException) BlockValidateExceptions() {}
func (BlockTooOldException) Code() ExcTypes           { return 3030006 }
func (BlockTooOldException) What() string             { return "Block is too old to push" }

type BlockFromTheFuture struct{ ELog }

func (BlockFromTheFuture) ChainExceptions()         {}
func (BlockFromTheFuture) BlockValidateExceptions() {}
func (BlockFromTheFuture) Code() ExcTypes           { return 3030007 }
func (BlockFromTheFuture) What() string             { return "Block is from the future" }

type WrongSigningKey struct{ ELog }

func (WrongSigningKey) ChainExceptions()         {}
func (WrongSigningKey) BlockValidateExceptions() {}
func (WrongSigningKey) Code() ExcTypes           { return 3030008 }
func (WrongSigningKey) What() string             { return "Block is not signed with expected key" }

type WrongProducer struct{ ELog }

func (WrongProducer) ChainExceptions()         {}
func (WrongProducer) BlockValidateExceptions() {}
func (WrongProducer) Code() ExcTypes           { return 3030009 }
func (WrongProducer) What() string             { return "Block is not signed by expected producer" }
