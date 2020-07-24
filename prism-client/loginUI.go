package main

import (
	"log"
	"github.com/marcusolsson/tui-go"
)


var logo = `$$$$$$$\            $$\
$$  __$$\           \__|
$$ |  $$ | $$$$$$\  $$\  $$$$$$$\ $$$$$$\$$$$\
$$$$$$$  |$$  __$$\ $$ |$$  _____|$$  _$$  _$$\
$$  ____/ $$ |  \__|$$ |\$$$$$$\  $$ / $$ / $$ |
$$ |      $$ |      $$ | \____$$\ $$ | $$ | $$ |
$$ |      $$ |      $$ |$$$$$$$  |$$ | $$ | $$ |
\__|      \__|      \__|\_______/ \__| \__| \__|
 by joshua-isak                     client ` + VERSION


func loginUI() (string, string, []byte){

	server := tui.NewEntry()
	server.SetFocused(true)

	user := tui.NewEntry()

	key := tui.NewEntry()
	key.SetEchoMode(tui.EchoModePassword)

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("IP"), tui.NewLabel("Username"), tui.NewLabel("32-Byte Key"))
	form.AppendRow(server, user, key)

	status := tui.NewStatusBar("Ready. Press esc to quit")

	login := tui.NewButton("[Login]")

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, login),
	)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(25, 1, tui.NewLabel("Connect to a server")),
		tui.NewPadder(5, 1, form),
		buttons,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(
		content,
		status,
	)

	tui.DefaultFocusChain.Set(server, user, key, login)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	login.OnActivated(func(b *tui.Button) {
		// Make sure the username is not longer than 20 characters
		if len(user.Text()) > 20 {
			status.SetText("Username cannot exceed 20 characters in length!")
			user.SetText("")
			user.SetFocused(true)
			login.SetFocused(false)
			return
		}

		// Make sure the key is 32 bytes long
		if len(key.Text()) != 32 {
			status.SetText("Key must be exactly 32 bytes in length!")
			key.SetText("")
			key.SetFocused(true)
			login.SetFocused(false)
			return
		}

		status.SetText("Logged in!")
		ui.Quit()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

	return server.Text(), user.Text(), []byte(key.Text())

}