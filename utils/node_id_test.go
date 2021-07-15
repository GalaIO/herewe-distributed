package utils

import "testing"

func Test_GenNodeId(t *testing.T) {
	t.Log(genNodeId("test"))
}

func Test_RandGenNodeId(t *testing.T) {
	t.Log(randGenNodeId())
}
