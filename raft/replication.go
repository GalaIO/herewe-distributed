package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GalaIO/herewe-distributed/logger"
	"strings"
	"time"
)

// replication states
type RepSate int

func (r RepSate) String() string {
	return repStateStringMap[r]
}

const (
	// who not in current cluster config
	Invalid RepSate = iota
	// do not voting in pre stage when cluster config chg
	NonVoting

	Follower
	Candidate
	Leader
)

var repStateStringMap = map[RepSate]string{
	Invalid:   "Invalid",
	NonVoting: "NonVoting",
	Follower:  "Follower",
	Candidate: "Candidate",
	Leader:    "Leader",
}

var repLog = logger.GetLogger("replication")

type Replication struct {
	// rep info
	conf  RepConfig
	state RepSate

	// persistent data
	RepData
	logs []int

	// mem data
	commitIndex      int64
	lastApplied      int64
	electionTimer    *time.Timer
	heartbeatTimer   *time.Timer
	electStart       bool // when electionTimer timeout, will set true
	recentlyLeaderId string
	stopCh           chan bool

	// service
	rpcClient RpcClient

	// mem data in leader
	// reinitialized after election
	// while clean in state transfer
	nextIndex  map[string]int64
	matchIndex map[string]int64

	// mem data in candidate
	// while clean in state transfer
	votedFrom       map[string]bool
	retryElectCount int
}

func NewRep(storage Storage, rpcClient RpcClient, conf RepConfig) (*Replication, error) {

	repData, err := storage.GetRepData()
	if err != nil && err != ErrNotFound {
		return nil, err
	}
	conRepData := conf.RepData
	majorityNum := (len(conRepData.Cluster.RepPeers) + 1) / 2
	rep := &Replication{
		conf:        conf,
		state:       Follower,
		rpcClient:   rpcClient,
		stopCh:      make(chan bool),
		logs:        nil,
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   nil,
		matchIndex:  nil,
	}
	rep.RepData = conRepData
	rep.RepData.MajorityNum = majorityNum

	if repData != nil {
		rep.RepData = *repData
	}

	if len(rep.Cluster.RepPeers) < 3 {
		return nil, fmt.Errorf("replication cluster must contain more 3 reps, "+
			"now %v", len(rep.Cluster.RepPeers))
	}

	if rep.conf.HeartbeatTimeout >= rep.conf.MinElectionTimeout {
		return nil, fmt.Errorf("replication HeartbeatTimeout need less MinElectionTimeout, "+
			"%v:%v", rep.conf.HeartbeatTimeout, rep.conf.MinElectionTimeout)
	}

	repLog.Debugf("rep %v init as %v", conf.RepId, rep)
	return rep, nil
}

// connect and listen other reps
func (r *Replication) Start() error {
	repLog.Debugf("rep %v has start...", r.conf.RepId)

	// start from follower
	r.initAsFollower()
	r.resetHeartbeatTimer()

	// run some loop
	go r.timeoutLoop()
	return nil
}

// connect and listen other reps
func (r *Replication) Stop() error {
	repLog.Debugf("rep %v will stop...", r.conf.RepId)
	r.stopElectionTimeout()
	r.stopHeartbeatTimeout()
	close(r.stopCh)
	return nil
}

func (r *Replication) initAsLeader() {
	r.state = Leader
	repCount := len(r.Cluster.RepPeers)
	r.nextIndex = make(map[string]int64, repCount-1)
	for _, peer := range r.Cluster.RepPeers {
		if strings.EqualFold(peer.RepId, r.conf.RepId) {
			continue
		}
		r.nextIndex[peer.RepId] = r.commitIndex
	}
	r.matchIndex = make(map[string]int64, repCount-1)
	for _, peer := range r.Cluster.RepPeers {
		if strings.EqualFold(peer.RepId, r.conf.RepId) {
			continue
		}
		r.matchIndex[peer.RepId] = r.commitIndex
	}

	// stop timeout
	r.stopElectionTimeout()
}

func (r *Replication) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

// transfer2Leader only candidate -> leader
func (r *Replication) transfer2Leader() error {

	if r.state != Candidate {
		return errors.New("only candidate could transfer to leader")
	}

	if len(r.votedFrom) < r.MajorityNum {
		return fmt.Errorf("haven not collect enough majority, "+
			"now %v, expect %v", len(r.votedFrom), r.MajorityNum)
	}

	repLog.Debugf("rep %v win to be leader", r.conf.RepId)
	r.initAsLeader()
	// reset send heartbeat
	r.resetHeartbeatTimer()
	return nil
}

func (r *Replication) HandleVote(repId string) error {
	if r.state != Candidate {
		return fmt.Errorf("only candidate could receive vote, now %v", r.state)
	}

	if _, exist := r.votedFrom[repId]; exist {
		repLog.Debugf("rep %v has vote before", repId)
	}
	r.votedFrom[repId] = true
	return nil
}

func (r *Replication) initAsFollower() {
	r.state = Follower
	r.VoteFor = ""

	// reset timeout
	r.resetElectionTimer()
}

