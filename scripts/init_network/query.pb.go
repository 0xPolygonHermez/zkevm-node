//*
// Bridge service.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: query.proto

package main

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Deposit message
type Deposit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrigNet    uint32 `protobuf:"varint,1,opt,name=orig_net,json=origNet,proto3" json:"orig_net,omitempty"`
	TokenAddr  string `protobuf:"bytes,2,opt,name=token_addr,json=tokenAddr,proto3" json:"token_addr,omitempty"`
	Amount     string `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	DestNet    uint32 `protobuf:"varint,4,opt,name=dest_net,json=destNet,proto3" json:"dest_net,omitempty"`
	DestAddr   string `protobuf:"bytes,5,opt,name=dest_addr,json=destAddr,proto3" json:"dest_addr,omitempty"`
	BlockNum   uint64 `protobuf:"varint,6,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
	DepositCnt uint64 `protobuf:"varint,7,opt,name=deposit_cnt,json=depositCnt,proto3" json:"deposit_cnt,omitempty"`
}

func (x *Deposit) Reset() {
	*x = Deposit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deposit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deposit) ProtoMessage() {}

func (x *Deposit) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deposit.ProtoReflect.Descriptor instead.
func (*Deposit) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{0}
}

func (x *Deposit) GetOrigNet() uint32 {
	if x != nil {
		return x.OrigNet
	}
	return 0
}

func (x *Deposit) GetTokenAddr() string {
	if x != nil {
		return x.TokenAddr
	}
	return ""
}

func (x *Deposit) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Deposit) GetDestNet() uint32 {
	if x != nil {
		return x.DestNet
	}
	return 0
}

func (x *Deposit) GetDestAddr() string {
	if x != nil {
		return x.DestAddr
	}
	return ""
}

func (x *Deposit) GetBlockNum() uint64 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

func (x *Deposit) GetDepositCnt() uint64 {
	if x != nil {
		return x.DepositCnt
	}
	return 0
}

