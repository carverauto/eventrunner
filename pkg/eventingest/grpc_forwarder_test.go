package eventingest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGRPCEventForwarder_ForwardEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockServiceClient(ctrl)
	forwarder := &GRPCEventForwarder{
		client: mockClient,
	}

	ctx := &gofr.Context{}
	tenantID := uuid.New()
	customerID := uuid.New()
	eventData := []byte("test event data")

	expectedRequest := &IngestEventRequest{
		TenantID:   tenantID,
		CustomerID: customerID,
		EventData:  eventData,
	}

	// Test successful forward
	mockClient.EXPECT().
		IngestEvent(gomock.Any(), gomock.Eq(expectedRequest), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Times(1)

	err := forwarder.ForwardEvent(ctx, tenantID, customerID, eventData)
	require.NoError(t, err)

	// Test error case
	mockClient.EXPECT().
		IngestEvent(gomock.Any(), gomock.Eq(expectedRequest), gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	err = forwarder.ForwardEvent(ctx, tenantID, customerID, eventData)
	assert.Error(t, err)
}
