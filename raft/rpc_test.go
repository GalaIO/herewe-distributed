package raft

import (
	"context"
	"errors"
)

type MockRpcClient struct {
}

func NewMockRpcClient() RpcClient {
	return &MockRpcClient{}
}

func (m MockRpcClient) SendRequestVote(ctx context.Context, params *ReqVoteParams, peer RepPeer) (*ReqVoteResult, error) {
	return nil, errors.New("send error in mock")
}

func (m MockRpcClient) SendAppendEntries(ctx context.Context, params *AppendEntriesParams, peer RepPeer) (*AppendEntriesResult, error) {
	return nil, errors.New("send error in mock")
}
