package authentication

import (
	"errors"
	"strings"

	"github.com/tortuga-softworks/cerberus/internal/session"
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

func (as *AuthenticationService) LogIn(email string) (string, error) {
	session, err := as.sessionStore.CreateSession(email)

	if err != nil {
		return "", err
	}

	return session.ID, nil
}

// This functions returns whether the email
func (as *AuthenticationService) ValidateEmail(email string) bool {
	return strings.Contains(email, "@")
}

func (as *AuthenticationService) Refresh(sessionID string) error {
	return as.sessionStore.RefreshSession(sessionID)
}
