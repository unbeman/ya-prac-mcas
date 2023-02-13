package handlers

//TODO: rename file
import (
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/parser"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"log"
	"net/http"
	"time"
)

type CollectorServer struct {
	storage storage.Repository
}

func (cs *CollectorServer) UpdateHandler(w http.ResponseWriter, req *http.Request) { //TODO: move http logic to middleware; sep final handlers by metric type
	if req.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("expect method POST at /update/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}
	metric, err := parser.ParseMetric(fmt.Sprint(req.URL)) //TODO: parse separately
	if err != nil {
		switch err {
		case parser.ErrNotEnoughParams:
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case parser.ErrParse:
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case parser.ErrUnknownType:
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}
	}
	err = cs.storage.Update(metric.GetName(), metric)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		http.Error(w, "Storage error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func NewCollectorServer(repository storage.Repository) *CollectorServer {
	return &CollectorServer{storage: repository}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
