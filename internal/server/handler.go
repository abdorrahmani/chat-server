package server

import (
	"bufio"
	"chat-server/internal/config"
	"fmt"
	"strconv"
	"strings"
)

// HandleInputs handles incoming messages from a client
func HandleInputs(scanner *bufio.Scanner, client *Client, server *ChatServer, cfg *config.Config) {
	for scanner.Scan() {
		message := scanner.Text()

		if len(message) > cfg.MaxMessageLength {
			fmt.Printf("❌message too long (max %d chars)\n", cfg.MaxMessageLength)
			client.Send("❌ message too long (max: " + strconv.Itoa(cfg.MaxMessageLength) + " chars)")
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
