package database

type DataBase interface {
	Close()

	Insert(in interface{}) error

	Find(tagName string, in interface{}, out interface{}) error

	Empty(begin, end, fieldName []byte) bool

	GetIndex(tagName string, in interface{}) (*multiIndex, error)

	GetMutableIndex(tagName string, in interface{}) (*multiIndex, error)

	Modify(data interface{}, fn interface{}) error

	Remove(data interface{}) error

	Undo()

	UndoAll()

	StartSession() *Session

	Commit(revision int64)

	SetRevision(revision int64)

	Revision() int64

	lowerBound(key, value, typeName []byte, in interface{}, greater bool) (*DbIterator, error)

	upperBound(key, value, typeName []byte, in interface{}, greater bool) (*DbIterator, error)

	squash()
}
