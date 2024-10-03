// Package eventrunner pkg/eventrunner/consumer_manager.go
package eventrunner

import (
	"context"
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

func (cm *ConsumerManager) ConsumeEvent(ctx context.Context, event *cloudevents.Event) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	gofrCtx := &gofr.Context{
		Context: ctx,
	}

	var errors []error
	for name, consumer := range cm.consumers {
		if err := consumer.ConsumeEvent(gofrCtx, event); err != nil {
			errors = append(errors, fmt.Errorf("consumer %s failed: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while consuming event: %v", errors)
	}

	return nil
}
