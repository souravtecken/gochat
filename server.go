package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func encodeData(data interface{}, messageType string) []byte {
	messageObj := make(map[string]interface{})
	messageObj["data"] = data
	messageObj["type"] = messageType
	jsonBytes, _ := json.Marshal(&messageObj)
	return jsonBytes
}

func broadcastData(userConMap map[string]net.Conn, data []byte) {
	for _, conn := range userConMap {
		fmt.Fprintf(conn, "%s\n", data)
	}
}

func sendAllUsers(conn *net.Conn, userConMap map[string]net.Conn) {
	users := make([]string, 0, len(userConMap))
	for user := range userConMap {
		users = append(users, user)
	}
	data := encodeData(users, "users")
	fmt.Fprintf(*conn, "%s\n", data)
	fmt.Println(string(data))
}

func sendUserJoin(user string, userConMap map[string]net.Conn) {
	data := encodeData(user, "userJoin")
	broadcastData(userConMap, data)
}

func sendUserLeave(user string, userConMap map[string]net.Conn) {
	data := encodeData(user, "userLeave")
	broadcastData(userConMap, data)
}

func sendMessage(userConMap map[string]net.Conn, message string, username string) {
	messageBody := make(map[string]interface{})
	messageBody["message"] = message
	messageBody["username"] = username
	data := encodeData(messageBody, "message")
	broadcastData(userConMap, data)
}

func handleConnection(conn net.Conn, userConMap map[string]net.Conn) {
	fmt.Println("Connection open:", conn.RemoteAddr())

	buf := make([]byte, 512)

	// Initial signup message containing username
	n, _ := conn.Read(buf[0:])
	username := string(buf[:n])

	// Send list of all current users
	sendAllUsers(&conn, userConMap)

	// Add user to connection map
	userConMap[username] = conn

	// Inform all users about joining
	sendUserJoin(username, userConMap)

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

	// Inform all users about leaving
	sendUserLeave(username, userConMap)
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
