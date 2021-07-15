package raft

import (
	context "context"
	"fmt"
	"github.com/GalaIO/herewe-distributed/logger"
	"google.golang.org/grpc"
	"net"
	"time"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  replica_service.proto
//go:generate mockgen -package=raft -destination=rpc_mock.go . RpcClient,RpcServer

var rpcLog = logger.GetLogger("rpc")

type RpcServer interface {
	ReplicaServiceServer
	Start() error
	Stop() error
}

type RpcServerImpl struct {
	UnimplementedReplicaServiceServer
	rep    *Replica
	server *grpc.Server
}

func NewRepServer(rep *Replica) RpcServer {
	r := &RpcServerImpl{
		rep: rep,
	}
	s := grpc.NewServer()
	RegisterReplicaServiceServer(s, r)
	r.server = s
	return r
}

func (r *RpcServerImpl) Start() error {
	conn, err := net.Listen("tcp", r.rep.conf.Addr)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}
	rpcLog.Infof("rpc server start at %v", r.rep.conf.Addr)
	if err := r.server.Serve(conn); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
	return nil
}

func (r *RpcServerImpl) Stop() error {
	rpcLog.Infof("rpc server stop...")
	r.server.Stop()
	return nil
}

func (r *RpcServerImpl) RequestVote(ctx context.Context, params *ReqVoteParams) (*ReqVoteResult, error) {
	return r.rep.RequestVote(params)
}

func (r *RpcServerImpl) AppendEntries(ctx context.Context, params *AppendEntriesParams) (*AppendEntriesResult, error) {
	return r.rep.AppendEntries(params)
}

type RpcClient interface {
	SendRequestVote(ctx context.Context, params *ReqVoteParams, peer RepPeer) (*ReqVoteResult, error)
	SendAppendEntries(ctx context.Context, params *AppendEntriesParams, peer RepPeer) (*AppendEntriesResult, error)
}

type RpcClientImpl struct {
}

func NewRpcClient() RpcClient {
	return &RpcClientImpl{}
}

func (r *RpcClientImpl) SendRequestVote(ctx context.Context, params *ReqVoteParams, peer RepPeer) (*ReqVoteResult, error) {
	conn, err := grpc.Dial(peer.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := NewReplicaServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	result, err := c.RequestVote(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RpcClientImpl) SendAppendEntries(ctx context.Context, params *AppendEntriesParams, peer RepPeer) (*AppendEntriesResult, error) {
	conn, err := grpc.Dial(peer.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := NewReplicaServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	result, err := c.AppendEntries(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
