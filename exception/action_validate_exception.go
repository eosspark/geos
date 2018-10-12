package exception

type ActionValidateException struct{ logMessage }

func (e *ActionValidateException) ChainExceptions()          {}
func (e *ActionValidateException) ActionValidateExceptions() {}
func (e *ActionValidateException) Code() ExcTypes            { return 3050000 }
func (e *ActionValidateException) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}

type AccountNameExistsException struct{ logMessage }

func (e *AccountNameExistsException) ChainExceptions()          {}
func (e *AccountNameExistsException) ActionValidateExceptions() {}
func (e *AccountNameExistsException) Code() ExcTypes            { return 3050001 }
func (e *AccountNameExistsException) What() string {
	return "Account name already exists"
}

type InvalidActionArgsException struct{ logMessage }

func (e *InvalidActionArgsException) ChainExceptions()          {}
func (e *InvalidActionArgsException) ActionValidateExceptions() {}
func (e *InvalidActionArgsException) Code() ExcTypes            { return 3050002 }
func (e *InvalidActionArgsException) What() string              { return "Invalid Action Arguments" }

type EosioAssertMessageException struct{ logMessage }

func (e *EosioAssertMessageException) ChainExceptions()          {}
func (e *EosioAssertMessageException) ActionValidateExceptions() {}
func (e *EosioAssertMessageException) Code() ExcTypes            { return 3050003 }
func (e *EosioAssertMessageException) What() string {
	return "eosio_assert_message assertion failure"
}

type EosioAssertCodeException struct{ logMessage }

func (e *EosioAssertCodeException) ChainExceptions()          {}
func (e *EosioAssertCodeException) ActionValidateExceptions() {}
func (e *EosioAssertCodeException) Code() ExcTypes            { return 3050004 }
func (e *EosioAssertCodeException) What() string {
	return "eosio_assert_code assertion failure"
}

type ActionNotFoundException struct{ logMessage }

func (e *ActionNotFoundException) ChainExceptions()          {}
func (e *ActionNotFoundException) ActionValidateExceptions() {}
func (e *ActionNotFoundException) Code() ExcTypes            { return 3050005 }
func (e *ActionNotFoundException) What() string {
	return "Action can not be found"
}

type ActionDataAndStructMismatch struct{ logMessage }

func (e *ActionDataAndStructMismatch) ChainExceptions()          {}
func (e *ActionDataAndStructMismatch) ActionValidateExceptions() {}
func (e *ActionDataAndStructMismatch) Code() ExcTypes            { return 3050006 }
func (e *ActionDataAndStructMismatch) What() string {
	return "Mismatch between action data and its struct"
}

type UnaccessibleApi struct{ logMessage }

func (e *UnaccessibleApi) ChainExceptions()          {}
func (e *UnaccessibleApi) ActionValidateExceptions() {}
func (e *UnaccessibleApi) Code() ExcTypes            { return 3050007 }
func (e *UnaccessibleApi) What() string {
	return "Attempt to use unaccessible API"
}

type AbortCalled struct{ logMessage }

func (e *AbortCalled) ChainExceptions()          {}
func (e *AbortCalled) ActionValidateExceptions() {}
func (e *AbortCalled) Code() ExcTypes            { return 3050008 }
func (e *AbortCalled) What() string              { return "Abort Called" }

type InlineActionTooBig struct{ logMessage }

func (e *InlineActionTooBig) ChainExceptions()          {}
func (e *InlineActionTooBig) ActionValidateExceptions() {}
func (e *InlineActionTooBig) Code() ExcTypes            { return 3050009 }
func (e *InlineActionTooBig) What() string {
	return "Inline Action exceeds maximum size limit"
}
