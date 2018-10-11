package exception

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

type PublicKeyTypeException struct{ logMessage }

func (e *PublicKeyTypeException) ChainExceptions()     {}
func (e *PublicKeyTypeException) ChainTypeExceptions() {}
func (e *PublicKeyTypeException) Code() ExcTypes       { return 3010002 }
func (e *PublicKeyTypeException) What() string         { return "Invalid public key" }

type PrivateKeyTypeException struct{ logMessage }

func (e *PrivateKeyTypeException) ChainExceptions()     {}
func (e *PrivateKeyTypeException) ChainTypeExceptions() {}
func (e *PrivateKeyTypeException) Code() ExcTypes       { return 3010003 }
func (e *PrivateKeyTypeException) What() string         { return "Invalid private key" }

type AuthorityTypeException struct{ logMessage }

func (e *AuthorityTypeException) ChainExceptions()     {}
func (e *AuthorityTypeException) ChainTypeExceptions() {}
func (e *AuthorityTypeException) Code() ExcTypes       { return 3010004 }
func (e *AuthorityTypeException) What() string         { return "Invalid authority" }

type ActionTypeException struct{ logMessage }

func (e *ActionTypeException) ChainExceptions()     {}
func (e *ActionTypeException) ChainTypeExceptions() {}
func (e *ActionTypeException) Code() ExcTypes       { return 3010005 }
func (e *ActionTypeException) What() string         { return "Invalid action" }

type TransactionTypeException struct{ logMessage }

func (e *TransactionTypeException) ChainExceptions()     {}
func (e *TransactionTypeException) ChainTypeExceptions() {}
func (e *TransactionTypeException) Code() ExcTypes       { return 3010006 }
func (e *TransactionTypeException) What() string         { return "Invalid transaction" }

type AbiTypeException struct{ logMessage }

func (e *AbiTypeException) ChainExceptions()     {}
func (e *AbiTypeException) ChainTypeExceptions() {}
func (e *AbiTypeException) Code() ExcTypes       { return 3010007 }
func (e *AbiTypeException) What() string         { return "Invalid ABI" }

type BlockIdTypeException struct{ logMessage }

func (e *BlockIdTypeException) ChainExceptions()     {}
func (e *BlockIdTypeException) ChainTypeExceptions() {}
func (e *BlockIdTypeException) Code() ExcTypes       { return 3010008 }
func (e *BlockIdTypeException) What() string         { return "Invalid block ID" }

type TransactionIdTypeException struct{ logMessage }

func (e *TransactionIdTypeException) ChainExceptions()     {}
func (e *TransactionIdTypeException) ChainTypeExceptions() {}
func (e *TransactionIdTypeException) Code() ExcTypes       { return 3010009 }
func (e *TransactionIdTypeException) What() string         { return "Invalid transaction ID" }

type PackedTransactionTypeException struct{ logMessage }

func (e *PackedTransactionTypeException) ChainExceptions()     {}
func (e *PackedTransactionTypeException) ChainTypeExceptions() {}
func (e *PackedTransactionTypeException) Code() ExcTypes       { return 3010010 }
func (e *PackedTransactionTypeException) What() string         { return "Invalid packed transaction" }

type AssetTypeException struct{ logMessage }

func (e *AssetTypeException) ChainExceptions()     {}
func (e *AssetTypeException) ChainTypeExceptions() {}
func (e *AssetTypeException) Code() ExcTypes       { return 3010011 }
func (e *AssetTypeException) What() string         { return "Invalid asset" }

type ChainIdTypeException struct{ logMessage }

func (e *ChainIdTypeException) ChainExceptions()     {}
func (e *ChainIdTypeException) ChainTypeExceptions() {}
func (e *ChainIdTypeException) Code() ExcTypes       { return 3010012 }
func (e *ChainIdTypeException) What() string         { return "Invalid chain ID" }

type FixedKeyTypeException struct{ logMessage }

func (e *FixedKeyTypeException) ChainExceptions()     {}
func (e *FixedKeyTypeException) ChainTypeExceptions() {}
func (e *FixedKeyTypeException) Code() ExcTypes       { return 3010013 }
func (e *FixedKeyTypeException) What() string         { return "Invalid fixed key" }

type SymbolTypeException struct{ logMessage }

func (e *SymbolTypeException) ChainExceptions()     {}
func (e *SymbolTypeException) ChainTypeExceptions() {}
func (e *SymbolTypeException) Code() ExcTypes       { return 30100014 }
func (e *SymbolTypeException) What() string         { return "Invalid symbol" }
