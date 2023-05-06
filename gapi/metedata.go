package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	userAgentHeader      = "user-agent"
	xForwardedForHeader  = "x-forwarded-host"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

// function to get clinetIP and userAgent from MetatData for grpc and http gateway
func (server *Server) extraMetaData(ctx context.Context) *MetaData {
	mtdt := &MetaData{}
	if data, ok := metadata.FromIncomingContext(ctx); ok {
		// get agent for api gateway
		if userAgents := data.Get(grpcGatewayUserAgent); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// get client for api gateway
		if ClientIPS := data.Get(xForwardedForHeader); len(ClientIPS) > 0 {
			mtdt.ClientIP = ClientIPS[0]
		}

		// get agent for grpc
		if userAgents := data.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
	}

	// get clientIP for grpc
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
