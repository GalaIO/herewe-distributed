package raft

import (
	"encoding/json"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
)

var (
	ErrNotFound = errors.New("data not found")
)

// Storage provide a simple k/v to stable storage
type Storage interface {
	io.Closer
	SaveRepData(*RepData) error
	GetRepData() (*RepData, error)
}

var RepDataKey = []byte("rep_data")

type LevelDBStorage struct {
	db *leveldb.DB
}

func (l *LevelDBStorage) Close() error {
	return l.db.Close()
}

func (l *LevelDBStorage) SaveRepData(data *RepData) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return l.db.Put(RepDataKey, val, nil)
}

func (l *LevelDBStorage) GetRepData() (*RepData, error) {
	val, err := l.db.Get(RepDataKey, nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	data := new(RepData)
	if err = json.Unmarshal(val, data); err != nil {
		return nil, err
	}
	return data, nil
}

func NewLevelDB(conf StorageConfig) (Storage, error) {
	db, err := leveldb.OpenFile(conf.DataDir, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBStorage{
		db: db,
	}, nil
}
