package main


import (
	"fmt"
	"net"
)

// PORT : Port to listen for new connections on
var PORT string = "14296"


// Code for client variable organization TODO: make this all structs and stuff later
var countID int = 1 //make this a random number selector for finding client ids
var clients = make(map[string]net.Conn)

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
	i, err := ReadSocket(connection)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection if the client doesn't respond with the "Initial" packet type
	if i.data[0] != 1 {
		connection.Close()
		fmt.Println("client sent bad data (Initial)")
	}

	l :=  int(i.data[1]) 			// read in the username length
	name := string(i.data[2:2+l])	// read l bytes from start of name turn that into a string

	// Check if the client's username is already in use

	// Send the Welcome packet to the client
	w := NewPacket(Welcome)
	w.PrepWelcome(clients)
	w.Send(connection)

	// Add the client to the clients map
	clients[name] = connection

	// Tell all clients a user has connected
	c := NewPacket(ClientConnect)
	c.PrepClientConnect(name)
	c.Broadcast(clients)

	fmt.Println(name + " has connected from " + connection.RemoteAddr().String() )

	// Listen for messages from the client and broadcast them to all other connected clients
	for {
		// Read in data from tcp socket and put it in a Packet object
		m, err := ReadSocket(connection)

		if err != nil {
			if err.Error() == "EOF"{
				break
			}
			fmt.Println("other err", err)
			break
		}

		// Close the connection if the client doesn't respond with the "GeneralMessage" packet type
		if PacketType(m.data[0]) != GeneralMessage {
			connection.Close()
			fmt.Println("client sent bad data (GeneralMessage)")
			break
		}

		netData := NewPacket(Received)
		netData.data = m.data

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

	// Close the tcp connection and remove the client from the clients map
	connection.Close()
	delete(clients, name)

	// Tell all clients a user has disconnected
	d := NewPacket(ClientDisconnect)
	d.PrepClientConnect(name)
	d.Broadcast(clients)

	fmt.Println(name + " has disconnected")


}


func main() {
	// Listen for new TCP connections
	listener, err := net.Listen("tcp", ":" + PORT)
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

	}

}
