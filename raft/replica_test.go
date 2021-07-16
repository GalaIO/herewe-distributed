package raft

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	mockRepIds = []string{
		"alice",
		"bob",
		"henk",
	}
)

var mockCluster = ClusterConfig{
	RepPeers: []RepPeer{
		{
			RepId: mockRepIds[0],
			Addr:  "127.0.0.1:6701",
		},
		{
			RepId: mockRepIds[1],
			Addr:  "127.0.0.1:6702",
		},
		{
			RepId: mockRepIds[2],
			Addr:  "127.0.0.1:6703",
		},
	},
}

var mockRepConfig = RepConfig{
	RepId:              "El0+6hRV9eKQfkJrMfViDP2qPzA=",
	Addr:               "127.0.0.1:6701",
	MinElectionTimeout: 300,
	MaxElectionTimeout: 500,
	HeartbeatTimeout:   200,
	RepData: RepData{
		Cluster:      mockCluster,
		CurrentTerm:  0,
		VoteFor:      "",
		MajorityNum:  0,
		LastLogIndex: 0,
		LastLogTerm:  0,
		CommitIndex:  0,
		LastApplied:  0,
	},
}

func TestReplica_transfer2Leader(t *testing.T) {

	rep := initMockRep(t)

	rep.initAsCandidate()
	for _, repId := range mockRepIds {
		err := rep.OnReceiveVote(repId)
		assert.NoError(t, err)
	}
	err := rep.transfer2Leader()
	assert.NoError(t, err)
	assert.Equal(t, Leader, rep.state)

}

func TestReplica_transfer2LeaderFromFollower(t *testing.T) {

	rep := initMockRep(t)

	for _, repId := range mockRepIds {
		err := rep.OnReceiveVote(repId)
		assert.Error(t, err)
	}
	err := rep.transfer2Leader()
	assert.Error(t, err)

}

func TestReplica_transfer2CandidateFromFollower(t *testing.T) {

	rep := initMockRep(t)
	oldTerm := rep.CurrentTerm

	rep.electStart = true
	err := rep.transfer2Candidate()
	assert.NoError(t, err)
	assert.Equal(t, Candidate, rep.state)
	assert.Equal(t, oldTerm+1, rep.CurrentTerm)
}

func TestReplica_transfer2CandidateFromCandidate(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsCandidate()
	oldTerm := rep.CurrentTerm

	rep.electStart = true
	err := rep.transfer2Candidate()
	assert.NoError(t, err)
	assert.Equal(t, Candidate, rep.state)
	assert.Equal(t, oldTerm+1, rep.CurrentTerm)
}

func TestReplica_transfer2CandidateFromLeader(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsLeader()

	rep.electStart = true
	err := rep.transfer2Candidate()
	assert.Error(t, err)
}

func TestReplica_electionTimeout(t *testing.T) {

	rep := initMockRep(t)

	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, true, rep.electStart)
	assert.Equal(t, Candidate, rep.state)
	count := rep.retryElectCount
	t.Log("retryElectCount:", rep.retryElectCount)

	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, true, rep.electStart)
	assert.Equal(t, Candidate, rep.state)
	assert.True(t, rep.retryElectCount > count)
	t.Log("retryElectCount:", rep.retryElectCount)
}

func TestReplica_transfer2FollowerFromLeader(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsLeader()

	var incomeTerm int64 = 10
	err := rep.transfer2Follower(incomeTerm, Follower)
	assert.NoError(t, err)
	assert.Equal(t, Follower, rep.state)
}

func TestReplica_transfer2FollowerFromCandidate(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsCandidate()

	var incomeTerm int64 = 10
	err := rep.transfer2Follower(incomeTerm, Follower)
	assert.NoError(t, err)
	assert.Equal(t, Follower, rep.state)
}

