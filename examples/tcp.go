package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")

	for {
		conn, _ := listener.Accept()
		buf := make([]byte, 1024)
		conn.Read(buf)
		fmt.Println(string(buf))
	}
}
