package main


import (
	"fmt"
	"net"
	"os"
)


// Code for client variable organization TODO: make this all structs and stuff later
var countID int = 1 //make this a random number selector for finding client ids
var clients = make(map[int]net.Conn)

type client struct {
	conn net.Conn
	name string
	id   int
}


func broadcast(msg string) {
	for _, v := range clients {
		fmt.Fprintf(v, msg+"\n")
	}
}


func handleConnection(connection net.Conn, id int) {
	//Read in the client's username (handle "Initial" packet)
	buf := make([]byte, 256)	// read up to 256 bytes into buf
	_, err := connection.Read(buf[0:])	// read up to size of buf
	if err != nil {
		fmt.Println(err)
		return
	}
	l :=  int(buf[1]) 			// read in the username length
	name := string(buf[2:2+l])	// read l bytes from start of name turn that into a string

	// Tell all clients a user has connected
	msg := name + " has connected"
	p := NewPacket(GeneralMessage)
	p.PrepGeneralMessage("", []byte(msg), false)
	p.Broadcast(clients)
	//p.PrintDataHex()

	fmt.Println(msg)

	// Listen for messages from the client and broadcast them to all other connected clients
	for {
		// Read in data from tcp socket and put it in a Packet object
		buf := make([]byte, 512)	// read up to 512 bytes into buf
		_, err := connection.Read(buf[0:])	// read up to size of buf
		if err != nil {
			if err.Error() == "EOF"{
				break
			}
			fmt.Println(err)
			return
		}
		netData := NewPacket(Received)
		netData.data = buf

		// Read in the possibly encrypted message
		netData.seek = 24		// Move reader to byte 24 of netData.data
		messageLen := netData.ReadUint8()
		message := netData.ReadBytes(int(messageLen))

		// Broadcast this received message to all other connected clients
		p := NewPacket(GeneralMessage)
		p.PrepGeneralMessage(name, message, true)
		p.Broadcast(clients)

		fmt.Println(name + ": --ENCRYPTED--")

	}

	// Handle the TCP connection closing
	msg = name + " has disconnected" // := not used because var msg previously declared
	fmt.Println(msg)

	p2 := NewPacket(GeneralMessage)
	p2.PrepGeneralMessage("", []byte(msg), false)
	p2.Broadcast(clients)

	connection.Close()
	delete(clients, id)
}


func main() {
	// Read in command line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: prism-server [port]")
		return
	}

	// Listen for new TCP connections
	PORT := ":" + arguments[1]
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server setup successful, listening for connections...")
	defer listener.Close()

	// Handle new TCP connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(connection, countID)
		// Add client id info
		clients[countID] = connection
		countID++

	}

}
