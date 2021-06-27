package main

import (
	"github.com/rivo/tview"
	"github.com/souravtecken/gochat/client/ui"
)

func run(app *tview.Application, username *string, host *string) {
	ui.LoginPage(app, username, host)
}

func main() {
	app := tview.NewApplication()

	username := ""
	host := ""

	run(app, &username, &host)
}
