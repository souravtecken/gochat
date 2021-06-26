package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/souravtecken/gochat/client/chat"
)

func LoginPage(app *tview.Application, username *string, host *string) {

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
			ChatPage(app, *username, *host)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Login")

	app.SetRoot(form, true).SetFocus(form).Run()
}

func ChatPage(app *tview.Application, username string, host string) {

	message := ""
	conn, err := chat.ConnectToServer(host)
	chat.SendMessage(&conn, username)
	if err != nil {
		panic(err)
	}

	chatPane := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	usersPane := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	inputField := tview.NewInputField()

	go chat.ReceiveMessages(&conn, chatPane, usersPane)

	handleSendMessage := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			chat.SendMessage(&conn, message)
			inputField.SetText("")
			message = ""
		}
	}

	handleInputMessage := func(text string) {
		message = text
	}

	inputField.SetDoneFunc(handleSendMessage)
	inputField.SetChangedFunc(handleInputMessage)

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
