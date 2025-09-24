package main

import (
	"chat-server/internal/config"
	"chat-server/internal/server"
	"chat-server/internal/server/network"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}

	chatServer := server.NewChatServer()

	if cfg.Server.Type == "tcp" {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
		if err != nil {
			fmt.Printf("Error listening on port 8080: %d, Err:%v\n", cfg.Server.Port, err)
			return
		}

		fmt.Printf("TCP Chat server listening on port :%d \n", cfg.Server.Port)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v\n", err)
				continue
			}
			go server.HandleConnection(network.NewTCPConnection(conn), chatServer, cfg)
		}
	} else if cfg.Server.Type == "websocket" {
		upgrader := websocket.Upgrader{}
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("Error upgrading websocket: %s\n", err)
				return
			}
			go server.HandleConnection(network.NewWSConnection(wsConn), chatServer, cfg)
		})

		log.Printf("Websocket chat server listening on port %d\n", cfg.Server.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil))
	} else {
		log.Printf("Unknown type: %s\n", cfg.Server.Type)
	}
}
