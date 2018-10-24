package client

import (
	"net"
	"fmt"
	"time"
)

func connect() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("Error net dial", err)
		return
	}

	time.Sleep(time.Second * 60)

	//_, werr := conn.Write([]byte("hello"))
	//
	//if werr != nil {
	//	fmt.Println("Error write", werr)
	//	return
	//}

	defer conn.Close()
}