// Claim message
type Claim struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index     uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	OrigNet   uint32 `protobuf:"varint,2,opt,name=orig_net,json=origNet,proto3" json:"orig_net,omitempty"`
	TokenAddr string `protobuf:"bytes,3,opt,name=token_addr,json=tokenAddr,proto3" json:"token_addr,omitempty"`
	Amount    string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	DestNet   uint32 `protobuf:"varint,5,opt,name=dest_net,json=destNet,proto3" json:"dest_net,omitempty"`
	DestAddr  string `protobuf:"bytes,6,opt,name=dest_addr,json=destAddr,proto3" json:"dest_addr,omitempty"`
	BlockNum  uint64 `protobuf:"varint,7,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
}

func (x *Claim) Reset() {
	*x = Claim{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Claim) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Claim) ProtoMessage() {}

func (x *Claim) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Claim.ProtoReflect.Descriptor instead.
func (*Claim) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{1}
}

func (x *Claim) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Claim) GetOrigNet() uint32 {
	if x != nil {
		return x.OrigNet
	}
	return 0
}

func (x *Claim) GetTokenAddr() string {
	if x != nil {
		return x.TokenAddr
	}
	return ""
}

func (x *Claim) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *Claim) GetDestNet() uint32 {
	if x != nil {
		return x.DestNet
	}
	return 0
}

func (x *Claim) GetDestAddr() string {
	if x != nil {
		return x.DestAddr
	}
	return ""
}

func (x *Claim) GetBlockNum() uint64 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

// Merkle Proof message
type Proof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MerkleProof    []string `protobuf:"bytes,1,rep,name=merkle_proof,json=merkleProof,proto3" json:"merkle_proof,omitempty"`
	ExitRootNum    uint64   `protobuf:"varint,2,opt,name=exit_root_num,json=exitRootNum,proto3" json:"exit_root_num,omitempty"`
	MainExitRoot   string   `protobuf:"bytes,3,opt,name=main_exit_root,json=mainExitRoot,proto3" json:"main_exit_root,omitempty"`
	RollupExitRoot string   `protobuf:"bytes,4,opt,name=rollup_exit_root,json=rollupExitRoot,proto3" json:"rollup_exit_root,omitempty"`
}

func (x *Proof) Reset() {
	*x = Proof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proof) ProtoMessage() {}

func (x *Proof) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Proof.ProtoReflect.Descriptor instead.
func (*Proof) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{2}
}

func (x *Proof) GetMerkleProof() []string {
	if x != nil {
		return x.MerkleProof
	}
	return nil
}

func (x *Proof) GetExitRootNum() uint64 {
	if x != nil {
		return x.ExitRootNum
	}
	return 0
}

func (x *Proof) GetMainExitRoot() string {
	if x != nil {
		return x.MainExitRoot
	}
	return ""
}

func (x *Proof) GetRollupExitRoot() string {
	if x != nil {
		return x.RollupExitRoot
	}
	return ""
}

type CheckAPIRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckAPIRequest) Reset() {
	*x = CheckAPIRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckAPIRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckAPIRequest) ProtoMessage() {}

func (x *CheckAPIRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckAPIRequest.ProtoReflect.Descriptor instead.
func (*CheckAPIRequest) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{3}
}

type GetBridgesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EtherAddr string `protobuf:"bytes,1,opt,name=ether_addr,json=etherAddr,proto3" json:"ether_addr,omitempty"`
	Offset    uint64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetBridgesRequest) Reset() {
	*x = GetBridgesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBridgesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBridgesRequest) ProtoMessage() {}

func (x *GetBridgesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBridgesRequest.ProtoReflect.Descriptor instead.
func (*GetBridgesRequest) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{4}
}

func (x *GetBridgesRequest) GetEtherAddr() string {
	if x != nil {
		return x.EtherAddr
	}
	return ""
}

func (x *GetBridgesRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetProofRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrigNet    uint32 `protobuf:"varint,1,opt,name=orig_net,json=origNet,proto3" json:"orig_net,omitempty"`
	DepositCnt uint64 `protobuf:"varint,2,opt,name=deposit_cnt,json=depositCnt,proto3" json:"deposit_cnt,omitempty"`
}

func (x *GetProofRequest) Reset() {
	*x = GetProofRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProofRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProofRequest) ProtoMessage() {}

func (x *GetProofRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProofRequest.ProtoReflect.Descriptor instead.
func (*GetProofRequest) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{5}
}

func (x *GetProofRequest) GetOrigNet() uint32 {
	if x != nil {
		return x.OrigNet
	}
	return 0
}

func (x *GetProofRequest) GetDepositCnt() uint64 {
	if x != nil {
		return x.DepositCnt
	}
	return 0
}

type GetClaimStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrigNet    uint32 `protobuf:"varint,1,opt,name=orig_net,json=origNet,proto3" json:"orig_net,omitempty"`
	DepositCnt uint64 `protobuf:"varint,2,opt,name=deposit_cnt,json=depositCnt,proto3" json:"deposit_cnt,omitempty"`
}

func (x *GetClaimStatusRequest) Reset() {
	*x = GetClaimStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimStatusRequest) ProtoMessage() {}

func (x *GetClaimStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimStatusRequest.ProtoReflect.Descriptor instead.
func (*GetClaimStatusRequest) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{6}
}

func (x *GetClaimStatusRequest) GetOrigNet() uint32 {
	if x != nil {
		return x.OrigNet
	}
	return 0
}

func (x *GetClaimStatusRequest) GetDepositCnt() uint64 {
	if x != nil {
		return x.DepositCnt
	}
	return 0
}

type GetClaimsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EtherAddr string `protobuf:"bytes,1,opt,name=ether_addr,json=etherAddr,proto3" json:"ether_addr,omitempty"`
	Offset    uint64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetClaimsRequest) Reset() {
	*x = GetClaimsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsRequest) ProtoMessage() {}

func (x *GetClaimsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsRequest.ProtoReflect.Descriptor instead.
func (*GetClaimsRequest) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{7}
}

func (x *GetClaimsRequest) GetEtherAddr() string {
	if x != nil {
		return x.EtherAddr
	}
	return ""
}

func (x *GetClaimsRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type CheckAPIResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Api string `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
}

func (x *CheckAPIResponse) Reset() {
	*x = CheckAPIResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckAPIResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckAPIResponse) ProtoMessage() {}

func (x *CheckAPIResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckAPIResponse.ProtoReflect.Descriptor instead.
func (*CheckAPIResponse) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{8}
}

func (x *CheckAPIResponse) GetApi() string {
	if x != nil {
		return x.Api
	}
	return ""
}

type GetBridgesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deposits []*Deposit `protobuf:"bytes,1,rep,name=deposits,proto3" json:"deposits,omitempty"`
}

func (x *GetBridgesResponse) Reset() {
	*x = GetBridgesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBridgesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBridgesResponse) ProtoMessage() {}

func (x *GetBridgesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBridgesResponse.ProtoReflect.Descriptor instead.
func (*GetBridgesResponse) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{9}
}

func (x *GetBridgesResponse) GetDeposits() []*Deposit {
	if x != nil {
		return x.Deposits
	}
	return nil
}

type GetProofResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proof *Proof `protobuf:"bytes,1,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *GetProofResponse) Reset() {
	*x = GetProofResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProofResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProofResponse) ProtoMessage() {}

func (x *GetProofResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProofResponse.ProtoReflect.Descriptor instead.
func (*GetProofResponse) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{10}
}

func (x *GetProofResponse) GetProof() *Proof {
	if x != nil {
		return x.Proof
	}
	return nil
}

type GetClaimStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ready bool `protobuf:"varint,1,opt,name=ready,proto3" json:"ready,omitempty"`
}

func (x *GetClaimStatusResponse) Reset() {
	*x = GetClaimStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimStatusResponse) ProtoMessage() {}

func (x *GetClaimStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimStatusResponse.ProtoReflect.Descriptor instead.
func (*GetClaimStatusResponse) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{11}
}

func (x *GetClaimStatusResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

type GetClaimsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Claims []*Claim `protobuf:"bytes,1,rep,name=claims,proto3" json:"claims,omitempty"`
}

func (x *GetClaimsResponse) Reset() {
	*x = GetClaimsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_query_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClaimsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClaimsResponse) ProtoMessage() {}

func (x *GetClaimsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClaimsResponse.ProtoReflect.Descriptor instead.
func (*GetClaimsResponse) Descriptor() ([]byte, []int) {
	return file_query_proto_rawDescGZIP(), []int{12}
}

func (x *GetClaimsResponse) GetClaims() []*Claim {
	if x != nil {
		return x.Claims
	}
	return nil
}

var File_query_proto protoreflect.FileDescriptor

var file_query_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x68,
	0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd1, 0x01,
	0x0a, 0x07, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x69,
	0x67, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x69,
	0x67, 0x4e, 0x65, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x41,
	0x64, 0x64, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x64,
	0x65, 0x73, 0x74, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x64,
	0x65, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x73, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d,
	0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x5f, 0x63, 0x6e, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6e,
	0x74, 0x22, 0xc4, 0x01, 0x0a, 0x05, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x69, 0x67, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x69, 0x67, 0x4e, 0x65, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x6e, 0x65, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x64, 0x65, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x64, 0x65, 0x73, 0x74, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x22, 0x9e, 0x01, 0x0a, 0x05, 0x50, 0x72, 0x6f,
	0x6f, 0x66, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x72, 0x6f,
	0x6f, 0x66, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65,
	0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x22, 0x0a, 0x0d, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x72, 0x6f,
	0x6f, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x65, 0x78,
	0x69, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x61, 0x69,
	0x6e, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x6d, 0x61, 0x69, 0x6e, 0x45, 0x78, 0x69, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12,
	0x28, 0x0a, 0x10, 0x72, 0x6f, 0x6c, 0x6c, 0x75, 0x70, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x72,
	0x6f, 0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x6f, 0x6c, 0x6c, 0x75,
	0x70, 0x45, 0x78, 0x69, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x41, 0x50, 0x49, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x74, 0x68, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x4d, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50,
	0x72, 0x6f, 0x6f, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x69, 0x67, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f,
	0x72, 0x69, 0x67, 0x4e, 0x65, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x5f, 0x63, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x43, 0x6e, 0x74, 0x22, 0x53, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x6c,
	0x61, 0x69, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x69, 0x67, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x69, 0x67, 0x4e, 0x65, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x64,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x5f, 0x63, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x43, 0x6e, 0x74, 0x22, 0x49, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x74, 0x68, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x24, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x41, 0x50, 0x49, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61,
	0x70, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x69, 0x22, 0x4b, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62,
	0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x52, 0x08, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x73, 0x22, 0x41, 0x0a, 0x10, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d,
	0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x2e, 0x0a,
	0x16, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x22, 0x44, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x52, 0x06, 0x63, 0x6c, 0x61,
	0x69, 0x6d, 0x73, 0x32, 0xc3, 0x04, 0x0a, 0x0d, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x08, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x50,
	0x49, 0x12, 0x21, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x50, 0x49, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72,
	0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x50, 0x49,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x06,
	0x12, 0x04, 0x2f, 0x61, 0x70, 0x69, 0x12, 0x76, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x72, 0x69,
	0x64, 0x67, 0x65, 0x73, 0x12, 0x23, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72,
	0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x72, 0x69, 0x64, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x68, 0x65, 0x72, 0x6d,
	0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65,
	0x73, 0x2f, 0x7b, 0x65, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x7d, 0x12, 0x69,
	0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x21, 0x2e, 0x68, 0x65, 0x72,
	0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x6d, 0x65, 0x72, 0x6b,
	0x6c, 0x65, 0x2d, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x73, 0x12, 0x7a, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x43, 0x6c, 0x61, 0x69, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x27, 0x2e, 0x68, 0x65,
	0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72,
	0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x2d, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x72, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69,
	0x6d, 0x73, 0x12, 0x22, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e, 0x62, 0x72, 0x69, 0x64,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2e,
	0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x61,
	0x69, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x2f, 0x7b, 0x65, 0x74,
	0x68, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x7d, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x7a, 0x2d, 0x62, 0x72, 0x69,
	0x64, 0x67, 0x65, 0x2f, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_query_proto_rawDescOnce sync.Once
	file_query_proto_rawDescData = file_query_proto_rawDesc
)

func file_query_proto_rawDescGZIP() []byte {
	file_query_proto_rawDescOnce.Do(func() {
		file_query_proto_rawDescData = protoimpl.X.CompressGZIP(file_query_proto_rawDescData)
	})
	return file_query_proto_rawDescData
}

var file_query_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_query_proto_goTypes = []interface{}{
	(*Deposit)(nil),                // 0: hermez.bridge.v1.Deposit
	(*Claim)(nil),                  // 1: hermez.bridge.v1.Claim
	(*Proof)(nil),                  // 2: hermez.bridge.v1.Proof
	(*CheckAPIRequest)(nil),        // 3: hermez.bridge.v1.CheckAPIRequest
	(*GetBridgesRequest)(nil),      // 4: hermez.bridge.v1.GetBridgesRequest
	(*GetProofRequest)(nil),        // 5: hermez.bridge.v1.GetProofRequest
	(*GetClaimStatusRequest)(nil),  // 6: hermez.bridge.v1.GetClaimStatusRequest
	(*GetClaimsRequest)(nil),       // 7: hermez.bridge.v1.GetClaimsRequest
	(*CheckAPIResponse)(nil),       // 8: hermez.bridge.v1.CheckAPIResponse
	(*GetBridgesResponse)(nil),     // 9: hermez.bridge.v1.GetBridgesResponse
	(*GetProofResponse)(nil),       // 10: hermez.bridge.v1.GetProofResponse
	(*GetClaimStatusResponse)(nil), // 11: hermez.bridge.v1.GetClaimStatusResponse
	(*GetClaimsResponse)(nil),      // 12: hermez.bridge.v1.GetClaimsResponse
}
var file_query_proto_depIdxs = []int32{
	0,  // 0: hermez.bridge.v1.GetBridgesResponse.deposits:type_name -> hermez.bridge.v1.Deposit
	2,  // 1: hermez.bridge.v1.GetProofResponse.proof:type_name -> hermez.bridge.v1.Proof
	1,  // 2: hermez.bridge.v1.GetClaimsResponse.claims:type_name -> hermez.bridge.v1.Claim
	3,  // 3: hermez.bridge.v1.BridgeService.CheckAPI:input_type -> hermez.bridge.v1.CheckAPIRequest
	4,  // 4: hermez.bridge.v1.BridgeService.GetBridges:input_type -> hermez.bridge.v1.GetBridgesRequest
	5,  // 5: hermez.bridge.v1.BridgeService.GetProof:input_type -> hermez.bridge.v1.GetProofRequest
	6,  // 6: hermez.bridge.v1.BridgeService.GetClaimStatus:input_type -> hermez.bridge.v1.GetClaimStatusRequest
	7,  // 7: hermez.bridge.v1.BridgeService.GetClaims:input_type -> hermez.bridge.v1.GetClaimsRequest
	8,  // 8: hermez.bridge.v1.BridgeService.CheckAPI:output_type -> hermez.bridge.v1.CheckAPIResponse
	9,  // 9: hermez.bridge.v1.BridgeService.GetBridges:output_type -> hermez.bridge.v1.GetBridgesResponse
	10, // 10: hermez.bridge.v1.BridgeService.GetProof:output_type -> hermez.bridge.v1.GetProofResponse
	11, // 11: hermez.bridge.v1.BridgeService.GetClaimStatus:output_type -> hermez.bridge.v1.GetClaimStatusResponse
	12, // 12: hermez.bridge.v1.BridgeService.GetClaims:output_type -> hermez.bridge.v1.GetClaimsResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_query_proto_init() }
func file_query_proto_init() {
	if File_query_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_query_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deposit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Claim); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Proof); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckAPIRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBridgesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProofRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimStatusRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckAPIResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBridgesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProofResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimStatusResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_query_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClaimsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_query_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_query_proto_goTypes,
		DependencyIndexes: file_query_proto_depIdxs,
		MessageInfos:      file_query_proto_msgTypes,
	}.Build()
	File_query_proto = out.File
	file_query_proto_rawDesc = nil
	file_query_proto_goTypes = nil
	file_query_proto_depIdxs = nil
}
