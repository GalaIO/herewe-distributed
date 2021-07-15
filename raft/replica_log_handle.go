package raft

import (
	"fmt"
)

func (r *Replica) requestLogEntry(command []byte) error {

	if r.state != Leader {
		return fmt.Errorf("I am not leader, state %v", r.state)
	}

	r.LastLogIndex++
	r.LastLogTerm = r.CurrentTerm
	entry := &LogEntry{
		Command: command,
		Index:   r.LastLogIndex,
		Term:    r.LastLogTerm,
	}

	r.saveLogEntry(r.LastLogIndex, entry)
	return r.SendAppendEntries()
}
