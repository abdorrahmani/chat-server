package server

import (
	"bufio"
	"chat-server/internal/config"
	"fmt"
	"net"
	"sync"
	"time"
)

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients map[string]*Client
	mutex   sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect Add a new client to the chat server
func (s *ChatServer) Connect(username string, maxClients, rateLimit int) (*Client, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	}

	if len(s.clients) >= maxClients {
		return nil, ErrServerFull
	}

	refillRate := time.Second / time.Duration(rateLimit)
	client := &Client{
		Username:  username,
		Message:   make(chan string, 10),
		connected: true,
		limiter:   NewTokenBucket(rateLimit, refillRate),
	}

	s.clients[username] = client
	return client, nil
}

// Disconnect removes a client form the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.clients, client.Username)
	client.connected = false
	close(client.Message)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		if client.Username != sender.Username {
			client.Send(fmt.Sprintf("[%s]: %s", sender.Username, message))
		}
	}
}

func (s *ChatServer) PrivateMessage(sender *Client, recipient, message string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !sender.connected {
		return ErrClientDisconnected
	}

	if client, exists := s.clients[recipient]; !exists {
		return ErrRecipientNotFound
	} else if !client.connected {
		return ErrClientDisconnected
	} else {
		client.Send(fmt.Sprintf("[Private] %s : %s", sender.Username, message))
		return nil
	}
}

// HandleConnection handles a new client connection to the chat server
func HandleConnection(conn net.Conn, server *ChatServer, cfg *config.Config) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString("Enter your username: ")
	writer.Flush()

	scanner.Scan()
	username := scanner.Text()

	client, err := server.Connect(username, cfg.MaxClients, cfg.RateLimit)
	if err != nil {
		writer.WriteString(fmt.Sprintln(err))
		writer.Flush()
		return
	}

	writer.WriteString(fmt.Sprintf("%s Welcome to the Anophel Chat service. (anophel\n", username))
	writer.Flush()

	server.Broadcast(client, fmt.Sprintf("%s has joined the chat\n", username))
	defer func() {
		server.Broadcast(client, fmt.Sprintf("%s has left the chat\n", username))
		server.Disconnect(client)
	}()

	go func() {
		for msg := range client.Message {
			writer.WriteString(msg + "\n")
			writer.Flush()
		}
	}()

	HandleInputs(scanner, client, server, cfg)
}
