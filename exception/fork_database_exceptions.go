package exception

import _ "github.com/eosspark/eos-go/log"

type ForkDatabaseException struct{ ELog }

func (ForkDatabaseException) ChainExceptions()        {}
func (ForkDatabaseException) ForkDatabaseExceptions() {}
func (ForkDatabaseException) Code() ExcTypes          { return 3020000 }
func (ForkDatabaseException) What() string            { return "Fork database exception" }

type ForkDbBlockNotFound struct{ ELog }

func (ForkDbBlockNotFound) ChainExceptions()        {}
func (ForkDbBlockNotFound) ForkDatabaseExceptions() {}
func (ForkDbBlockNotFound) Code() ExcTypes          { return 3020001 }
func (ForkDbBlockNotFound) What() string            { return "Block can not be found" }
