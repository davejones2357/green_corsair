package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/davejones2357/pigeon_post/config"
	"github.com/davejones2357/pigeon_post/protocol"
	"github.com/davejones2357/pigeon_post/types"
)

// This is a simple TCP client that connects to a server, sends a message, and waits for a reply.
// Based on https://medium.com/@viktordev/socket-programming-in-go-write-a-simple-tcp-client-server-c9609edf3671

func main() {
	start := flag.Bool("start", false, "start the client")
	clientID := flag.String("client", "", "client ID to use when connecting")

	flag.Parse()

	if !*start {
		flag.Usage()
		return
	}

	if *clientID == "" {
		log.Fatal("--client is required when using --start")
	}

	startClient(*clientID)
}

func startClient(clientID string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("CONFIG ERROR: %v", err)
	}

	conn, err := net.Dial(cfg.Protocol, cfg.Host+":"+cfg.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	connectCmd := types.NewConnect(clientID)

	// Convert to protocol message and serialize.
	msg := connectCmd.ToMessage()
	data, n := protocol.Serialize(msg)

	// Send the serialized bytes.
	if _, err := conn.Write(data[:n]); err != nil {
		fmt.Println("write error:", err)
		return
	}

	// Wait for a server response.
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		fmt.Println("read deadline error:", err)
		return
	}

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

	fmt.Printf("server reply: type=%q client=%q body=%q\n",
		resp.MessageType,
		resp.ClientID,
		resp.Message,
	)
}
