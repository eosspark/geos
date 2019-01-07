package fasthttp

import (
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type AsyncServer struct {
	*Server
	ctx          *asio.IoContext
	LogAllErrors bool
}

func ListenAndAsyncServe(ctx *asio.IoContext, addr string, handler RequestHandler) error {
	s := &AsyncServer{
		ctx:          ctx,
		LogAllErrors: true,
		Server: &Server{
			Handler: handler,
		},
	}

	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	if s.TCPKeepalive {
		if tcpln, ok := ln.(*net.TCPListener); ok {
			return s.Serve(tcpKeepaliveListener{
				TCPListener:     tcpln,
				keepalivePeriod: s.TCPKeepalivePeriod,
			})
		}
	}
	return s.Serve(ln)
}

func (s *AsyncServer) AsyncAccept(ln net.Listener, lastPerIPErrorTime time.Time, f func(c net.Conn)) {
	go s.accept(ln, lastPerIPErrorTime, f)
}

func (s *AsyncServer) accept(ln net.Listener, lastPerIPErrorTime time.Time, f func(c net.Conn)) {
	var c net.Conn
	var err error

	if c, err = acceptConn(s.Server, ln, &lastPerIPErrorTime); err != nil {
		if err == io.EOF {
			return
		}
		return
	}
	s.ctx.Post(func(err error) {
		f(c)
	})
}

func (s *AsyncServer) Serve(ln net.Listener) error {
	//var lastOverflowErrorTime time.Time
	var lastPerIPErrorTime time.Time
	s.mu.Lock()
	{
		if s.ln != nil {
			s.mu.Unlock()
			return ErrAlreadyServing
		}

		s.ln = ln
	}
	s.mu.Unlock()

	maxWorkersCount := s.getConcurrency()
	s.concurrencyCh = make(chan struct{}, maxWorkersCount)
	//wp := &workerPool{
	//	WorkerFunc:      s.serveConn,
	//	MaxWorkersCount: maxWorkersCount,
	//	LogAllErrors:    s.LogAllErrors,
	//	Logger:          s.logger(),
	//	connState:       s.setState,
	//}
	//wp.Start()

	// Count our waiting to accept a connection as an open connection.
	// This way we can't get into any weird state where just after accepting
	// a connection Shutdown is called which reads open as 0 because it isn't
	// incremented yet.
	atomic.AddInt32(&s.open, 1)
	defer atomic.AddInt32(&s.open, -1)

	var handleFunc func()
	handleFunc = func() {
		s.AsyncAccept(ln, lastPerIPErrorTime, func(c net.Conn) {
			s.setState(c, StateNew)
			atomic.AddInt32(&s.open, 1)

			//if !wp.Serve(c) {
			//if !wp.Serve(c) {
			//	atomic.AddInt32(&s.open, -1)
			//	s.writeFastError(c, StatusServiceUnavailable,
			//		"The connection cannot be served because Server.Concurrency limit exceeded")
			//	c.Close()
			//	s.setState(c, StateClosed)
			//	if time.Since(lastOverflowErrorTime) > time.Minute {
			//		s.logger().Printf("The incoming connection cannot be served, because %d concurrent connections are served. " +
			//			"Try increasing Server.Concurrency")
			//		lastOverflowErrorTime = time.Now()
			//	}
			//
			//	// The current server reached concurrency limit,
			//	// so give other concurrently running servers a chance
			//	// accepting incoming connections on the same address.
			//	//
			//	// There is a hope other servers didn't reach their
			//	// concurrency limits yet :)
			//	time.Sleep(100 * time.Millisecond)
			//}

			con := c
			s.ctx.Post(func(err error) {
				if con == nil {
					return
				}

				if err = s.serveConn(con); err != nil {
					if err != errHijacked {
						errStr := err.Error()
						if s.LogAllErrors || !(strings.Contains(errStr, "broken pipe") ||
							strings.Contains(errStr, "reset by peer") ||
							strings.Contains(errStr, "request headers: small read buffer") ||
							strings.Contains(errStr, "i/o timeout")) {
							s.Logger.Printf("error when serving connection %q<->%q: %s", con.LocalAddr(), con.RemoteAddr(), err)
							con.Close()
							s.setState(con, StateClosed)
						}
					} else if err == errHijacked {
						s.setState(con, StateHijacked)
					}
				}
			})

			c = nil
			handleFunc()
		})
	}

	handleFunc()

	return nil
}
