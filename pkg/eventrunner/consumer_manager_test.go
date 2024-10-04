package eventrunner

import (
	"fmt"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/logging"
)

func TestNewConsumerManager(t *testing.T) {
	app := &gofr.App{}
	logger := logging.NewMockLogger(logging.DEBUG)
	cm := NewConsumerManager(app, logger)

	assert.NotNil(t, cm)
	assert.Equal(t, app, cm.app)
	assert.Equal(t, logger, cm.logger)
	assert.Empty(t, cm.consumers)
}

func TestConsumerManager_AddConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := &gofr.App{}
	logger := logging.NewMockLogger(logging.DEBUG)
	cm := NewConsumerManager(app, logger)

	consumer := NewMockConsumer(ctrl)
	cm.AddConsumer("test", consumer)

	assert.Len(t, cm.consumers, 1)
	assert.Equal(t, consumer, cm.consumers["test"])
}

func TestConsumerManager_ConsumeEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := &gofr.App{}
	logger := logging.NewMockLogger(logging.DEBUG)
	cm := NewConsumerManager(app, logger)

	mockContext := &gofr.Context{}
	mockEvent := cloudevents.NewEvent()

	t.Run("No consumers", func(t *testing.T) {
		err := cm.ConsumeEvent(mockContext, &mockEvent)
		assert.NoError(t, err)
	})

	t.Run("Single successful consumer", func(t *testing.T) {
		consumer := NewMockConsumer(ctrl)
		consumer.EXPECT().ConsumeEvent(mockContext, &mockEvent).Return(nil)
		cm.AddConsumer("test", consumer)

		err := cm.ConsumeEvent(mockContext, &mockEvent)
		assert.NoError(t, err)

		cm.consumers = make(map[string]Consumer) // Reset consumers
	})

	t.Run("Multiple consumers, one fails", func(t *testing.T) {
		successConsumer := NewMockConsumer(ctrl)
		successConsumer.EXPECT().ConsumeEvent(mockContext, &mockEvent).Return(nil)

		failConsumer := NewMockConsumer(ctrl)
		failConsumer.EXPECT().ConsumeEvent(mockContext, &mockEvent).Return(fmt.Errorf("consumer failed"))

		cm.AddConsumer("success", successConsumer)
		cm.AddConsumer("fail", failConsumer)

		err := cm.ConsumeEvent(mockContext, &mockEvent)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "consumer fail failed: consumer failed")

		cm.consumers = make(map[string]Consumer) // Reset consumers
	})

	t.Run("Nil consumer", func(t *testing.T) {
		cm.AddConsumer("nil", nil)

		err := cm.ConsumeEvent(mockContext, &mockEvent)
		assert.NoError(t, err)

		cm.consumers = make(map[string]Consumer) // Reset consumers
	})

	t.Run("Nil context", func(t *testing.T) {
		err := cm.ConsumeEvent(nil, &mockEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil context provided")
	})

	t.Run("Nil event", func(t *testing.T) {
		err := cm.ConsumeEvent(mockContext, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil event provided")
	})
}
