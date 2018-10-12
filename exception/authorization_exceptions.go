package exception

type AuthorizationException struct{ logMessage }

func (e *AuthorizationException) ChainExceptions()         {}
func (e *AuthorizationException) AuthorizationExceptions() {}
func (e *AuthorizationException) Code() ExcTypes           { return 3090000 }
func (e *AuthorizationException) What() string             { return "Authorization exception" }

type TxDuplicateSig struct{ logMessage }

func (e *TxDuplicateSig) ChainExceptions()         {}
func (e *TxDuplicateSig) AuthorizationExceptions() {}
func (e *TxDuplicateSig) Code() ExcTypes           { return 3090001 }
func (e *TxDuplicateSig) What() string             { return "Duplicate signature included" }

type TxIrrelevantSig struct{ logMessage }

func (e *TxIrrelevantSig) ChainExceptions()         {}
func (e *TxIrrelevantSig) AuthorizationExceptions() {}
func (e *TxIrrelevantSig) Code() ExcTypes           { return 3090002 }
func (e *TxIrrelevantSig) What() string             { return "Irrelevant signature included" }

type UnsatisfiedAuthorization struct{ logMessage }

func (e *UnsatisfiedAuthorization) ChainExceptions()         {}
func (e *UnsatisfiedAuthorization) AuthorizationExceptions() {}
func (e *UnsatisfiedAuthorization) Code() ExcTypes           { return 3090003 }
func (e *UnsatisfiedAuthorization) What() string             {
	return "Provided keys, permissions, and delays do not satisfy declared authorizations"
}

type MissingAuthException struct{ logMessage }

func (e *MissingAuthException) ChainExceptions()         {}
func (e *MissingAuthException) AuthorizationExceptions() {}
func (e *MissingAuthException) Code() ExcTypes           { return 3090004 }
func (e *MissingAuthException) What() string             { return "Missing required authority" }

type IrrelevantAuthException struct{ logMessage }

func (e *IrrelevantAuthException) ChainExceptions()         {}
func (e *IrrelevantAuthException) AuthorizationExceptions() {}
func (e *IrrelevantAuthException) Code() ExcTypes           { return 3090005 }
func (e *IrrelevantAuthException) What() string             { return "Irrelevant authority included" }

type InsufficientDelayException struct{ logMessage }

func (e *InsufficientDelayException) ChainExceptions()         {}
func (e *InsufficientDelayException) AuthorizationExceptions() {}
func (e *InsufficientDelayException) Code() ExcTypes           { return 3090006 }
func (e *InsufficientDelayException) What() string             { return "Insufficient delay" }

type InvalidPermission struct{ logMessage }

func (e *InvalidPermission) ChainExceptions()         {}
func (e *InvalidPermission) AuthorizationExceptions() {}
func (e *InvalidPermission) Code() ExcTypes           { return 3090007 }
func (e *InvalidPermission) What() string             { return "Invalid Permission" }

type UnlinkableMinPermissionAction struct{ logMessage }

func (e *UnlinkableMinPermissionAction) ChainExceptions()         {}
func (e *UnlinkableMinPermissionAction) AuthorizationExceptions() {}
func (e *UnlinkableMinPermissionAction) Code() ExcTypes           { return 3090008 }
func (e *UnlinkableMinPermissionAction) What() string             {
	return "The action is not allowed to be linked with minimum permission"
}

type InvalidParentPermission struct{ logMessage }

func (e *InvalidParentPermission) ChainExceptions()         {}
func (e *InvalidParentPermission) AuthorizationExceptions() {}
func (e *InvalidParentPermission) Code() ExcTypes           { return 3090009 }
func (e *InvalidParentPermission) What() string             { return "The parent permission is invalid" }
