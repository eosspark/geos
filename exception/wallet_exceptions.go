package exception

type WalletException struct{ logMessage }

func (e *WalletException) ChainExceptions()  {}
func (e *WalletException) WalletExceptions() {}
func (e *WalletException) Code() ExcTypes    { return 3120000 }
func (e *WalletException) What() string {
	return "Invalid contract vm version"
}
