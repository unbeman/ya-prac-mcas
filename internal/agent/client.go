package agent

import (
	"context"
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/parser"
	"io"
	"log"
	"net/http"
)

type ClientMetric interface { //TODO: rename
	SendMetric(ctx context.Context, m metrics.Metric)
}

type clientMetric struct {
	addr   string
	client http.Client
}

func NewClientMetric(addr string, cli http.Client) *clientMetric {
	return &clientMetric{addr: addr, client: cli}
}

func (cs clientMetric) SendMetric(ctx context.Context, m metrics.Metric) { // TODO: write http connector
	url := parser.FormatURL(cs.addr, m)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	request.Header.Set("Content-Type", "text/plain")
	if err != nil {
		log.Fatalln(err)
	}
	response, err := cs.client.Do(request)
	if err != nil {
		log.Println(err) //TODO: retry request?
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Received status code: %v for post request to %v", response.StatusCode, url)
}
