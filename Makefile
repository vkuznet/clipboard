.PHONY: all

GO := go
BINDIR := .

# certs parts
CERT_DIR := $(HOME)/.clipboard/certs
SERVER_KEY := $(CERT_DIR)/server.key
SERVER_CRT := $(CERT_DIR)/server.crt

# config parts
CONFIG_DIR=$(HOME)/.clipboard
CONFIG_FILE=$(CONFIG_DIR)/config.json
HISTORY_FILE=$(CONFIG_DIR)/shistory.json
CERTS_DIR=$(CONFIG_DIR)/certs


all: clean clipboard scopy spaste shistory

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
	@echo "Generating self-signed certificates in $(CERT_DIR) with duration 36500 days ..."
	openssl genrsa -out $(SERVER_KEY) 2048

$(SERVER_CRT): $(SERVER_KEY)
	@mkdir -p $(CERT_DIR)
	openssl req -x509 -new -nodes -key $(SERVER_KEY) -sha256 -days 36500 \
		-out $(SERVER_CRT) -subj "/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,DNS:localhost/"/

# generate $HOME/.clipboard/config.json file
config:
	@mkdir -p $(CONFIG_DIR)
	@mkdir -p $(CERTS_DIR)
	@echo '{' > $(CONFIG_FILE)
	@echo '  "port": 14443,' >> $(CONFIG_FILE)
	@echo '  "clipboard_size": 10,' >> $(CONFIG_FILE)
	@echo '  "history_file": "$(HISTORY_FILE)",' >> $(CONFIG_FILE)
	@echo '  "server_key": "$(CERTS_DIR)/server.key",' >> $(CONFIG_FILE)
	@echo '  "server_cert": "$(CERTS_DIR)/server.crt"' >> $(CONFIG_FILE)
	@echo '}' >> $(CONFIG_FILE)
	@echo "Configuration created at $(CONFIG_FILE)"

# Clean up build files
clean:
	rm -f $(BINDIR)/clipboard $(BINDIR)/scopy $(BINDIR)/spaste $(BINDIR)/shistory
	rm -f server.key server.crt
