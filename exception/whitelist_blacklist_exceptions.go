package exception

type WhitelistBlacklistException struct{ logMessage }

func (e *WhitelistBlacklistException) ChainExceptions()              {}
func (e *WhitelistBlacklistException) WhitelistBlacklistExceptions() {}
func (e *WhitelistBlacklistException) Code() ExcTypes                { return 3130000 }
func (e *WhitelistBlacklistException) What() string {
	return "Actor or contract whitelist/blacklist exception"
}
