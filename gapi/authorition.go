package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/bacnx/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationBearer = "bearer"
)

func (server *Server) authorizeRequest(ctx context.Context) (*token.Payload, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing meetadata")
	}

	vals := metadata.Get(AuthorizationHeader)
	if len(vals) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := strings.Fields(vals[0])
	if len(authHeader) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	tokenType := authHeader[0]
	if strings.ToLower(tokenType) != AuthorizationBearer {
		return nil, fmt.Errorf("unsuppported authorization type: %s", tokenType)
	}

	accessToken := authHeader[1]
	return server.tokenMaker.VerifyToken(accessToken)
}
