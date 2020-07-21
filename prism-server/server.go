package main


import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
)


// Code for client variable organization TODO: make this all structs and stuff later
var count_id int = 1	//make this a random number selector for finding client ids
var clients = make(map[int]net.Conn)


type client struct {
	conn net.Conn
	name string
	id int
}


func broadcast(msg string){
	for _, v := range clients {
		fmt.Fprintf(v, msg + "\n")
	}
}


func handleConnection(connection net.Conn, id int) {
	//Read in the client's username
	netData, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
			fmt.Println(err)
			return
	}
	var name string = strings.TrimSpace(string(netData))
	msg := name + " has connected"
	fmt.Println(msg)
	broadcast(msg)

	// Listen for messages from the client and broadcast them to all other connected clients
	for {
		netData, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {	
			if err.Error() == "EOF" {
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		//fmt.Fprintf(connection, "LOL")
		msg = name + ": " + strings.TrimSuffix(string(netData), "\n") 
		fmt.Println(msg)
		broadcast(msg)
	}

	msg = name + " has disconnected"	// := not used because var msg previously declared
	fmt.Println(msg)
	broadcast(msg)
	connection.Close()
	delete(clients, id)
}



func main() {
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

		go handleConnection(connection, count_id)
		// Add client id info
		clients[count_id] = connection
		count_id++
		


	}


 
}

