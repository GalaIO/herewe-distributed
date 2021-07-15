package raft

import (
	"strconv"
	"testing"
	"time"
)

var mockConfig = Config{
	Rep: mockRepConfig,
	Storage: StorageConfig{
		DataDir: "./db",
	},
}

var mockRafts []*Raft

func init() {
	for i, peer := range mockCluster.RepPeers {
		mockConfig.Storage.DataDir = "./db/" + strconv.Itoa(i)

		mockConfig.Rep.RepId = peer.RepId
		mockConfig.Rep.Addr = peer.Addr

		r, _ := NewRaft(mockConfig)
		mockRafts = append(mockRafts, r)
	}
}

func TestRaft_ClusterLeaderCrash(t *testing.T) {
	for _, r := range mockRafts {
		go r.Start()
	}

	time.Sleep(1 * time.Second)
	var leader *Raft
	for _, r := range mockRafts {
		t.Log("rep state", r.conf.Rep.RepId, r.rep.state)
		if r.rep.state == Leader {
			leader = r
		}
	}

	_ = leader.Stop()
	time.Sleep(1 * time.Second)
	for _, r := range mockRafts {
		t.Log("rep state", r.conf.Rep.RepId, r.rep.state)
	}

	leader, _ = NewRaft(leader.conf)
	go leader.Start()
	time.Sleep(1 * time.Second)
	for _, r := range mockRafts {
		t.Log("rep state", r.conf.Rep.RepId, r.rep.state)
	}
}
