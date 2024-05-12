// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/service/securedata/service.go

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

// Delete mocks base method.
func (m *MockdataStorage) Delete(ctx context.Context, ids []int, userID uint, kinds []domain.Kind) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, ids, userID, kinds)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockdataStorageMockRecorder) Delete(ctx, ids, userID, kinds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockdataStorage)(nil).Delete), ctx, ids, userID, kinds)
}

// GetAll mocks base method.
func (m *MockdataStorage) GetAll(ctx context.Context, userID uint) ([]domain.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, userID)
	ret0, _ := ret[0].([]domain.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockdataStorageMockRecorder) GetAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockdataStorage)(nil).GetAll), ctx, userID)
}

// GetDataByID mocks base method.
func (m *MockdataStorage) GetDataByID(ctx context.Context, id int) (domain.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataByID", ctx, id)
	ret0, _ := ret[0].(domain.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataByID indicates an expected call of GetDataByID.
func (mr *MockdataStorageMockRecorder) GetDataByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataByID", reflect.TypeOf((*MockdataStorage)(nil).GetDataByID), ctx, id)
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

// UpdateByID mocks base method.
func (m *MockdataStorage) UpdateByID(ctx context.Context, id int, dto domain.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateByID", ctx, id, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateByID indicates an expected call of UpdateByID.
func (mr *MockdataStorageMockRecorder) UpdateByID(ctx, id, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateByID", reflect.TypeOf((*MockdataStorage)(nil).UpdateByID), ctx, id, dto)
}