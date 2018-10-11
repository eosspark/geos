package exception

type MiscException struct{ logMessage }

func (e *MiscException) ChainExceptions() {}
func (e *MiscException) MiscExceptions()  {}
func (e *MiscException) Code() ExcTypes   { return 3100000 }
func (e *MiscException) What() string {
	return "Miscellaneous exception"
}
