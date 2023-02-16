package authentication

import (
	"cerberus/internal/session"
	"errors"
	"strings"
)

type AuthenticationService struct {
	sessionStore session.SessionStore
}

func NewAuthenticationService(sessionStore session.SessionStore) (*AuthenticationService, error) {
	if sessionStore == nil {
		return nil, errors.New("the session store can not be nil")
	}
	return &AuthenticationService{sessionStore}, nil
}

func (as *AuthenticationService) Login(username string) (string, error) {
	session, err := as.sessionStore.CreateSession(username)

	if err != nil {
		return "", err
	}

	return session.ID, nil
}

// This functions returns whether the username
func (as *AuthenticationService) ValidateUsername(username string) bool {
	return strings.Contains(username, "@")
}

func (as *AuthenticationService) Refresh(sessionID string) (bool, error) {
	return as.sessionStore.RefreshSession(sessionID)
}
