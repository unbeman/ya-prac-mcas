package handlers

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	pb "github.com/unbeman/ya-prac-mcas/proto"
)

type GRPCService struct {
	pb.UnimplementedMetricsCollectorServer
	control *controller.Controller
}

func NewGRPCService(control *controller.Controller) *GRPCService {
	return &GRPCService{control: control}
}

func (g *GRPCService) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	params := metrics.Params{Name: in.Name, Type: in.Type}

	m, err := g.control.GetMetric(ctx, params)
	if err != nil {
		return nil, g.processedError(err)
	}

	out := &pb.GetMetricResponse{Metric: m.ToProto()}
	out.Metric.Hash = g.control.GetHash(m)
	return out, nil
}
func (g *GRPCService) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	ms, err := g.control.GetAll(ctx)
	if err != nil {
		return nil, g.processedError(err)
	}

	protoMetrics := make([]*pb.Metric, 0, len(ms))
	for _, m := range ms {
		mp := m.ToProto()
		mp.Hash = g.control.GetHash(m)
		protoMetrics = append(protoMetrics, mp)
	}

	out := &pb.GetMetricsResponse{Metrics: protoMetrics}
	return out, nil
}
func (g *GRPCService) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	params, err := metrics.ParseProto(in.Metric, metrics.PName, metrics.PType, metrics.PValue)
	if err != nil {
		return nil, g.processedError(err)
	}

	m, err := g.control.UpdateMetric(ctx, params)
	if err != nil {
		return nil, g.processedError(err)
	}

	out := &pb.UpdateMetricResponse{Metric: m.ToProto()}
	out.Metric.Hash = g.control.GetHash(m)
	return out, nil
}
func (g *GRPCService) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var metricsParams metrics.ParamsSlice
	err := metricsParams.ParseProto(in.Metrics)
	if err != nil {
		return nil, g.processedError(err)
	}

	metricsParams, err = g.control.UpdateMetrics(ctx, metricsParams)
	if err != nil {
		return nil, g.processedError(err)
	}

	out := &pb.UpdateMetricsResponse{Metrics: metricsParams.ToProto()}
	return out, nil
}
func (g *GRPCService) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	err := g.control.Ping(ctx)
	if err != nil {
		return nil, g.processedError(err)
	}
	return &pb.PingResponse{}, nil
}

func (g *GRPCService) processedError(err error) error {
	var grpcCode codes.Code
	switch {
	case errors.Is(err, controller.ErrInvalidHash):
		grpcCode = codes.InvalidArgument
	case errors.Is(err, metrics.ErrInvalidType):
		grpcCode = codes.Unimplemented
	case errors.Is(err, metrics.ErrInvalidValue):
		grpcCode = codes.InvalidArgument
	case errors.Is(err, storage.ErrNotFound):
		grpcCode = codes.NotFound
	default:
		grpcCode = codes.Internal
	}
	return status.Error(grpcCode, err.Error())
}
