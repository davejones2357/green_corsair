package main

import (
	"fmt"
	"net"
	"github.com/davejones2357/pigeon_post/protocol"
)

// This is a simple TCP client that connects to a server, sends a message, and waits for a reply.
// Based on https://medium.com/@viktordev/socket-programming-in-go-write-a-simple-tcp-client-server-c9609edf3671

func main() {

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Send some data to the server
	msg:= protocol.SimpleMessage{"CONNECTION","0123456789","Hello new server\n"}
	bytes,_ := protocol.Serialize(msg)
	_, err = conn.Write([]byte(bytes))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for and read the reply
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	fmt.Println("Server reply:", string(buffer[:n]))
}
