// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GalaIO/herewe-distributed/raft (interfaces: Storage)

// Package raft is a generated GoMock package.
package raft

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetRepData mocks base method.
func (m *MockStorage) GetRepData() (*RepData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepData")
	ret0, _ := ret[0].(*RepData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepData indicates an expected call of GetRepData.
func (mr *MockStorageMockRecorder) GetRepData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepData", reflect.TypeOf((*MockStorage)(nil).GetRepData))
}

// SaveRepData mocks base method.
func (m *MockStorage) SaveRepData(arg0 *RepData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRepData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRepData indicates an expected call of SaveRepData.
func (mr *MockStorageMockRecorder) SaveRepData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRepData", reflect.TypeOf((*MockStorage)(nil).SaveRepData), arg0)
}
