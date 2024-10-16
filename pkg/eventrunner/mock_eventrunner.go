// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go
//
// Generated by this command:
//
//	mockgen -destination=mock_eventrunner.go -package=eventrunner -source=./interfaces.go EventSink,EventConsumer,NATSClient,AppInterface
//

// Package eventrunner is a generated GoMock package.
package eventrunner

import (
	context "context"
	reflect "reflect"

	v2 "github.com/cloudevents/sdk-go/v2"
	gomock "go.uber.org/mock/gomock"
	gofr "gofr.dev/pkg/gofr"
	container "gofr.dev/pkg/gofr/container"
	datasource "gofr.dev/pkg/gofr/datasource"
	pubsub "gofr.dev/pkg/gofr/datasource/pubsub"
	logging "gofr.dev/pkg/gofr/logging"
	metrics "gofr.dev/pkg/gofr/metrics"
	migration "gofr.dev/pkg/gofr/migration"
)

// MockEventSink is a mock of EventSink interface.
type MockEventSink struct {
	ctrl     *gomock.Controller
	recorder *MockEventSinkMockRecorder
}

// MockEventSinkMockRecorder is the mock recorder for MockEventSink.
type MockEventSinkMockRecorder struct {
	mock *MockEventSink
}

// NewMockEventSink creates a new mock instance.
func NewMockEventSink(ctrl *gomock.Controller) *MockEventSink {
	mock := &MockEventSink{ctrl: ctrl}
	mock.recorder = &MockEventSinkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventSink) EXPECT() *MockEventSinkMockRecorder {
	return m.recorder
}

// LogEvent mocks base method.
func (m *MockEventSink) LogEvent(arg0 context.Context, arg1 *v2.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogEvent indicates an expected call of LogEvent.
func (mr *MockEventSinkMockRecorder) LogEvent(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogEvent", reflect.TypeOf((*MockEventSink)(nil).LogEvent), arg0, arg1)
}

// MockEventConsumer is a mock of EventConsumer interface.
type MockEventConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockEventConsumerMockRecorder
}

// MockEventConsumerMockRecorder is the mock recorder for MockEventConsumer.
type MockEventConsumerMockRecorder struct {
	mock *MockEventConsumer
}

// NewMockEventConsumer creates a new mock instance.
func NewMockEventConsumer(ctrl *gomock.Controller) *MockEventConsumer {
	mock := &MockEventConsumer{ctrl: ctrl}
	mock.recorder = &MockEventConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventConsumer) EXPECT() *MockEventConsumerMockRecorder {
	return m.recorder
}

// ConsumeEvent mocks base method.
func (m *MockEventConsumer) ConsumeEvent(ctx *gofr.Context, event *v2.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeEvent indicates an expected call of ConsumeEvent.
func (mr *MockEventConsumerMockRecorder) ConsumeEvent(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeEvent", reflect.TypeOf((*MockEventConsumer)(nil).ConsumeEvent), ctx, event)
}

// MockBuffer is a mock of Buffer interface.
type MockBuffer struct {
	ctrl     *gomock.Controller
	recorder *MockBufferMockRecorder
}

// MockBufferMockRecorder is the mock recorder for MockBuffer.
type MockBufferMockRecorder struct {
	mock *MockBuffer
}

// NewMockBuffer creates a new mock instance.
func NewMockBuffer(ctrl *gomock.Controller) *MockBuffer {
	mock := &MockBuffer{ctrl: ctrl}
	mock.recorder = &MockBufferMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuffer) EXPECT() *MockBufferMockRecorder {
	return m.recorder
}

// Bytes mocks base method.
func (m *MockBuffer) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes.
func (mr *MockBufferMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockBuffer)(nil).Bytes))
}

// Reset mocks base method.
func (m *MockBuffer) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockBufferMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockBuffer)(nil).Reset))
}

// Write mocks base method.
func (m *MockBuffer) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockBufferMockRecorder) Write(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockBuffer)(nil).Write), p)
}

// MockNATSClient is a mock of NATSClient interface.
type MockNATSClient struct {
	ctrl     *gomock.Controller
	recorder *MockNATSClientMockRecorder
}

// MockNATSClientMockRecorder is the mock recorder for MockNATSClient.
type MockNATSClientMockRecorder struct {
	mock *MockNATSClient
}

// NewMockNATSClient creates a new mock instance.
func NewMockNATSClient(ctrl *gomock.Controller) *MockNATSClient {
	mock := &MockNATSClient{ctrl: ctrl}
	mock.recorder = &MockNATSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNATSClient) EXPECT() *MockNATSClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockNATSClient) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockNATSClientMockRecorder) Close(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockNATSClient)(nil).Close), ctx)
}

// Connect mocks base method.
func (m *MockNATSClient) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockNATSClientMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockNATSClient)(nil).Connect))
}

