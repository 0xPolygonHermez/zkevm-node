// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StateDBServiceClient is the client API for StateDBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StateDBServiceClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	SetProgram(ctx context.Context, in *SetProgramRequest, opts ...grpc.CallOption) (*SetProgramResponse, error)
	GetProgram(ctx context.Context, in *GetProgramRequest, opts ...grpc.CallOption) (*GetProgramResponse, error)
	Flush(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type stateDBServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStateDBServiceClient(cc grpc.ClientConnInterface) StateDBServiceClient {
	return &stateDBServiceClient{cc}
}

func (c *stateDBServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/statedb.v1.StateDBService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateDBServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/statedb.v1.StateDBService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateDBServiceClient) SetProgram(ctx context.Context, in *SetProgramRequest, opts ...grpc.CallOption) (*SetProgramResponse, error) {
	out := new(SetProgramResponse)
	err := c.cc.Invoke(ctx, "/statedb.v1.StateDBService/SetProgram", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateDBServiceClient) GetProgram(ctx context.Context, in *GetProgramRequest, opts ...grpc.CallOption) (*GetProgramResponse, error) {
	out := new(GetProgramResponse)
	err := c.cc.Invoke(ctx, "/statedb.v1.StateDBService/GetProgram", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateDBServiceClient) Flush(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/statedb.v1.StateDBService/Flush", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateDBServiceServer is the server API for StateDBService service.
// All implementations must embed UnimplementedStateDBServiceServer
// for forward compatibility
type StateDBServiceServer interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	SetProgram(context.Context, *SetProgramRequest) (*SetProgramResponse, error)
	GetProgram(context.Context, *GetProgramRequest) (*GetProgramResponse, error)
	Flush(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedStateDBServiceServer()
}

// UnimplementedStateDBServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStateDBServiceServer struct {
}

func (UnimplementedStateDBServiceServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedStateDBServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedStateDBServiceServer) SetProgram(context.Context, *SetProgramRequest) (*SetProgramResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetProgram not implemented")
}
func (UnimplementedStateDBServiceServer) GetProgram(context.Context, *GetProgramRequest) (*GetProgramResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProgram not implemented")
}
func (UnimplementedStateDBServiceServer) Flush(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (UnimplementedStateDBServiceServer) mustEmbedUnimplementedStateDBServiceServer() {}

// UnsafeStateDBServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StateDBServiceServer will
// result in compilation errors.
type UnsafeStateDBServiceServer interface {
	mustEmbedUnimplementedStateDBServiceServer()
}

func RegisterStateDBServiceServer(s grpc.ServiceRegistrar, srv StateDBServiceServer) {
	s.RegisterService(&StateDBService_ServiceDesc, srv)
}

func _StateDBService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateDBServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statedb.v1.StateDBService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateDBServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateDBService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateDBServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statedb.v1.StateDBService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateDBServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateDBService_SetProgram_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetProgramRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateDBServiceServer).SetProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statedb.v1.StateDBService/SetProgram",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateDBServiceServer).SetProgram(ctx, req.(*SetProgramRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateDBService_GetProgram_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProgramRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateDBServiceServer).GetProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statedb.v1.StateDBService/GetProgram",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateDBServiceServer).GetProgram(ctx, req.(*GetProgramRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateDBService_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateDBServiceServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/statedb.v1.StateDBService/Flush",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateDBServiceServer).Flush(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StateDBService_ServiceDesc is the grpc.ServiceDesc for StateDBService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StateDBService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statedb.v1.StateDBService",
	HandlerType: (*StateDBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _StateDBService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _StateDBService_Get_Handler,
		},
		{
			MethodName: "SetProgram",
			Handler:    _StateDBService_SetProgram_Handler,
		},
		{
			MethodName: "GetProgram",
			Handler:    _StateDBService_GetProgram_Handler,
		},
		{
			MethodName: "Flush",
			Handler:    _StateDBService_Flush_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statedb.proto",
}
