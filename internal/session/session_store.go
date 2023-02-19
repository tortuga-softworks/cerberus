package session

import "context"

type SessionStore interface {
	CreateSession(ctx context.Context, userID string) (*Session, error)
	FindSessionByID(ctx context.Context, sessionID string) (*Session, error)
	RefreshSession(ctx context.Context, sessionID string) error
	DeleteSession(ctx context.Context, sessionID string) error
}
