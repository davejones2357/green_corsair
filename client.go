package main

import (
	"fmt"
	"net"
	"log"
	"time"

	"pigeon_post/config"
	"pigeon_post/protocol"
	"pigeon_post/types"
)

// This is a simple TCP client that connects to a server, sends a message, and waits for a reply.
// Based on https://medium.com/@viktordev/socket-programming-in-go-write-a-simple-tcp-client-server-c9609edf3671

func main() {

	// Connect to the server
	
	cfg, err := config.Load()
    if err != nil {
        log.Fatalf("CONFIG ERROR: %v", err)
    }

	conn, err := net.Dial(cfg.Protocol, cfg.Host + ":" + cfg.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Build a CONNECT command with an arbitrary client ID
	clientID := "client1234"
	connectCmd := types.NewConnect(clientID)

	// Convert to protocol message and serialize
	msg := connectCmd.ToMessage()
	data, n := protocol.Serialize(msg)

	// Send the serialized bytes
	if _, err := conn.Write(data[:n]); err != nil {
		fmt.Println("write error:", err)
		return
	}

	// Wait for a server response
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	nr, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error (or no reply):", err)
		return
	}

	resp, err := protocol.Deserialize(buf[:nr])
	if err != nil {
		fmt.Println("deserialize error:", err)
		return
	}

	fmt.Printf("server reply: type=%q client=%q body=%q\n", resp.MessageType, resp.ClientID, resp.Message)

}
