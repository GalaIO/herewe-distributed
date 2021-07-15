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

// ReplicationServiceClient is the client API for ReplicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationServiceClient interface {
	RequestVote(ctx context.Context, in *ReqVoteParams, opts ...grpc.CallOption) (*ReqVoteResult, error)
	AppendEntries(ctx context.Context, in *AppendEntriesParams, opts ...grpc.CallOption) (*AppendEntriesResult, error)
}

type replicationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationServiceClient(cc grpc.ClientConnInterface) ReplicationServiceClient {
	return &replicationServiceClient{cc}
}

func (c *replicationServiceClient) RequestVote(ctx context.Context, in *ReqVoteParams, opts ...grpc.CallOption) (*ReqVoteResult, error) {
	out := new(ReqVoteResult)
	err := c.cc.Invoke(ctx, "/raft.ReplicationService/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationServiceClient) AppendEntries(ctx context.Context, in *AppendEntriesParams, opts ...grpc.CallOption) (*AppendEntriesResult, error) {
	out := new(AppendEntriesResult)
	err := c.cc.Invoke(ctx, "/raft.ReplicationService/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServiceServer is the server API for ReplicationService service.
// All implementations must embed UnimplementedReplicationServiceServer
// for forward compatibility
type ReplicationServiceServer interface {
	RequestVote(context.Context, *ReqVoteParams) (*ReqVoteResult, error)
	AppendEntries(context.Context, *AppendEntriesParams) (*AppendEntriesResult, error)
	mustEmbedUnimplementedReplicationServiceServer()
}

// UnimplementedReplicationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServiceServer struct {
}

func (UnimplementedReplicationServiceServer) RequestVote(context.Context, *ReqVoteParams) (*ReqVoteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedReplicationServiceServer) AppendEntries(context.Context, *AppendEntriesParams) (*AppendEntriesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedReplicationServiceServer) mustEmbedUnimplementedReplicationServiceServer() {}

// UnsafeReplicationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServiceServer will
// result in compilation errors.
type UnsafeReplicationServiceServer interface {
	mustEmbedUnimplementedReplicationServiceServer()
}

func RegisterReplicationServiceServer(s grpc.ServiceRegistrar, srv ReplicationServiceServer) {
	s.RegisterService(&ReplicationService_ServiceDesc, srv)
}

func _ReplicationService_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqVoteParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServiceServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.ReplicationService/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServiceServer).RequestVote(ctx, req.(*ReqVoteParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplicationService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.ReplicationService/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServiceServer).AppendEntries(ctx, req.(*AppendEntriesParams))
	}
	return interceptor(ctx, in, info, handler)
}

// ReplicationService_ServiceDesc is the grpc.ServiceDesc for ReplicationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReplicationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft.ReplicationService",
	HandlerType: (*ReplicationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _ReplicationService_RequestVote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _ReplicationService_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft/replication_service.proto",
}