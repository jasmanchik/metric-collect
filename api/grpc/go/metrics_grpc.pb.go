// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: metrics.proto

package metricsv1

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

// MetricClient is the client API for Metric service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	RequestsCounter(ctx context.Context, in *RequestsCounterRequest, opts ...grpc.CallOption) (*RequestsCounterResponse, error)
}

type metricClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricClient(cc grpc.ClientConnInterface) MetricClient {
	return &metricClient{cc}
}

func (c *metricClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/metrics.Metric/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricClient) RequestsCounter(ctx context.Context, in *RequestsCounterRequest, opts ...grpc.CallOption) (*RequestsCounterResponse, error) {
	out := new(RequestsCounterResponse)
	err := c.cc.Invoke(ctx, "/metrics.Metric/RequestsCounter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricServer is the server API for Metric service.
// All implementations must embed UnimplementedMetricServer
// for forward compatibility
type MetricServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	RequestsCounter(context.Context, *RequestsCounterRequest) (*RequestsCounterResponse, error)
	mustEmbedUnimplementedMetricServer()
}

// UnimplementedMetricServer must be embedded to have forward compatible implementations.
type UnimplementedMetricServer struct {
}

func (UnimplementedMetricServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMetricServer) RequestsCounter(context.Context, *RequestsCounterRequest) (*RequestsCounterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestsCounter not implemented")
}
func (UnimplementedMetricServer) mustEmbedUnimplementedMetricServer() {}

// UnsafeMetricServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricServer will
// result in compilation errors.
type UnsafeMetricServer interface {
	mustEmbedUnimplementedMetricServer()
}

func RegisterMetricServer(s grpc.ServiceRegistrar, srv MetricServer) {
	s.RegisterService(&Metric_ServiceDesc, srv)
}

func _Metric_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/metrics.Metric/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metric_RequestsCounter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestsCounterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServer).RequestsCounter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/metrics.Metric/RequestsCounter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServer).RequestsCounter(ctx, req.(*RequestsCounterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metric_ServiceDesc is the grpc.ServiceDesc for Metric service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metric_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metrics.Metric",
	HandlerType: (*MetricServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Metric_Ping_Handler,
		},
		{
			MethodName: "RequestsCounter",
			Handler:    _Metric_RequestsCounter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}