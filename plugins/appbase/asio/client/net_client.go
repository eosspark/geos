package main

import (
	"net"
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	)

func main() {
	for i:=0; i<100; i++ {
		go doDial()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	select {
	case <- c:
		return
	}
}

func doDial() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("Error net dial", err)
		return
	}

	defer conn.Close()

	for ;; {
		time.Sleep(time.Second)


		_, werr := conn.Write([]byte("hello: " + conn.LocalAddr().String()))

		if werr != nil {
			fmt.Println("Error write", werr)
			return
		}
	}
}