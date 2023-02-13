package main

import (
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage/ram"
	"log"
	"net/http"
)

const (
	ServerAddress = "localhost:8080"
)

func main() { //TODO: more logs
	mux := http.NewServeMux()
	storage := ram.NewRAMStorage()
	server := handlers.NewCollectorServer(storage)
	mux.HandleFunc("/update/", server.UpdateHandler)
	handler := handlers.Logging(mux)
	log.Fatal(http.ListenAndServe(ServerAddress, handler))
}
