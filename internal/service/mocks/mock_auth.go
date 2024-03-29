// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/auth.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	service "tigerhall_kittens/internal/service"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockAuthService) LoginUser(ctx context.Context, req service.LoginUserReq) (*service.LoginUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, req)
	ret0, _ := ret[0].(*service.LoginUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockAuthServiceMockRecorder) LoginUser(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockAuthService)(nil).LoginUser), ctx, req)
}
