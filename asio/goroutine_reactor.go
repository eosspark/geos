package asio

import (
	"reflect"
	"fmt"
	"os"
	)

type GoroutineReactor struct {
	opq 	 chan operation
	//notifies chan os.Signal
	down chan struct{}
}

type operation struct {
	function  interface{}
	argument  []interface{}
}

func NewGouroutineReactor() *GoroutineReactor {
	r := new(GoroutineReactor)
	r.opq = make(chan operation, 128)
	r.down = make(chan struct{}, 1)
	//r.notifies = make(chan os.Signal, 1)
	return r
}

func (g *GoroutineReactor) run() {
	for ;; {
		select {
		case <-g.down:
			return
		case op := <-g.opq:
			g.doReactor(op.function, op.argument)
			break
		}
	}
}

func (g *GoroutineReactor) stop() {
	g.down <- struct{}{}
}

func (g *GoroutineReactor) push(op interface{}, args ...interface{}) {
	g.opq <- operation{op, args}
}

func (g *GoroutineReactor) notify(sig ...os.Signal) {
	//signal.Notify(g.notifies, sig...)
}

func (g *GoroutineReactor) doReactor(op interface{}, args []interface{}) {
	opv := reflect.ValueOf(op)
	opt := reflect.TypeOf(op)

	if opt.Kind() != reflect.Func {
		fmt.Println("opt must be a callback function")
		return
	}

	opNum := opt.NumIn()
	if opNum != len(args) {
		fmt.Println("invalid arguments", "opNum:", opNum)
		return
	}

	opArgs := make([]reflect.Value, opNum)

	for i:=0; i<opt.NumIn(); i++ {
		if args[i] != nil {
			opArgs[i] = reflect.ValueOf(args[i])
		}
	}

	opv.Call(opArgs)
}

