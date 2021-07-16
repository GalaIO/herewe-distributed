package raft

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (r *Replica) RequestVote(params *ReqVoteParams) (*ReqVoteResult, error) {
	if params == nil {
		return nil, errors.New("wrong params")
	}

	err := r.VoteToU(params.CandidateId, params.Term)
	if err != nil {
		repLog.Debugf("rep %v vote to fail, err %v", r.conf.RepId, err)
		return &ReqVoteResult{
			Term:        r.CurrentTerm,
			VoteGranted: false,
		}, nil
	}

	repLog.Debugf("rep %v vote to success", r.conf.RepId)
	return &ReqVoteResult{
		Term:        r.CurrentTerm,
		VoteGranted: true,
	}, nil
}

func (r *Replica) AppendEntries(params *AppendEntriesParams) (*AppendEntriesResult, error) {

	if params == nil {
		return nil, errors.New("wrong params")
	}

	err := r.checkTerm(params.Term, Leader)
	if err != nil {
		return nil, err
	}
	// accept the leader
	r.recentlyLeaderId = params.LeaderId

	// check lastLog
	success := false
	if params.PrevLogIndex == r.LastLogIndex && params.PrevLogTerm == r.LastLogTerm {
		if err := r.saveLogEntries(r.LastLogIndex+1, params.Entries); err != nil {
			return nil, err
		}

		// commit by LeaderCommit
		if err := r.commitLogEntries(params.LeaderCommit); err != nil {
			return nil, err
		}
		success = true
		repLog.Debugf("rep %v follower append to params %v:%v, new %v, self %v:%v",
			r.conf.RepId, params.PrevLogIndex, params.PrevLogTerm,
			len(params.Entries), r.LastLogIndex, r.LastLogTerm)
	} else {
		// just return self state
		repLog.Debugf("rep %v not match leader, params %v:%v, self %v:%v",
			r.conf.RepId, params.PrevLogIndex, params.PrevLogTerm,
			r.LastLogIndex, r.LastLogTerm)
	}
	return &AppendEntriesResult{
		Term:                 r.CurrentTerm,
		Success:              success,
		LastLogTerm:          r.LastLogTerm,
		LastLogIndex:         r.LastLogIndex,
		FirstIndexInLastTerm: r.FirstIndexInTerm,
		CommitIndex:          r.CommitIndex,
	}, nil
}

func (r *Replica) checkTerm(incomeTerm int64, remoteState RepSate) error {
	if incomeTerm < r.CurrentTerm {
		return fmt.Errorf("old term %v:%v", incomeTerm, r.CurrentTerm)
	}

	if incomeTerm > r.CurrentTerm {
		err := r.transfer2Follower(incomeTerm, remoteState)
		if err != nil {
			return err
		}
	}

	// if self is follower, also need contain follower
	if Follower == r.state {
		err := r.transfer2Follower(incomeTerm, remoteState)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Replica) SendAppendEntries(heartbeat bool) error {
	repLog.Debugf("rep %v start send log entries, term %v", r.conf.RepId, r.CurrentTerm)

	// heartbeat periodically
	r.resetHeartbeatTimer()

	ctx := context.Background()
	for _, peer := range r.Cluster.RepPeers {
		if strings.EqualFold(peer.RepId, r.conf.RepId) {
			continue
		}
		nextEntry, err := r.queryEntryByIndex(r.nextIndex[peer.RepId] - 1)
		if err != nil {
			repLog.Debugf("rep %v query entry err %v", r.conf.RepId, err)
			return err
		}
		var entries []*LogEntry
		if !heartbeat {
			if entries, err = r.queryRangeEntryByIndex(r.nextIndex[peer.RepId], r.LastLogIndex); err != nil {
				repLog.Debugf("rep %v query range entry err %v", r.conf.RepId, err)
				return err
			}
		}
		params := &AppendEntriesParams{
			Term:         r.CurrentTerm,
			LeaderId:     r.conf.RepId,
			PrevLogIndex: nextEntry.Index,
			PrevLogTerm:  nextEntry.Term,
			Entries:      entries,
			LeaderCommit: r.CommitIndex,
		}

		repLog.Debugf("rep %v send to %v, params %v, next %v", r.conf.RepId,
			peer.RepId, params, r.nextIndex[peer.RepId])
		aeResult, err := r.rpcClient.SendAppendEntries(ctx, params, peer)
		if err != nil {
			repLog.Debugf("rep %v send AppendEntries to %v, err %v", r.conf.RepId, peer, err)
			continue
		}
		if err := r.handleAppendEntriesResult(aeResult, peer); err != nil {
			repLog.Debugf("rep %v append entry to %v err %v", r.conf.RepId, peer.RepId, err)
		}
	}
	// cal if retch majority append new ones
	reachCount := 1
	for _, next := range r.nextIndex {
		if next == r.LastLogIndex+1 {
			reachCount++
		}
	}
	repLog.Debugf("rep %v reach new entry count %v, majority %v", r.conf.RepId, reachCount, r.MajorityNum)
	if reachCount < r.MajorityNum {
		return fmt.Errorf("majority have not accept the entries, try it again")
	}

	return nil
}
