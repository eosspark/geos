package exception

import . "github.com/eosspark/eos-go/log"

type ActionValidateException struct{ LogMessage }

func (ActionValidateException) ChainExceptions()          {}
func (ActionValidateException) ActionValidateExceptions() {}
func (ActionValidateException) Code() ExcTypes            { return 3050000 }
func (ActionValidateException) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}

type AccountNameExistsException struct{ LogMessage }

func (AccountNameExistsException) ChainExceptions()          {}
func (AccountNameExistsException) ActionValidateExceptions() {}
func (AccountNameExistsException) Code() ExcTypes            { return 3050001 }
func (AccountNameExistsException) What() string {
	return "Account name already exists"
}

type InvalidActionArgsException struct{ LogMessage }

func (InvalidActionArgsException) ChainExceptions()          {}
func (InvalidActionArgsException) ActionValidateExceptions() {}
func (InvalidActionArgsException) Code() ExcTypes            { return 3050002 }
func (InvalidActionArgsException) What() string              { return "Invalid Action Arguments" }

type EosioAssertMessageException struct{ LogMessage }

func (EosioAssertMessageException) ChainExceptions()          {}
func (EosioAssertMessageException) ActionValidateExceptions() {}
func (EosioAssertMessageException) Code() ExcTypes            { return 3050003 }
func (EosioAssertMessageException) What() string {
	return "eosio_assert_message assertion failure"
}

type EosioAssertCodeException struct{ LogMessage }

func (EosioAssertCodeException) ChainExceptions()          {}
func (EosioAssertCodeException) ActionValidateExceptions() {}
func (EosioAssertCodeException) Code() ExcTypes            { return 3050004 }
func (EosioAssertCodeException) What() string {
	return "eosio_assert_code assertion failure"
}

type ActionNotFoundException struct{ LogMessage }

func (ActionNotFoundException) ChainExceptions()          {}
func (ActionNotFoundException) ActionValidateExceptions() {}
func (ActionNotFoundException) Code() ExcTypes            { return 3050005 }
func (ActionNotFoundException) What() string {
	return "Action can not be found"
}

type ActionDataAndStructMismatch struct{ LogMessage }

func (ActionDataAndStructMismatch) ChainExceptions()          {}
func (ActionDataAndStructMismatch) ActionValidateExceptions() {}
func (ActionDataAndStructMismatch) Code() ExcTypes            { return 3050006 }
func (ActionDataAndStructMismatch) What() string {
	return "Mismatch between action data and its struct"
}

type UnaccessibleApi struct{ LogMessage }

func (UnaccessibleApi) ChainExceptions()          {}
func (UnaccessibleApi) ActionValidateExceptions() {}
func (UnaccessibleApi) Code() ExcTypes            { return 3050007 }
func (UnaccessibleApi) What() string {
	return "Attempt to use unaccessible API"
}

type AbortCalled struct{ LogMessage }

func (AbortCalled) ChainExceptions()          {}
func (AbortCalled) ActionValidateExceptions() {}
func (AbortCalled) Code() ExcTypes            { return 3050008 }
func (AbortCalled) What() string              { return "Abort Called" }

type InlineActionTooBig struct{ LogMessage }

func (InlineActionTooBig) ChainExceptions()          {}
func (InlineActionTooBig) ActionValidateExceptions() {}
func (InlineActionTooBig) Code() ExcTypes            { return 3050009 }
func (InlineActionTooBig) What() string {
	return "Inline Action exceeds maximum size limit"
}
