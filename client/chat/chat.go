package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rivo/tview"
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

func handleReceiveMessage(buf string, conn *net.Conn, chatDisplay io.Writer,
	usersDisplay *tview.TextView, users map[string]bool, wg *sync.WaitGroup, mut *sync.Mutex) {

	data := decodeMessage([]byte(buf))

	if data["type"] == "users" {
		usersList := data["data"].([]interface{})
		for _, user := range usersList {
			users[user.(string)] = true
			fmt.Fprintf(usersDisplay, "%s\n", user)
		}
		wg.Done()
	}

	wg.Wait()
	mut.Lock()

	if data["type"] == "message" {
		// There's probably a cleaner way to do this
		messageBody := data["data"].(map[string]interface{})
		fmt.Fprintf(chatDisplay, "%s: %s\n", messageBody["username"].(string), messageBody["message"].(string))
	} else if data["type"] == "userJoin" {
		users[data["data"].(string)] = true
		usersDisplay.Clear()
		for user := range users {
			fmt.Fprintf(usersDisplay, "%s\n", user)
		}
	} else if data["type"] == "userLeave" {
		delete(users, data["data"].(string))
		usersDisplay.Clear()
		for user := range users {
			fmt.Fprintf(usersDisplay, "%s\n", user)
		}
	}

	mut.Unlock()
}

func ReceiveMessages(conn *net.Conn, chatDisplay io.Writer, usersDisplay *tview.TextView) {

	users := make(map[string]bool)
	reader := bufio.NewReader(*conn)
	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(1)
	for {
		buf, _ := reader.ReadString('\n')
		go handleReceiveMessage(buf, conn, chatDisplay, usersDisplay, users, &wg, &mut)
	}
}
