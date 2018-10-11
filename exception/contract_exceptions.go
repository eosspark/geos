package exception

type ContractException struct{ logMessage }

func (e *ContractException) ChainExceptions()    {}
func (e *ContractException) ContractExceptions() {}
func (e *ContractException) Code() ExcTypes      { return 3160000 }
func (e *ContractException) What() string        { return "Contract exception" }

type InvalidContractVmVersion struct{ logMessage }

func (e *InvalidContractVmVersion) ChainExceptions()        {}
func (e *InvalidContractVmVersion) ForkDatabaseExceptions() {}
func (e *InvalidContractVmVersion) Code() ExcTypes          { return 3160007 }
func (e *InvalidContractVmVersion) What() string            { return "Invalid contract vm version" }

type SetExactCode struct{ logMessage }

func (e *SetExactCode) ChainExceptions()        {}
func (e *SetExactCode) ForkDatabaseExceptions() {}
func (e *SetExactCode) Code() ExcTypes          { return 3160008 }
func (e *SetExactCode) What() string {
	return "Contract is already running this version of code"
}
