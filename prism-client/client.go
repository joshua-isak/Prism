package main

import (
	"fmt"
	"net"
	"github.com/marcusolsson/tui-go"
	"time"
	"errors"
)

// VERSION :  version number of build
var VERSION string = "v0.4"

// PORT : Port to listen for new connection on
var PORT string = "14296"

// Clients connected to the server
var clients = make(map[string]string)


// Handle the Welcome packet type
func handleWelcome(p Packet, clients map[string]string, u *uiThing, clientList *tui.List) {
	p.seek = 1
	x := int(p.ReadUint8())

	// Skip all of this if this client is the first one in the server
	if x == 0 {
		return
	}

	// Loop through every username
	for i := 0; i < x; i++ {
		// Read in the username
		l := int(p.ReadUint8())
		name := p.ReadString(l)

		// Add the username to the clients map
		clients[name] = name	// I really only want the map to be a array I can remove values from it by name :^)
	}

	// Update the connected clients list in the chatUI
	u.ui.Update(func(){
		// Iterate over the map of clients and add their names to the clientList
		for _, y := range clients {
			clientList.AddItems(y)
		}
	})

}


// Handle the ClientConnect packet type
func handleClientConnect(p Packet, clients map[string]string, u *uiThing, clientList *tui.List, username string, history *tui.Box) {
	p.seek = 1

	// Read in the username
	l := int(p.ReadUint8())
	name := p.ReadString(l)

	// Prepare a message to print in chat that someone has connected
	message := name + " has connected"

	// If this is this client's name, append a "(You)" to it
	if name == username {
		name = username + " (You)"
	}

	// Add the client name to the map of clients
	clients[name] = name

	// Update the connected clients list in the chatUI
	u.ui.Update(func(){
		// Refresh list by removing all its items then adding them back (this is a workaround to the tui library :/)
		clientList.RemoveItems()

		// Iterate over the map of clients and add their names to the clientList
		for _, y := range clients {
			clientList.AddItems(y)
		}
	})

	// Display var message in chat
	u.ui.Update( func(){
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(message))),
			tui.NewSpacer(),
		))
	})

}


// Handle the ClientDisconnect packet type
func handleClientDisconnect(p Packet, clients map[string]string, u *uiThing, clientList *tui.List, username string, history *tui.Box) {
	p.seek = 1

	// Read in the username
	l := int(p.ReadUint8())
	name := p.ReadString(l)

	// Prepare a message to print in chat that someone has disconnected
	message := name + " has disconnected"

	// Remove the client name from the map of clients
	delete(clients, name)

	// Update the connected clients list in the chatUI
	u.ui.Update(func(){
		// Refresh list by removing all its items then adding them back again from the now updated map clients
		clientList.RemoveItems()

		// Iterate over the map of clients and add them to the clientList
		for _, y := range clients {
			clientList.AddItems(y)
		}
	})

	// Display var message in chat
	u.ui.Update( func(){
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(message))),
			tui.NewSpacer(),
		))
	})

}


// Handle the GeneralMessage packet type
func handleGeneralMessage(p Packet, key []byte, history *tui.Box, u *uiThing) {

	p.seek = 1

	l := p.ReadUint8()						// read in the senderName length
	senderName := p.ReadString(int(l))		// read in the senderName as a string

	// Check if the message is encrypted
	p.seek = 23
	isEncrypted := p.ReadBool()

	// Read in the message as a byte array
	messageSize := p.ReadUint8()
	message := p.ReadBytes(int(messageSize))

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


// Handle the ServerDisconnect packet type
func handleServerDisconnect(p Packet) error {
	// be a good client and forcibly close when the server tells you to
	p.seek = 1

	errCode := int(p.ReadUint8())
	len := p.ReadUint8()
	reason := p.ReadString(int(len))

	return errors.New("Server forced this client to disconnect. Reason " + string(errCode) + ": " + reason)
}



// Connection : Listen for packets from the server and handle them
func Connection(conn net.Conn, clients map[string]string, k []byte, username string, h *tui.Box, c *tui.List, u *uiThing) error {
	defer conn.Close()

	// Send the Initial packet to the server
	init := NewPacket(Initial)
	init.PrepInitial(username)
	init.Send(conn)

	// Listen for and handle packets received from the server
	for {
		// Read in data from tcp socket and put it in a Packet object
		p, err := ReadSocket(conn)
		if err != nil {
			return nil	// this error is mostlikely a "use of closed network connection" anyway right....
			//return err
		}

		// Handle this packet based on the packet type
		pType := PacketType(p.data[0])

		switch pType {

		case Welcome:
			handleWelcome(p, clients, u, c)

		case ServerDisconnect:
			return handleServerDisconnect(p)

		case ClientConnect:
			handleClientConnect(p, clients, u, c, username, h)

		case ClientDisconnect:
			handleClientDisconnect(p, clients, u, c, username, h)

		case GeneralMessage:
			handleGeneralMessage(p, k, h, u)

		default:
			return errors.New("Received an invalid packet type")

		}
	}

}



func main() {
	// Get login information from the login UI
	address, username, key := loginUI()

	// Open a TCP connection to the server 			//TODOmaybe put this in the loginUI?
	//fmt.Println("Attemping to connect to:", address)
	connection, err := net.Dial("tcp", address + ":" + PORT)
	if err != nil {
		fmt.Println(err)
		time.Sleep(10 * time.Second)	// Let that error really sink in...
		return
	}

	// Init some widgets widget... and ui pointer thingy
	history := tui.NewVBox()
	clientList := tui.NewList()
	var u uiThing

	// Init the chat UI
	go chatUI(username, connection, key, history, clientList, address, &u)

	// Give some time for chatUI to initialize
	// chat UI not finishing init before printServerMessage runs WILL CAUSE A PANIC
	// TODO add some real blocking here with channels!
	time.Sleep(1 * time.Second)

	// Handle GeneralMessage packets from the server
	err = Connection(connection, clients, key, username, history, clientList, &u)

	// Close the chat UI
	u.ui.Quit()

	if err != nil {
		fmt.Println(err)
		time.Sleep(10 * time.Second) 	// Let that error really sink in...
	}

}


// Ui : part of nasty fix for updating ui when a new message is received from the server
type uiThing struct {
	ui tui.UI
}