package main

import (
	"chat-server/internal/config"
	"chat-server/internal/server"
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

	if cfg.Type == "tcp" {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PORT))
		if err != nil {
			fmt.Printf("Error listening on port 8080: %s, Err:%v\n", cfg.PORT, err)
			return
		}

		fmt.Printf("TCP Chat server listening on port :%s \n", cfg.PORT)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v\n", err)
				continue
			}
			go server.HandleConnection(server.NewTCPConnection(conn), chatServer, cfg)
		}
	} else if cfg.Type == "websocket" {
		upgrader := websocket.Upgrader{}
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("Error upgrading websocket: %s\n", err)
				return
			}
			go server.HandleConnection(server.NewWSConnection(wsConn), chatServer, cfg)
		})

		log.Printf("Websocket chat server listening on port %s\n", cfg.PORT)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.PORT), nil))
	} else {
		log.Printf("Unknown type: %s\n", cfg.Type)
	}
}
