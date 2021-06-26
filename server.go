package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

func encodeData(data interface{}, messageType string) []byte {
	messageObj := make(map[string]interface{})
	messageObj["data"] = data
	messageObj["type"] = messageType
	jsonBytes, _ := json.Marshal(&messageObj)
	jsonBytes = bytes.Trim(jsonBytes, "\x00")
	return jsonBytes
}

func broadcastMessage(userConMap map[string]net.Conn, message []byte) {
	for _, conn := range userConMap {
		conn.Write(message)
	}
}

func sendMessage(useConMap map[string]net.Conn, message string, username string) {
	messageBody := make(map[string]interface{})
	messageBody["message"] = message
	messageBody["username"] = username
	data := encodeData(messageBody, "message")
	fmt.Println(string(data))
	broadcastMessage(useConMap, data)
}

func handleConnection(conn net.Conn, userConMap map[string]net.Conn) {
	fmt.Println("Connection open:", conn.RemoteAddr())

	buf := make([]byte, 512)

	// Initial signup message containing username
	n, _ := conn.Read(buf[0:])
	username := string(buf[:n])

	// Add user to connection map
	userConMap[username] = conn

	fmt.Println(conn.RemoteAddr(), "is", username)
	for {
		n, _ := conn.Read(buf[0:])
		if n == 0 {
			// Connection has closed
			break
		}
		message := string(buf[:n])
		go sendMessage(userConMap, message, username)
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
