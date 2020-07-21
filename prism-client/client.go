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

func printServerMessage(connection net.Conn) {
	for {
		netData, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Print("\n") // needed cause of "\033[F" ANSI code to move cursor up
		fmt.Print("\033[F" + netData)
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

	// Open a TCP connection to the server
	ADDRESS := arguments[2] + ":" + arguments[3]

	connection, err := net.Dial("tcp", ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send our username to the server
	fmt.Fprintf(connection, username+"\n")

	// Goroutine to print messages from the server
	go printServerMessage(connection)

	for {
		// Send user inputted messages to the server
		reader := bufio.NewReader(os.Stdin)
		//fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")
	}

	// Display messages received from the server

}
