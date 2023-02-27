package main

import (
	"log"
	"net/http"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

// TODO: wrap to init server
func main() { //TODO: more logs, add signals and context
	cfg := configs.NewServerConfig().FromEnv()
	ramRepo := storage.NewRAMRepository()
	ch := handlers.NewCollectorHandler(ramRepo)
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(cfg.Address, ch))
}
