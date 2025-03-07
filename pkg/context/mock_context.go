// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/carverauto/eventrunner/pkg/context (interfaces: Context)
//
// Generated by this command:
//
//	mockgen -destination=mock_context.go -package=context github.com/carverauto/eventrunner/pkg/context Context
//

// Package context is a generated GoMock package.
package context

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
	gofr "gofr.dev/pkg/gofr"
)

// MockContext is a mock of Context interface.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// Bind mocks base method.
func (m *MockContext) Bind(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bind", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Bind indicates an expected call of Bind.
func (mr *MockContextMockRecorder) Bind(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bind", reflect.TypeOf((*MockContext)(nil).Bind), arg0)
}

// Context mocks base method.
func (m *MockContext) Context() *gofr.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(*gofr.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockContextMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockContext)(nil).Context))
}

// GetAPIKey mocks base method.
func (m *MockContext) GetAPIKey() (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAPIKey")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetAPIKey indicates an expected call of GetAPIKey.
func (mr *MockContextMockRecorder) GetAPIKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKey", reflect.TypeOf((*MockContext)(nil).GetAPIKey))
}

// GetClaim mocks base method.
func (m *MockContext) GetClaim(arg0 string) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaim", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetClaim indicates an expected call of GetClaim.
func (mr *MockContextMockRecorder) GetClaim(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaim", reflect.TypeOf((*MockContext)(nil).GetClaim), arg0)
}

// GetStringClaim mocks base method.
func (m *MockContext) GetStringClaim(arg0 string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringClaim", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetStringClaim indicates an expected call of GetStringClaim.
func (mr *MockContextMockRecorder) GetStringClaim(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringClaim", reflect.TypeOf((*MockContext)(nil).GetStringClaim), arg0)
}

// GetUUIDClaim mocks base method.
func (m *MockContext) GetUUIDClaim(arg0 string) (uuid.UUID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUUIDClaim", arg0)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetUUIDClaim indicates an expected call of GetUUIDClaim.
func (mr *MockContextMockRecorder) GetUUIDClaim(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUUIDClaim", reflect.TypeOf((*MockContext)(nil).GetUUIDClaim), arg0)
}

// SetClaim mocks base method.
func (m *MockContext) SetClaim(arg0 string, arg1 any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetClaim", arg0, arg1)
}

// SetClaim indicates an expected call of SetClaim.
func (mr *MockContextMockRecorder) SetClaim(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClaim", reflect.TypeOf((*MockContext)(nil).SetClaim), arg0, arg1)
}