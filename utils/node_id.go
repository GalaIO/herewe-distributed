package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

func genNodeId(src string) string {
	hash := sha1.Sum([]byte(src))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func randGenNodeId() string {
	rand.Seed(time.Now().UnixNano())
	src := strconv.FormatUint(rand.Uint64(), 10)
	return genNodeId(src)
}
