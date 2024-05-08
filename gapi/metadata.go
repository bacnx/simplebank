package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	GrpcGatewayUserAgent = "grpcgateway-user-agent"
	XForwardedHost       = "x-forwarded-host"
	UserAgent            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := Metadata{}

	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := headers.Get(GrpcGatewayUserAgent); len(vals) > 0 {
			mtdt.UserAgent = vals[0]
		}

		if vals := headers.Get(UserAgent); len(vals) > 0 {
			mtdt.UserAgent = vals[0]
		}

		if vals := headers.Get(XForwardedHost); len(vals) > 0 {
			mtdt.ClientIP = vals[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return &mtdt
}
