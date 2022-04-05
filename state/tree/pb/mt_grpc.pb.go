// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: mt.proto

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

// MTServiceClient is the client API for MTService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MTServiceClient interface {
	// Getters
	/// Get balance for a given address at a given root
	GetBalance(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	/// Get nonce for a given address at a given root
	GetNonce(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetNonceResponse, error)
	/// Get code for a given address at a given root
	GetCode(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetCodeResponse, error)
	/// Get code hash for a given address at a given root
	GetCodeHash(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetCodeHashResponse, error)
	/// Get smart contract storage for a given address and position at a given root
	GetStorageAt(ctx context.Context, in *GetStorageAtRequest, opts ...grpc.CallOption) (*GetStorageAtResponse, error)
	/// Get set of reads and writes that updated the state at a given root
	GetStateTransitionNodes(ctx context.Context, in *GetStateTransitionNodesRequest, opts ...grpc.CallOption) (*GetStateTransitionNodesResponse, error)
	/// Reverse a hash of an exisiting Merkletree node
	ReverseHash(ctx context.Context, in *ReverseHashRequest, opts ...grpc.CallOption) (*ReverseHashResponse, error)
	// Setters
	/// Set the balance for an account at a root
	SetBalance(ctx context.Context, in *SetBalanceRequest, opts ...grpc.CallOption) (*CommonSetResponse, error)
	/// Set the nonce of an account at a root
	SetNonce(ctx context.Context, in *SetNonceRequest, opts ...grpc.CallOption) (*CommonSetResponse, error)
	/// Set the code for an account at a root
	SetCode(ctx context.Context, in *SetCodeRequest, opts ...grpc.CallOption) (*CommonSetResponse, error)
	/// Set smart contract storage for an account and position at a root
	SetStorageAt(ctx context.Context, in *SetStorageAtRequest, opts ...grpc.CallOption) (*CommonSetResponse, error)
	/// Set an entry of the reverse hash table
	SetHashValue(ctx context.Context, in *HashValuePair, opts ...grpc.CallOption) (*SetHashValueResponse, error)
	/// Set many entries of the reverse hash table
	SetStateTransitionNodes(ctx context.Context, in *SetStateTransitionNodesRequest, opts ...grpc.CallOption) (*SetStateTransitionNodesResponse, error)
}

type mTServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMTServiceClient(cc grpc.ClientConnInterface) MTServiceClient {
	return &mTServiceClient{cc}
}

