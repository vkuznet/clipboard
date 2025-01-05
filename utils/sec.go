package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
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
