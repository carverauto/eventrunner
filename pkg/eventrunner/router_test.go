package eventrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/datasource"
	cassandraPkg "gofr.dev/pkg/gofr/datasource/cassandra"
	"gofr.dev/pkg/gofr/datasource/pubsub"
	"gofr.dev/pkg/gofr/logging"
	"gofr.dev/pkg/gofr/metrics"
	"gofr.dev/pkg/gofr/migration"
)

// MockConsumerManager is a mock of ConsumerManager interface
type MockConsumerManager struct {
	ConsumeEventFunc func(*gofr.Context, *cloudevents.Event) error
}

func (m *MockConsumerManager) ConsumeEvent(c *gofr.Context, e *cloudevents.Event) error {
	if m.ConsumeEventFunc != nil {
		return m.ConsumeEventFunc(c, e)
	}
	return nil
}

type MockPubSubWrapper struct {
	PublishFunc     func(ctx context.Context, topic string, message []byte) error
	SubscribeFunc   func(ctx context.Context, topic string) (*pubsub.Message, error)
	CreateTopicFunc func(ctx context.Context, name string) error
	DeleteTopicFunc func(ctx context.Context, name string) error
	CloseFunc       func() error
	HealthFunc      func() datasource.Health
	ConnectFunc     func()
	UseLoggerFunc   func(logger any)
	UseMetricsFunc  func(metrics any)
	UseTracerFunc   func(tracer any)
}

func (m *MockPubSubWrapper) Publish(ctx context.Context, topic string, message []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, topic, message)
	}
	return nil
}

type MockRequest struct {
	ctx    context.Context
	params map[string][]string // Adjusted to hold slices of strings
	body   []byte
}

func (r *MockRequest) Context() context.Context {
	return r.ctx
}

func (r *MockRequest) Param(key string) string {
	if vals, ok := r.params[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (r *MockRequest) Params(key string) []string {
	return r.params[key]
}

func (r *MockRequest) PathParam(key string) string {
	if vals, ok := r.params[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (r *MockRequest) Bind(i interface{}) error {
	return json.Unmarshal(r.body, i)
}

func (r *MockRequest) HostName() string {
	return "localhost"
}

func TestEventRouter_handleEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	t.Run("Valid CloudEvent", func(t *testing.T) {
		event := cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetType("test.event")
		event.SetSource("test")
		eventJSON, _ := json.Marshal(event)

		mockRequest := &MockRequest{
			ctx:    context.Background(),
			body:   eventJSON,
			params: make(map[string][]string),
		}

		mockContext := &gofr.Context{
			Context:   context.Background(),
			Request:   mockRequest,
			Container: mockContainer,
		}

		err := er.handleEvent(mockContext)
		assert.NoError(t, err)
	})

}

func TestEventRouter_routeEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")

	mockRequest := &MockRequest{
		ctx: context.Background(),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	mockConsumerManager.ConsumeEventFunc = func(c *gofr.Context, e *cloudevents.Event) error {
		return nil
	}

	err := er.routeEvent(mockContext, &event)
	assert.NoError(t, err)
}

func TestNewEventRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Use the generated mock
	mockNatsClient := NewMockNATSClient(ctrl)

	// Set up expectations on the mock NATS client
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockNatsClient.EXPECT().Connect().Return(nil).AnyTimes()

	// Create a mock Cassandra client or pass nil
	var mockCassandraClient cassandraPkg.Client = nil // Or implement a mock if needed

	// Create the EventRouter with mock clients
	er := NewEventRouter(mockNatsClient, mockCassandraClient)

	// Assertions
	assert.NotNil(t, er)
	assert.NotNil(t, er.app)
	assert.NotNil(t, er.natsClient)
	assert.NotNil(t, er.bufferPool)
	assert.NotNil(t, er.consumerManager)
}

func TestEventRouter_Use(t *testing.T) {
	er := &EventRouter{
		middlewares: []Middleware{},
	}

	// Define a sample middleware
	sampleMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(c *gofr.Context, e *cloudevents.Event) error {
			// Middleware logic here
			return next(c, e)
		}
	}

	// Use the middleware
	er.Use(sampleMiddleware)

	// Assert that the middleware was added
	assert.Len(t, er.middlewares, 1)
	assert.Equal(t, sampleMiddleware, er.middlewares[0])
}

// Define MockApp that implements AppInterface
type MockApp struct {
	subscribeCalled bool
	runCalled       bool
	subscribeFunc   func(topic string, handler HandlerFunc)
	runFunc         func()
}

// func (m *MockApp) Subscribe(topic string, handler gofr.Handler) {
func (m *MockApp) Subscribe(topic string, handler HandlerFunc) {
	m.subscribeCalled = true
	if m.subscribeFunc != nil {
		m.subscribeFunc(topic, handler)
	}
}

func (m *MockApp) Run() {
	m.runCalled = true
	if m.runFunc != nil {
		m.runFunc()
	}
}

func (m *MockApp) Logger() logging.Logger {
	// Return a mock or real logger
	return logging.NewLogger(logging.INFO)
}

func (m *MockApp) Metrics() metrics.Manager {
	// Return a mock or real metrics manager
	return nil
}

func (m *MockApp) AddPubSub(pubsubClient pubsub.Client) {
	// Implement if necessary
}

// have the method: AddCassandra(cassandraClient *cassandra. Client)
func (m *MockApp) AddCassandra(cassandraClient container.CassandraProvider) {
	// Implement if necessary
}

func (m *MockApp) Migrate(migrationsMap map[int64]migration.Migrate) {
	// Implement if necessary
}

func TestEventRouter_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := NewMockAppInterface(ctrl)

	er := &EventRouter{
		app: mockApp,
	}

	// Set up expectations on mockApp
	mockApp.EXPECT().Subscribe("events.products", er.handleEvent).Times(1)
	mockApp.EXPECT().Run().Times(1)

	// Call Start
	er.Start()
}