func (c *mTServiceClient) GetBalance(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) GetNonce(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetNonceResponse, error) {
	out := new(GetNonceResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetNonce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) GetCode(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetCodeResponse, error) {
	out := new(GetCodeResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) GetCodeHash(ctx context.Context, in *CommonGetRequest, opts ...grpc.CallOption) (*GetCodeHashResponse, error) {
	out := new(GetCodeHashResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetCodeHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) GetStorageAt(ctx context.Context, in *GetStorageAtRequest, opts ...grpc.CallOption) (*GetStorageAtResponse, error) {
	out := new(GetStorageAtResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetStorageAt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) GetStateTransitionNodes(ctx context.Context, in *GetStateTransitionNodesRequest, opts ...grpc.CallOption) (*GetStateTransitionNodesResponse, error) {
	out := new(GetStateTransitionNodesResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/GetStateTransitionNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) ReverseHash(ctx context.Context, in *ReverseHashRequest, opts ...grpc.CallOption) (*ReverseHashResponse, error) {
	out := new(ReverseHashResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/ReverseHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetBalance(ctx context.Context, in *SetBalanceRequest, opts ...grpc.CallOption) (*CommonSetResponse, error) {
	out := new(CommonSetResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetNonce(ctx context.Context, in *SetNonceRequest, opts ...grpc.CallOption) (*CommonSetResponse, error) {
	out := new(CommonSetResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetNonce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetCode(ctx context.Context, in *SetCodeRequest, opts ...grpc.CallOption) (*CommonSetResponse, error) {
	out := new(CommonSetResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetStorageAt(ctx context.Context, in *SetStorageAtRequest, opts ...grpc.CallOption) (*CommonSetResponse, error) {
	out := new(CommonSetResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetStorageAt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetHashValue(ctx context.Context, in *HashValuePair, opts ...grpc.CallOption) (*SetHashValueResponse, error) {
	out := new(SetHashValueResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetHashValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mTServiceClient) SetStateTransitionNodes(ctx context.Context, in *SetStateTransitionNodesRequest, opts ...grpc.CallOption) (*SetStateTransitionNodesResponse, error) {
	out := new(SetStateTransitionNodesResponse)
	err := c.cc.Invoke(ctx, "/mt.v1.MTService/SetStateTransitionNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MTServiceServer is the server API for MTService service.
// All implementations must embed UnimplementedMTServiceServer
// for forward compatibility
type MTServiceServer interface {
	// Getters
	/// Get balance for a given address at a given root
	GetBalance(context.Context, *CommonGetRequest) (*GetBalanceResponse, error)
	/// Get nonce for a given address at a given root
	GetNonce(context.Context, *CommonGetRequest) (*GetNonceResponse, error)
	/// Get code for a given address at a given root
	GetCode(context.Context, *CommonGetRequest) (*GetCodeResponse, error)
	/// Get code hash for a given address at a given root
	GetCodeHash(context.Context, *CommonGetRequest) (*GetCodeHashResponse, error)
	/// Get smart contract storage for a given address and position at a given root
	GetStorageAt(context.Context, *GetStorageAtRequest) (*GetStorageAtResponse, error)
	/// Get set of reads and writes that updated the state at a given root
	GetStateTransitionNodes(context.Context, *GetStateTransitionNodesRequest) (*GetStateTransitionNodesResponse, error)
	/// Reverse a hash of an exisiting Merkletree node
	ReverseHash(context.Context, *ReverseHashRequest) (*ReverseHashResponse, error)
	// Setters
	/// Set the balance for an account at a root
	SetBalance(context.Context, *SetBalanceRequest) (*CommonSetResponse, error)
	/// Set the nonce of an account at a root
	SetNonce(context.Context, *SetNonceRequest) (*CommonSetResponse, error)
	/// Set the code for an account at a root
	SetCode(context.Context, *SetCodeRequest) (*CommonSetResponse, error)
	/// Set smart contract storage for an account and position at a root
	SetStorageAt(context.Context, *SetStorageAtRequest) (*CommonSetResponse, error)
	/// Set an entry of the reverse hash table
	SetHashValue(context.Context, *HashValuePair) (*SetHashValueResponse, error)
	/// Set many entries of the reverse hash table
	SetStateTransitionNodes(context.Context, *SetStateTransitionNodesRequest) (*SetStateTransitionNodesResponse, error)
	mustEmbedUnimplementedMTServiceServer()
}

// UnimplementedMTServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMTServiceServer struct {
}

func (UnimplementedMTServiceServer) GetBalance(context.Context, *CommonGetRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedMTServiceServer) GetNonce(context.Context, *CommonGetRequest) (*GetNonceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNonce not implemented")
}
func (UnimplementedMTServiceServer) GetCode(context.Context, *CommonGetRequest) (*GetCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCode not implemented")
}
func (UnimplementedMTServiceServer) GetCodeHash(context.Context, *CommonGetRequest) (*GetCodeHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCodeHash not implemented")
}
func (UnimplementedMTServiceServer) GetStorageAt(context.Context, *GetStorageAtRequest) (*GetStorageAtResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStorageAt not implemented")
}
func (UnimplementedMTServiceServer) GetStateTransitionNodes(context.Context, *GetStateTransitionNodesRequest) (*GetStateTransitionNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStateTransitionNodes not implemented")
}
func (UnimplementedMTServiceServer) ReverseHash(context.Context, *ReverseHashRequest) (*ReverseHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReverseHash not implemented")
}
func (UnimplementedMTServiceServer) SetBalance(context.Context, *SetBalanceRequest) (*CommonSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetBalance not implemented")
}
func (UnimplementedMTServiceServer) SetNonce(context.Context, *SetNonceRequest) (*CommonSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNonce not implemented")
}
func (UnimplementedMTServiceServer) SetCode(context.Context, *SetCodeRequest) (*CommonSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetCode not implemented")
}
func (UnimplementedMTServiceServer) SetStorageAt(context.Context, *SetStorageAtRequest) (*CommonSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStorageAt not implemented")
}
func (UnimplementedMTServiceServer) SetHashValue(context.Context, *HashValuePair) (*SetHashValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetHashValue not implemented")
}
func (UnimplementedMTServiceServer) SetStateTransitionNodes(context.Context, *SetStateTransitionNodesRequest) (*SetStateTransitionNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStateTransitionNodes not implemented")
}
func (UnimplementedMTServiceServer) mustEmbedUnimplementedMTServiceServer() {}

// UnsafeMTServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MTServiceServer will
// result in compilation errors.
type UnsafeMTServiceServer interface {
	mustEmbedUnimplementedMTServiceServer()
}

func RegisterMTServiceServer(s grpc.ServiceRegistrar, srv MTServiceServer) {
	s.RegisterService(&MTService_ServiceDesc, srv)
}

func _MTService_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetBalance(ctx, req.(*CommonGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_GetNonce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetNonce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetNonce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetNonce(ctx, req.(*CommonGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_GetCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetCode(ctx, req.(*CommonGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_GetCodeHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetCodeHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetCodeHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetCodeHash(ctx, req.(*CommonGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_GetStorageAt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStorageAtRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetStorageAt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetStorageAt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetStorageAt(ctx, req.(*GetStorageAtRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_GetStateTransitionNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateTransitionNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).GetStateTransitionNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/GetStateTransitionNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).GetStateTransitionNodes(ctx, req.(*GetStateTransitionNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_ReverseHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReverseHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).ReverseHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/ReverseHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).ReverseHash(ctx, req.(*ReverseHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetBalance(ctx, req.(*SetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetNonce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetNonceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetNonce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetNonce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetNonce(ctx, req.(*SetNonceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetCode(ctx, req.(*SetCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetStorageAt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStorageAtRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetStorageAt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetStorageAt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetStorageAt(ctx, req.(*SetStorageAtRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetHashValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HashValuePair)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetHashValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetHashValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetHashValue(ctx, req.(*HashValuePair))
	}
	return interceptor(ctx, in, info, handler)
}

func _MTService_SetStateTransitionNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStateTransitionNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MTServiceServer).SetStateTransitionNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mt.v1.MTService/SetStateTransitionNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MTServiceServer).SetStateTransitionNodes(ctx, req.(*SetStateTransitionNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MTService_ServiceDesc is the grpc.ServiceDesc for MTService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MTService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mt.v1.MTService",
	HandlerType: (*MTServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBalance",
			Handler:    _MTService_GetBalance_Handler,
		},
		{
			MethodName: "GetNonce",
			Handler:    _MTService_GetNonce_Handler,
		},
		{
			MethodName: "GetCode",
			Handler:    _MTService_GetCode_Handler,
		},
		{
			MethodName: "GetCodeHash",
			Handler:    _MTService_GetCodeHash_Handler,
		},
		{
			MethodName: "GetStorageAt",
			Handler:    _MTService_GetStorageAt_Handler,
		},
		{
			MethodName: "GetStateTransitionNodes",
			Handler:    _MTService_GetStateTransitionNodes_Handler,
		},
		{
			MethodName: "ReverseHash",
			Handler:    _MTService_ReverseHash_Handler,
		},
		{
			MethodName: "SetBalance",
			Handler:    _MTService_SetBalance_Handler,
		},
		{
			MethodName: "SetNonce",
			Handler:    _MTService_SetNonce_Handler,
		},
		{
			MethodName: "SetCode",
			Handler:    _MTService_SetCode_Handler,
		},
		{
			MethodName: "SetStorageAt",
			Handler:    _MTService_SetStorageAt_Handler,
		},
		{
			MethodName: "SetHashValue",
			Handler:    _MTService_SetHashValue_Handler,
		},
		{
			MethodName: "SetStateTransitionNodes",
			Handler:    _MTService_SetStateTransitionNodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mt.proto",
}
