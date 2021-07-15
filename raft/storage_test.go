package raft

type mockStorage struct {
}

func (m mockStorage) Close() error {
	return nil
}

func (m mockStorage) SaveRepData(data *RepData) error {
	return nil
}

func (m mockStorage) GetRepData() (*RepData, error) {
	return nil, nil
}

func NewMockStorage() Storage {
	return mockStorage{}
}
