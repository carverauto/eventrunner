// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go
//
// Generated by this command:
//
//	mockgen -destination=mock_eventingest.go -package=eventingest -source=./interfaces.go ServiceClient,EventForwarder
//

// Package eventingest is a generated GoMock package.
package eventingest

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
	gofr "gofr.dev/pkg/gofr"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockServiceClient is a mock of ServiceClient interface.
type MockServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockServiceClientMockRecorder
}

// MockServiceClientMockRecorder is the mock recorder for MockServiceClient.
type MockServiceClientMockRecorder struct {
	mock *MockServiceClient
}

// NewMockServiceClient creates a new mock instance.
func NewMockServiceClient(ctrl *gomock.Controller) *MockServiceClient {
	mock := &MockServiceClient{ctrl: ctrl}
	mock.recorder = &MockServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceClient) EXPECT() *MockServiceClientMockRecorder {
	return m.recorder
}

// IngestEvent mocks base method.
func (m *MockServiceClient) IngestEvent(ctx context.Context, in *IngestEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IngestEvent", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IngestEvent indicates an expected call of IngestEvent.
func (mr *MockServiceClientMockRecorder) IngestEvent(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IngestEvent", reflect.TypeOf((*MockServiceClient)(nil).IngestEvent), varargs...)
}

// MockEventForwarder is a mock of EventForwarder interface.
type MockEventForwarder struct {
	ctrl     *gomock.Controller
	recorder *MockEventForwarderMockRecorder
}

// MockEventForwarderMockRecorder is the mock recorder for MockEventForwarder.
type MockEventForwarderMockRecorder struct {
	mock *MockEventForwarder
}

// NewMockEventForwarder creates a new mock instance.
func NewMockEventForwarder(ctrl *gomock.Controller) *MockEventForwarder {
	mock := &MockEventForwarder{ctrl: ctrl}
	mock.recorder = &MockEventForwarderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventForwarder) EXPECT() *MockEventForwarderMockRecorder {
	return m.recorder
}

// ForwardEvent mocks base method.
func (m *MockEventForwarder) ForwardEvent(c *gofr.Context, tenantID, customerID uuid.UUID, eventData []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForwardEvent", c, tenantID, customerID, eventData)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForwardEvent indicates an expected call of ForwardEvent.
func (mr *MockEventForwarderMockRecorder) ForwardEvent(c, tenantID, customerID, eventData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardEvent", reflect.TypeOf((*MockEventForwarder)(nil).ForwardEvent), c, tenantID, customerID, eventData)
}