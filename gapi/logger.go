package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	statusCode := codes.Unknown
	if s, ok := status.FromError(err); ok {
		statusCode = s.Code()
	}

	log.Err(err).
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_msg", statusCode.String()).
		Dur("duration", time.Since(start)).
		Msg("received gRPC request")

	return resp, err
}
