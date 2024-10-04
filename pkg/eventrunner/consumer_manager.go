// Package eventrunner pkg/eventrunner/consumer_manager.go
package eventrunner

import (
	"fmt"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/pubsub"
)

type ConsumerManager struct {
	app       AppInterface
	consumers map[string]EventConsumer
	mu        sync.RWMutex
	logger    pubsub.Logger
}

func NewConsumerManager(app AppInterface, logger pubsub.Logger) *ConsumerManager {
	return &ConsumerManager{
		app:       app,
		consumers: make(map[string]EventConsumer),
		logger:    logger,
	}
}

func (cm *ConsumerManager) AddConsumer(name string, consumer EventConsumer) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.consumers[name] = consumer
}

func (cm *ConsumerManager) ConsumeEvent(c *gofr.Context, event *cloudevents.Event) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c == nil {
		return fmt.Errorf("nil context provided to ConsumerManager")
	}
	if event == nil {
		return fmt.Errorf("nil event provided to ConsumerManager")
	}

	var errors []error
	for name, consumer := range cm.consumers {
		if consumer == nil {
			cm.logger.Logf("EventConsumer %s is nil, skipping", name)
			continue
		}
		if err := consumer.ConsumeEvent(c, event); err != nil {
			cm.logger.Errorf("EventConsumer %s failed: %v", name, err)
			errors = append(errors, fmt.Errorf("consumer %s failed: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while consuming event: %v", errors)
	}

	return nil
}
