package cerberus

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedAuthenticationServer
}

func (s *Server) Login(ctx context.Context, in *LoginRequest) (*LoginResponse, error) {
	var username = in.Username

	if username != "" {
		return &LoginResponse{SessionId: createSessionID()}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "the username cannot be empty")
	}
}

func createSessionID() string { // TODO securize
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
