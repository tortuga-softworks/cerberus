package authentication

import (
	"context"
	"errors"
	"strings"

	"github.com/tortuga-softworks/cerberus/internal/session"
	"github.com/tortuga-softworks/hestia/pkg/account"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	accountStore account.AccountStore
	sessionStore session.SessionStore
}

func NewAuthenticationService(accountStore account.AccountStore, sessionStore session.SessionStore) (*AuthenticationService, error) {
	if accountStore == nil {
		return nil, errors.New("could not create a registration service: account store is nil")
	}

	if sessionStore == nil {
		return nil, errors.New("could not create a registration service: session store is nil")
	}
	return &AuthenticationService{accountStore, sessionStore}, nil
}

func (service *AuthenticationService) LogIn(ctx context.Context, email string, password string) (string, error) {
	if !validateEmail(email) {
		return "", &EmailFormatError{Email: email}
	}

	account, err := service.accountStore.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = checkPassword(account.PasswordHash, password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", &PasswordMismatchError{Message: err.Error()}
		} else {
			return "", err
		}
	}

	session, err := service.sessionStore.CreateSession(ctx, email)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

// This functions returns whether the email
func validateEmail(email string) bool {
	return strings.Contains(email, "@")
}

func checkPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func (service *AuthenticationService) Refresh(ctx context.Context, sessionID string) error {
	return service.sessionStore.RefreshSession(ctx, sessionID)
}

func (service *AuthenticationService) LogOut(ctx context.Context, sessionID string) error {
	return service.sessionStore.DeleteSession(ctx, sessionID)
}
