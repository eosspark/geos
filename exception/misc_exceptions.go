package exception

import _ "github.com/eosspark/eos-go/log"

type MiscException struct{ ELog }

func (MiscException) ChainExceptions() {}
func (MiscException) MiscExceptions()  {}
func (MiscException) Code() ExcTypes   { return 3100000 }
func (MiscException) What() string     { return "Miscellaneous exception" }

type RateLimitingStateInconsistent struct{ ELog }

func (RateLimitingStateInconsistent) ChainExceptions() {}
func (RateLimitingStateInconsistent) MiscExceptions()  {}
func (RateLimitingStateInconsistent) Code() ExcTypes   { return 3100001 }
func (RateLimitingStateInconsistent) What() string {
	return "Internal state is no longer consistent"
}

type UnknownBlockException struct{ ELog }

func (UnknownBlockException) ChainExceptions() {}
func (UnknownBlockException) MiscExceptions()  {}
func (UnknownBlockException) Code() ExcTypes   { return 3100002 }
func (UnknownBlockException) What() string     { return "Unknown block" }

type UnknownTransactionException struct{ ELog }

func (UnknownTransactionException) ChainExceptions() {}
func (UnknownTransactionException) MiscExceptions()  {}
func (UnknownTransactionException) Code() ExcTypes   { return 3100003 }
func (UnknownTransactionException) What() string     { return "Unknown transaction" }

type FixedReversibleDbException struct{ ELog }

func (FixedReversibleDbException) ChainExceptions() {}
func (FixedReversibleDbException) MiscExceptions()  {}
func (FixedReversibleDbException) Code() ExcTypes   { return 3100004 }
func (FixedReversibleDbException) What() string {
	return "Corrupted reversible block database was fixed"
}

type ExtractGenesisStateException struct{ ELog }

func (ExtractGenesisStateException) ChainExceptions() {}
func (ExtractGenesisStateException) MiscExceptions()  {}
func (ExtractGenesisStateException) Code() ExcTypes   { return 3100005 }
func (ExtractGenesisStateException) What() string {
	return "Extracted genesis state from blocks.log"
}

type SubjectiveBlockProductionException struct{ ELog }

func (SubjectiveBlockProductionException) ChainExceptions() {}
func (SubjectiveBlockProductionException) MiscExceptions()  {}
func (SubjectiveBlockProductionException) Code() ExcTypes   { return 3100006 }
func (SubjectiveBlockProductionException) What() string {
	return "Subjective exception thrown during block production"
}

type MultipleVoterInfo struct{ ELog }

func (MultipleVoterInfo) ChainExceptions() {}
func (MultipleVoterInfo) MiscExceptions()  {}
func (MultipleVoterInfo) Code() ExcTypes   { return 3100007 }
func (MultipleVoterInfo) What() string {
	return "Multiple voter info detected"
}

type UnsupportedFeature struct{ ELog }

func (UnsupportedFeature) ChainExceptions() {}
func (UnsupportedFeature) MiscExceptions()  {}
func (UnsupportedFeature) Code() ExcTypes   { return 3100008 }
func (UnsupportedFeature) What() string {
	return "Feature is currently unsupported"
}

type NodeManagementSuccess struct{ ELog }

func (NodeManagementSuccess) ChainExceptions() {}
func (NodeManagementSuccess) MiscExceptions()  {}
func (NodeManagementSuccess) Code() ExcTypes   { return 3100009 }
func (NodeManagementSuccess) What() string {
	return "Node management operation successfully executed"
}
