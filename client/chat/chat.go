package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"math/rand"

	"github.com/rivo/tview"
)

var (
	colors = [...]string{"red", "orange", "yellow", "green"}
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
	usersDisplay *tview.TextView, users map[string]int, wg *sync.WaitGroup, mut *sync.Mutex) {

	data := decodeMessage([]byte(buf))

	if data["type"] == "users" {
		usersList := data["data"].([]interface{})
		for _, user := range usersList {
			users[user.(string)] = rand.Intn(len(colors))
			fmt.Fprintf(usersDisplay, "[%s]%s[white]\n", users[user.(string)], user)
		}
		wg.Done()
	}

	wg.Wait()
	mut.Lock()

	if data["type"] == "message" {
		// There's probably a cleaner way to do this
		messageBody := data["data"].(map[string]interface{})
		username := messageBody["username"].(string)
		message := messageBody["message"].(string)
		fmt.Fprintf(chatDisplay, "[%s]%s[white]: %s\n", colors[users[username]], username, message)
	} else if data["type"] == "userJoin" {
		users[data["data"].(string)] = rand.Intn(len(colors))
		usersDisplay.Clear()
		for user, color := range users {
			fmt.Fprintf(usersDisplay, "[%s]%s[white]\n", colors[color], user)
		}
	} else if data["type"] == "userLeave" {
		delete(users, data["data"].(string))
		usersDisplay.Clear()
		for user, color := range users {
			fmt.Fprintf(usersDisplay, "[%s]%s[white]\n", colors[color], user)
		}
	}

	mut.Unlock()
}

func ReceiveMessages(conn *net.Conn, chatDisplay io.Writer, usersDisplay *tview.TextView) {

	users := make(map[string]int)
	reader := bufio.NewReader(*conn)
	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(1)
	for {
		buf, _ := reader.ReadString('\n')
		go handleReceiveMessage(buf, conn, chatDisplay, usersDisplay, users, &wg, &mut)
	}
}
