// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package raft

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

// ReplicaServiceClient is the client API for ReplicaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicaServiceClient interface {
	RequestVote(ctx context.Context, in *ReqVoteParams, opts ...grpc.CallOption) (*ReqVoteResult, error)
	AppendEntries(ctx context.Context, in *AppendEntriesParams, opts ...grpc.CallOption) (*AppendEntriesResult, error)
}

type replicaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicaServiceClient(cc grpc.ClientConnInterface) ReplicaServiceClient {
	return &replicaServiceClient{cc}
}

func (c *replicaServiceClient) RequestVote(ctx context.Context, in *ReqVoteParams, opts ...grpc.CallOption) (*ReqVoteResult, error) {
	out := new(ReqVoteResult)
	err := c.cc.Invoke(ctx, "/raft.ReplicaService/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaServiceClient) AppendEntries(ctx context.Context, in *AppendEntriesParams, opts ...grpc.CallOption) (*AppendEntriesResult, error) {
	out := new(AppendEntriesResult)
	err := c.cc.Invoke(ctx, "/raft.ReplicaService/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicaServiceServer is the server API for ReplicaService service.
// All implementations must embed UnimplementedReplicaServiceServer
// for forward compatibility
type ReplicaServiceServer interface {
	RequestVote(context.Context, *ReqVoteParams) (*ReqVoteResult, error)
	AppendEntries(context.Context, *AppendEntriesParams) (*AppendEntriesResult, error)
	mustEmbedUnimplementedReplicaServiceServer()
}

// UnimplementedReplicaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReplicaServiceServer struct {
}

func (UnimplementedReplicaServiceServer) RequestVote(context.Context, *ReqVoteParams) (*ReqVoteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedReplicaServiceServer) AppendEntries(context.Context, *AppendEntriesParams) (*AppendEntriesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedReplicaServiceServer) mustEmbedUnimplementedReplicaServiceServer() {}

// UnsafeReplicaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicaServiceServer will
// result in compilation errors.
type UnsafeReplicaServiceServer interface {
	mustEmbedUnimplementedReplicaServiceServer()
}

func RegisterReplicaServiceServer(s grpc.ServiceRegistrar, srv ReplicaServiceServer) {
	s.RegisterService(&ReplicaService_ServiceDesc, srv)
}

func _ReplicaService_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqVoteParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServiceServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.ReplicaService/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServiceServer).RequestVote(ctx, req.(*ReqVoteParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplicaService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.ReplicaService/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServiceServer).AppendEntries(ctx, req.(*AppendEntriesParams))
	}
	return interceptor(ctx, in, info, handler)
}

// ReplicaService_ServiceDesc is the grpc.ServiceDesc for ReplicaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReplicaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft.ReplicaService",
	HandlerType: (*ReplicaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _ReplicaService_RequestVote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _ReplicaService_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft/replica_service.proto",
}