package asio

import (
	"net"
	"io"
)

type ReactiveSocket struct {
	ctx *IoContext
}

func NewReactiveSocket(ctx *IoContext) *ReactiveSocket {
	d := new(ReactiveSocket)
	d.ctx = ctx
	return d
}

func (r *ReactiveSocket) AsyncRead(reader io.Reader, b []byte, op func(n int, ec ErrorCode)) {
	go func() {
		n, err := reader.Read(b)
		r.ctx.GetService().push(op, n, NewErrorCode(err))
	}()
}


func (r *ReactiveSocket) AsyncWrite(writer io.Writer, b []byte, op func(n int, ec ErrorCode)) {
	go func() {
		n, err := writer.Write(b)
		r.ctx.GetService().push(op, n, NewErrorCode(err))
	}()

}

func (r *ReactiveSocket) AsyncAccept(listen net.Listener, op func(conn net.Conn, ec ErrorCode)) {
	go func() {
		connect, err := listen.Accept()
		r.ctx.GetService().push(op, connect, NewErrorCode(err))
	}()
}



