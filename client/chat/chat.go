package chat

import (
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

func ReceiveMessages(conn *net.Conn, chatDisplay io.Writer) {
	buf := make([]byte, 2048)
	for {
		n, _ := (*conn).Read(buf[0:])
		message := string(buf[:n])
		fmt.Fprintf(chatDisplay, "%s\n", message)
	}
}
