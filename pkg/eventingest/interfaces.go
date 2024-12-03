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

package eventingest

import (
	"context"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -destination=mock_eventingest.go -package=eventingest -source=./interfaces.go ServiceClient,EventForwarder

// ServiceClient is the gRPC client interface.
type ServiceClient interface {
	IngestEvent(ctx context.Context, in *IngestEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// EventForwarder interface defines the contract for forwarding events.
type EventForwarder interface {
	ForwardEvent(c *gofr.Context, tenantID, customerID uuid.UUID, eventData []byte) error
}
