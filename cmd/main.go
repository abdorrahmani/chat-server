package main

import (
	"chat-server/internal/server"
	"fmt"
	"log"
	"net"
)

func main() {
	s := server.NewChatServer()

	listener, err := net.Listen("tcp", ":8080")
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

		go server.HandleConnection(conn, s)
	}
}
