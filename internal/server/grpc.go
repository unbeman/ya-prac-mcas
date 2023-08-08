package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"

	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	pb "github.com/unbeman/ya-prac-mcas/proto"
)

type GRPCServer struct {
	address string
	server  *grpc.Server
	service *handlers.GRPCService
}

func NewGRPCServer(addr string, control *controller.Controller) *GRPCServer {
	server := grpc.NewServer()
	service := handlers.NewGRPCService(control)

	return &GRPCServer{address: addr, server: server, service: service}
}

func (g *GRPCServer) GetAddress() string {
	return g.address
}

func (g *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("can't bind address: %w", err)
	}

	pb.RegisterMetricsCollectorServer(g.server, g.service)

	log.Info("starting gRPC server")
	return g.server.Serve(listen)
}

func (g *GRPCServer) Close() error {
	g.server.GracefulStop()
	return nil
}
