package multiindex

type IndexType interface {
	GetSuperIndex() interface{}
	GetFinalIndex() interface{}
}

type NodeType interface {
	GetSuperNode() interface{}
	GetFinalNode() interface{}
}

type IteratorType interface {
	IsEnd() bool
	HasNext() bool
}

type ReverseIteratorType interface {
	IsBegin() bool
	HasPrev() bool
}