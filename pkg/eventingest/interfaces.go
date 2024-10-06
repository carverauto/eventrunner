package eventingest

import (
	"context"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServiceClient is the gRPC client interface.
type ServiceClient interface {
	IngestEvent(ctx context.Context, in *IngestEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// EventForwarder interface defines the contract for forwarding events.
type EventForwarder interface {
	ForwardEvent(c *gofr.Context, tenantID, customerID uuid.UUID, eventData []byte) error
}
