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
