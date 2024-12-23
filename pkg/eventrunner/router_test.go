/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

package eventrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	cassandraPkg "gofr.dev/pkg/gofr/datasource/cassandra"
	"gofr.dev/pkg/gofr/logging"
)

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

func (*MockRequest) HostName() string {
	return "localhost"
}

func TestEventRouter_handleEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.DEBUG)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		// app:             gofr.New(),
		app:             mockApp,
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
		logger:     mockLogger,
	}
	er.getBufferFunc = er.defaultGetBuffer

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
	done := make(chan struct{})

	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Do(func(ctx *gofr.Context, event *cloudevents.Event) {
		close(done)
	}).Return(nil).Times(1)

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.DEBUG)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             mockApp,
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
		logger:     mockLogger,
	}
	er.getBufferFunc = er.defaultGetBuffer

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

	// Wait for ConsumeEvent to be called
	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("ConsumeEvent was not called")
	}
}

func TestNewEventRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := NewMockAppInterface(ctrl)
	mockNatsClient := NewMockNATSClient(ctrl)
	mockCassandraClient := &cassandraPkg.Client{}

	mockApp.EXPECT().AddCassandra(gomock.Any()).Times(1)
	mockApp.EXPECT().Migrate(gomock.Any()).Times(1)
	mockApp.EXPECT().Logger().Return(logging.NewLogger(logging.INFO)).AnyTimes()
	mockApp.EXPECT().Metrics().Return(nil).AnyTimes()

	er := NewEventRouter(mockApp, mockNatsClient, mockCassandraClient)

	assert.NotNil(t, er)
	assert.Equal(t, mockApp, er.app)
	assert.Equal(t, mockNatsClient, er.natsClient)
	assert.NotNil(t, er.bufferPool)
	assert.NotNil(t, er.consumerManager)
}

func TestEventRouter_Use(t *testing.T) {
	er := &EventRouter{
		middlewares: []Middleware{},
	}

	sampleMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(c *gofr.Context, e *cloudevents.Event) error {
			return next(c, e)
		}
	}

	er.Use(sampleMiddleware)

	assert.Len(t, er.middlewares, 1)

	// Check if the function types match
	actualType := reflect.TypeOf(er.middlewares[0])
	expectedType := reflect.TypeOf(Middleware(nil))
	assert.Equal(t, expectedType, actualType)
}

func TestEventRouter_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := NewMockAppInterface(ctrl)

	er := &EventRouter{
		app: mockApp,
	}

	// Set up expectations on mockApp
	mockApp.EXPECT().Subscribe("events.products", gomock.Any()).Times(1)
	mockApp.EXPECT().Run().Times(1)

	// Call Start
	er.Start()
}

func TestEventRouter_handleEvent_NonCloudEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.DEBUG)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             mockApp,
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
		logger:     mockLogger,
	}
	er.getBufferFunc = er.defaultGetBuffer

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
	mockNatsClient := NewMockNATSClient(ctrl)

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
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse message")
}

func TestEventRouter_applyMiddleware_NoMiddlewares(t *testing.T) {
	er := &EventRouter{
		middlewares: []Middleware{},
	}

	handlerCalled := false

	handler := func(*gofr.Context, *cloudevents.Event) error {
		handlerCalled = true
		return nil
	}

	wrappedHandler := er.applyMiddleware(handler)

	// Call the handler
	err := wrappedHandler(&gofr.Context{}, &cloudevents.Event{})
	require.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestEventRouter_applyMiddleware_WithMiddlewares(t *testing.T) {
	er := &EventRouter{}

	var callSequence []string

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

	handler := func(*gofr.Context, *cloudevents.Event) error {
		callSequence = append(callSequence, "handler")
		return nil
	}

	wrappedHandler := er.applyMiddleware(handler)

	// Call the handler
	err := wrappedHandler(&gofr.Context{}, &cloudevents.Event{})
	require.NoError(t, err)

	// Assert the call sequence
	expectedSequence := []string{
		"middleware1 before",
		"middleware2 before",
		"handler",
		"middleware2 after",
		"middleware1 after",
	}
	assert.Equal(t, expectedSequence, callSequence)
}

func TestEventRouter_routeEvent_ConsumeEventError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(errConsumeEventError).AnyTimes()

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.INFO)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	er := &EventRouter{
		app:             mockApp,
		natsClient:      mockNatsClient,
		consumerManager: mockConsumerManager,
		bufferPool:      &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		logger:          mockLogger,
	}
	er.getBufferFunc = er.defaultGetBuffer // Add this line

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")

	mockContext := &gofr.Context{
		Context: context.Background(),
	}

	err := er.routeEvent(mockContext, &event)
	assert.NoError(t, err) // The error from ConsumeEvent is logged but not returned
}

func TestEventRouter_routeEvent_PublishError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.DEBUG)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)
	mockNatsClient.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errPublishError).Times(1)

	er := &EventRouter{
		app:             gofr.New(),
		consumerManager: mockConsumerManager,
		bufferPool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
		natsClient: mockNatsClient,
		logger:     mockLogger,
	}
	er.getBufferFunc = er.defaultGetBuffer // Add this line

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
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish event")
}

func TestEventRouter_routeEvent_EncodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumerManager := NewMockEventConsumer(ctrl)
	mockConsumerManager.EXPECT().ConsumeEvent(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockApp := NewMockAppInterface(ctrl)
	mockLogger := logging.NewMockLogger(logging.DEBUG)
	mockApp.EXPECT().Logger().Return(mockLogger).AnyTimes()

	mockNatsClient := NewMockNATSClient(ctrl)

	er := &EventRouter{
		app:             mockApp,
		natsClient:      mockNatsClient,
		consumerManager: mockConsumerManager,
		bufferPool:      &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		logger:          mockLogger,
	}

	// Override getBufferFunc to return a failing buffer
	er.getBufferFunc = func() Buffer {
		return &failingBuffer{}
	}

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("test.event")
	event.SetSource("test")
	err := event.SetData(cloudevents.ApplicationJSON, map[string]string{"key": "value"})
	require.NoError(t, err)

	mockContext := &gofr.Context{
		Context: context.Background(),
	}

	err = er.routeEvent(mockContext, &event)
	require.Error(t, err, "Expected an error but got none")
	t.Logf("Actual error: %v", err)
	assert.Contains(t, err.Error(), "failed to encode event", "Error message does not contain the expected substring")
}

// failingBuffer is a custom buffer that always fails on Write.
type failingBuffer struct{}

func (*failingBuffer) Write([]byte) (n int, err error) {
	fmt.Println("Forced write error")
	return 0, errBufferForcedWriteError
}

func (*failingBuffer) Bytes() []byte {
	return []byte{}
}

func (*failingBuffer) Reset() {}
