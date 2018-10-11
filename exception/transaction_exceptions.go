package exception

type TransactionException struct{ logMessage }

func (e *TransactionException) ChainExceptions()       {}
func (e *TransactionException) TransactionExceptions() {}
func (e *TransactionException) Code() ExcTypes         { return 3040000 }
func (e *TransactionException) What() string           { return "Transaction exception" }
