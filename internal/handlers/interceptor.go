package handlers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"

	"github.com/unbeman/ya-prac-mcas/internal/utils"
)

//todo: add logging interceptor

func IPCheckerServerInterceptor(trustedSubnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if trustedSubnet != nil {
			meta, ok := metadata.FromIncomingContext(ctx)
			log.Info(meta)
			if ok {
				clientIP := meta.Get("X-Real-IP")[0]

				if err := utils.CheckIPBelongsNetwork(clientIP, trustedSubnet); err != nil {
					return nil, status.Error(codes.PermissionDenied, err.Error())
				}
			}
		}
		return handler(ctx, req)
	}
}
