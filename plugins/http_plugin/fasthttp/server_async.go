package fasthttp

import (
	"net"
	"time"
	"io"
	"fmt"
	"sync/atomic"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
)

type AsyncServer struct {
	ctx *asio.IoContext
	*Server
}

func ListenAndAsyncServe(ctx *asio.IoContext, addr string, handler RequestHandler) error {
	s := &AsyncServer{
		ctx: ctx,
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

func (s *AsyncServer) AsyncAccept(ln net.Listener, lastPerIPErrorTime time.Time, wp *workerPool, maxWorkersCount int, f func(c net.Conn)) {
	go s.accept(ln, lastPerIPErrorTime, wp, maxWorkersCount, f)
}

func (s *AsyncServer) accept(ln net.Listener, lastPerIPErrorTime time.Time, wp *workerPool, maxWorkersCount int, f func(c net.Conn)) {
	var c net.Conn
	var err error

	if c, err = acceptConn(s.Server, ln, &lastPerIPErrorTime); err != nil {
		wp.Stop()
		if err == io.EOF {
			fmt.Println("accept err", err.Error())
		}
		fmt.Println("accept err", err.Error())
	}
	s.ctx.Post(func(err error) {
		f(c)
	})
}

func (s *AsyncServer) Serve(ln net.Listener) error {
	var lastOverflowErrorTime time.Time
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
	wp := &workerPool{
		WorkerFunc:      s.serveConn,
		MaxWorkersCount: maxWorkersCount,
		LogAllErrors:    s.LogAllErrors,
		Logger:          s.logger(),
		connState:       s.setState,
	}
	wp.Start()

	// Count our waiting to accept a connection as an open connection.
	// This way we can't get into any weird state where just after accepting
	// a connection Shutdown is called which reads open as 0 because it isn't
	// incremented yet.
	atomic.AddInt32(&s.open, 1)
	defer atomic.AddInt32(&s.open, -1)

	var handleFunc func()
	handleFunc = func() {
		s.AsyncAccept(ln, lastPerIPErrorTime, wp, maxWorkersCount, func(c net.Conn) {
			s.setState(c, StateNew)
			atomic.AddInt32(&s.open, 1)
			if !wp.Serve(c) {
				atomic.AddInt32(&s.open, -1)
				s.writeFastError(c, StatusServiceUnavailable,
					"The connection cannot be served because Server.Concurrency limit exceeded")
				c.Close()
				s.setState(c, StateClosed)
				if time.Since(lastOverflowErrorTime) > time.Minute {
					s.logger().Printf("The incoming connection cannot be served, because %d concurrent connections are served. "+
						"Try increasing Server.Concurrency", maxWorkersCount)
					lastOverflowErrorTime = time.Now()
				}

				// The current server reached concurrency limit,
				// so give other concurrently running servers a chance
				// accepting incoming connections on the same address.
				//
				// There is a hope other servers didn't reach their
				// concurrency limits yet :)
				time.Sleep(100 * time.Millisecond)
			}
			c = nil
			handleFunc()
		})
	}

	handleFunc()

	return nil
}
