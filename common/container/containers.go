package container

import "github.com/eosspark/eos-go/log"

type Container interface {
	Empty() bool
	Size() int
	Clear()

	Serializer
}

type Serializer interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type Set interface {
	Container
}

type Map interface {
	Container
}

var Logger = log.NewWithHandle("container::log", log.DiscardHandler())