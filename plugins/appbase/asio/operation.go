package asio

import "net"

type operation interface {
	call()
}

type deadlineTimerOp struct {
	function func(err error)
	err 	 error
}

func (op deadlineTimerOp) call() {
	op.function(op.err)
}

type socketAcceptOp struct {
	function func(conn net.Conn, err error)
	conn	 net.Conn
	err 	 error
}

func (op socketAcceptOp) call() {
	op.function(op.conn, op.err)
}

type socketReadOp struct {
	function func(n int, err error)
	n	 	 int
	err 	 error
}

func (op socketReadOp) call() {
	op.function(op.n, op.err)
}

type socketReadFullOp = socketReadOp
type socketWriteOp = socketReadOp
type signalSetOp = deadlineTimerOp
type postOp = deadlineTimerOp
