package eventingest

import (
	"context"

	"github.com/google/uuid"

	"gofr.dev/pkg/gofr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// IngestEventRequest is the request structure for the IngestEvent RPC.
type IngestEventRequest struct {
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	EventData  []byte
}

// GRPCEventForwarder implements the EventForwarder interface.
type GRPCEventForwarder struct {
	client ServiceClient
}

// NewGRPCEventForwarder creates a new GRPCEventForwarder.
func NewGRPCEventForwarder(conn *grpc.ClientConn) *GRPCEventForwarder {
	return &GRPCEventForwarder{
		client: NewEventIngestServiceClient(conn),
	}
}

// NewEventIngestServiceClient creates a new ServiceClient.
func NewEventIngestServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &eventIngestServiceClient{cc}
}

type eventIngestServiceClient struct {
	cc grpc.ClientConnInterface
}

// IngestEvent calls the IngestEvent RPC.
func (c *eventIngestServiceClient) IngestEvent(
	ctx context.Context, in *IngestEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)

	if err := c.cc.Invoke(ctx, "/eventingest.IngestService/IngestEvent", in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

// ForwardEvent implements the EventForwarder interface.
func (f *GRPCEventForwarder) ForwardEvent(c *gofr.Context, tenantID, customerID uuid.UUID, eventData []byte) error {
	_, err := f.client.IngestEvent(c, &IngestEventRequest{
		TenantID:   tenantID,
		CustomerID: customerID,
		EventData:  eventData,
	})

	return err
}