func TestReplica_transfer2FollowerFromCandidateInSameTerm(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsCandidate()

	incomeTerm := rep.CurrentTerm
	err := rep.transfer2Follower(incomeTerm, Leader)
	assert.NoError(t, err)
	assert.Equal(t, Follower, rep.state)
}

func TestReplica_transfer2FollowerFromCandidateInSameTerm2(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsCandidate()

	incomeTerm := rep.CurrentTerm
	err := rep.transfer2Follower(incomeTerm, Candidate)
	assert.Error(t, err)
}

func TestReplica_VoteToUInFollower(t *testing.T) {

	rep := initMockRep(t)

	incomeTerm := rep.CurrentTerm
	err := rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.NoError(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)
	assert.Equal(t, incomeTerm, rep.CurrentTerm)

	err = rep.VoteToU(mockRepIds[1], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)

	// wait timeout chg candidate
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, rep.VoteFor, rep.conf.RepId)
	assert.Equal(t, Candidate, rep.state)

	// higher term
	incomeTerm = rep.CurrentTerm + 1
	err = rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.NoError(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)
	assert.Equal(t, incomeTerm, rep.CurrentTerm)

	err = rep.VoteToU(mockRepIds[1], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)

}

func TestReplica_VoteToUInCandidate(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsCandidate()

	incomeTerm := rep.CurrentTerm
	err := rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, rep.conf.RepId, rep.VoteFor)
	assert.Equal(t, Candidate, rep.state)

	// wait timeout chg candidate
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, rep.VoteFor, rep.conf.RepId)
	assert.Equal(t, Candidate, rep.state)

	// higher term
	incomeTerm = rep.CurrentTerm + 1
	err = rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.NoError(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)
	assert.Equal(t, incomeTerm, rep.CurrentTerm)

	err = rep.VoteToU(mockRepIds[1], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)

}

func TestReplica_VoteToUInLeader(t *testing.T) {

	rep := initMockRep(t)
	rep.initAsLeader()

	incomeTerm := rep.CurrentTerm
	err := rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, Leader, rep.state)

	// higher term
	incomeTerm = rep.CurrentTerm + 1
	err = rep.VoteToU(mockRepIds[0], incomeTerm)
	assert.NoError(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)
	assert.Equal(t, incomeTerm, rep.CurrentTerm)

	err = rep.VoteToU(mockRepIds[1], incomeTerm)
	assert.Error(t, err)
	assert.Equal(t, rep.VoteFor, mockRepIds[0])
	assert.Equal(t, Follower, rep.state)
}

func TestReplica_HandelHeartBeat(t *testing.T) {

	rep := initMockRep(t)

	old := rep.CurrentTerm
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		rep.AppendEntries(&AppendEntriesParams{
			Term:         old,
			LeaderId:     mockRepIds[0],
			PrevLogIndex: 0,
			PrevLogTerm:  0,
			Entries:      nil,
			LeaderCommit: 0,
		})
		assert.Equal(t, Follower, rep.state)
		assert.Equal(t, old, rep.CurrentTerm)
	}
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, Candidate, rep.state)
	assert.True(t, old < rep.CurrentTerm)
}

func TestReplica_RequestLogEntry(t *testing.T) {

	ctrl := gomock.NewController(t)
	rpcClient := NewMockRpcClient(ctrl)
	rpcClient.EXPECT().SendRequestVote(gomock.Any(), gomock.Any(),
		gomock.Any()).Return(nil, ErrNotFound).AnyTimes()

	rep := initRpcClientMockRep(t, rpcClient)
	rep.initAsLeader()

	result := &AppendEntriesResult{
		Term:                 rep.CurrentTerm,
		Success:              true,
		LastLogTerm:          rep.LastLogTerm,
		LastLogIndex:         rep.LastLogIndex + 1,
		FirstIndexInLastTerm: 0,
		CommitIndex:          0,
	}
	rpcClient.EXPECT().SendAppendEntries(gomock.Any(), gomock.Any(),
		gomock.Any()).Return(result, nil).AnyTimes()
	err := rep.requestLogEntry([]byte("hello"))
	assert.NoError(t, err)
	for _, index := range rep.nextIndex {
		assert.Equal(t, rep.LastLogIndex+1, index)
	}
	// TODO
	//for _, index := range rep.matchIndex {
	//	assert.Equal(t, rep.LastLogIndex+1, index)
	//}
}

