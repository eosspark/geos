package exception

type MongoDbException struct{ logMessage }

func (e *MongoDbException) ChainExceptions()   {}
func (e *MongoDbException) MongoDbExceptions() {}
func (e *MongoDbException) Code() ExcTypes     { return 3220000 }
func (e *MongoDbException) What() string {
	return "Mongo DB exception"
}
