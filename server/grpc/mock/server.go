// Code generated by MockGen. DO NOT EDIT.
// Source: arcadium.dev/core/server/grpc (interfaces: Server)

// Package mockgrpc is a generated GoMock package.
package mockgrpc

import (
	reflect "reflect"

	grpc "arcadium.dev/core/server/grpc"
	gomock "github.com/golang/mock/gomock"
)

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockServer) Register(arg0 []grpc.Service) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0)
}

// Register indicates an expected call of Register.
func (mr *MockServerMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockServer)(nil).Register), arg0)
}

// Serve mocks base method.
func (m *MockServer) Serve(arg0 chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Serve", arg0)
}

// Serve indicates an expected call of Serve.
func (mr *MockServerMockRecorder) Serve(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serve", reflect.TypeOf((*MockServer)(nil).Serve), arg0)
}

// Stop mocks base method.
func (m *MockServer) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServer)(nil).Stop))
}