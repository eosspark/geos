package exception

type ForkDatabaseException struct{ logMessage }

func (e *ForkDatabaseException) ChainExceptions()        {}
func (e *ForkDatabaseException) ForkDatabaseExceptions() {}
func (e *ForkDatabaseException) Code() ExcTypes          { return 3020000 }
func (e *ForkDatabaseException) What() string            { return "Fork database exception" }

type ForkDbBlockNotFound struct{ logMessage }

func (e *ForkDbBlockNotFound) ChainExceptions()        {}
func (e *ForkDbBlockNotFound) ForkDatabaseExceptions() {}
func (e *ForkDbBlockNotFound) Code() ExcTypes          { return 3020001 }
func (e *ForkDbBlockNotFound) What() string            { return "Block can not be found" }
