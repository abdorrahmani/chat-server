package network

import (
	"chat-server/internal/config"
	"crypto/tls"
	"fmt"
	"net"
)

func parseMinVersion(ver string) uint16 {
	switch ver {
	case "TLS13":
		return tls.VersionTLS13
	case "TLS12":
		return tls.VersionTLS12
	default:
		return tls.VersionTLS12
	}
}

func tlsConfig(cfg config.TLSConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   parseMinVersion(cfg.MinVersion),
	}
	return tlsCfg, nil
}

func NewTLS(ln net.Listener, tlsCfg config.TLSConfig) (net.Listener, error) {
	cfg, err := tlsConfig(tlsCfg)
	if err != nil {
		return nil, err
	}
	tlsListener := tls.NewListener(ln, cfg)
	return tlsListener, nil
}
