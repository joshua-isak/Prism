package main

import (
	"github.com/marcusolsson/tui-go"
	"log"
	"net"
)


func chatUI(username string, c net.Conn, key []byte, history *tui.Box, users *tui.List, address string, u *uiThing) {

	sidebar := tui.NewVBox(
		tui.NewPadder(1, 0, tui.NewLabel(" CONNECTED: ")),
		tui.NewPadder(1, 1, users),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	topInfo := tui.NewStatusBar("  Prism goClient " + VERSION + "  |  Server: " + address + "  |  Username: " + username + "  |  " +  "Press esc to quit")

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
		// Truncate messages longer than 200
		msg := []byte(entry.Text())
		if len(msg) > 200 {
			msg = msg[0:199]
		}

		// Encrypt user text
		msg = encrypt(msg , key)

		// Send message to server
		p := NewPacket(GeneralMessage)
		p.PrepGeneralMessage(username, msg, true)
		p.Send(c)

		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)
	root = tui.NewVBox(topInfo, root)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// Pointers :(
	u.ui = ui

	ui.SetKeybinding("Esc", func() {
		c.Close()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}