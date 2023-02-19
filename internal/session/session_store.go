package session

import "context"

type SessionStore interface {
	CreateSession(ctx context.Context, username string) (*Session, error)
	// RetrieveSession(ctx context.Context, sessionID string) (bool, error)
	RefreshSession(ctx context.Context, sessionID string) error
	// DeleteSession(ctx context.Context, sessionID string) (string, error)
}
