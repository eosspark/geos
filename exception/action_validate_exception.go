package exception

type ActionValidateException struct{ logMessage }

func (e *ActionValidateException) ChainExceptions()          {}
func (e *ActionValidateException) ActionValidateExceptions() {}
func (e *ActionValidateException) Code() ExcTypes            { return 3050000 }
func (e *ActionValidateException) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}
