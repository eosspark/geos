package exception

type WhitelistBlacklistException struct{ logMessage }

func (WhitelistBlacklistException) ChainExceptions()              {}
func (WhitelistBlacklistException) WhitelistBlacklistExceptions() {}
func (WhitelistBlacklistException) Code() ExcTypes                { return 3130000 }
func (WhitelistBlacklistException) What() string {
	return "Actor or contract whitelist/blacklist exception"
}

type ActorWhitelistException struct{ logMessage }

func (ActorWhitelistException) ChainExceptions()              {}
func (ActorWhitelistException) WhitelistBlacklistExceptions() {}
func (ActorWhitelistException) Code() ExcTypes                { return 3130001 }
func (ActorWhitelistException) What() string {
	return "Authorizing actor of transaction is not on the whitelist"
}

type ActorBlacklistException struct{ logMessage }

func (ActorBlacklistException) ChainExceptions()              {}
func (ActorBlacklistException) WhitelistBlacklistExceptions() {}
func (ActorBlacklistException) Code() ExcTypes                { return 3130002 }
func (ActorBlacklistException) What() string {
	return "Authorizing actor of transaction is on the blacklist"
}

type ContractWhitelistException struct{ logMessage }

func (ContractWhitelistException) ChainExceptions()              {}
func (ContractWhitelistException) WhitelistBlacklistExceptions() {}
func (ContractWhitelistException) Code() ExcTypes                { return 3130003 }
func (ContractWhitelistException) What() string {
	return "Contract to execute is not on the whitelist"
}

type ContractBlacklistException struct{ logMessage }

func (ContractBlacklistException) ChainExceptions()              {}
func (ContractBlacklistException) WhitelistBlacklistExceptions() {}
func (ContractBlacklistException) Code() ExcTypes                { return 3130004 }
func (ContractBlacklistException) What() string {
	return "Contract to execute is on the blacklist"
}

type ActionBlacklistException struct{ logMessage }

func (ActionBlacklistException) ChainExceptions()              {}
func (ActionBlacklistException) WhitelistBlacklistExceptions() {}
func (ActionBlacklistException) Code() ExcTypes                { return 3130005 }
func (ActionBlacklistException) What() string {
	return "Action to execute is on the blacklist"
}

type KeyBlacklistException struct{ logMessage }

func (KeyBlacklistException) ChainExceptions()              {}
func (KeyBlacklistException) WhitelistBlacklistExceptions() {}
func (KeyBlacklistException) Code() ExcTypes                { return 3130006 }
func (KeyBlacklistException) What() string {
	return "Public key in authority is on the blacklist"
}
