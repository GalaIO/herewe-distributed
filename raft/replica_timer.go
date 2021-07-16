package raft

import (
	"fmt"
	"math/rand"
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
			err := r.whenElectionTimeout()
			if err != nil {
				repLog.Errorf("rep %v whenElectionTimeout err %v", r.conf.RepId, err)
			}
		case <-r.heartbeatTimer.C:
			err := r.whenHeartBeatTimeout()
			if err != nil {
				repLog.Errorf("rep %v whenHeartBeatTimeout err %v", r.conf.RepId, err)
			}
		case <-r.stopCh:
			repLog.Debugf("rep %v stop timer loop", r.conf.RepId)
			return
		}
	}
}

func (r *Replica) whenElectionTimeout() error {
	r.electStart = true
	r.transfer2Candidate()
	return nil
}

func (r *Replica) whenHeartBeatTimeout() error {
	if r.state != Leader {
		r.stopHeartbeatTimeout()
		return nil
	}
	repLog.Debugf("rep %v start send heartbeat, term %v", r.conf.RepId, r.CurrentTerm)
	return r.SendAppendEntries(true)
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
