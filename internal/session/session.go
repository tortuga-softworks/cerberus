package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
)

type Session struct {
	ID           string
	Username     string
	CreationTime time.Time
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, s)
}

// Generates a random session ID.
// The current implementation needs to be refactored to a more secure implementation.
// TODO securize
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
