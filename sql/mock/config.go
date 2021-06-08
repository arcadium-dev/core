// Code generated by MockGen. DO NOT EDIT.
// Source: arcadium.dev/core/sql (interfaces: Config)

// Package mocksql is a generated GoMock package.
package mocksql

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// DSN mocks base method.
func (m *MockConfig) DSN() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DSN")
	ret0, _ := ret[0].(string)
	return ret0
}

// DSN indicates an expected call of DSN.
func (mr *MockConfigMockRecorder) DSN() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DSN", reflect.TypeOf((*MockConfig)(nil).DSN))
}

// DriverName mocks base method.
func (m *MockConfig) DriverName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DriverName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DriverName indicates an expected call of DriverName.
func (mr *MockConfigMockRecorder) DriverName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DriverName", reflect.TypeOf((*MockConfig)(nil).DriverName))
}
