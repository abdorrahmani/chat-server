package server

import "sync"

// Client represent a connected chat client
type Client struct {
	Username  string
	Message   chan string
	connected bool
	mutex     sync.RWMutex
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.connected {
		c.Message <- message
	}
	return
}

// Receive returns the next message for the client
func (c *Client) Receive() string {
	if message, ok := <-c.Message; ok {
		return message
	}
	return ""
}
