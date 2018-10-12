package exception

type AuthorizationException struct{ logMessage }

func (e *AuthorizationException) ChainExceptions()         {}
func (e *AuthorizationException) AuthorizationExceptions() {}
func (e *AuthorizationException) Code() ExcTypes           { return 3090000 }
func (e *AuthorizationException) What() string             { return "Authorization exception" }
