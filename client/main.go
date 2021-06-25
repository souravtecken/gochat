package main

import (
	"github.com/rivo/tview"
	"github.com/souravtecken/gochat/client/ui"
)

func run(app *tview.Application, username *string, host *string) {
	ui.LoginPage(app, username, host)
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
