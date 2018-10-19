package database

type DataBase interface {

	Close()

	Insert(data interface{}) error

	Find(fieldName string, data interface{}, out interface{}) error

	GetIndex(fieldName string, data interface{}) (*multiIndex, error)

	Modify(data interface{}, fn interface{}) error

	Remove(data interface{}) error

	Undo()

	UndoAll()

	StartSession() *Session

	Commit(revision int64)

	SetRevision(revision int64)

	Revision() int64

	lowerBound(key,value,typeName []byte,data interface{},greater bool) (*DbIterator,error)

	upperBound(key,value,typeName []byte,data interface{},greater bool) (*DbIterator,error)

	squash()
}
