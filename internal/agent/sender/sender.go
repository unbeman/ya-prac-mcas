package sender

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

const defaultRate = 1 * time.Second

type Sender interface {
	SendMetrics(ctx context.Context, slice metrics.ParamsSlice)
}

func GetSender(cfg configs.ConnectionConfig, pubKey *rsa.PublicKey) (Sender, error) {
	switch cfg.Protocol {
	case configs.GRPCProtocol:
		return NewGRPCSender(cfg)
	default:
		return NewHTTPSender(cfg, pubKey)
	}
}
