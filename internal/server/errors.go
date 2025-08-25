package server

import "errors"

// Common Errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrClientDisconnected   = errors.New("client disconnected")
	ErrServerFull           = errors.New("server full")
	ErrInvalidCommand       = errors.New("invalid command")
	ErrRecipientNotFound    = errors.New("recipient not found")
)
