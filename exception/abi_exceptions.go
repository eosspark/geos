package exception

type AbiException struct{ ELog }

func (AbiException) ChainExceptions() {}
func (AbiException) AbiExceptions()   {}
func (AbiException) Code() ExcTypes   { return 3015000 }
func (AbiException) What() string     { return "ABI exception" }

type AbiNotFoundException struct{ ELog }

func (AbiNotFoundException) ChainExceptions() {}
func (AbiNotFoundException) AbiExceptions()   {}
func (AbiNotFoundException) Code() ExcTypes   { return 3015001 }
func (AbiNotFoundException) What() string     { return "No ABI Found" }

type InvalidRicardianClauseException struct{ ELog }

func (InvalidRicardianClauseException) ChainExceptions() {}
func (InvalidRicardianClauseException) AbiExceptions()   {}
func (InvalidRicardianClauseException) Code() ExcTypes   { return 3015002 }
func (InvalidRicardianClauseException) What() string     { return "Invalid Ricardian Clause" }

type InvalidActionClauseException struct{ ELog }

func (InvalidActionClauseException) ChainExceptions() {}
func (InvalidActionClauseException) AbiExceptions()   {}
func (InvalidActionClauseException) Code() ExcTypes   { return 3015003 }
func (InvalidActionClauseException) What() string     { return "Invalid Ricardian Action" }

type InvalidTypeInsideAbi struct{ ELog }

func (InvalidTypeInsideAbi) ChainExceptions() {}
func (InvalidTypeInsideAbi) AbiExceptions()   {}
func (InvalidTypeInsideAbi) Code() ExcTypes   { return 3015004 }
func (InvalidTypeInsideAbi) What() string     { return "The type defined in the ABI is invalid" } // Not to be confused with abi_type_exception

type DuplicateAbiTypeDefException struct{ ELog }

func (DuplicateAbiTypeDefException) ChainExceptions() {}
func (DuplicateAbiTypeDefException) AbiExceptions()   {}
func (DuplicateAbiTypeDefException) Code() ExcTypes   { return 3015005 }
func (DuplicateAbiTypeDefException) What() string     { return "Duplicate type definition in the ABI" }

type DuplicateAbiStructDefException struct{ ELog }

func (DuplicateAbiStructDefException) ChainExceptions() {}
func (DuplicateAbiStructDefException) AbiExceptions()   {}
func (DuplicateAbiStructDefException) Code() ExcTypes   { return 3015006 }
func (DuplicateAbiStructDefException) What() string     { return "Duplicate struct definition in the ABI" }

type DuplicateAbiActionDefException struct{ ELog }

func (DuplicateAbiActionDefException) ChainExceptions() {}
func (DuplicateAbiActionDefException) AbiExceptions()   {}
func (DuplicateAbiActionDefException) Code() ExcTypes   { return 3015007 }
func (DuplicateAbiActionDefException) What() string     { return "Duplicate action definition in the ABI" }

type DuplicateAbiTableDefException struct{ ELog }

func (DuplicateAbiTableDefException) ChainExceptions() {}
func (DuplicateAbiTableDefException) AbiExceptions()   {}
func (DuplicateAbiTableDefException) Code() ExcTypes   { return 3015008 }
func (DuplicateAbiTableDefException) What() string     { return "Duplicate table definition in the ABI" }

type DuplicateAbiErrMsgDefException struct{ ELog }

func (DuplicateAbiErrMsgDefException) ChainExceptions() {}
func (DuplicateAbiErrMsgDefException) AbiExceptions()   {}
func (DuplicateAbiErrMsgDefException) Code() ExcTypes   { return 3015009 }
func (DuplicateAbiErrMsgDefException) What() string     { return "Duplicate error message definition in the ABI" }

type AbiSerializationDeadlineException struct{ ELog }

func (AbiSerializationDeadlineException) ChainExceptions() {}
func (AbiSerializationDeadlineException) AbiExceptions()   {}
func (AbiSerializationDeadlineException) Code() ExcTypes   { return 3015010 }
func (AbiSerializationDeadlineException) What() string {
	return "ABI serialization time has exceeded the deadline"
}

type AbiRecursionDepthException struct{ ELog }

func (AbiRecursionDepthException) ChainExceptions() {}
func (AbiRecursionDepthException) AbiExceptions()   {}
func (AbiRecursionDepthException) Code() ExcTypes   { return 3015011 }
func (AbiRecursionDepthException) What() string {
	return "ABI recursive definition has exceeded the max recursion depth"
}

type AbiCircularDefException struct{ ELog }

func (AbiCircularDefException) ChainExceptions() {}
func (AbiCircularDefException) AbiExceptions()   {}
func (AbiCircularDefException) Code() ExcTypes   { return 3015012 }
func (AbiCircularDefException) What() string     { return "Circular definition is detected in the ABI" }

type UnpackException struct{ ELog }

func (UnpackException) ChainExceptions() {}
func (UnpackException) AbiExceptions()   {}
func (UnpackException) Code() ExcTypes   { return 3015013 }
func (UnpackException) What() string     { return "Unpack data exception" }

type PackException struct{ ELog }

func (PackException) ChainExceptions() {}
func (PackException) AbiExceptions()   {}
func (PackException) Code() ExcTypes   { return 3015014 }
func (PackException) What() string     { return "Pack data exception" }

type DuplicateAbiVariantDefException struct{ ELog }

func (DuplicateAbiVariantDefException) ChainExceptions() {}
func (DuplicateAbiVariantDefException) AbiExceptions()   {}
func (DuplicateAbiVariantDefException) Code() ExcTypes   { return 3015015 }
func (DuplicateAbiVariantDefException) What() string     { return "Duplicate variant definition in the ABI" }

type UnsupportedAbiVersionException struct{ ELog }

func (UnsupportedAbiVersionException) ChainExceptions() {}
func (UnsupportedAbiVersionException) AbiExceptions()   {}
func (UnsupportedAbiVersionException) Code() ExcTypes   { return 3015016 }
func (UnsupportedAbiVersionException) What() string     { return "ABI has an unsupported version" }
