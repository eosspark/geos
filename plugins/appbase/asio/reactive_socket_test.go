package asio

import (
	"testing"
	"net"
	"fmt"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestReactiveSocket_AsyncAccept(t *testing.T) {
	ioc := NewIoContext()
	socket := NewReactiveSocket(ioc)
	connects := make([]net.Conn, 0)

	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err)
	}

	startAcceptLoop(t, socket, listen, &connects)

	const COUNT = 1000
	var doWrite func()
	var (
		stop 	= false
		dialogs = make([]net.Conn, COUNT)
		index 	= 0
	)

	go func() {
		for i:=0; i<COUNT && !stop; i++ {
			conn, err := net.Dial("tcp", ":8888")
			if err != nil {
				t.Fatal(err)
				return
			}

			dialogs[index] = conn
			index ++

			doWrite = func () {
				time.Sleep(time.Second)
				conn.Write([]byte("hello"))
				if !stop {
					doWrite()
				}
			}

			go doWrite()
		}
	}()

	shutdown := NewDeadlineTimer(ioc)
	shutdown.ExpiresFromNow(time.Second * 10)
	shutdown.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		stop = true
		ioc.Stop()
	})

	ioc.Run()

	for _,c := range connects {
		c.Close()
	}

	for i:=0; i<index && i<COUNT; i++ {
		dialogs[i].Close()
	}

	assert.Equal(t, COUNT, len(connects))
}

func startAcceptLoop(t *testing.T, socket *ReactiveSocket, listen net.Listener, connects *[]net.Conn) {
	socket.AsyncAccept(listen, func(conn net.Conn, err error) {
		//defer conn.Close()

		if conn == nil {
			fmt.Println("Error connect, nil")
			startAcceptLoop(t, socket, listen, connects)
			return
		}

		if err != nil {
			conn.Close()
			fmt.Println("Error connect", err)
			startAcceptLoop(t, socket, listen, connects)
			return
		}

		fmt.Println(conn.RemoteAddr().String())

		*connects = append(*connects, conn)

		startRead(t, socket, conn)

		startAcceptLoop(t, socket, listen, connects)
	})
}



func startRead(t *testing.T, socket *ReactiveSocket, conn net.Conn) {
	buf := make([]byte, 64)
	socket.AsyncRead(conn, buf, func(n int, err error) {
		if err != nil {
			t.Fatal(err)
			return
		}
		if n > 0 {
			//msg := string(buf[:n])
			//fmt.Println(msg)
		}
		startRead(t, socket, conn)
	})
}


