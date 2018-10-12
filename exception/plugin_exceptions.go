package exception

type PluginException struct{ logMessage }

func (PluginException) ChainExceptions()  {}
func (PluginException) PluginExceptions() {}
func (PluginException) Code() ExcTypes    { return 3110000 }
func (PluginException) What() string {
	return "Plugin exception"
}

type MissingChainApiPluginException struct{ logMessage }

func (MissingChainApiPluginException) ChainExceptions()  {}
func (MissingChainApiPluginException) PluginExceptions() {}
func (MissingChainApiPluginException) Code() ExcTypes    { return 3110001 }
func (MissingChainApiPluginException) What() string {
	return "Missing Chain API Plugin"
}

type MissingWalletApiPluginException struct{ logMessage }

func (MissingWalletApiPluginException) ChainExceptions()  {}
func (MissingWalletApiPluginException) PluginExceptions() {}
func (MissingWalletApiPluginException) Code() ExcTypes    { return 3110002 }
func (MissingWalletApiPluginException) What() string {
	return "Missing Wallet API Plugin"
}

type MissingHistoryApiPluginException struct{ logMessage }

func (MissingHistoryApiPluginException) ChainExceptions()  {}
func (MissingHistoryApiPluginException) PluginExceptions() {}
func (MissingHistoryApiPluginException) Code() ExcTypes    { return 3110003 }
func (MissingHistoryApiPluginException) What() string {
	return "Missing History API Plugin"
}

type MissingNetApiPluginException struct{ logMessage }

func (MissingNetApiPluginException) ChainExceptions()  {}
func (MissingNetApiPluginException) PluginExceptions() {}
func (MissingNetApiPluginException) Code() ExcTypes    { return 3110004 }
func (MissingNetApiPluginException) What() string {
	return "Missing Net API Plugin"
}

type MissingChainPluginException struct{ logMessage }

func (MissingChainPluginException) ChainExceptions()  {}
func (MissingChainPluginException) PluginExceptions() {}
func (MissingChainPluginException) Code() ExcTypes    { return 3110005 }
func (MissingChainPluginException) What() string {
	return "Missing Chain Plugin"
}

type PluginConfigException struct{ logMessage }

func (PluginConfigException) ChainExceptions()  {}
func (PluginConfigException) PluginExceptions() {}
func (PluginConfigException) Code() ExcTypes    { return 3110006 }
func (PluginConfigException) What() string {
	return "Incorrect plugin configuration"
}
