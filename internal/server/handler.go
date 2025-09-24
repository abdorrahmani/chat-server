package server

import (
	"chat-server/internal/config"
	"chat-server/internal/server/network"
	"fmt"
	"strconv"
	"strings"
)

// HandleInputs handles incoming messages from a client
func HandleInputs(conn network.Connection, client *Client, server *ChatServer, cfg *config.Config) {
	for {
		message, err := conn.ReadLine()
		if err != nil {
			break
		}

		if len(message) > cfg.Message.MaxLength {
			fmt.Printf("❌message too long (max %d chars)\n", cfg.Message.MaxLength)
			client.Send("❌ message too long (max: " + strconv.Itoa(cfg.Message.MaxLength) + " chars)")
			continue
		}

		if !client.limiter.Allow() {
			client.Send("You are sending message too fast! slow down.")
			continue
		}

		if strings.TrimSpace(message) == "/quit" {
			client.Send("You have left the chat.")
			return
		}

		if strings.HasPrefix(message, "/pm") {
			parts := strings.SplitN(message, " ", 3)
			if len(parts) < 3 {
				client.Send("ERROR: Invalid private message format. Use /pm <username> <message>\n")
				continue
			}
			recipient, msg := parts[1], parts[2]
			err := server.PrivateMessage(client, recipient, msg)
			if err != nil {
				client.Send("ERROR: Invalid private message " + err.Error())
			} else {
				client.Send("ME: " + msg)
			}
		} else {
			server.Broadcast(client, message)
			client.Send("ME: " + message)
		}
	}
}
