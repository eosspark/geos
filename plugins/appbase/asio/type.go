package asio

type Path = string

type socketReadFullOp = socketReadOp
type socketWriteOp = socketReadOp
type signalSetOp = deadlineTimerOp
type postOp = deadlineTimerOp