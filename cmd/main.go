package main

import (
	"encoding/json"
	"fmt"
	"github.com/GalaIO/herewe-distributed/logger"
	"github.com/GalaIO/herewe-distributed/raft"
	"io/ioutil"
	"os"
)

var mainLog = logger.GetLogger("main")

func main() {
	file, err := os.Open("./raft/config.json")
	if err != nil {
		panic(fmt.Errorf("init open config file err %v", err))
	}

	bytes, _ := ioutil.ReadAll(file)
	conf := new(raft.Config)
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		panic(fmt.Errorf("parse config json err %v", err))
	}
	r, err := raft.NewRaft(*conf)
	if err != nil {
		panic(fmt.Errorf("NewRaft err %v", err))
	}

	err = r.Start()
	if err != nil {
		panic(fmt.Errorf("start raft err %v", err))
	}
}
