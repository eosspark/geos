package exception

import . "github.com/eosspark/eos-go/log"

type ChainTypeException struct{ LogMessage }

func (ChainTypeException) ChainExceptions()     {}
func (ChainTypeException) ChainTypeExceptions() {}
func (ChainTypeException) Code() ExcTypes       { return 3010000 }
func (ChainTypeException) What() string         { return "chain type exception" }

type NameTypeException struct{ LogMessage }

func (NameTypeException) ChainExceptions()     {}
func (NameTypeException) ChainTypeExceptions() {}
func (NameTypeException) Code() ExcTypes       { return 3010001 }
func (NameTypeException) What() string         { return "Invalid name" }

type PublicKeyTypeException struct{ LogMessage }

func (PublicKeyTypeException) ChainExceptions()     {}
func (PublicKeyTypeException) ChainTypeExceptions() {}
func (PublicKeyTypeException) Code() ExcTypes       { return 3010002 }
func (PublicKeyTypeException) What() string         { return "Invalid public key" }

type PrivateKeyTypeException struct{ LogMessage }

func (PrivateKeyTypeException) ChainExceptions()     {}
func (PrivateKeyTypeException) ChainTypeExceptions() {}
func (PrivateKeyTypeException) Code() ExcTypes       { return 3010003 }
func (PrivateKeyTypeException) What() string         { return "Invalid private key" }

type AuthorityTypeException struct{ LogMessage }

func (AuthorityTypeException) ChainExceptions()     {}
func (AuthorityTypeException) ChainTypeExceptions() {}
func (AuthorityTypeException) Code() ExcTypes       { return 3010004 }
func (AuthorityTypeException) What() string         { return "Invalid authority" }

type ActionTypeException struct{ LogMessage }

func (ActionTypeException) ChainExceptions()     {}
func (ActionTypeException) ChainTypeExceptions() {}
func (ActionTypeException) Code() ExcTypes       { return 3010005 }
func (ActionTypeException) What() string         { return "Invalid action" }

type TransactionTypeException struct{ LogMessage }

func (TransactionTypeException) ChainExceptions()     {}
func (TransactionTypeException) ChainTypeExceptions() {}
func (TransactionTypeException) Code() ExcTypes       { return 3010006 }
func (TransactionTypeException) What() string         { return "Invalid transaction" }

type AbiTypeException struct{ LogMessage }

func (AbiTypeException) ChainExceptions()     {}
func (AbiTypeException) ChainTypeExceptions() {}
func (AbiTypeException) Code() ExcTypes       { return 3010007 }
func (AbiTypeException) What() string         { return "Invalid ABI" }

type BlockIdTypeException struct{ LogMessage }

func (BlockIdTypeException) ChainExceptions()     {}
func (BlockIdTypeException) ChainTypeExceptions() {}
func (BlockIdTypeException) Code() ExcTypes       { return 3010008 }
func (BlockIdTypeException) What() string         { return "Invalid block ID" }

type TransactionIdTypeException struct{ LogMessage }

func (TransactionIdTypeException) ChainExceptions()     {}
func (TransactionIdTypeException) ChainTypeExceptions() {}
func (TransactionIdTypeException) Code() ExcTypes       { return 3010009 }
func (TransactionIdTypeException) What() string         { return "Invalid transaction ID" }

type PackedTransactionTypeException struct{ LogMessage }

func (PackedTransactionTypeException) ChainExceptions()     {}
func (PackedTransactionTypeException) ChainTypeExceptions() {}
func (PackedTransactionTypeException) Code() ExcTypes       { return 3010010 }
func (PackedTransactionTypeException) What() string         { return "Invalid packed transaction" }

type AssetTypeException struct{ LogMessage }

func (AssetTypeException) ChainExceptions()     {}
func (AssetTypeException) ChainTypeExceptions() {}
func (AssetTypeException) Code() ExcTypes       { return 3010011 }
func (AssetTypeException) What() string         { return "Invalid asset" }

type ChainIdTypeException struct{ LogMessage }

func (ChainIdTypeException) ChainExceptions()     {}
func (ChainIdTypeException) ChainTypeExceptions() {}
func (ChainIdTypeException) Code() ExcTypes       { return 3010012 }
func (ChainIdTypeException) What() string         { return "Invalid chain ID" }

type FixedKeyTypeException struct{ LogMessage }

func (FixedKeyTypeException) ChainExceptions()     {}
func (FixedKeyTypeException) ChainTypeExceptions() {}
func (FixedKeyTypeException) Code() ExcTypes       { return 3010013 }
func (FixedKeyTypeException) What() string         { return "Invalid fixed key" }

type SymbolTypeException struct{ LogMessage }

func (SymbolTypeException) ChainExceptions()     {}
func (SymbolTypeException) ChainTypeExceptions() {}
func (SymbolTypeException) Code() ExcTypes       { return 30100014 }
func (SymbolTypeException) What() string         { return "Invalid symbol" }
