// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetUserByFullName mocks base method.
func (m *MockRepositoryInterface) GetUserByFullName(ctx context.Context, input GetUserByFullNameInput) (UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByFullName", ctx, input)
	ret0, _ := ret[0].(UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByFullName indicates an expected call of GetUserByFullName.
func (mr *MockRepositoryInterfaceMockRecorder) GetUserByFullName(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByFullName", reflect.TypeOf((*MockRepositoryInterface)(nil).GetUserByFullName), ctx, input)
}

// GetUserByID mocks base method.
func (m *MockRepositoryInterface) GetUserByID(ctx context.Context, input GetUserByIDInput) (UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, input)
	ret0, _ := ret[0].(UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockRepositoryInterfaceMockRecorder) GetUserByID(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetUserByID), ctx, input)
}

// GetUserByPhoneNumber mocks base method.
func (m *MockRepositoryInterface) GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByPhoneNumber", ctx, input)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByPhoneNumber indicates an expected call of GetUserByPhoneNumber.
func (mr *MockRepositoryInterfaceMockRecorder) GetUserByPhoneNumber(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByPhoneNumber", reflect.TypeOf((*MockRepositoryInterface)(nil).GetUserByPhoneNumber), ctx, input)
}

// InsertUser mocks base method.
func (m *MockRepositoryInterface) InsertUser(ctx context.Context, input User) (InsertUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", ctx, input)
	ret0, _ := ret[0].(InsertUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockRepositoryInterfaceMockRecorder) InsertUser(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockRepositoryInterface)(nil).InsertUser), ctx, input)
}

// UpdateLastLoginAndSuccessfullyLogin mocks base method.
func (m *MockRepositoryInterface) UpdateLastLoginAndSuccessfullyLogin(ctx context.Context, input UpdateLastLoginAndSuccessfullyLoginInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastLoginAndSuccessfullyLogin", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastLoginAndSuccessfullyLogin indicates an expected call of UpdateLastLoginAndSuccessfullyLogin.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateLastLoginAndSuccessfullyLogin(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastLoginAndSuccessfullyLogin", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateLastLoginAndSuccessfullyLogin), ctx, input)
}

// UpdateUser mocks base method.
func (m *MockRepositoryInterface) UpdateUser(ctx context.Context, input UpdateUserInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateUser(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateUser), ctx, input)
}
