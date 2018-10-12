package exception

type ContractException struct{ logMessage }

func (ContractException) ChainExceptions()    {}
func (ContractException) ContractExceptions() {}
func (ContractException) Code() ExcTypes      { return 3160000 }
func (ContractException) What() string        { return "Contract exception" }

type InvalidTablePayer struct{ logMessage }

func (InvalidTablePayer) ChainExceptions()   {}
func (InvalidTablePayer) ContractException() {}
func (InvalidTablePayer) Code() ExcTypes     { return 3160001 }
func (InvalidTablePayer) What() string       { return "The payer of the table data is invalid" }

type TableAccessViolation struct{ logMessage }

func (TableAccessViolation) ChainExceptions()   {}
func (TableAccessViolation) ContractException() {}
func (TableAccessViolation) Code() ExcTypes     { return 3160002 }
func (TableAccessViolation) What() string       { return "Table access violation" }

type InvalidTableTterator struct{ logMessage }

func (InvalidTableTterator) ChainExceptions()   {}
func (InvalidTableTterator) ContractException() {}
func (InvalidTableTterator) Code() ExcTypes     { return 3160003 }
func (InvalidTableTterator) What() string       { return "Invalid table iterator" }

type TableNotInCache struct{ logMessage }

func (TableNotInCache) ChainExceptions()   {}
func (TableNotInCache) ContractException() {}
func (TableNotInCache) Code() ExcTypes     { return 3160004 }
func (TableNotInCache) What() string       { return "Table can not be found inside the cache" }

type TableOperationNotPermitted struct{ logMessage }

func (TableOperationNotPermitted) ChainExceptions()   {}
func (TableOperationNotPermitted) ContractException() {}
func (TableOperationNotPermitted) Code() ExcTypes     { return 3160005 }
func (TableOperationNotPermitted) What() string       { return "The table operation is not allowed" }

type InvalidContractVmType struct{ logMessage }

func (InvalidContractVmType) ChainExceptions()   {}
func (InvalidContractVmType) ContractException() {}
func (InvalidContractVmType) Code() ExcTypes     { return 3160006 }
func (InvalidContractVmType) What() string       { return "Invalid contract vm type" }

type InvalidContractVmVersion struct{ logMessage }

func (InvalidContractVmVersion) ChainExceptions()   {}
func (InvalidContractVmVersion) ContractException() {}
func (InvalidContractVmVersion) Code() ExcTypes     { return 3160007 }
func (InvalidContractVmVersion) What() string       { return "Invalid contract vm version" }

type SetExactCode struct{ logMessage }

func (SetExactCode) ChainExceptions()   {}
func (SetExactCode) ContractException() {}
func (SetExactCode) Code() ExcTypes     { return 3160008 }
func (SetExactCode) What() string {
	return "Contract is already running this version of code"
}

type WastFileNotFound struct{ logMessage }

func (WastFileNotFound) ChainExceptions()   {}
func (WastFileNotFound) ContractException() {}
func (WastFileNotFound) Code() ExcTypes     { return 3160009 }
func (WastFileNotFound) What() string       { return "No wast file found" }

type AbiFileNotFound struct{ logMessage }

func (AbiFileNotFound) ChainExceptions()   {}
func (AbiFileNotFound) ContractException() {}
func (AbiFileNotFound) Code() ExcTypes     { return 3160010 }
func (AbiFileNotFound) What() string       { return "No abi file found" }
