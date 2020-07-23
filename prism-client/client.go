package main

import (
	"fmt"
	"net"
	"github.com/marcusolsson/tui-go"
	"time"
)



// Basically Handle the GeneralMessage packet type
func printServerMessage(connection net.Conn, key []byte, history *tui.Box, u *uiThing) {
	defer connection.Close()
	for {
		// Read in data from tcp socket and put it in a Packet object
		buf := make([]byte, 1024)	// read up to 1024 bytes into buf
		_, err := connection.Read(buf[0:])	// read up to size of buf
		if err != nil {
			fmt.Println(err)
			return
		}
		netData := NewPacket(Received)
		netData.data = buf

		// Close connection if server didn't send a GeneralMessage PacketType
		pType := netData.ReadUint8()
		if pType != 5 { return }

		l := netData.ReadUint8()						// read in the senderName length
		senderName := netData.ReadString(int(l))		// read in the senderName as a string

		// Check if the message is encrypted
		netData.seek = 23
		isEncrypted := netData.ReadBool()

		// Read in the message as a byte array
		messageSize := netData.ReadUint8()
		message := netData.ReadBytes(int(messageSize))

		// Decrypt the message if it is encrypted
		if isEncrypted {
			message = decrypt(message, key)
		}

		// Print message using textUI ;)
		u.ui.Update( func(){
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format("15:04")),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", senderName))),
				tui.NewLabel(string(message)),
				tui.NewSpacer(),
			))
		})

	}
}



func main() {
	// Get login information from the login UI
	address, username, key := loginUI()

	// Open a TCP connection to the server 			//TODOmaybe put this in the loginUI?
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Init the chat history widget... and ui pointer thingy
	history := tui.NewVBox()
	var u uiThing

	// Start the goroutine to print received messages from the server
	go printServerMessage(connection, key, history, &u)

	// Init the chat UI
	go chatUI(username, connection, key, history, address, &u)

	// Allot some time for the textUI to finish initializing... TODO CHANGE THIS
	// Getting data from the server WILL CAUSE A PANIC if chatUI init does not finish before then!
	time.Sleep(1 * time.Second)

	// Send our username to the server
	p := NewPacket(Initial)
	p.PrepInitial(username)
	p.Send(connection)

	// TODO ADD SOME LEGITIMATE BLOCKING SO WE EXIT WHEN USER PRESSES ESCAPE! (or textui goroutine ends)
	time.Sleep(1 * time.Hour)
}


// Ui : part of nasty fix for updating ui when a new message is received from the server
type uiThing struct {
	ui tui.UI
}