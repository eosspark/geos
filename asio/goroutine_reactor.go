package gosio

import (
	"reflect"
	"fmt"
)

type GoroutineReactor struct {
	opQueue events
}

type event struct {
	op   interface{}
	args []interface{}
}

type events chan event

func NewGouroutineReactor() *GoroutineReactor {
	r := new(GoroutineReactor)
	r.opQueue = make(events, 128)
	return r
}

func (g *GoroutineReactor) run() {
	for {
		select {
		case op := <-g.opQueue:
			g.doReactor(op.op, op.args)
			break
		}
	}
}

func (g *GoroutineReactor) push(op interface{}, args ...interface{}) {
	g.opQueue <- event{op, args}
}

func (g *GoroutineReactor) doReactor(op interface{}, args []interface{}) {
	opv := reflect.ValueOf(op)
	opt := reflect.TypeOf(op)

	if opt.Kind() != reflect.Func {
		fmt.Println("opt is not a function")
		return
	}

	opNum := opt.NumIn()
	if opNum != len(args) {
		fmt.Println("invalid arguments", "opNum:", opNum)
		return
	}

	opArgs := make([]reflect.Value, 0, opNum)

	for i:=0; i<opt.NumIn(); i++ {
		if args[i] != nil {
			opArgs = append(opArgs, reflect.ValueOf(args[i]))
		}
	}

	opv.Call(opArgs)
}

