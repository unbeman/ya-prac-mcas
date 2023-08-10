package sender

import (
	"context"
	"crypto/rsa"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
	pb "github.com/unbeman/ya-prac-mcas/proto"
)

type GRPCSender struct {
	client      pb.MetricsCollectorClient
	timeout     time.Duration
	rateLimiter *rate.Limiter
	publicKey   *rsa.PublicKey
}

func NewGRPCSender(cfg configs.ConnectionConfig) (*GRPCSender, error) {
	rl := rate.NewLimiter(rate.Every(defaultRate), cfg.RateTokensCount)

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("NewGRPCSender: can't dial to %s: %w", cfg.Address, err)
	}
	client := pb.NewMetricsCollectorClient(conn)
	return &GRPCSender{client: client, timeout: cfg.ReportTimeout, rateLimiter: rl}, nil
}

func (gs *GRPCSender) SendMetrics(ctx context.Context, slice metrics.ParamsSlice) {
	err := gs.rateLimiter.Wait(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	ctx2, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	ip, err := utils.GetOutboundIP()
	if err != nil {
		log.Error(err)
		return
	}

	meta := metadata.New(map[string]string{"x-real-ip": ip})
	ctx2 = metadata.NewOutgoingContext(ctx2, meta)

	protMetrics := slice.ToProto()

	_, err = gs.client.UpdateMetrics(ctx2, &pb.UpdateMetricsRequest{Metrics: protMetrics}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		if e, ok := status.FromError(err); ok {
			log.Errorf("SendMetrics: status code %d, msg: %s", e.Code(), e.Message())
		} else {
			log.Errorf("SendMetrics: %s", err)
		}
	}

	log.Info("Metrics send")
}
