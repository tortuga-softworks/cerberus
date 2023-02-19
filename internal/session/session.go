package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Session struct {
	ID           string
	UserID       string
	CreationTime time.Time
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// Generates a random session ID.
func generateSessionID() (string, error) {
	// Generate 16 random bytes
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}

	// Append a timestamp to the random bytes
	ts := time.Now().UnixNano()
	tsBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		tsBytes[i] = byte(ts >> (i * 8))
	}
	idBytes := append(randBytes, tsBytes...)

	// Hash the ID bytes using SHA-256
	hash := sha256.Sum256(idBytes)

	// Convert the hash to a string and return it
	return hex.EncodeToString(hash[:]), nil
}
