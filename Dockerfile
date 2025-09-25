# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates openssl

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/config.yml .

COPY --from=builder /app/tls ./tls

EXPOSE 8080

CMD ["./main"]