package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
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
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_msg", statusCode.String()).
		Dur("duration", time.Since(start)).
		Msg("received a gRPC request")

	return resp, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		res := &ResponseRecorder{w, http.StatusOK, []byte{}}
		next.ServeHTTP(res, r)

		logging := log.Info()
		if res.statusCode < 200 || res.statusCode > 299 {
			logging = log.Error().Bytes("body", res.body)
		}

		duration := time.Since(start)

		logging.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status_code", res.statusCode).
			Str("status_text", http.StatusText(res.statusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