func TestReplica_AppendEntries(t *testing.T) {
	rep := initMockRep(t)
	samples := []struct {
		input  *AppendEntriesParams
		state  RepData
		output *AppendEntriesResult
	}{
		{
			input: &AppendEntriesParams{
				Term:         0,
				LeaderId:     mockRepIds[0],
				PrevLogIndex: 0,
				PrevLogTerm:  0,
				Entries:      []*LogEntry{dummyEntry(1, 0), dummyEntry(2, 0)},
				LeaderCommit: 0,
			},
			state: RepData{
				Cluster:          rep.Cluster,
				VoteFor:          "",
				MajorityNum:      rep.MajorityNum,
				CommitIndex:      0,
				LastApplied:      0,
				CurrentTerm:      0,
				LastLogIndex:     0,
				LastLogTerm:      0,
				FirstIndexInTerm: 0,
			},
			output: &AppendEntriesResult{
				Term:                 rep.CurrentTerm,
				Success:              true,
				LastLogTerm:          0,
				LastLogIndex:         2,
				FirstIndexInLastTerm: 0,
				CommitIndex:          0,
			},
		},
		{
			input: &AppendEntriesParams{
				Term:         8,
				LeaderId:     mockRepIds[0],
				PrevLogIndex: 10,
				PrevLogTerm:  6,
				Entries:      []*LogEntry{dummyEntry(1, 0), dummyEntry(2, 0)},
				LeaderCommit: 0,
			},
			state: RepData{
				Cluster:          rep.Cluster,
				VoteFor:          "",
				MajorityNum:      rep.MajorityNum,
				CommitIndex:      0,
				LastApplied:      0,
				CurrentTerm:      3,
				LastLogIndex:     11,
				LastLogTerm:      3,
				FirstIndexInTerm: 7,
			},
			output: &AppendEntriesResult{
				Term:                 8,
				Success:              false,
				LastLogTerm:          3,
				LastLogIndex:         11,
				FirstIndexInLastTerm: 7,
				CommitIndex:          0,
			},
		},
	}
	for _, sample := range samples {
		rep.RepData = sample.state
		tmpResult, err := rep.AppendEntries(sample.input)
		assert.NoError(t, err)
		assert.Equal(t, sample.output, tmpResult)
	}
}

//TODO add more safety tests
// election commitId tests

func initMockRep(t *testing.T) *Replica {
	ctrl := gomock.NewController(t)
	storage := NewMockStorage(ctrl)
	storage.EXPECT().GetRepData().Return(nil, ErrNotFound).AnyTimes()
	storage.EXPECT().SaveRepData(gomock.Any()).Return(nil).AnyTimes()

	rpcClient := NewMockRpcClient(ctrl)
	rpcClient.EXPECT().SendAppendEntries(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrNotFound).AnyTimes()
	rpcClient.EXPECT().SendRequestVote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrNotFound).AnyTimes()

	rep, _ := NewRep(storage, rpcClient, mockRepConfig)

	err := rep.Start()
	assert.NoError(t, err)
	return rep
}

func initRpcClientMockRep(t *testing.T, client RpcClient) *Replica {
	ctrl := gomock.NewController(t)
	storage := NewMockStorage(ctrl)
	storage.EXPECT().GetRepData().Return(nil, ErrNotFound).AnyTimes()
	storage.EXPECT().SaveRepData(gomock.Any()).Return(nil).AnyTimes()
	rep, _ := NewRep(storage, client, mockRepConfig)

	err := rep.Start()
	assert.NoError(t, err)
	return rep
}
