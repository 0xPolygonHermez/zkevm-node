// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: aggregator.proto

package pb

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

// AggregatorServiceClient is the client API for AggregatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AggregatorServiceClient interface {
	Channel(ctx context.Context, opts ...grpc.CallOption) (AggregatorService_ChannelClient, error)
}

type aggregatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAggregatorServiceClient(cc grpc.ClientConnInterface) AggregatorServiceClient {
	return &aggregatorServiceClient{cc}
}

func (c *aggregatorServiceClient) Channel(ctx context.Context, opts ...grpc.CallOption) (AggregatorService_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &AggregatorService_ServiceDesc.Streams[0], "/aggregator.v1.AggregatorService/Channel", opts...)
	if err != nil {
		return nil, err
	}
	x := &aggregatorServiceChannelClient{stream}
	return x, nil
}

type AggregatorService_ChannelClient interface {
	Send(*ProverMessage) error
	Recv() (*AggregatorMessage, error)
	grpc.ClientStream
}

type aggregatorServiceChannelClient struct {
	grpc.ClientStream
}

func (x *aggregatorServiceChannelClient) Send(m *ProverMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aggregatorServiceChannelClient) Recv() (*AggregatorMessage, error) {
	m := new(AggregatorMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AggregatorServiceServer is the server API for AggregatorService service.
// All implementations must embed UnimplementedAggregatorServiceServer
// for forward compatibility
type AggregatorServiceServer interface {
	Channel(AggregatorService_ChannelServer) error
	mustEmbedUnimplementedAggregatorServiceServer()
}

// UnimplementedAggregatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAggregatorServiceServer struct {
}

func (UnimplementedAggregatorServiceServer) Channel(AggregatorService_ChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method Channel not implemented")
}
func (UnimplementedAggregatorServiceServer) mustEmbedUnimplementedAggregatorServiceServer() {}

// UnsafeAggregatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AggregatorServiceServer will
// result in compilation errors.
type UnsafeAggregatorServiceServer interface {
	mustEmbedUnimplementedAggregatorServiceServer()
}

func RegisterAggregatorServiceServer(s grpc.ServiceRegistrar, srv AggregatorServiceServer) {
	s.RegisterService(&AggregatorService_ServiceDesc, srv)
}

func _AggregatorService_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AggregatorServiceServer).Channel(&aggregatorServiceChannelServer{stream})
}

type AggregatorService_ChannelServer interface {
	Send(*AggregatorMessage) error
	Recv() (*ProverMessage, error)
	grpc.ServerStream
}

type aggregatorServiceChannelServer struct {
	grpc.ServerStream
}

func (x *aggregatorServiceChannelServer) Send(m *AggregatorMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aggregatorServiceChannelServer) Recv() (*ProverMessage, error) {
	m := new(ProverMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AggregatorService_ServiceDesc is the grpc.ServiceDesc for AggregatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AggregatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aggregator.v1.AggregatorService",
	HandlerType: (*AggregatorServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Channel",
			Handler:       _AggregatorService_Channel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "aggregator.proto",
}
