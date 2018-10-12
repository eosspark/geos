package database



type DataBase interface {
	Insert(data interface{})

	Find(fieldName string, data interface{}) (DbIterator, error)

	Get(fieldName string, data interface{}) (Iterator, error)

	Modify(data interface{}, fn interface{}) error

	Remove(data interface{}) error
}
