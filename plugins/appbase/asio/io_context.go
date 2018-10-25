package asio

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

func (i *IoContext) Stop() {
	i.GetService().stop()
}

func (i *IoContext) Post(op func()) {
	i.GetService().post(op)
}

