package main

import (
	"net"
	"fmt"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("Error net dial", err)
		return
	}

	defer conn.Close()

	for ;; {
		time.Sleep(time.Second)

		_, werr := conn.Write([]byte("hello"))

		if werr != nil {
			fmt.Println("Error write", werr)
			return
		}
	}


}