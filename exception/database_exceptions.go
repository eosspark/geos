package exception

type DatabaseException struct{ logMessage }

func (DatabaseException) ChainExceptions()    {}
func (DatabaseException) DatabaseExceptions() {}
func (DatabaseException) Code() ExcTypes      { return 3060000 }
func (DatabaseException) What() string        { return "Database exception" }

type PermissionQueryException struct{ logMessage }

func (PermissionQueryException) ChainExceptions()    {}
func (PermissionQueryException) DatabaseExceptions() {}
func (PermissionQueryException) Code() ExcTypes      { return 3060001 }
func (PermissionQueryException) What() string        { return "Permission Query Exception" }

type AccountQueryException struct{ logMessage }

func (AccountQueryException) ChainExceptions()    {}
func (AccountQueryException) DatabaseExceptions() {}
func (AccountQueryException) Code() ExcTypes      { return 3060002 }
func (AccountQueryException) What() string        { return "Account Query Exception" }

type ContractTableQueryException struct{ logMessage }

func (ContractTableQueryException) ChainExceptions()    {}
func (ContractTableQueryException) DatabaseExceptions() {}
func (ContractTableQueryException) Code() ExcTypes      { return 3060003 }
func (ContractTableQueryException) What() string        { return "Contract Table Query Exception" }

type ContractQueryException struct{ logMessage }

func (ContractQueryException) ChainExceptions()    {}
func (ContractQueryException) DatabaseExceptions() {}
func (ContractQueryException) Code() ExcTypes      { return 3060004 }
func (ContractQueryException) What() string        { return "Contract Query Exception" }

// implements GuardExceptions
type GuardException struct{ logMessage }

func (GuardException) ChainExceptions()    {}
func (GuardException) GuardExceptions()    {}
func (GuardException) DatabaseExceptions() {}
func (GuardException) Code() ExcTypes      { return 3060100 }
func (GuardException) What() string        { return "Database exception" }

type DatabaseGuardException struct{ logMessage }

func (DatabaseGuardException) ChainExceptions()    {}
func (DatabaseGuardException) GuardExceptions()    {}
func (DatabaseGuardException) DatabaseExceptions() {}
func (DatabaseGuardException) Code() ExcTypes      { return 3060101 }
func (DatabaseGuardException) What() string        { return "Database usage is at unsafe levels" }

type ReversibleGuardException struct{ logMessage }

func (ReversibleGuardException) ChainExceptions()    {}
func (ReversibleGuardException) GuardExceptions()    {}
func (ReversibleGuardException) DatabaseExceptions() {}
func (ReversibleGuardException) Code() ExcTypes      { return 3060102 }
func (ReversibleGuardException) What() string {
	return "Reversible block log usage is at unsafe levels"
}
