package main

import (
	"fmt"
	"net"
)

func main() {

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Send some data to the server
	_, err = conn.Write([]byte("CONNECTION0123456789Hello, server!\n"))
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
