package main

import (
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func loginPage(app *tview.Application, username *string, host *string) {

	handleUsernameInput := func(text string) {
		*username = text
	}

	handleHostInput := func(text string) {
		*host = text
	}

	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, handleUsernameInput).
		AddInputField("Host", "", 20, nil, handleHostInput).
		AddButton("Login", func() {
			chatPage(app, *username, *host)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Login")

	app.SetRoot(form, true).SetFocus(form).Run()
}

func connectToServer(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	return conn, err
}

func sendMessage(conn *net.Conn, message string) {
	fmt.Fprintf(*conn, message)
}

func chatPage(app *tview.Application, username string, host string) {

	message := ""
	conn, err := connectToServer(host)
	if err != nil {
		panic(err)
	}

	handleSendMessage := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			sendMessage(&conn, message)
		}
	}

	handleInputMessage := func(text string) {
		message = text
	}

	chatPane := tview.NewBox().SetBorder(true).SetTitle("Messages")
	usersPane := tview.NewBox().SetBorder(true).SetTitle("Users")
	inputField := tview.NewInputField().
		SetDoneFunc(handleSendMessage).
		SetChangedFunc(handleInputMessage)

	leftPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatPane, 0, 1, false).
		AddItem(inputField, 2, 1, true)

	flex := tview.NewFlex().
		AddItem(leftPane, 0, 5, false).
		AddItem(usersPane, 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
}

func run(app *tview.Application, username *string, host *string) {
	loginPage(app, username, host)
}

func main() {
	// conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	// fmt.Println("Hello")
	// fmt.Fprintf(conn, "Hello there")
	// message, _ := bufio.NewReader(conn).ReadString('\n')
	// fmt.Println(message)
	app := tview.NewApplication()

	username := ""
	host := ""

	run(app, &username, &host)

}
