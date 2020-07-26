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


func handleConnection(conn net.Conn, id int) {
	//Read in the client's username (handle "Initial" packet)
	i, err := ReadSocket(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection if the client doesn't respond with the "Initial" packet type
	if i.data[0] != 1 {
		d := NewPacket(ServerDisconnect)
		d.PrepServerDisconnect(0, "client sent bad data")
		d.Send(conn)
		conn.Close()
		fmt.Println(conn.RemoteAddr().String() + " client sent bad data (Initial)")
	}

	l :=  int(i.data[1]) 			// read in the username length
	name := string(i.data[2:2+l])	// read l bytes from start of name turn that into a string

	// Check if the client's username is already in use
	if _, ok := clients[name]; ok {
		// Tell the client this username is already in use and disconnect it
		d := NewPacket(ServerDisconnect)
		d.PrepServerDisconnect(5, "username already in use")
		d.Send(conn)
		conn.Close()
		fmt.Println(conn.RemoteAddr().String() + " tried to connect with a username already in use")
		return
	}

	// Send the Welcome packet to the client
	w := NewPacket(Welcome)
	w.PrepWelcome(clients)
	w.Send(conn)

	// Add the client to the clients map
	clients[name] = conn

	// Tell all clients a user has connected
	c := NewPacket(ClientConnect)
	c.PrepClientConnect(name)
	c.Broadcast(clients)

	fmt.Println(name + " has connected from " + conn.RemoteAddr().String() )

	// Listen for messages from the client and broadcast them to all other connected clients
	for {
		// Read in data from tcp socket and put it in a Packet object
		m, err := ReadSocket(conn)

		if err != nil {
			if err.Error() == "EOF"{
				break
			}
			fmt.Println("other err", err)
			break
		}

		// Close the connection if the client doesn't respond with the "GeneralMessage" packet type
		if PacketType(m.data[0]) != GeneralMessage {
			d := NewPacket(ServerDisconnect)
			d.PrepServerDisconnect(0, "client sent bad data")
			d.Send(conn)
			conn.Close()
			fmt.Println("client sent bad data (GeneralMessage)")
			break
		}

		netData := NewPacket(Received)
		netData.data = m.data

		// Read in the hopefully encrypted message
		netData.seek = 24		// Move reader to byte 24 of netData.data
		messageLen := netData.ReadUint8()
		message := netData.ReadBytes(int(messageLen))

		// Broadcast this received message to all other connected clients
		p := NewPacket(GeneralMessage)
		p.PrepGeneralMessage(name, message, true)
		p.Broadcast(clients)

		//fmt.Println(name + ": --ENCRYPTED--")

	}

	// Close the tcp connection and remove the client from the clients map
	conn.Close()
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
