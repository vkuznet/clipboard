package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// SecureClipboard represents secure clipboard
type SecureClipboard struct {
	history    []string
	lock       sync.Mutex
	key        []byte
	maxHistory int
}

// NewSecureClipboard provides instance of secure clipboard
func NewSecureClipboard(config *Config) (*SecureClipboard, error) {
	skey := make([]byte, 32) // AES-256 key
	if _, err := rand.Read(skey); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	if config.Secret != "" {
		skey = Generate32ByteKey(config.Secret)
	}
	return &SecureClipboard{
		history:    make([]string, 0, config.ClipboardSize),
		key:        skey,
		maxHistory: config.ClipboardSize,
	}, nil
}

func (sc *SecureClipboard) encrypt(data string) (string, error) {
	return Encrypt(sc.key, data)
}

func (sc *SecureClipboard) decrypt(hexData string) (string, error) {
	return Decrypt(sc.key, hexData)
}

// Copy function copies given data to secure clipboard
func (sc *SecureClipboard) Copy(data string) error {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	encrypted, err := sc.encrypt(data)
	if err != nil {
		return err
	}

	// Add to history
	if len(sc.history) >= sc.maxHistory {
		sc.history = sc.history[1:]
	}
	sc.history = append(sc.history, encrypted)

	return nil
}

// Paste function pastes data from secure clipboard
func (sc *SecureClipboard) Paste(idx int) (string, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if len(sc.history) == 0 {
		return "", errors.New("clipboard is empty")
	}
	if idx != 0 && len(sc.history) > idx-1 {
		//         return sc.decrypt(sc.history[idx-1])
		return sc.history[idx-1], nil
	}

	//     return sc.decrypt(sc.history[len(sc.history)-1])
	return sc.history[len(sc.history)-1], nil
}

func (sc *SecureClipboard) GetHistory() ([]string, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	var history []string
	for _, encrypted := range sc.history {
		/*
			decrypted, err := sc.decrypt(encrypted)
			if err != nil {
				return nil, err
			}
			history = append(history, decrypted)
		*/

		history = append(history, encrypted)
	}
	return history, nil
}

// Save history to a file
func (sc *SecureClipboard) SaveHistory(filename string) error {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(sc.history)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// Load history from a file
func (sc *SecureClipboard) LoadHistory(filename string) error {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var history []string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&history); err != nil {
		return err
	}

	sc.history = history
	return nil
}

// HTTP Handlers
func (sc *SecureClipboard) HandleCopy(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "data query parameter is required", http.StatusBadRequest)
		return
	}

	err := sc.Copy(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data copied to clipboard")
}

func (sc *SecureClipboard) HandlePaste(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	var idx int
	if item != "" {
		if val, err := strconv.Atoi(item); err == nil {
			idx = val
		}
	}

	data, err := sc.Paste(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, data)
}

func (sc *SecureClipboard) HandleHistory(w http.ResponseWriter, r *http.Request) {
	history, err := sc.GetHistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
