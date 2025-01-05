.PHONY: all clean

GO := go
BINDIR := .
CERT_DIR := certs
SERVER_KEY := $(CERT_DIR)/server.key
SERVER_CRT := $(CERT_DIR)/server.crt


all: clipboard scopy spaste shistory certs

clipboard:
	$(GO) build -o $(BINDIR)/clipboard ./server

scopy:
	$(GO) build -o $(BINDIR)/scopy ./clients/scopy

spaste:
	$(GO) build -o $(BINDIR)/spaste ./clients/spaste

shistory:
	$(GO) build -o $(BINDIR)/shistory ./clients/shistory


certs: $(SERVER_KEY) $(SERVER_CRT)

$(SERVER_KEY):
	@mkdir -p $(CERT_DIR)
	openssl genrsa -out $(SERVER_KEY) 2048

$(SERVER_CRT): $(SERVER_KEY)
	@mkdir -p $(CERT_DIR)
	openssl req -x509 -new -nodes -key $(SERVER_KEY) -sha256 -days 365 \
		-out $(SERVER_CRT) -subj "/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,DNS:localhost/"/

# Generate self-signed certificates
certs_orig:
	@echo "Generating self-signed certificates..."
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout server.key -out server.crt \
		-subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=localhost"
	@echo "Certificates generated: server.key and server.crt"

# Clean up build files
clean:
	rm -rf certs
	rm -f $(BINDIR)/clipboard $(BINDIR)/scopy $(BINDIR)/spaste $(BINDIR)/shistory
	rm -f server.key server.crt