func (r *Replication) initAsCandidate() {
	if r.state != Candidate {
		r.retryElectCount = 0
	}
	r.state = Candidate
	r.CurrentTerm++
	r.retryElectCount++
	r.votedFrom = make(map[string]bool, len(r.Cluster.RepPeers))

	// vote for self
	r.VoteFor = r.conf.RepId
	r.votedFrom[r.conf.RepId] = true

	// reset timeout
	r.resetElectionTimer()
}

// transfer2Candidate, only when electionTimeout cloud -> Candidate
// may Follower -> Candidate or Candidate -> Candidate
func (r *Replication) transfer2Candidate() error {

	if !r.electStart {
		return errors.New("there no election timeout to start elect")
	}

	switch r.state {
	case Follower:
		repLog.Debugf("rep %v transfer candidate from follower", r.conf.RepId)
	case Candidate:
		repLog.Debugf("rep %v transfer candidate from candidate", r.conf.RepId)
	default:
		return fmt.Errorf("rep %v wrong state %v to candidate", r.conf.RepId, r.state)
	}

	r.initAsCandidate()
	r.startElect()
	return nil
}

func (r *Replication) transfer2Follower(incomeTerm int64, remoteState RepSate) error {

	if incomeTerm < r.CurrentTerm {
		return fmt.Errorf("low term %v:%v from %v",
			incomeTerm, r.CurrentTerm, remoteState)
	}

	if r.CurrentTerm < incomeTerm {
		repLog.Debugf("rep %v find a higher term, term %v:%v from %v",
			r.conf.RepId, incomeTerm, r.CurrentTerm, remoteState)
	}

	switch r.state {
	case Follower:
		repLog.Debugf("rep %v already a follower", r.conf.RepId)
	case Leader:
		repLog.Debugf("rep %v transfer follower from leader", r.conf.RepId)
		if incomeTerm <= r.CurrentTerm {
			return fmt.Errorf("low term %v:%v from %v",
				incomeTerm, r.CurrentTerm, remoteState)
		}
	case Candidate:
		repLog.Debugf("rep %v transfer follower from candidate", r.conf.RepId)
		if remoteState != Leader && incomeTerm <= r.CurrentTerm {
			return fmt.Errorf("low term %v:%v from %v",
				incomeTerm, r.CurrentTerm, remoteState)
		}
	default:
		return fmt.Errorf("wrong state %v to follower", r.state)
	}

	r.initAsFollower()
	return nil
}

func (r *Replication) VoteToU(repId string, incomeTerm int64) error {
	if incomeTerm < r.CurrentTerm {
		return fmt.Errorf("old term %v now %v", incomeTerm, r.CurrentTerm)
	}

	if incomeTerm == r.CurrentTerm && r.VoteFor != "" && !strings.EqualFold(repId, r.VoteFor) {
		return fmt.Errorf("term %v has voteFor %v, sorry %v", incomeTerm, r.VoteFor, repId)
	}

	repLog.Debugf("rep %v vote to %v, with incomeTerm %v now %v",
		r.conf.RepId, repId, incomeTerm, r.CurrentTerm)
	if err := r.transfer2Follower(incomeTerm, Candidate); err != nil {
		return err
	}
	r.VoteFor = repId
	r.CurrentTerm = incomeTerm
	return nil
}

func (r *Replication) HandleHeartBeat(incomeTerm int64, leaderId string) error {
	r.recentlyLeaderId = leaderId
	return r.transfer2Follower(incomeTerm, Leader)
}

func (r *Replication) startElect() {
	ctx := context.Background()
	voteParams := &ReqVoteParams{
		Term:         r.CurrentTerm,
		CandidateId:  r.conf.RepId,
		LastLogIndex: 0,
		LastLogTerm:  0,
	}
	for _, peer := range r.Cluster.RepPeers {
		if strings.EqualFold(peer.RepId, r.conf.RepId) {
			continue
		}
		voteResult, err := r.rpcClient.SendRequestVote(ctx, voteParams, peer)
		if err != nil {
			repLog.Debugf("rep %v send RequestVote to %v, err %v", r.conf.RepId, peer, err)
			continue
		}
		r.handleRequestVoteResult(voteResult, peer)
	}
}

func (r *Replication) handleRequestVoteResult(voteResult *ReqVoteResult,
	peer RepPeer) {
	// if not voted, try transfer to follower
	if !voteResult.VoteGranted {
		if err := r.transfer2Follower(voteResult.Term, Follower); err != nil {
			repLog.Debugf("rep %v not voted, transfer2Follower err %v", r.conf.RepId, err)
		}
		return
	}
	if err := r.HandleVote(peer.RepId); err != nil {
		repLog.Debugf("rep %v get vote, HandleVote err %v", r.conf.RepId, err)
		return
	}
	if err := r.transfer2Leader(); err != nil {
		repLog.Debugf("rep %v get vote, transfer2Leader err %v", r.conf.RepId, err)
	}
}

func (r *Replication) handleAppendEntriesResult(result *AppendEntriesResult,
	peer RepPeer) {
	// if not success, try transfer to follower
	if !result.Success {
		if err := r.transfer2Follower(result.Term, Follower); err != nil {
			repLog.Debugf("rep %v not voted, transfer2Follower err %v", r.conf.RepId, err)
		}
		return
	}

	// other just retry
}
