package exception

import . "github.com/eosspark/eos-go/log"

type WalletException struct{ LogMessage }

func (WalletException) ChainExceptions()  {}
func (WalletException) WalletExceptions() {}
func (WalletException) Code() ExcTypes    { return 3120000 }
func (WalletException) What() string {
	return "Invalid contract vm version"
}

type WalletExistException struct{ LogMessage }

func (WalletExistException) ChainExceptions()  {}
func (WalletExistException) WalletExceptions() {}
func (WalletExistException) Code() ExcTypes    { return 3120001 }
func (WalletExistException) What() string      { return "Wallet already exists" }

type WalletNonexistentException struct{ LogMessage }

func (WalletNonexistentException) ChainExceptions()  {}
func (WalletNonexistentException) WalletExceptions() {}
func (WalletNonexistentException) Code() ExcTypes    { return 3120002 }
func (WalletNonexistentException) What() string      { return "Nonexistent wallet" }

type WalletLockedException struct{ LogMessage }

func (WalletLockedException) ChainExceptions()  {}
func (WalletLockedException) WalletExceptions() {}
func (WalletLockedException) Code() ExcTypes    { return 3120003 }
func (WalletLockedException) What() string      { return "Locked wallet" }

type WalletMissingPubKeyException struct{ LogMessage }

func (WalletMissingPubKeyException) ChainExceptions()  {}
func (WalletMissingPubKeyException) WalletExceptions() {}
func (WalletMissingPubKeyException) Code() ExcTypes    { return 3120004 }
func (WalletMissingPubKeyException) What() string      { return "Missing public key" }

type WalletInvalidPasswordException struct{ LogMessage }

func (WalletInvalidPasswordException) ChainExceptions()  {}
func (WalletInvalidPasswordException) WalletExceptions() {}
func (WalletInvalidPasswordException) Code() ExcTypes    { return 3120005 }
func (WalletInvalidPasswordException) What() string      { return "Invalid wallet password" }

type WalletNotAvailableException struct{ LogMessage }

func (WalletNotAvailableException) ChainExceptions()  {}
func (WalletNotAvailableException) WalletExceptions() {}
func (WalletNotAvailableException) Code() ExcTypes    { return 3120006 }
func (WalletNotAvailableException) What() string      { return "No available wallet" }

type WalletUnlockedException struct{ LogMessage }

func (WalletUnlockedException) ChainExceptions()  {}
func (WalletUnlockedException) WalletExceptions() {}
func (WalletUnlockedException) Code() ExcTypes    { return 3120007 }
func (WalletUnlockedException) What() string      { return "Already unlocked" }

type KeyExistException struct{ LogMessage }

func (KeyExistException) ChainExceptions()  {}
func (KeyExistException) WalletExceptions() {}
func (KeyExistException) Code() ExcTypes    { return 3120008 }
func (KeyExistException) What() string      { return "Key already exists" }

type KeyNonexistentException struct{ LogMessage }

func (KeyNonexistentException) ChainExceptions()  {}
func (KeyNonexistentException) WalletExceptions() {}
func (KeyNonexistentException) Code() ExcTypes    { return 3120009 }
func (KeyNonexistentException) What() string      { return "Nonexistent key" }

type UnsupportedKeyTypeException struct{ LogMessage }

func (UnsupportedKeyTypeException) ChainExceptions()  {}
func (UnsupportedKeyTypeException) WalletExceptions() {}
func (UnsupportedKeyTypeException) Code() ExcTypes    { return 3120010 }
func (UnsupportedKeyTypeException) What() string      { return "Unsupported key type" }

type InvalidLockTimeoutException struct{ LogMessage }

func (InvalidLockTimeoutException) ChainExceptions()  {}
func (InvalidLockTimeoutException) WalletExceptions() {}
func (InvalidLockTimeoutException) Code() ExcTypes    { return 3120011 }
func (InvalidLockTimeoutException) What() string {
	return "Wallet lock timeout is invalid"
}

type SecureEnclaveException struct{ LogMessage }

func (SecureEnclaveException) ChainExceptions()  {}
func (SecureEnclaveException) WalletExceptions() {}
func (SecureEnclaveException) Code() ExcTypes    { return 3120012 }
func (SecureEnclaveException) What() string {
	return "Secure Enclave Exception"
}
