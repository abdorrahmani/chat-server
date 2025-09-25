#!/bin/bash

mkdir -p tls

openssl genrsa -out tls/server.key 2048

openssl req -new -key tls/server.key -out tls/server.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

openssl x509 -req -days 365 -in tls/server.csr -signkey tls/server.key -out tls/server.crt

rm tls/server.csr

echo "TLS certificates generated successfully!"
echo "Certificate: tls/server.crt"
echo "Private Key: tls/server.key"
