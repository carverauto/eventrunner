// Package eventrunner pkg/eventrunner/consumer_manager.go
package eventrunner

import (
	"fmt"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
)

type ConsumerManager struct {
	app       *gofr.App
	consumers map[string]Consumer
	mu        sync.RWMutex
}

func NewConsumerManager(app *gofr.App) *ConsumerManager {
	return &ConsumerManager{
		app:       app,
		consumers: make(map[string]Consumer),
	}
}

func (cm *ConsumerManager) AddConsumer(name string, consumer Consumer) {
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
			cm.app.Logger().Warnf("Consumer %s is nil, skipping", name)
			continue
		}
		if err := consumer.ConsumeEvent(c, event); err != nil {
			cm.app.Logger().Errorf("Consumer %s failed: %v", name, err)
			errors = append(errors, fmt.Errorf("consumer %s failed: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while consuming event: %v", errors)
	}

	return nil
}
