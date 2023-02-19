package authentication

import (
	"errors"
	"reflect"

	"github.com/tortuga-softworks/cerberus/proto"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tortuga-softworks/cerberus/internal/session"
)

type AuthenticationServer struct {
	proto.UnimplementedAuthenticationServer
	authenticationService *AuthenticationService
}

func NewAuthenticationServer(authenticationService *AuthenticationService) (*AuthenticationServer, error) {
	if authenticationService == nil {
		return nil, errors.New("the session store can not be nil")
	}
	return &AuthenticationServer{authenticationService: authenticationService}, nil
}

func (as *AuthenticationServer) LogIn(ctx context.Context, in *proto.LogInRequest) (*proto.LogInResponse, error) {
	var email = in.Email

	if as.authenticationService.ValidateEmail(email) {
		sessionID, err := as.authenticationService.LogIn(email)

		if err != nil {
			return nil, err
		}

		return &proto.LogInResponse{SessionId: sessionID}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "email")
	}
}

func (as *AuthenticationServer) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	var sessionID = in.SessionId

	err := as.authenticationService.Refresh(sessionID)
	if err == nil {
		return &proto.RefreshResponse{}, nil
	} else {
		switch err.(type) {
		case *session.SessionNotFoundError:
			return nil, status.Error(codes.NotFound, sessionID)
		default:
			return nil, status.Errorf(codes.Internal, "%v: %v", reflect.TypeOf(err), err)
		}
	}
}
