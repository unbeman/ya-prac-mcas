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

func main() { //TODO: more logs, add signals and context
	repo := ram.NewRAMStorage()
	ch := handlers.NewCollectorHandler(repo)
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(ServerAddress, ch))
}
