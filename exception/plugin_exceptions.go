package exception

type PluginException struct{ logMessage }

func (e *PluginException) ChainExceptions()  {}
func (e *PluginException) PluginExceptions() {}
func (e *PluginException) Code() ExcTypes    { return 3110000 }
func (e *PluginException) What() string {
	return "Plugin exception"
}

type MissingChainApiPluginException struct{ logMessage }

func (e *MissingChainApiPluginException) ChainExceptions()  {}
func (e *MissingChainApiPluginException) PluginExceptions() {}
func (e *MissingChainApiPluginException) Code() ExcTypes    { return 3110001 }
func (e *MissingChainApiPluginException) What() string {
	return "Missing Chain API Plugin"
}

type MissingWalletApiPluginException struct{ logMessage }

func (e *MissingWalletApiPluginException) ChainExceptions()  {}
func (e *MissingWalletApiPluginException) PluginExceptions() {}
func (e *MissingWalletApiPluginException) Code() ExcTypes    { return 3110002 }
func (e *MissingWalletApiPluginException) What() string {
	return "Missing Wallet API Plugin"
}

type MissingHistoryApiPluginException struct{ logMessage }

func (e *MissingHistoryApiPluginException) ChainExceptions()  {}
func (e *MissingHistoryApiPluginException) PluginExceptions() {}
func (e *MissingHistoryApiPluginException) Code() ExcTypes    { return 3110003 }
func (e *MissingHistoryApiPluginException) What() string {
	return "Missing History API Plugin"
}

type MissingNetApiPluginException struct{ logMessage }

func (e *MissingNetApiPluginException) ChainExceptions()  {}
func (e *MissingNetApiPluginException) PluginExceptions() {}
func (e *MissingNetApiPluginException) Code() ExcTypes    { return 3110004 }
func (e *MissingNetApiPluginException) What() string {
	return "Missing Net API Plugin"
}

type MissingChainPluginException struct{ logMessage }

func (e *MissingChainPluginException) ChainExceptions()  {}
func (e *MissingChainPluginException) PluginExceptions() {}
func (e *MissingChainPluginException) Code() ExcTypes    { return 3110005 }
func (e *MissingChainPluginException) What() string {
	return "Missing Chain Plugin"
}

type PluginConfigException struct{ logMessage }

func (e *PluginConfigException) ChainExceptions()  {}
func (e *PluginConfigException) PluginExceptions() {}
func (e *PluginConfigException) Code() ExcTypes    { return 3110006 }
func (e *PluginConfigException) What() string {
	return "Incorrect plugin configuration"
}
