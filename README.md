 ## Anophel Chat Server (Go)

 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?logo=opensourceinitiative&logoColor=white)](LICENSE)
 [![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://github.com/abdorrahmani/chat-server)
 [![GitHub stars](https://img.shields.io/github/stars/abdorrahmani/chat-server?logo=github)](https://github.com/abdorrahmani/chat-server/stargazers)
 [![GitHub issues](https://img.shields.io/github/issues/abdorrahmani/chat-server?logo=github)](https://github.com/abdorrahmani/chat-server/issues)
 [![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](https://github.com/abdorrahmani/chat-server)
 [![Release](https://img.shields.io/github/v/release/abdorrahmani/chat-server?include_prereleases&sort=semver)](https://github.com/abdorrahmani/chat-server/releases)

Lightweight chat server supporting TCP and WebSocket transports, with optional TLS for both (TLS over TCP and WSS). Includes configurable rate limiting, message size limits, optional password gate, and Docker/Compose deployment.

### Features
- TCP and WebSocket transports
- Optional TLS for both transports (TLS1.2+)
- Simple chat with broadcast and private messages
- Rate limiting per client and message length limit
- Optional password prompt before joining
- Dockerfile and docker-compose included

---

## Table of Contents
- Setup
- Configuration
- Generate TLS certificates
- Run (TCP / WebSocket)
- Use the server (commands & examples)
- TLS usage (TCP + WSS)
- Docker and Docker Compose
- Logs and files
- Troubleshooting

---

## Setup

Prerequisites:
- Go 1.24+
- OpenSSL (for local TLS cert generation)
- Optional: Docker 24+ and Docker Compose v2+

Clone and build:

```bash
git clone https://github.com/abdorrahmani/chat-server.git
cd chat-server
go build -o bin/chat ./cmd/main.go
```

Run with default config (non-TLS TCP on port 8080):

```bash
./bin/chat
```

The server reads `config.yml` from the current working directory.

---

## Configuration

All settings live in `config.yml` and are loaded via Viper at startup.

```yaml
server:
  host: "0.0.0.0"   # (currently informational; the server binds to :port)
  port: 8080         # listening port for both TCP and WebSocket
  type: "tcp"        # "tcp" or "websocket"
  maxClients: 100
  readTimeout: 5     # seconds (not all timeouts may be enforced in handlers)
  writeTimeout: 5    # seconds

security:
  requirePassword: false
  password: "1234"
  hashMessage: true          # reserved; message hashing helper exists in utils
  hashAlgorithm: "hmac-sha256"
  hashKey: "supersecret"

rateLimit:
  messagePerSecond: 5        # token bucket rate per client
  burst: 5                   # not all fields may be used directly

message:
  maxLength: 1000            # maximum allowed message length

log:
  enableLogging: false
  file: "chat.log"

tls:
  tlsRequire: false          # enable TLS for TCP or WebSocket
  certFile: "tls/server.crt"
  keyFile: "tls/server.key"
  minVersion: "TLS12"        # TLS12 or TLS13
```

Key notes:
- Set `server.type` to `tcp` for raw TCP, or `websocket` for WS/WSS.
- Set `tls.tlsRequire: true` to enable TLS (affects both TCP and WebSocket depending on `server.type`).
- Files referenced in `tls` must exist and be readable by the process.

---

## Generate TLS certificates

Self-signed certificates for local development are supported. Scripts are provided:

- Windows PowerShell / cmd:
  ```bash
  .\generate-certs.bat
  ```

- Linux/macOS:
  ```bash
  bash generate-certs.sh
  ```

These create:
- `tls/server.key` (private key)
- `tls/server.crt` (self-signed cert)

Update `config.yml`:
```yaml
tls:
  tlsRequire: true
  certFile: "tls/server.crt"
  keyFile: "tls/server.key"
  minVersion: "TLS12"
```

---

## Run

### TCP mode (plaintext)
```bash
sed -i.bak 's/type: ".*"/type: "tcp"/' config.yml   # or edit manually
sed -i.bak 's/tlsRequire: .*/tlsRequire: false/' config.yml
./bin/chat
```

### TCP mode (TLS)
```bash
sed -i.bak 's/type: ".*"/type: "tcp"/' config.yml
sed -i.bak 's/tlsRequire: .*/tlsRequire: true/' config.yml
./bin/chat
```

### WebSocket mode (WS)
```bash
sed -i.bak 's/type: ".*"/type: "websocket"/' config.yml
sed -i.bak 's/tlsRequire: .*/tlsRequire: false/' config.yml
./bin/chat
```
Listens on `ws://localhost:<port>/ws`.

### WebSocket mode (WSS)
```bash
sed -i.bak 's/type: ".*"/type: "websocket"/' config.yml
sed -i.bak 's/tlsRequire: .*/tlsRequire: true/' config.yml
./bin/chat
```
Listens on `wss://localhost:<port>/ws`.

---

## Use the server

When a client connects, the server prompts:
1) `Enter your username:`
2) If `security.requirePassword` is true, prompts: `Enter password:` (must match `security.password`).

Commands and behavior:
- `/quit`: leave the chat
- `/pm <username> <message>`: send a private message
- Any other text: broadcast to all other users
- Echo: the server sends `ME: <your message>` back to the sender
- Rate limit: if you send too quickly, you’ll receive a slowdown message
- Max length: messages exceeding `message.maxLength` are rejected

### Example clients

#### Raw TCP
- Plaintext TCP:
  ```bash
  nc 127.0.0.1 8080
  # or on Windows (Git Bash/WSL): ncat 127.0.0.1 8080
  ```

- TLS over TCP (using OpenSSL):
  ```bash
  openssl s_client -connect 127.0.0.1:8080 -quiet
  # for self-signed certs add: -verify 0
  ```

#### WebSocket
- WS using `wscat`:
  ```bash
  npx wscat@latest -c ws://localhost:8080/ws
  ```

- WSS using `wscat` (self-signed):
  ```bash
  npx wscat@latest -c wss://localhost:8080/ws --no-check
  ```

---

## TLS usage details

The server enables TLS when `tls.tlsRequire: true`.

- TCP + TLS: wraps the TCP listener with `crypto/tls` using `tls/server.crt` and `tls/server.key`.
- WebSocket + TLS: uses `http.ListenAndServeTLS` to serve WSS at `/ws`.
- Minimum TLS version is controlled by `tls.minVersion` (`TLS12` or `TLS13`).
- For self-signed certs, clients must disable verification or trust the cert.

Security considerations:
- Self-signed certs are for development only. Use proper CA-issued certificates in production.
- Keep private keys secure and set file permissions appropriately.

### Test TLS over TCP inside Docker container (Windows dev mode example)
- Ensure `server.type: "tcp"` and `tls.tlsRequire: true` in `config.yml`.
- Start the container (Docker or Compose).
- Run inside the container:
  ```bash
  docker exec -it chat-server-chat-server-1 openssl s_client -connect localhost:8080 -quiet -verify 0
  ```
  Notes:
  - `localhost:8080` targets the service from inside the container.
  - `-verify 0` skips verification for self-signed certs (dev only).
  - If you run from the host instead, use `127.0.0.1:8080` (or the mapped port) and ensure the port is published.

---

## Docker

### Build and run (Dockerfile)
```bash
docker build -t chat-server:local .
docker run --rm -p 8080:8080 \
  -v %cd%/config.yml:/root/config.yml \
  -v %cd%/tls:/root/tls \
  chat-server:local
```

Notes:
- The image copies `config.yml` and `tls/` at build time, but the command above also mounts your local `config.yml` and `tls/` to override at runtime. Adjust paths for your shell/OS.

### Docker Compose
`docker-compose.yml` already exposes port 8080 and mounts `config.yml` and `tls/`:

```bash
docker compose up --build
```

Edit `config.yml` to switch between `tcp`/`websocket` and to enable/disable TLS. Recreate the container after changes.

---

## Logs and files
- TLS assets: `tls/server.crt`, `tls/server.key`
- Configuration: `config.yml`
- Binary (local build): `bin/chat`
- Optional logs (if enabled): `logs/` or file specified in `log.file`

---

## Troubleshooting
- Port already in use: change `server.port` or stop the conflicting service.
- Self-signed cert errors: for testing clients, disable verification (`--no-check` in wscat, `-verify 0` in `openssl s_client`) or trust the cert.
- Cannot connect over WSS: confirm `tls.tlsRequire: true`, cert/key paths exist in the container/host process, and the port is reachable.
- Messages dropped or slow: you may be hitting the per-client rate limit; adjust `rateLimit.messagePerSecond`.
- Username in use: each username must be unique per server instance.

---

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

- GitHub: [@abdorrahmani](https://github.com/abdorrahmani)