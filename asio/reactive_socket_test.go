package asio

import (
	"testing"
	"net"
	"fmt"
	)

func TestNet_Accept(t *testing.T) {
	listen, _ := net.Listen("tcp", "127.0.0.1:8888")
	for ;; {
		con, _ := listen.Accept()

		fmt.Println(con.RemoteAddr().String())

		buf := make([]byte, 10)
		con.Read(buf)
		fmt.Println(string(buf))
	}

}

func startRead(socket *ReactiveSocket, conn net.Conn) {
	buf := make([]byte, 64)
	socket.AsyncRead(conn, buf, func(n int, ec ErrorCode) {
		if ec.Valid {
			fmt.Println("Error read", ec.Error)
			return
		}
		if n > 0 {
			msg := string(buf[:n])
			fmt.Println(msg)
		}
		startRead(socket, conn)
	})
}

func startWrite(socket *ReactiveSocket, conn net.Conn) {
	msg := "i've received message"
	socket.AsyncWrite(conn, []byte(msg), func(n int, ec ErrorCode) {
		if ec.Valid {
			fmt.Println("Error write", ec.Error)
			return
		}
		startWrite(socket, conn)
	})
}

func startAcceptLoop(socket *ReactiveSocket, listen net.Listener) {
	socket.AsyncAccept(listen, func(conn net.Conn, ec ErrorCode) {
		//defer conn.Close()

		if ec.Valid {
			fmt.Println("Error connect", ec.Error)
		}

		fmt.Println(conn.RemoteAddr().String())

		startRead(socket, conn)
		startWrite(socket, conn)

		startAcceptLoop(socket, listen)
	})
}


func TestReactiveSocket_AsyncAccept(t *testing.T) {
	iosv := NewIoContext()
	socket := NewReactiveSocket(iosv)
	
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err)
	}

	startAcceptLoop(socket, listen)

	iosv.Run()
}
