package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func ConnectToServer(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	return conn, err
}

func SendMessage(conn *net.Conn, message string) {
	fmt.Fprintf(*conn, "%s", message)
}

func decodeMessage(jsonBytes []byte) map[string]interface{} {
	messageObj := make(map[string]interface{})
	json.Unmarshal(jsonBytes, &messageObj)
	return messageObj
}

func ReceiveMessages(conn *net.Conn, chatDisplay io.Writer) {
	buf := make([]byte, 2048)
	for {
		n, _ := (*conn).Read(buf[0:])
		data := decodeMessage(buf[:n])

		if data["type"] == "message" {
			// There's probably a cleaner way to do this
			messageBody := data["data"].(map[string]interface{})
			fmt.Fprintf(chatDisplay, "%s: %s\n", messageBody["username"].(string), messageBody["message"].(string))
		}
	}
}
