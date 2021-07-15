package raft

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (r *Replica) resetElectionTimer() {
	r.stopElectionTimeout()
	n := r.conf.MaxElectionTimeout - r.conf.MinElectionTimeout
	if n <= 0 {
		panic(fmt.Errorf("wrong election timeout, MaxElectionTimeout: %v, "+
			"MinElectionTimeout: %v", r.conf.MaxElectionTimeout, r.conf.MinElectionTimeout))
	}
	rt := rand.Intn(n) + r.conf.MinElectionTimeout
	timeout := time.Duration(rt) * time.Millisecond

	if r.electionTimer == nil {
		r.electionTimer = time.NewTimer(timeout)
	} else {
		r.electionTimer.Reset(timeout)
	}
}

func (r *Replica) resetHeartbeatTimer() {
	r.stopHeartbeatTimeout()
	timeout := time.Duration(r.conf.HeartbeatTimeout) * time.Millisecond

	if r.heartbeatTimer == nil {
		r.heartbeatTimer = time.NewTimer(timeout)
	} else {
		r.heartbeatTimer.Reset(timeout)
	}
}

func (r *Replica) timeoutLoop() {
	for {
		select {
		case <-r.electionTimer.C:
			err := r.handleElectionTimeout()
			if err != nil {
				repLog.Errorf("rep %v handleElectionTimeout err %v", r.conf.RepId, err)
			}
		case <-r.heartbeatTimer.C:
			err := r.handleHeartBeatTimeout()
			if err != nil {
				repLog.Errorf("rep %v handleHeartBeatTimeout err %v", r.conf.RepId, err)
			}
		case <-r.stopCh:
			repLog.Debugf("rep %v stop timer loop", r.conf.RepId)
			return
		}
	}
}

func (r *Replica) handleElectionTimeout() error {
	r.electStart = true
	r.transfer2Candidate()
	return nil
}

func (r *Replica) handleHeartBeatTimeout() error {
	if r.state != Leader {
		r.stopHeartbeatTimeout()
		return nil
	}
	repLog.Debugf("rep %v start send heartbeat, term %v", r.conf.RepId, r.CurrentTerm)
	ctx := context.Background()
	params := &AppendEntriesParams{
		Term:         r.CurrentTerm,
		LeaderId:     r.conf.RepId,
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries:      nil,
		LeaderCommit: 0,
	}
	for _, peer := range r.Cluster.RepPeers {
		if strings.EqualFold(peer.RepId, r.conf.RepId) {
			continue
		}
		aeResult, err := r.rpcClient.SendAppendEntries(ctx, params, peer)
		if err != nil {
			repLog.Debugf("rep %v send AppendEntries to %v, err %v", r.conf.RepId, peer, err)
			continue
		}
		r.handleAppendEntriesResult(aeResult, peer)
	}

	// heartbeat periodically
	r.resetHeartbeatTimer()
	return nil
}

func (r *Replica) stopElectionTimeout() {
	if r.electionTimer == nil {
		return
	}
	stop := r.electionTimer.Stop()
	if !stop {
		select {
		case <-r.electionTimer.C:
		default:
		}
	}
}

func (r *Replica) stopHeartbeatTimeout() {
	if r.heartbeatTimer == nil {
		return
	}
	stop := r.heartbeatTimer.Stop()
	if !stop {
		select {
		case <-r.heartbeatTimer.C:
		default:
		}
	}
}
