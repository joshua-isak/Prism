package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


func connect() {
	return
}


func printServerMessage(connection net.Conn, key []byte) {
	defer connection.Close()
	for {
		var output string

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

		// Append username of the message sender to output
		l := netData.ReadUint8()						// read in the senderName length
		senderName := netData.ReadString(int(l))		// read in the senderName as a string
		if l > 0 {
			output += senderName + ": "					// don't add senderName if it is empty
		}

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

		// Cast the messge to a string and add it to the output
		output += string(message)

		// Print out the message!
		// Some "fancy" terminal formatting, should really make this with curses...
		fmt.Print("\n") // needed cause of "\033[F" ANSI code to move cursor up
		fmt.Println( output)//"\033[F" + output) //message)
	}
}



func main() {
	// Get command line arguments
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: prism-client [username] [server] [port]")
		fmt.Println("Example: prism-client Anon 192.168.0.1 201")
		return
	}
	var username string = arguments[1]
	var k string = arguments[4]
	key := []byte(k)

	// Open a TCP connection to the server
	ADDRESS := arguments[2] + ":" + arguments[3]

	connection, err := net.Dial("tcp", ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send our username to the server
	p := NewPacket(Initial)
	p.PrepInitial(username)
	p.Send(connection)

	// Goroutine to print received messages from the server
	go printServerMessage(connection, key)

	for {
		// Read in user text
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// Encrypt user text
		msg := encrypt([]byte(text), key)

		// Send message to server
		p := NewPacket(GeneralMessage)
		p.PrepGeneralMessage(username, msg, true)
		p.Send(connection)
	}

}
