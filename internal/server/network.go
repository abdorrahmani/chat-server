package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/gorilla/websocket"
)

type Connection interface {
	ReadLine() (string, error)
	WriteLine(msg string) error
	Close() error
}

// ------------TCP-------------

type TCPConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTCPConnection(conn net.Conn) *TCPConnection {
	return &TCPConnection{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *TCPConnection) ReadLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// remove trailing \r and \n for Windows and Unix newlines
	return strings.TrimRight(line, "\r\n"), nil
}

func (c *TCPConnection) WriteLine(msg string) error {
	_, err := c.writer.WriteString(msg + "\n")
	if err != nil {
		return err
	}

	return c.writer.Flush()
}
func (c *TCPConnection) Close() error {
	return c.conn.Close()
}

// ----------WEBSOCKET---------

type WSConnection struct {
	conn *websocket.Conn
}

func NewWSConnection(conn *websocket.Conn) *WSConnection {
	return &WSConnection{
		conn: conn,
	}
}

func (c *WSConnection) ReadLine() (string, error) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}

	return string(msg), nil
}

func (c *WSConnection) WriteLine(msg string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *WSConnection) Close() error {
	return c.conn.Close()
}
