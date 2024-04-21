// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/service/auth/service.go

// Package auth_mock is a generated GoMock package.
package auth_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/ilya372317/pass-keeper/internal/server/domain"
)

// MockTokenManager is a mock of TokenManager interface.
type MockTokenManager struct {
	ctrl     *gomock.Controller
	recorder *MockTokenManagerMockRecorder
}

// MockTokenManagerMockRecorder is the mock recorder for MockTokenManager.
type MockTokenManagerMockRecorder struct {
	mock *MockTokenManager
}

// NewMockTokenManager creates a new mock instance.
func NewMockTokenManager(ctrl *gomock.Controller) *MockTokenManager {
	mock := &MockTokenManager{ctrl: ctrl}
	mock.recorder = &MockTokenManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenManager) EXPECT() *MockTokenManagerMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockTokenManager) Generate(user *domain.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockTokenManagerMockRecorder) Generate(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockTokenManager)(nil).Generate), user)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetUserByEmail mocks base method.
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), ctx, email)
}