package raft

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	mockRepIds = []string{
		"El0+6hRV9eKQfkJrMfViDP2qPzA=",
		"9d25d/FlqJQrBSiwvE3hUnupM3A=",
		"YJpacOKiEjzyHRBm3CWoPGXI+Jc=",
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
	MinElectionTimeout: 150,
	MaxElectionTimeout: 300,
	HeartbeatTimeout:   50,
	RepData: RepData{
		Cluster:     mockCluster,
		CurrentTerm: 0,
		VoteFor:     "",
	},
}

func TestReplica_transfer2Leader(t *testing.T) {

	rep := initMockRep(t)

	rep.initAsCandidate()
	for _, repId := range mockRepIds {
		err := rep.HandleVote(repId)
		assert.NoError(t, err)
	}
	err := rep.transfer2Leader()
	assert.NoError(t, err)
	assert.Equal(t, Leader, rep.state)

}

func TestReplica_transfer2LeaderFromFollower(t *testing.T) {

	rep := initMockRep(t)

	for _, repId := range mockRepIds {
		err := rep.HandleVote(repId)
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
		rep.HandleHeartBeat(old, mockRepIds[0])
		assert.Equal(t, Follower, rep.state)
		assert.Equal(t, old, rep.CurrentTerm)
	}
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, Candidate, rep.state)
	assert.True(t, old < rep.CurrentTerm)
}

func initMockRep(t *testing.T) *Replica {
	storage := NewMockStorage()
	rpcClient := NewMockRpcClient()
	rep, _ := NewRep(storage, rpcClient, mockRepConfig)

	err := rep.Start()
	assert.NoError(t, err)
	return rep
}
