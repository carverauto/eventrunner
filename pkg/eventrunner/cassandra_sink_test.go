package eventrunner

import (
	"context"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

type MockCassandra struct {
	ExecFunc            func(stmt string, values ...any) error
	QueryFunc           func(dest any, stmt string, values ...any) error
	ExecCASFunc         func(dest any, stmt string, values ...any) (bool, error)
	NewBatchFunc        func(name string, batchType int) error
	BatchQueryFunc      func(name, stmt string, values ...any)
	ExecuteBatchFunc    func(name string) error
	ExecuteBatchCASFunc func(name string, dest ...any) (bool, error)
	HealthCheckFunc     func(context.Context) (any, error)
}

func (m *MockCassandra) Exec(stmt string, values ...any) error {
	if m.ExecFunc != nil {
		return m.ExecFunc(stmt, values...)
	}
	return nil
}

func (m *MockCassandra) Query(dest any, stmt string, values ...any) error {
	if m.QueryFunc != nil {
		return m.QueryFunc(dest, stmt, values...)
	}
	return nil
}

func (m *MockCassandra) ExecCAS(dest any, stmt string, values ...any) (bool, error) {
	if m.ExecCASFunc != nil {
		return m.ExecCASFunc(dest, stmt, values...)
	}
	return true, nil
}

func (m *MockCassandra) NewBatch(name string, batchType int) error {
	if m.NewBatchFunc != nil {
		return m.NewBatchFunc(name, batchType)
	}
	return nil
}

func (m *MockCassandra) BatchQuery(name, stmt string, values ...any) error {
	if m.BatchQueryFunc != nil {
		m.BatchQueryFunc(name, stmt, values...)
	}
	return nil
}

func (m *MockCassandra) ExecuteBatch(name string) error {
	if m.ExecuteBatchFunc != nil {
		return m.ExecuteBatchFunc(name)
	}
	return nil
}

func (m *MockCassandra) ExecuteBatchCAS(name string, dest ...any) (bool, error) {
	if m.ExecuteBatchCASFunc != nil {
		return m.ExecuteBatchCASFunc(name, dest...)
	}
	return true, nil
}

func (m *MockCassandra) HealthCheck(ctx context.Context) (any, error) {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return "OK", nil
}

// MockContext attempts to mimic the structure of gofr.Context
type MockContext struct {
	*container.Container
}

func NewMockContext() *gofr.Context {
	mockContainer := &container.Container{}
	return &gofr.Context{Container: mockContainer}
}

func TestCassandraEventSink_ConsumeEvent(t *testing.T) {
	sink := NewCassandraEventSink()

	event := cloudevents.NewEvent()
	event.SetID("123")
	event.SetSource("test-source")
	event.SetType("test-type")
	event.SetSubject("test-subject")
	event.SetTime(time.Now())
	event.SetDataContentType("application/json")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"key": "value"})

	mockCassandra := &MockCassandra{}
	mockContext := NewMockContext()
	mockContext.Container.Cassandra = mockCassandra

	err := sink.ConsumeEvent(mockContext, &event)
	require.NoError(t, err)
}

func TestCassandraEventSink_ConsumeEvent_NilContext(t *testing.T) {
	sink := NewCassandraEventSink()
	err := sink.ConsumeEvent(nil, &cloudevents.Event{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil context provided")
}

func TestCassandraEventSink_ConsumeEvent_NilEvent(t *testing.T) {
	sink := NewCassandraEventSink()
	mockContext := NewMockContext()
	err := sink.ConsumeEvent(mockContext, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil event provided")
}

func TestCassandraEventSink_ConsumeEvent_NilCassandra(t *testing.T) {
	sink := NewCassandraEventSink()
	mockContext := NewMockContext()
	mockContext.Container.Cassandra = nil
	err := sink.ConsumeEvent(mockContext, &cloudevents.Event{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cassandra client is nil")
}

func TestCassandraEventSink_ConsumeEvent_InvalidJSON(t *testing.T) {
	sink := NewCassandraEventSink()
	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, []byte("invalid json"))
	mockCassandra := &MockCassandra{}
	mockContext := NewMockContext()
	mockContext.Container.Cassandra = mockCassandra
	err := sink.ConsumeEvent(mockContext, &event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event data is not valid JSON")
}

func TestCassandraEventSink_ConsumeEvent_CassandraError(t *testing.T) {
	sink := NewCassandraEventSink()
	event := cloudevents.NewEvent()
	event.SetID("123")
	event.SetSource("test-source")
	event.SetType("test-type")
	event.SetSubject("test-subject")
	event.SetTime(time.Now())
	event.SetDataContentType("application/json")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"key": "value"})

	mockCassandra := &MockCassandra{
		ExecFunc: func(stmt string, values ...any) error {
			return assert.AnError
		},
	}
	mockContext := NewMockContext()
	mockContext.Container.Cassandra = mockCassandra

	err := sink.ConsumeEvent(mockContext, &event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert event into Cassandra")
}
