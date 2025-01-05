# Secure Clipboard Manager

The Secure Clipboard Manager is a Go-based application that provides a secure and persistent clipboard for storing and transferring sensitive data. It encrypts clipboard content to protect against unauthorized access and supports an HTTP(S) API for inter-application communication. The clipboard includes history management and configurable settings for enhanced usability.

## Features

- **Secure Copy/Paste**: Encrypts clipboard data using AES-256 encryption.
- **Persistent Clipboard**: Stores up to a configurable number of clipboard items for future access.
- **Clipboard History**: Maintains a configurable history of past clipboard entries.
- **Inter-Application Communication**: Provides an HTTP(S) server for securely copying and pasting data between applications.
- **Configurable Settings**: Allows users to specify the clipboard size, history size, server port, and TLS certificates for secure communication.
- **Client Utilities**: Includes `scopy`, `spaste`, and `shistory` clients for interacting with the clipboard server.

## Configuration

The application uses a configuration file to define its behavior. Example `config.json`:

```json
{
  "port": 14443,
  "clipboard_size": 10,
  "clipboard_secret": "secret",
  "history_size": 10,
  "history_file": "clipboard_history.json",
  "server_key": "certs/server.key",
  "server_cert": "certs/server.crt"
}
```

- `port`: The port number for the server.
- `clipboard_size`: The maximum number of items in the clipboard.
- `history_size`: The maximum number of items in the clipboard history.
- `history_file`: The file for persisting clipboard history.
- `server_key`: Path to the server's private key for TLS.
- `server_cert`: Path to the server's certificate for TLS.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/vkuznet/clipboard.git
   cd clipboard
   ```

2. Build the project:
   ```bash
   make
   ```

3. Generate self-signed certificates (if needed):
   ```bash
   make certs
   ```

## Usage

1. **Start the Server**:
   ```bash
   ./clipboard -config config.json
   ```

2. **Copy Data**:
   Use the `scopy` client to copy data to the clipboard:
   ```bash
   ./scopy "Some secure content"
   ```

3. **Paste Data**:
   Use the `spaste` client to retrieve data from the clipboard:
   ```bash
   ./spaste
   ```

4. **View Clipboard History**:
   Use the `shistory` client to view the clipboard history:
   ```bash
   ./shistory
   ```

5. **Secure Communication**:
   The server runs in HTTPS mode if valid `serverKey` and `serverCert` are provided in the configuration. Use `curl` or custom clients to interact with the API securely.

## API Endpoints

- `POST /copy`: Copy data to the clipboard.
- `GET /paste`: Retrieve the latest clipboard content.
- `GET /history`: View clipboard history.

## Example Usage

Start the server:
```bash
./clipboard -config config.json
```

Copy data:
```bash
./scopy "Hello, Secure Clipboard!"
```

Paste data:
```bash
./spaste
# Output: Hello, Secure Clipboard!
```

View history:
```bash
./shistory
# Output:
# 1. Hello, Secure Clipboard!
```

## Security

- Clipboard content and history are encrypted using AES-256 encryption.
- HTTPS ensures secure communication between clients and the server.
- Certificates can be generated using `make certs` for local testing.

## Contribution

Feel free to submit issues or pull requests for improvements and additional features.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
