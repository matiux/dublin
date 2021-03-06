// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/matiux/dublin/dublin/commandhandling (interfaces: CommandHandler)

// Package mock_commandhandling is a generated GoMock package.
package commandhandling

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommandHandler is a mock of CommandHandler interface.
type MockCommandHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCommandHandlerMockRecorder
}

// MockCommandHandlerMockRecorder is the mock recorder for MockCommandHandler.
type MockCommandHandlerMockRecorder struct {
	mock *MockCommandHandler
}

// NewMockCommandHandler creates a new mock instance.
func NewMockCommandHandler(ctrl *gomock.Controller) *MockCommandHandler {
	mock := &MockCommandHandler{ctrl: ctrl}
	mock.recorder = &MockCommandHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandHandler) EXPECT() *MockCommandHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockCommandHandler) Handle(cmd Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", cmd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockCommandHandlerMockRecorder) Handle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockCommandHandler)(nil).Handle), arg0)
}
