package database

type DataBase interface {
	Close()

	Insert(in interface{}) error

	Find(tagName string, in interface{}, out interface{},skip... SkipSuffix) error

	Empty(begin, end, fieldName []byte) bool

	GetIndex(tagName string, in interface{}) (*MultiIndex, error)

	GetMutableIndex(tagName string, in interface{}) (*MultiIndex, error)

	Modify(data interface{}, fn interface{}) error

	Remove(data interface{}) error

	Undo()

	UndoAll()

	StartSession() *Session

	Commit(revision int64)

	SetRevision(revision int64)

	Revision() int64

	lowerBound(key, value, typeName []byte, in interface{},skip... SkipSuffix) (*DbIterator, error)

	upperBound(key, value, typeName []byte, in interface{},skip... SkipSuffix) (*DbIterator, error)

	IteratorTo(begin, end, fieldName []byte, in interface{},skip...SkipSuffix) (*DbIterator, error)

	BeginIterator(begin, end, typeName []byte) (*DbIterator, error)
	EndIterator(begin, end, typeName []byte) (*DbIterator, error)

	squash()
}
