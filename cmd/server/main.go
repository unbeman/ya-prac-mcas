package main

import (
	"log"
	"net/http"

	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

const (
	ServerAddress = "localhost:8080"
)

func main() { //TODO: more logs, add signals and context
	ramRepo := storage.NewRepository(storage.NewGaugeRAMStorage(), storage.NewCounterRAMStorage())
	ch, err := handlers.NewCollectorHandler(ramRepo)
	if err != nil {
		log.Fatalln("Unable to create handler")
	}
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(ServerAddress, ch))
}
