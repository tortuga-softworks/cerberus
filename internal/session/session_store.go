package session

type SessionStore interface {
	CreateSession(username string) (*Session, error)
	// RetrieveSession(sessionID string) (bool, error)
	RefreshSession(sessionID string) error
	// DeleteSession(sessionID string) (string, error)
}
