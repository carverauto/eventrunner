package eventrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/datasource"
	"gofr.dev/pkg/gofr/datasource/pubsub"
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
	er := &EventRouter{
		app: gofr.New(),
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		consumerManager: &MockConsumerManager{},
		natsClient: &MockPubSubWrapper{
			PublishFunc: func(ctx context.Context, topic string, message []byte) error {
				// Simulate publishing behavior; return nil to indicate success
				return nil
			},
		},
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
	mockConsumerManager := &MockConsumerManager{}

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: &MockPubSubWrapper{
			PublishFunc: func(ctx context.Context, topic string, message []byte) error {
				// Simulate publishing behavior; you can return an error to test error handling
				return nil
			},
		},
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