func TestEventRouter_handleEvent_NonCloudEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	// Simulate a raw JSON message
	rawMessage := map[string]interface{}{
		"foo": "bar",
	}
	messageJSON, _ := json.Marshal(rawMessage)

	mockRequest := &MockRequest{
		ctx:    context.Background(),
		body:   messageJSON,
		params: make(map[string][]string),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	err := er.handleEvent(mockContext)
	assert.NoError(t, err)
}

func TestEventRouter_handleEvent_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	// Simulate invalid JSON
	invalidJSON := []byte("{ invalid json }")

	mockRequest := &MockRequest{
		ctx:    context.Background(),
		body:   invalidJSON,
		params: make(map[string][]string),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	err := er.handleEvent(mockContext)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse message")
}

func TestEventRouter_applyMiddleware_NoMiddlewares(t *testing.T) {
	er := &EventRouter{
		middlewares: []Middleware{},
	}

	handlerCalled := false

	handler := func(c *gofr.Context, e *cloudevents.Event) error {
		handlerCalled = true
		return nil
	}

	wrappedHandler := er.applyMiddleware(handler)

	// Call the handler
	err := wrappedHandler(&gofr.Context{}, &cloudevents.Event{})
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestEventRouter_applyMiddleware_WithMiddlewares(t *testing.T) {
	er := &EventRouter{}

	callSequence := []string{}

	middleware1 := func(next HandlerFunc) HandlerFunc {
		return func(c *gofr.Context, e *cloudevents.Event) error {
			callSequence = append(callSequence, "middleware1 before")
			err := next(c, e)
			callSequence = append(callSequence, "middleware1 after")
			return err
		}
	}

	middleware2 := func(next HandlerFunc) HandlerFunc {
		return func(c *gofr.Context, e *cloudevents.Event) error {
			callSequence = append(callSequence, "middleware2 before")
			err := next(c, e)
			callSequence = append(callSequence, "middleware2 after")
			return err
		}
	}

	er.Use(middleware1)
	er.Use(middleware2)

	handler := func(c *gofr.Context, e *cloudevents.Event) error {
		callSequence = append(callSequence, "handler")
		return nil
	}

	wrappedHandler := er.applyMiddleware(handler)

	// Call the handler
	err := wrappedHandler(&gofr.Context{}, &cloudevents.Event{})
	assert.NoError(t, err)

	// Assert the call sequence
	expectedSequence := []string{
		"middleware2 before",
		"middleware1 before",
		"handler",
		"middleware1 after",
		"middleware2 after",
	}
	assert.Equal(t, expectedSequence, callSequence)
}

func TestEventRouter_routeEvent_ConsumeEventError(t *testing.T) {
	consumeEventCalled := false

	mockConsumerManager := &MockConsumerManager{
		ConsumeEventFunc: func(c *gofr.Context, e *cloudevents.Event) error {
			consumeEventCalled = true
			return fmt.Errorf("consume event error")
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")

	mockRequest := &MockRequest{
		ctx: context.Background(),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	err := er.routeEvent(mockContext, &event)
	assert.NoError(t, err)
	assert.True(t, consumeEventCalled)
}

func TestEventRouter_routeEvent_PublishError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")

	mockRequest := &MockRequest{
		ctx: context.Background(),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	err := er.routeEvent(mockContext, &event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish event")
}

type UnmarshalableType struct{}

func (u UnmarshalableType) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}

func TestEventRouter_routeEvent_EncodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
	}

	mockContainer, _ := container.NewMockContainer(t)

	// Create an event that cannot be marshaled
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")

	// Set event data to an unmarshalable value
	event.SetData(cloudevents.ApplicationJSON, UnmarshalableType{})

	mockRequest := &MockRequest{
		ctx: context.Background(),
	}

	mockContext := &gofr.Context{
		Context:   context.Background(),
		Request:   mockRequest,
		Container: mockContainer,
	}

	err := er.routeEvent(mockContext, &event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encode event")
}
