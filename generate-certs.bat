@echo off

REM Create tls directory if it doesn't exist
if not exist tls mkdir tls

REM Generate private key
openssl genrsa -out tls/server.key 2048

REM Generate certificate signing request
openssl req -new -key tls/server.key -out tls/server.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

REM Generate self-signed certificate
openssl x509 -req -days 365 -in tls/server.csr -signkey tls/server.key -out tls/server.crt

REM Clean up CSR file
del tls/server.csr

echo TLS certificates generated successfully!
echo Certificate: tls/server.crt
echo Private Key: tls/server.key
