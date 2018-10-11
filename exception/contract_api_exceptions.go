package exception

type ContractApiException struct{ logMessage }

func (e *ContractApiException) ChainExceptions()       {}
func (e *ContractApiException) ContractApiExceptions() {}
func (e *ContractApiException) Code() ExcTypes         { return 3230000 }
func (e *ContractApiException) What() string           { return "Contract API exception" }
