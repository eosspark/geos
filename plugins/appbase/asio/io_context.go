package asio

type IoContext struct {
	service ReactorService
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
	//use new goroutine for channel-blocking
	go i.GetService().post(op)
}

type ReactorService interface {
	run()
	stop()
	post(op interface{}, args ...interface{})
	notify(op interface{}, args ...interface{})
}

