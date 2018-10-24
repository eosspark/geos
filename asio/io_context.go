package gosio

type IoContext struct {
	reactor Reactor
}

func NewIoContext() *IoContext {
	i := new(IoContext)

	return i
}

func (i *IoContext) Run() {
	i.GetService().run()
}