// CreateTopic mocks base method.
func (m *MockNATSClient) CreateTopic(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTopic", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTopic indicates an expected call of CreateTopic.
func (mr *MockNATSClientMockRecorder) CreateTopic(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopic", reflect.TypeOf((*MockNATSClient)(nil).CreateTopic), ctx, name)
}

// DeleteTopic mocks base method.
func (m *MockNATSClient) DeleteTopic(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTopic", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTopic indicates an expected call of DeleteTopic.
func (mr *MockNATSClientMockRecorder) DeleteTopic(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopic", reflect.TypeOf((*MockNATSClient)(nil).DeleteTopic), ctx, name)
}

// Health mocks base method.
func (m *MockNATSClient) Health() datasource.Health {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health")
	ret0, _ := ret[0].(datasource.Health)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockNATSClientMockRecorder) Health() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockNATSClient)(nil).Health))
}

// Publish mocks base method.
func (m *MockNATSClient) Publish(ctx context.Context, topic string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockNATSClientMockRecorder) Publish(ctx, topic, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockNATSClient)(nil).Publish), ctx, topic, message)
}

// Subscribe mocks base method.
func (m *MockNATSClient) Subscribe(ctx context.Context, topic string) (*pubsub.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, topic)
	ret0, _ := ret[0].(*pubsub.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockNATSClientMockRecorder) Subscribe(ctx, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockNATSClient)(nil).Subscribe), ctx, topic)
}

// UseLogger mocks base method.
func (m *MockNATSClient) UseLogger(logger any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UseLogger", logger)
}

// UseLogger indicates an expected call of UseLogger.
func (mr *MockNATSClientMockRecorder) UseLogger(logger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseLogger", reflect.TypeOf((*MockNATSClient)(nil).UseLogger), logger)
}

// UseMetrics mocks base method.
func (m *MockNATSClient) UseMetrics(metrics any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UseMetrics", metrics)
}

// UseMetrics indicates an expected call of UseMetrics.
func (mr *MockNATSClientMockRecorder) UseMetrics(metrics any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseMetrics", reflect.TypeOf((*MockNATSClient)(nil).UseMetrics), metrics)
}

// UseTracer mocks base method.
func (m *MockNATSClient) UseTracer(tracer any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UseTracer", tracer)
}

// UseTracer indicates an expected call of UseTracer.
func (mr *MockNATSClientMockRecorder) UseTracer(tracer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UseTracer", reflect.TypeOf((*MockNATSClient)(nil).UseTracer), tracer)
}

// MockAppInterface is a mock of AppInterface interface.
type MockAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAppInterfaceMockRecorder
}

// MockAppInterfaceMockRecorder is the mock recorder for MockAppInterface.
type MockAppInterfaceMockRecorder struct {
	mock *MockAppInterface
}

// NewMockAppInterface creates a new mock instance.
func NewMockAppInterface(ctrl *gomock.Controller) *MockAppInterface {
	mock := &MockAppInterface{ctrl: ctrl}
	mock.recorder = &MockAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppInterface) EXPECT() *MockAppInterfaceMockRecorder {
	return m.recorder
}

// AddCassandra mocks base method.
func (m *MockAppInterface) AddCassandra(cassandraClient container.CassandraProvider) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCassandra", cassandraClient)
}

// AddCassandra indicates an expected call of AddCassandra.
func (mr *MockAppInterfaceMockRecorder) AddCassandra(cassandraClient any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCassandra", reflect.TypeOf((*MockAppInterface)(nil).AddCassandra), cassandraClient)
}

// AddPubSub mocks base method.
func (m *MockAppInterface) AddPubSub(pubsub container.PubSubProvider) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddPubSub", pubsub)
}

// AddPubSub indicates an expected call of AddPubSub.
func (mr *MockAppInterfaceMockRecorder) AddPubSub(pubsub any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPubSub", reflect.TypeOf((*MockAppInterface)(nil).AddPubSub), pubsub)
}

// Logger mocks base method.
func (m *MockAppInterface) Logger() logging.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logger")
	ret0, _ := ret[0].(logging.Logger)
	return ret0
}

// Logger indicates an expected call of Logger.
func (mr *MockAppInterfaceMockRecorder) Logger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*MockAppInterface)(nil).Logger))
}

// Metrics mocks base method.
func (m *MockAppInterface) Metrics() metrics.Manager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metrics")
	ret0, _ := ret[0].(metrics.Manager)
	return ret0
}

// Metrics indicates an expected call of Metrics.
func (mr *MockAppInterfaceMockRecorder) Metrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metrics", reflect.TypeOf((*MockAppInterface)(nil).Metrics))
}

// Migrate mocks base method.
func (m *MockAppInterface) Migrate(migrationsMap map[int64]migration.Migrate) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Migrate", migrationsMap)
}

// Migrate indicates an expected call of Migrate.
func (mr *MockAppInterfaceMockRecorder) Migrate(migrationsMap any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockAppInterface)(nil).Migrate), migrationsMap)
}

// Run mocks base method.
func (m *MockAppInterface) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockAppInterfaceMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockAppInterface)(nil).Run))
}

// Subscribe mocks base method.
func (m *MockAppInterface) Subscribe(topic string, handler gofr.SubscribeFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", topic, handler)
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockAppInterfaceMockRecorder) Subscribe(topic, handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockAppInterface)(nil).Subscribe), topic, handler)
}
