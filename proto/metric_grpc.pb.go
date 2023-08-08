// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/metric.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MetricsCollector_GetMetric_FullMethodName     = "/mcas.MetricsCollector/GetMetric"
	MetricsCollector_GetMetrics_FullMethodName    = "/mcas.MetricsCollector/GetMetrics"
	MetricsCollector_UpdateMetric_FullMethodName  = "/mcas.MetricsCollector/UpdateMetric"
	MetricsCollector_UpdateMetrics_FullMethodName = "/mcas.MetricsCollector/UpdateMetrics"
	MetricsCollector_Ping_FullMethodName          = "/mcas.MetricsCollector/Ping"
)

// MetricsCollectorClient is the client API for MetricsCollector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsCollectorClient interface {
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (*GetMetricsResponse, error)
	UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error)
	UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type metricsCollectorClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsCollectorClient(cc grpc.ClientConnInterface) MetricsCollectorClient {
	return &metricsCollectorClient{cc}
}

func (c *metricsCollectorClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, MetricsCollector_GetMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsCollectorClient) GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (*GetMetricsResponse, error) {
	out := new(GetMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsCollector_GetMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsCollectorClient) UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error) {
	out := new(UpdateMetricResponse)
	err := c.cc.Invoke(ctx, MetricsCollector_UpdateMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsCollectorClient) UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error) {
	out := new(UpdateMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsCollector_UpdateMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsCollectorClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, MetricsCollector_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsCollectorServer is the server API for MetricsCollector service.
// All implementations must embed UnimplementedMetricsCollectorServer
// for forward compatibility
type MetricsCollectorServer interface {
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	GetMetrics(context.Context, *GetMetricsRequest) (*GetMetricsResponse, error)
	UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error)
	UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedMetricsCollectorServer()
}

// UnimplementedMetricsCollectorServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsCollectorServer struct {
}

func (UnimplementedMetricsCollectorServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsCollectorServer) GetMetrics(context.Context, *GetMetricsRequest) (*GetMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedMetricsCollectorServer) UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMetricsCollectorServer) UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetrics not implemented")
}
func (UnimplementedMetricsCollectorServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMetricsCollectorServer) mustEmbedUnimplementedMetricsCollectorServer() {}

// UnsafeMetricsCollectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsCollectorServer will
// result in compilation errors.
type UnsafeMetricsCollectorServer interface {
	mustEmbedUnimplementedMetricsCollectorServer()
}

func RegisterMetricsCollectorServer(s grpc.ServiceRegistrar, srv MetricsCollectorServer) {
	s.RegisterService(&MetricsCollector_ServiceDesc, srv)
}

func _MetricsCollector_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsCollectorServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsCollector_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsCollectorServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsCollector_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsCollectorServer).GetMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsCollector_GetMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsCollectorServer).GetMetrics(ctx, req.(*GetMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsCollector_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsCollectorServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsCollector_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsCollectorServer).UpdateMetric(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsCollector_UpdateMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsCollectorServer).UpdateMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsCollector_UpdateMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsCollectorServer).UpdateMetrics(ctx, req.(*UpdateMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsCollector_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsCollectorServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsCollector_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsCollectorServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsCollector_ServiceDesc is the grpc.ServiceDesc for MetricsCollector service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsCollector_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mcas.MetricsCollector",
	HandlerType: (*MetricsCollectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMetric",
			Handler:    _MetricsCollector_GetMetric_Handler,
		},
		{
			MethodName: "GetMetrics",
			Handler:    _MetricsCollector_GetMetrics_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _MetricsCollector_UpdateMetric_Handler,
		},
		{
			MethodName: "UpdateMetrics",
			Handler:    _MetricsCollector_UpdateMetrics_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _MetricsCollector_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metric.proto",
}