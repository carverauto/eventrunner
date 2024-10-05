// Package eventrunner pkg/eventrunner/consumer_manager.go
package eventrunner

import (
	"fmt"
	"strings"
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

// ConsumerErrors is a custom error type that holds multiple errors.
type ConsumerErrors struct {
	Errors []error
}

// Error implements the error interface for ConsumerErrors.
func (ce *ConsumerErrors) Error() string {
	if len(ce.Errors) == 0 {
		return "no errors occurred"
	}

	errorStrings := make([]string, len(ce.Errors))

	for i, err := range ce.Errors {
		errorStrings[i] = err.Error()
	}

	return fmt.Sprintf("errors occurred while consuming event: %s", strings.Join(errorStrings, "; "))
}

func (cm *ConsumerManager) ConsumeEvent(c *gofr.Context, event *cloudevents.Event) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if c == nil {
		return errNilContext
	}

	if event == nil {
		return errNilEvent
	}

	var consumerErrors ConsumerErrors

	for name, consumer := range cm.consumers {
		if consumer == nil {
			cm.logger.Logf("EventConsumer %s is nil, skipping", name)
			continue
		}

		if err := consumer.ConsumeEvent(c, event); err != nil {
			cm.logger.Errorf("EventConsumer %s failed: %v", name, err)
			consumerErrors.Errors = append(consumerErrors.Errors, fmt.Errorf("consumer %s failed: %w", name, err))
		}
	}

	if len(consumerErrors.Errors) > 0 {
		return &consumerErrors
	}

	return nil
}
