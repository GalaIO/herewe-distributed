package raft

import (
	"errors"
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
		err := r.HandleHeartBeat(params.Term, params.LeaderId)
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
