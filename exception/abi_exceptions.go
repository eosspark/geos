package exception

type AbiException struct{ logMessage }

func (e *AbiException) ChainExceptions()                {}
func (e *AbiException) ControllerEmitSignalExceptions() {}
func (e *AbiException) Code() ExcTypes                  { return 3015000 }
func (e *AbiException) What() string                    { return "ABI exception" }
