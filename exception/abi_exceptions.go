package exception

type AbiException struct{ logMessage }

func (e *AbiException) ChainExceptions() {}
func (e *AbiException) AbiExceptions()   {}
func (e *AbiException) Code() ExcTypes   { return 3015000 }
func (e *AbiException) What() string     { return "ABI exception" }

type AbiNotFoundException struct{ logMessage }

func (e *AbiNotFoundException) ChainExceptions() {}
func (e *AbiNotFoundException) AbiExceptions()   {}
func (e *AbiNotFoundException) Code() ExcTypes   { return 3015001 }
func (e *AbiNotFoundException) What() string     { return "No ABI Found" }

type InvalidRicardianClauseException struct{ logMessage }

func (e *InvalidRicardianClauseException) ChainExceptions() {}
func (e *InvalidRicardianClauseException) AbiExceptions()   {}
func (e *InvalidRicardianClauseException) Code() ExcTypes   { return 3015002 }
func (e *InvalidRicardianClauseException) What() string     { return "Invalid Ricardian Clause" }

type InvalidActionClauseException struct{ logMessage }

func (e *InvalidActionClauseException) ChainExceptions() {}
func (e *InvalidActionClauseException) AbiExceptions()   {}
func (e *InvalidActionClauseException) Code() ExcTypes   { return 3015003 }
func (e *InvalidActionClauseException) What() string     { return "Invalid Ricardian Action" }

type InvalidTypeInsideAbi struct{ logMessage }

func (e *InvalidTypeInsideAbi) ChainExceptions() {}
func (e *InvalidTypeInsideAbi) AbiExceptions()   {}
func (e *InvalidTypeInsideAbi) Code() ExcTypes   { return 3015004 }
func (e *InvalidTypeInsideAbi) What() string     { return "The type defined in the ABI is invalid" } // Not to be confused with abi_type_exception

type DuplicateAbiTypeDefException struct{ logMessage }

func (e *DuplicateAbiTypeDefException) ChainExceptions() {}
func (e *DuplicateAbiTypeDefException) AbiExceptions()   {}
func (e *DuplicateAbiTypeDefException) Code() ExcTypes   { return 3015005 }
func (e *DuplicateAbiTypeDefException) What() string     { return "Duplicate type definition in the ABI" }

type DuplicateAbiStructDefException struct{ logMessage }

func (e *DuplicateAbiStructDefException) ChainExceptions() {}
func (e *DuplicateAbiStructDefException) AbiExceptions()   {}
func (e *DuplicateAbiStructDefException) Code() ExcTypes   { return 3015006 }
func (e *DuplicateAbiStructDefException) What() string     { return "Duplicate struct definition in the ABI" }

type DuplicateAbiActionDefException struct{ logMessage }

func (e *DuplicateAbiActionDefException) ChainExceptions() {}
func (e *DuplicateAbiActionDefException) AbiExceptions()   {}
func (e *DuplicateAbiActionDefException) Code() ExcTypes   { return 3015007 }
func (e *DuplicateAbiActionDefException) What() string     { return "Duplicate action definition in the ABI" }

type DuplicateAbiTableDefException struct{ logMessage }

func (e *DuplicateAbiTableDefException) ChainExceptions() {}
func (e *DuplicateAbiTableDefException) AbiExceptions()   {}
func (e *DuplicateAbiTableDefException) Code() ExcTypes   { return 3015008 }
func (e *DuplicateAbiTableDefException) What() string     { return "Duplicate table definition in the ABI" }

type DuplicateAbiErrMsgDefException struct{ logMessage }

func (e *DuplicateAbiErrMsgDefException) ChainExceptions() {}
func (e *DuplicateAbiErrMsgDefException) AbiExceptions()   {}
func (e *DuplicateAbiErrMsgDefException) Code() ExcTypes   { return 3015009 }
func (e *DuplicateAbiErrMsgDefException) What() string     { return "Duplicate error message definition in the ABI" }

type AbiSerializationDeadlineException struct{ logMessage }

func (e *AbiSerializationDeadlineException) ChainExceptions() {}
func (e *AbiSerializationDeadlineException) AbiExceptions()   {}
func (e *AbiSerializationDeadlineException) Code() ExcTypes   { return 3015010 }
func (e *AbiSerializationDeadlineException) What() string {
	return "ABI serialization time has exceeded the deadline"
}

type AbiRecursionDepthException struct{ logMessage }

func (e *AbiRecursionDepthException) ChainExceptions() {}
func (e *AbiRecursionDepthException) AbiExceptions()   {}
func (e *AbiRecursionDepthException) Code() ExcTypes   { return 3015011 }
func (e *AbiRecursionDepthException) What() string {
	return "ABI recursive definition has exceeded the max recursion depth"
}

type AbiCircularDefException struct{ logMessage }

func (e *AbiCircularDefException) ChainExceptions() {}
func (e *AbiCircularDefException) AbiExceptions()   {}
func (e *AbiCircularDefException) Code() ExcTypes   { return 3015012 }
func (e *AbiCircularDefException) What() string     { return "Circular definition is detected in the ABI" }

type UnpackException struct{ logMessage }

func (e *UnpackException) ChainExceptions() {}
func (e *UnpackException) AbiExceptions()   {}
func (e *UnpackException) Code() ExcTypes   { return 3015013 }
func (e *UnpackException) What() string     { return "Unpack data exception" }

type PackException struct{ logMessage }

func (e *PackException) ChainExceptions() {}
func (e *PackException) AbiExceptions()   {}
func (e *PackException) Code() ExcTypes   { return 3015014 }
func (e *PackException) What() string     { return "Pack data exception" }

type DuplicateAbiVariantDefException struct{ logMessage }

func (e *DuplicateAbiVariantDefException) ChainExceptions() {}
func (e *DuplicateAbiVariantDefException) AbiExceptions()   {}
func (e *DuplicateAbiVariantDefException) Code() ExcTypes   { return 3015015 }
func (e *DuplicateAbiVariantDefException) What() string     { return "Duplicate variant definition in the ABI" }

type UnsupportedAbiVersionException struct{ logMessage }

func (e *UnsupportedAbiVersionException) ChainExceptions() {}
func (e *UnsupportedAbiVersionException) AbiExceptions()   {}
func (e *UnsupportedAbiVersionException) Code() ExcTypes   { return 3015016 }
func (e *UnsupportedAbiVersionException) What() string     { return "ABI has an unsupported version" }
