package authentication

import (
	proto "cerberus/proto"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthenticationServer struct {
	proto.UnimplementedAuthenticationServer
	AuthenticationService *AuthenticationService
}

func (as *AuthenticationServer) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	var username = in.Username

	if as.AuthenticationService.ValidateUsername(username) {
		sessionID, err := as.AuthenticationService.Login(username)

		if err != nil {
			return nil, err
		}

		return &proto.LoginResponse{SessionId: sessionID}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "the username is not valid")
	}
}

func (as *AuthenticationServer) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	var sessionID = in.SessionId

	refreshed, err := as.AuthenticationService.Refresh(sessionID)
	if err == nil {
		if refreshed {
			return &proto.RefreshResponse{}, nil
		} else {
			return nil, status.Error(codes.NotFound, "the session does not exist")
		}
	} else {
		return nil, err
	}
}
