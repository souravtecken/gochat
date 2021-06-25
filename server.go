package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	buf := make([]byte, 512)
	n, _ := conn.Read(buf[0:])
	fmt.Println(string(buf), n)
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err == nil {
			handleConnection(conn)
		}
	}
}
