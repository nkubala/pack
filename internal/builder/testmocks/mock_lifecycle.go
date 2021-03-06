// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpacks/pack/internal/builder (interfaces: Lifecycle)

// Package testmocks is a generated GoMock package.
package testmocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	builder "github.com/buildpacks/pack/internal/builder"
)

// MockLifecycle is a mock of Lifecycle interface
type MockLifecycle struct {
	ctrl     *gomock.Controller
	recorder *MockLifecycleMockRecorder
}

// MockLifecycleMockRecorder is the mock recorder for MockLifecycle
type MockLifecycleMockRecorder struct {
	mock *MockLifecycle
}

// NewMockLifecycle creates a new mock instance
func NewMockLifecycle(ctrl *gomock.Controller) *MockLifecycle {
	mock := &MockLifecycle{ctrl: ctrl}
	mock.recorder = &MockLifecycleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLifecycle) EXPECT() *MockLifecycleMockRecorder {
	return m.recorder
}

// Descriptor mocks base method
func (m *MockLifecycle) Descriptor() builder.LifecycleDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Descriptor")
	ret0, _ := ret[0].(builder.LifecycleDescriptor)
	return ret0
}

// Descriptor indicates an expected call of Descriptor
func (mr *MockLifecycleMockRecorder) Descriptor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Descriptor", reflect.TypeOf((*MockLifecycle)(nil).Descriptor))
}

// Open mocks base method
func (m *MockLifecycle) Open() (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open
func (mr *MockLifecycleMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockLifecycle)(nil).Open))
}
