package main

import (
	"chat-server/internal/config"
	"chat-server/internal/server"
	"fmt"
	"log"
	"net"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}

	s := server.NewChatServer()

	listener, err := net.Listen("tcp", ":"+cfg.PORT)
	if err != nil {
		fmt.Printf("Error listening on port 8080: %s\n", err)
		return
	}

	defer listener.Close()

	fmt.Println("Listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}

		go server.HandleConnection(conn, s, cfg)
	}
}
