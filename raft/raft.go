package raft

import (
	"fmt"
)

type Raft struct {
	conf   Config
	db     Storage
	rep    *Replica
	server RpcServer
}

func NewRaft(conf Config) (*Raft, error) {
	r := &Raft{
		conf: conf,
	}
	var err error
	r.db, err = NewLevelDB(conf.Storage)
	if err != nil {
		panic(fmt.Errorf("init levelDB err %v", err))
	}

	rpcClient := NewRpcClient()
	r.rep, err = NewRep(r.db, rpcClient, conf.Rep)
	if err != nil {
		panic(fmt.Errorf("init replica err %v", err))
	}
	r.server = NewRepServer(r.rep)

	return r, nil
}

func (r *Raft) Start() error {
	err := r.rep.Start()
	if err != nil {
		panic(fmt.Errorf("start rep err %v", err))
	}
	err = r.server.Start()
	if err != nil {
		panic(fmt.Errorf("start rep err %v", err))
	}
	return nil
}

func (r *Raft) Stop() error {
	err := r.rep.Stop()
	if err != nil {
		panic(fmt.Errorf("stop rep err %v", err))
	}
	err = r.server.Stop()
	if err != nil {
		panic(fmt.Errorf("stop rpc server err %v", err))
	}
	err = r.db.Close()
	if err != nil {
		panic(fmt.Errorf("stop db err %v", err))
	}
	return nil
}
