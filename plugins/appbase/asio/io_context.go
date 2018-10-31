package asio

import "reflect"

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

type operation struct {
	function  interface{}
	argument  []interface{}
}

type ReactorService interface {
	run()
	stop()
	post(op interface{}, args ...interface{})
	notify(op interface{}, args ...interface{})
}

func doReactor(opr operation) {
	opv := reflect.ValueOf(opr.function)
	opt := reflect.TypeOf(opr.function)

	if opt == nil {
		println("invalid operation <nil>")
	}

	if opt.Kind() != reflect.Func {
		println("op must be a callback function")
		return
	}

	opNum := opt.NumIn()
	if opNum != len(opr.argument) {
		println("invalid arguments", "arguments needs:", opNum)
		return
	}

	opArgs := make([]reflect.Value, opNum)

	for i:=0; i<opt.NumIn(); i++ {
		if opr.argument[i] == nil {
			opArgs[i] = reflect.Zero(opt.In(i))
			continue
		}

		if !reflect.TypeOf(opr.argument[i]).AssignableTo(opt.In(i)) {
			println("invalid arguments", "wrong args#", i)
			return
		}

		opArgs[i] = reflect.ValueOf(opr.argument[i])
	}

	//fmt.Println("args", opArgs)

	opv.Call(opArgs)
}

