package server

import (
	"context"
	"crypto/rsa"
	"net"
	"net/http"

	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, control *controller.Controller, privateKey *rsa.PrivateKey, trustedSubnet *net.IPNet) *HTTPServer {
	handler := handlers.NewCollectorHandler(control, privateKey, trustedSubnet)
	return &HTTPServer{server: &http.Server{Addr: addr, Handler: handler}}
}

func (h *HTTPServer) GetAddress() string {
	return h.server.Addr
}

func (h *HTTPServer) Run() error {
	return h.server.ListenAndServe()
}

func (h *HTTPServer) Close() error {
	return h.server.Shutdown(context.TODO())
}
