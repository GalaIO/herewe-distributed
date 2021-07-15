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

	if len(params.Entries) == 0 {
		err := r.OnReceiveHeartBeat(params.Term, params.LeaderId)
		if err != nil {
			repLog.Debugf("rep %v heartbeat fail, err %v", r.conf.RepId, err)
			return &AppendEntriesResult{
				Term:    r.CurrentTerm,
				Success: false,
			}, nil
		}
		return &AppendEntriesResult{
			Term:    r.CurrentTerm,
			Success: true,
		}, nil
	}

	//TODO
	//r.HandleAppendEntries()
	repLog.Debugf("rep %v vote to success", r.conf.RepId)
	return &AppendEntriesResult{
		Term:    r.CurrentTerm,
		Success: false,
	}, nil
}

func (r *Replica) SendAppendEntries() error {
	repLog.Debugf("rep %v start send log entries, term %v", r.conf.RepId, r.CurrentTerm)
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
		entries, err := r.queryRangeEntryByIndex(r.nextIndex[peer.RepId], r.LastLogIndex)
		if err != nil {
			repLog.Debugf("rep %v query range entry err %v", r.conf.RepId, err)
			return err
		}
		params := &AppendEntriesParams{
			Term:         r.CurrentTerm,
			LeaderId:     r.conf.RepId,
			PrevLogIndex: nextEntry.Index,
			PrevLogTerm:  nextEntry.Term,
			Entries:      entries,
			LeaderCommit: r.CommitIndex,
		}
		aeResult, err := r.rpcClient.SendAppendEntries(ctx, params, peer)
		if err != nil {
			repLog.Debugf("rep %v send AppendEntries to %v, err %v", r.conf.RepId, peer, err)
			continue
		}
		if err := r.handleAppendEntriesResult(aeResult, peer); err != nil {
			repLog.Debugf("rep %v append entry to err %v", r.conf.RepId, peer.RepId, err)
		}
	}
	// cal if retch majority append new ones
	reachCount := 0
	for _, next := range r.nextIndex {
		if next == r.LastLogIndex {
			reachCount++
		}
	}
	repLog.Debugf("rep %v reach new entry count %v", r.conf.RepId, reachCount)
	if reachCount < r.MajorityNum {
		return fmt.Errorf("majority have not accept the entries, try it again")
	}

	return nil
}
