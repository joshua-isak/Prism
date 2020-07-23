package main

import (
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
)


func chatUI(username string, c net.Conn, key []byte, history *tui.Box, address string, u *uiThing) {
	sidebar := tui.NewVBox(
		tui.NewLabel("pRism v0.1   "),
		tui.NewLabel(""),
		tui.NewLabel("Server:"),
		tui.NewLabel(address + " "),
		tui.NewLabel(""),
		tui.NewLabel("Username:"),
		tui.NewLabel(username + " "),
		tui.NewLabel(""),
		tui.NewLabel(""),
		tui.NewLabel("Press esc"),
		tui.NewLabel("to quit"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(false)

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(entry *tui.Entry) {
		// Encrypt user text
		msg := encrypt([]byte(entry.Text()), key)

		// Send message to server
		p := NewPacket(GeneralMessage)
		p.PrepGeneralMessage(username, msg, true)
		p.Send(c)

		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// Pointers :(
	u.ui = ui

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}