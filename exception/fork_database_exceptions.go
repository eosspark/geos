package exception

type ForkDatabaseException struct{ logMessage }

func (ForkDatabaseException) ChainExceptions()        {}
func (ForkDatabaseException) ForkDatabaseExceptions() {}
func (ForkDatabaseException) Code() ExcTypes          { return 3020000 }
func (ForkDatabaseException) What() string            { return "Fork database exception" }

type ForkDbBlockNotFound struct{ logMessage }

func (ForkDbBlockNotFound) ChainExceptions()        {}
func (ForkDbBlockNotFound) ForkDatabaseExceptions() {}
func (ForkDbBlockNotFound) Code() ExcTypes          { return 3020001 }
func (ForkDbBlockNotFound) What() string            { return "Block can not be found" }
