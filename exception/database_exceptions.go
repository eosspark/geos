package exception


type DatabaseException struct{ logMessage }

func (e *DatabaseException) ChainExceptions()    {}
func (e *DatabaseException) DatabaseExceptions() {}
func (e *DatabaseException) Code() ExcTypes      { return 3060000 }
func (e *DatabaseException) What() string        { return "Database exception" }


type GuardException struct{ logMessage }

func (e *GuardException) ChainExceptions()    {}
func (e *GuardException) GuardExceptions()    {}
func (e *GuardException) DatabaseExceptions() {}
func (e *GuardException) Code() ExcTypes      { return 3060100 }
func (e *GuardException) What() string        { return "Database exception" }
