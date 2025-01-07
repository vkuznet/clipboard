package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
)

func decodeHex(s string) ([]byte, error) {
	dst := make([]byte, len(s)/2)
	_, err := fmt.Sscanf(s, "%x", &dst)
	return dst, err
}

// Encrypt function encrypts given data with our secure key
func Encrypt(skey []byte, data string) (string, error) {
	block, err := aes.NewCipher(skey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt function decrypts given data with our secure key
func Decrypt(skey []byte, hexData string) (string, error) {
	ciphertext, err := decodeHex(hexData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(skey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// Generate32ByteKey generates a 32-byte key from an input secret string
func Generate32ByteKey(secret string) []byte {
	// Hash the secret using SHA-256
	hash := sha256.Sum256([]byte(secret))
	// Return the hash as a 32-byte key
	return hash[:]
}

// GenerateRandomString generates a random string of the specified number of bytes.
func GenerateRandomString(numBytes int) (string, error) {
	if numBytes <= 0 {
		return "", fmt.Errorf("number of bytes must be greater than zero")
	}

	// Create a byte slice to hold random bytes
	bytes := make([]byte, numBytes)

	// Fill the slice with random bytes
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to a hexadecimal string
	return hex.EncodeToString(bytes), nil
}

// GenerateSecret generates a persistent secret string based on the machine's hardware.
// It accepts the number of bytes `size` to generate the secret.
// The secret will be consistent across restarts but unique to the machine.
func GenerateSecret(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("size must be greater than 0")
	}

	// Get hardware-specific information (e.g., MAC address)
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	var macAddresses string
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 {
			macAddresses += iface.HardwareAddr.String()
		}
	}

	if macAddresses == "" {
		return "", fmt.Errorf("no MAC address found on the system")
	}

	// Hash the MAC address using SHA256 to derive a deterministic secret
	hash := sha256.Sum256([]byte(macAddresses))
	hashHex := hex.EncodeToString(hash[:])

	// Trim the hash to the desired size
	if size > len(hashHex) {
		return "", fmt.Errorf("requested size exceeds maximum hash size of %d", len(hashHex))
	}

	return hashHex[:size], nil
}
