package exception

import _ "github.com/eosspark/eos-go/log"

type AuthorizationException struct{ ELog }

func (AuthorizationException) ChainExceptions()         {}
func (AuthorizationException) AuthorizationExceptions() {}
func (AuthorizationException) Code() ExcTypes           { return 3090000 }
func (AuthorizationException) What() string             { return "Authorization exception" }

type TxDuplicateSig struct{ ELog }

func (TxDuplicateSig) ChainExceptions()         {}
func (TxDuplicateSig) AuthorizationExceptions() {}
func (TxDuplicateSig) Code() ExcTypes           { return 3090001 }
func (TxDuplicateSig) What() string             { return "Duplicate signature included" }

type TxIrrelevantSig struct{ ELog }

func (TxIrrelevantSig) ChainExceptions()         {}
func (TxIrrelevantSig) AuthorizationExceptions() {}
func (TxIrrelevantSig) Code() ExcTypes           { return 3090002 }
func (TxIrrelevantSig) What() string             { return "Irrelevant signature included" }

type UnsatisfiedAuthorization struct{ ELog }

func (UnsatisfiedAuthorization) ChainExceptions()         {}
func (UnsatisfiedAuthorization) AuthorizationExceptions() {}
func (UnsatisfiedAuthorization) Code() ExcTypes           { return 3090003 }
func (UnsatisfiedAuthorization) What() string             {
	return "Provided keys, permissions, and delays do not satisfy declared authorizations"
}

type MissingAuthException struct{ ELog }

func (MissingAuthException) ChainExceptions()         {}
func (MissingAuthException) AuthorizationExceptions() {}
func (MissingAuthException) Code() ExcTypes           { return 3090004 }
func (MissingAuthException) What() string             { return "Missing required authority" }

type IrrelevantAuthException struct{ ELog }

func (IrrelevantAuthException) ChainExceptions()         {}
func (IrrelevantAuthException) AuthorizationExceptions() {}
func (IrrelevantAuthException) Code() ExcTypes           { return 3090005 }
func (IrrelevantAuthException) What() string             { return "Irrelevant authority included" }

type InsufficientDelayException struct{ ELog }

func (InsufficientDelayException) ChainExceptions()         {}
func (InsufficientDelayException) AuthorizationExceptions() {}
func (InsufficientDelayException) Code() ExcTypes           { return 3090006 }
func (InsufficientDelayException) What() string             { return "Insufficient delay" }

type InvalidPermission struct{ ELog }

func (InvalidPermission) ChainExceptions()         {}
func (InvalidPermission) AuthorizationExceptions() {}
func (InvalidPermission) Code() ExcTypes           { return 3090007 }
func (InvalidPermission) What() string             { return "Invalid Permission" }

type UnlinkableMinPermissionAction struct{ ELog }

func (UnlinkableMinPermissionAction) ChainExceptions()         {}
func (UnlinkableMinPermissionAction) AuthorizationExceptions() {}
func (UnlinkableMinPermissionAction) Code() ExcTypes           { return 3090008 }
func (UnlinkableMinPermissionAction) What() string             {
	return "The action is not allowed to be linked with minimum permission"
}

type InvalidParentPermission struct{ ELog }

func (InvalidParentPermission) ChainExceptions()         {}
func (InvalidParentPermission) AuthorizationExceptions() {}
func (InvalidParentPermission) Code() ExcTypes           { return 3090009 }
func (InvalidParentPermission) What() string             { return "The parent permission is invalid" }
