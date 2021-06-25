package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func encodeData(data interface{}, messageType string) []byte {
	messageObj := make(map[string]interface{})
	messageObj["type"] = messageType
	messageObj["data"] = data
	jsonString, _ := json.Marshal(messageObj)
	fmt.Println(jsonString)
	return []byte(jsonString)
}

func broadcastMessage(userConMap map[string]net.Conn, message string, username string) {
	for _, conn := range userConMap {
		fmt.Fprintf(conn, message)
	}
}

func handleConnection(conn net.Conn, userConMap map[string]net.Conn) {
	fmt.Println("Connection open:", conn.RemoteAddr())

	buf := make([]byte, 512)
	conn.Read(buf[0:]) // Initial signup message
	username := string(buf)
	userConMap[username] = conn
	fmt.Println(conn.RemoteAddr(), "is", username)
	for {
		n, _ := conn.Read(buf[0:])
		if n == 0 { // Connection has closed
			break
		}
		message := string(buf[:n])
		go broadcastMessage(userConMap, message, username)
		fmt.Println(message, n)
	}
	delete(userConMap, username)
	fmt.Println("Connection closed:", conn.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	userConMap := make(map[string]net.Conn)

	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := ln.Accept()
		if err == nil {
			go handleConnection(conn, userConMap)
		}
	}
}
