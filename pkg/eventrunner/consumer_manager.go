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

		// go func() {
		if err := consumer.ConsumeEvent(c, event); err != nil {
			cm.logger.Errorf("EventConsumer %s failed: %v", name, err)
			consumerErrors.Errors = append(consumerErrors.Errors, fmt.Errorf("consumer %s failed: %w", name, err))
		}
		// }()
	}

	if len(consumerErrors.Errors) > 0 {
		return &consumerErrors
	}

	return nil
}
