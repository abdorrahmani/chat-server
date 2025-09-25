package grpcserver

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"chat-server/internal/config"
	core "chat-server/internal/server"
	chatpb "chat-server/internal/server/network/grpc"

	"google.golang.org/grpc"
)

// ChatGRPCServer implements the generated gRPC service and bridges to the core ChatServer
type ChatGRPCServer struct {
	core *core.ChatServer
	cfg  *config.Config
	chatpb.UnimplementedChatServiceServer
}

func New(coreServer *core.ChatServer, cfg *config.Config) *ChatGRPCServer {
	return &ChatGRPCServer{core: coreServer, cfg: cfg}
}

func (s *ChatGRPCServer) SendMessage(ctx context.Context, req *chatpb.ChatMessage) (*chatpb.ChatResponse, error) {
	if req == nil {
		return &chatpb.ChatResponse{Status: "empty request"}, nil
	}
	user := req.GetUser()
	text := req.GetText()

	if len(text) > s.cfg.Message.MaxLength {
		return &chatpb.ChatResponse{Status: fmt.Sprintf("message too long (max %d chars)", s.cfg.Message.MaxLength)}, nil
	}

	// Create a lightweight sender representation to reuse Broadcast formatting
	sender := &core.Client{Username: user, Message: make(chan string, 1)}
	s.core.Broadcast(sender, text)

	return &chatpb.ChatResponse{Status: "ok"}, nil
}

// Chat implements bidirectional chat similar to TCP/WebSocket modes
func (s *ChatGRPCServer) Chat(stream grpc.BidiStreamingServer[chatpb.ClientEvent, chatpb.ServerEvent]) error {
	// Prompt for username
	if err := stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Prompt{Prompt: &chatpb.Prompt{Text: "Enter your username: "}}}); err != nil {
		return err
	}

	first, err := stream.Recv()
	if err != nil {
		return err
	}
	join := first.GetJoin()
	if join == nil || strings.TrimSpace(join.GetUsername()) == "" {
		_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: "ERROR: username required"}}})
		return nil
	}
	username := join.GetUsername()

	// Password if required
	if s.cfg.Security.RequirePassword {
		if strings.TrimSpace(join.GetPassword()) == "" {
			if err := stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Prompt{Prompt: &chatpb.Prompt{Text: "Enter password: "}}}); err != nil {
				return err
			}
			passEvt, err := stream.Recv()
			if err != nil {
				return err
			}
			if p := passEvt.GetJoin(); p != nil && p.GetPassword() != "" {
				join = p
			} else if t := passEvt.GetText(); t != nil {
				join = &chatpb.Join{Username: username, Password: t.GetMessage()}
			}
		}
		if join.GetPassword() != s.cfg.Security.Password {
			_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: "Passwords do not match"}}})
			return nil
		}
	}

	// Connect client
	client, err := s.core.Connect(username, s.cfg.Server.MaxClients, s.cfg.RateLimit.MessagePerSecond)
	if err != nil {
		_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: err.Error()}}})
		return nil
	}
	defer func() {
		s.core.Broadcast(client, fmt.Sprintf("%s has left the chat", username))
		s.core.Disconnect(client)
	}()

	// Welcome and join notice
	_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: username + ", Welcome to the Anophel Chat service"}}})
	s.core.Broadcast(client, "has joined the chat")

	// Forward outbound messages to the stream
	go func() {
		for msg := range client.Message {
			if strings.HasPrefix(msg, "ME: ") {
				_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Echo{Echo: &chatpb.Echo{Text: strings.TrimPrefix(msg, "ME: ")}}})
				continue
			}
			// Try to extract [from]: text
			from := ""
			text := msg
			re := regexp.MustCompile(`^\[(.+?)\]:\s*(.*)$`)
			if m := re.FindStringSubmatch(msg); len(m) == 3 {
				from = m[1]
				text = m[2]
			}
			_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Chat{Chat: &chatpb.Chat{From: from, Text: text}}})
		}
	}()

	// Read incoming client events
	for {
		evt, err := stream.Recv()
		if err != nil {
			return nil
		}
		if t := evt.GetText(); t != nil {
			message := t.GetMessage()

			if len(message) > s.cfg.Message.MaxLength {
				_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: fmt.Sprintf("❌ message too long (max: %d chars)", s.cfg.Message.MaxLength)}}})
				continue
			}
			// Rate limiting: reuse client.limiter if available via Send path
			// Fallback: rely on core handler policy

			// Commands
			if strings.TrimSpace(message) == "/quit" {
				_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: "You have left the chat."}}})
				return nil
			}

			if strings.HasPrefix(message, "/pm") {
				parts := strings.SplitN(message, " ", 3)
				if len(parts) < 3 {
					_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: "ERROR: Invalid private message format. Use /pm <username> <message>"}}})
					continue
				}
				recipient, pm := parts[1], parts[2]
				if err := s.core.PrivateMessage(client, recipient, pm); err != nil {
					_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Notice{Notice: &chatpb.Notice{Text: "ERROR: Invalid private message " + err.Error()}}})
				} else {
					_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Echo{Echo: &chatpb.Echo{Text: pm}}})
				}
			} else {
				s.core.Broadcast(client, message)
				_ = stream.Send(&chatpb.ServerEvent{Payload: &chatpb.ServerEvent_Echo{Echo: &chatpb.Echo{Text: message}}})
			}
		}
	}
}
