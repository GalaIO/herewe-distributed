package raft

// RepPeer peer info
// {
// 	"repId": "xxx",
//  "addr":  "127.0.0.1:6701",
// }
type RepPeer struct {
	RepId string `json:"repId"`
	Addr  string `json:"addr"`
}

// ClusterConfig peers info
type ClusterConfig struct {
	RepPeers []RepPeer `json:"repPeers"` // map the repId and address
}

// RepData save replication persistent data
type RepData struct {
	Cluster     ClusterConfig `json:"cluster"`
	CurrentTerm int64         `json:"currentTerm"`
	VoteFor     string        `json:"VoteFor"`
	MajorityNum int           `json:"MajorityNum"` // majority of cluster
}

// RepConfig replication load local config
type RepConfig struct {
	RepId              string  `json:"repId"`              // the identify repId
	Addr               string  `json:"addr"`               // address to listen for msg
	MinElectionTimeout int     `json:"minElectionTimeout"` // MinElectionTimeout as mill
	MaxElectionTimeout int     `json:"maxElectionTimeout"` // MaxElectionTimeout as mill
	HeartbeatTimeout   int     `json:"heartbeatTimeout"`   // HeartbeatTimeout as mill
	RepData            RepData `json:"repData"`            // when dataDir is empty, will load cluster from input config
}

type StorageConfig struct {
	DataDir string `json:"dataDir"`
}

type Config struct {
	Rep     RepConfig     `json:"rep"`
	Storage StorageConfig `json:"storage"`
}
