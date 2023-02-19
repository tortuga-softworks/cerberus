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
		return nil, errors.New("could not create a authentication server: authentication service is nil")
	}
	return &AuthenticationServer{authenticationService: authenticationService}, nil
}

func (as *AuthenticationServer) LogIn(ctx context.Context, in *proto.LogInRequest) (*proto.LogInResponse, error) {
	var email = in.Email
	var password = in.Password

	sessionID, err := as.authenticationService.LogIn(ctx, email, password)
	if err != nil {
		switch err.(type) {
		case *EmailFormatError:
			return nil, status.Error(codes.InvalidArgument, "email")
		case *PasswordMismatchError:
			return nil, status.Error(codes.Unauthenticated, email)
		default:
			return nil, status.Errorf(codes.Internal, "%v: %v", reflect.TypeOf(err), err)
		}
	}

	return &proto.LogInResponse{SessionId: sessionID}, nil
}

func (as *AuthenticationServer) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	var sessionID = in.SessionId

	err := as.authenticationService.Refresh(ctx, sessionID)
	if err != nil {
		switch err.(type) {
		case *session.SessionNotFoundError:
			return nil, status.Error(codes.NotFound, sessionID)
		default:
			return nil, status.Errorf(codes.Internal, "%v: %v", reflect.TypeOf(err), err)
		}
	}

	return &proto.RefreshResponse{}, nil
}

func (as *AuthenticationServer) LogOut(ctx context.Context, in *proto.LogOutRequest) (*proto.LogOutResponse, error) {
	var sessionID = in.SessionId

	err := as.authenticationService.LogOut(ctx, sessionID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", reflect.TypeOf(err), err)
	}

	return &proto.LogOutResponse{}, nil
}
