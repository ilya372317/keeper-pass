// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/service/data/service.go

// Package data_mock is a generated GoMock package.
package data_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/ilya372317/pass-keeper/internal/server/domain"
)

// Mockkeyring is a mock of keyring interface.
type Mockkeyring struct {
	ctrl     *gomock.Controller
	recorder *MockkeyringMockRecorder
}

// MockkeyringMockRecorder is the mock recorder for Mockkeyring.
type MockkeyringMockRecorder struct {
	mock *Mockkeyring
}

// NewMockkeyring creates a new mock instance.
func NewMockkeyring(ctrl *gomock.Controller) *Mockkeyring {
	mock := &Mockkeyring{ctrl: ctrl}
	mock.recorder = &MockkeyringMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockkeyring) EXPECT() *MockkeyringMockRecorder {
	return m.recorder
}

// GetGeneralKey mocks base method.
func (m *Mockkeyring) GetGeneralKey(arg0 context.Context) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGeneralKey", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGeneralKey indicates an expected call of GetGeneralKey.
func (mr *MockkeyringMockRecorder) GetGeneralKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGeneralKey", reflect.TypeOf((*Mockkeyring)(nil).GetGeneralKey), arg0)
}

// MockdataStorage is a mock of dataStorage interface.
type MockdataStorage struct {
	ctrl     *gomock.Controller
	recorder *MockdataStorageMockRecorder
}

// MockdataStorageMockRecorder is the mock recorder for MockdataStorage.
type MockdataStorageMockRecorder struct {
	mock *MockdataStorage
}

// NewMockdataStorage creates a new mock instance.
func NewMockdataStorage(ctrl *gomock.Controller) *MockdataStorage {
	mock := &MockdataStorage{ctrl: ctrl}
	mock.recorder = &MockdataStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdataStorage) EXPECT() *MockdataStorageMockRecorder {
	return m.recorder
}

// SaveData mocks base method.
func (m *MockdataStorage) SaveData(arg0 context.Context, arg1 domain.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveData", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveData indicates an expected call of SaveData.
func (mr *MockdataStorageMockRecorder) SaveData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveData", reflect.TypeOf((*MockdataStorage)(nil).SaveData), arg0, arg1)
}
