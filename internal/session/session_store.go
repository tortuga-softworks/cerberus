package session

import "context"

type SessionStore interface {
	Create(ctx context.Context, userID string) (*Session, error)
	FindByID(ctx context.Context, sessionID string) (*Session, error)
	Refresh(ctx context.Context, sessionID string) error
	Delete(ctx context.Context, sessionID string) error
}
