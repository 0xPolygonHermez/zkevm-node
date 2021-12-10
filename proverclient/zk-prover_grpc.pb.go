// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proverclient

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

// ZKProverClient is the client API for ZKProver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ZKProverClient interface {
	GetStatus(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*State, error)
	GenProof(ctx context.Context, opts ...grpc.CallOption) (ZKProver_GenProofClient, error)
	Cancel(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*State, error)
	GetProof(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*Proof, error)
}

type zKProverClient struct {
	cc grpc.ClientConnInterface
}

func NewZKProverClient(cc grpc.ClientConnInterface) ZKProverClient {
	return &zKProverClient{cc}
}

func (c *zKProverClient) GetStatus(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*State, error) {
	out := new(State)
	err := c.cc.Invoke(ctx, "/zkprover.ZKProver/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zKProverClient) GenProof(ctx context.Context, opts ...grpc.CallOption) (ZKProver_GenProofClient, error) {
	stream, err := c.cc.NewStream(ctx, &ZKProver_ServiceDesc.Streams[0], "/zkprover.ZKProver/GenProof", opts...)
	if err != nil {
		return nil, err
	}
	x := &zKProverGenProofClient{stream}
	return x, nil
}

type ZKProver_GenProofClient interface {
	Send(*Batch) error
	Recv() (*Proof, error)
	grpc.ClientStream
}

type zKProverGenProofClient struct {
	grpc.ClientStream
}

func (x *zKProverGenProofClient) Send(m *Batch) error {
	return x.ClientStream.SendMsg(m)
}

func (x *zKProverGenProofClient) Recv() (*Proof, error) {
	m := new(Proof)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *zKProverClient) Cancel(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*State, error) {
	out := new(State)
	err := c.cc.Invoke(ctx, "/zkprover.ZKProver/Cancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zKProverClient) GetProof(ctx context.Context, in *NoParams, opts ...grpc.CallOption) (*Proof, error) {
	out := new(Proof)
	err := c.cc.Invoke(ctx, "/zkprover.ZKProver/GetProof", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZKProverServer is the server API for ZKProver service.
// All implementations must embed UnimplementedZKProverServer
// for forward compatibility
type ZKProverServer interface {
	GetStatus(context.Context, *NoParams) (*State, error)
	GenProof(ZKProver_GenProofServer) error
	Cancel(context.Context, *NoParams) (*State, error)
	GetProof(context.Context, *NoParams) (*Proof, error)
	mustEmbedUnimplementedZKProverServer()
}

// UnimplementedZKProverServer must be embedded to have forward compatible implementations.
type UnimplementedZKProverServer struct {
}

func (UnimplementedZKProverServer) GetStatus(context.Context, *NoParams) (*State, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedZKProverServer) GenProof(ZKProver_GenProofServer) error {
	return status.Errorf(codes.Unimplemented, "method GenProof not implemented")
}
func (UnimplementedZKProverServer) Cancel(context.Context, *NoParams) (*State, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cancel not implemented")
}
func (UnimplementedZKProverServer) GetProof(context.Context, *NoParams) (*Proof, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProof not implemented")
}
func (UnimplementedZKProverServer) mustEmbedUnimplementedZKProverServer() {}

// UnsafeZKProverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ZKProverServer will
// result in compilation errors.
type UnsafeZKProverServer interface {
	mustEmbedUnimplementedZKProverServer()
}

func RegisterZKProverServer(s grpc.ServiceRegistrar, srv ZKProverServer) {
	s.RegisterService(&ZKProver_ServiceDesc, srv)
}

func _ZKProver_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZKProverServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zkprover.ZKProver/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZKProverServer).GetStatus(ctx, req.(*NoParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZKProver_GenProof_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ZKProverServer).GenProof(&zKProverGenProofServer{stream})
}

type ZKProver_GenProofServer interface {
	Send(*Proof) error
	Recv() (*Batch, error)
	grpc.ServerStream
}

type zKProverGenProofServer struct {
	grpc.ServerStream
}

func (x *zKProverGenProofServer) Send(m *Proof) error {
	return x.ServerStream.SendMsg(m)
}

func (x *zKProverGenProofServer) Recv() (*Batch, error) {
	m := new(Batch)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ZKProver_Cancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZKProverServer).Cancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zkprover.ZKProver/Cancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZKProverServer).Cancel(ctx, req.(*NoParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZKProver_GetProof_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZKProverServer).GetProof(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zkprover.ZKProver/GetProof",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZKProverServer).GetProof(ctx, req.(*NoParams))
	}
	return interceptor(ctx, in, info, handler)
}

// ZKProver_ServiceDesc is the grpc.ServiceDesc for ZKProver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ZKProver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zkprover.ZKProver",
	HandlerType: (*ZKProverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _ZKProver_GetStatus_Handler,
		},
		{
			MethodName: "Cancel",
			Handler:    _ZKProver_Cancel_Handler,
		},
		{
			MethodName: "GetProof",
			Handler:    _ZKProver_GetProof_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GenProof",
			Handler:       _ZKProver_GenProof_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/zk-prover.proto",
}
