// Code generated by MockGen. DO NOT EDIT.
// Source: challenge.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	domain "crypto-project-1/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockChallengeRepository is a mock of ChallengeRepository interface.
type MockChallengeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChallengeRepositoryMockRecorder
}

// MockChallengeRepositoryMockRecorder is the mock recorder for MockChallengeRepository.
type MockChallengeRepositoryMockRecorder struct {
	mock *MockChallengeRepository
}

// NewMockChallengeRepository creates a new mock instance.
func NewMockChallengeRepository(ctrl *gomock.Controller) *MockChallengeRepository {
	mock := &MockChallengeRepository{ctrl: ctrl}
	mock.recorder = &MockChallengeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChallengeRepository) EXPECT() *MockChallengeRepositoryMockRecorder {
	return m.recorder
}

// CreateChallenge mocks base method.
func (m *MockChallengeRepository) CreateChallenge(arg0, arg1 string, arg2 int64) (*domain.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChallenge", arg0, arg1, arg2)
	ret0, _ := ret[0].(*domain.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChallenge indicates an expected call of CreateChallenge.
func (mr *MockChallengeRepositoryMockRecorder) CreateChallenge(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChallenge", reflect.TypeOf((*MockChallengeRepository)(nil).CreateChallenge), arg0, arg1, arg2)
}

// GetChallenges mocks base method.
func (m *MockChallengeRepository) GetChallenges(arg0, arg1 string) ([]*domain.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallenges", arg0, arg1)
	ret0, _ := ret[0].([]*domain.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallenges indicates an expected call of GetChallenges.
func (mr *MockChallengeRepositoryMockRecorder) GetChallenges(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallenges", reflect.TypeOf((*MockChallengeRepository)(nil).GetChallenges), arg0, arg1)
}
