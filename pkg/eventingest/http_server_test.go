package eventingest

import (
	"errors"
	"testing"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

func TestHandleEvent(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*gomock.Controller) (customctx.Interface, EventForwarder)
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "Success",
			setupMocks: func(ctrl *gomock.Controller) (customctx.Interface, EventForwarder) {
				mockCtx := customctx.NewMockInterface(ctrl)
				mockEF := NewMockEventForwarder(ctrl)

				tenantID := uuid.New()
				customerID := uuid.New()
				eventData := []byte(`{"key": "value"}`)

				mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(tenantID, true)
				mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(customerID, true)
				mockCtx.EXPECT().Bind(gomock.Any()).SetArg(0, eventData).Return(nil)
				mockCtx.EXPECT().Context().Return(&gofr.Context{})

				mockEF.EXPECT().ForwardEvent(gomock.Any(), tenantID, customerID, eventData).Return(nil)

				return mockCtx, mockEF
			},
			expectedResult: map[string]string{"status": "received"},
			expectedError:  nil,
		},
		{
			name: "Missing tenant ID",
			setupMocks: func(ctrl *gomock.Controller) (customctx.Interface, EventForwarder) {
				mockCtx := customctx.NewMockInterface(ctrl)
				mockEF := NewMockEventForwarder(ctrl)

				mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.Nil, false)

				return mockCtx, mockEF
			},
			expectedResult: nil,
			expectedError:  NewAuthError("Missing tenant ID"),
		},
		{
			name: "Missing customer ID",
			setupMocks: func(ctrl *gomock.Controller) (customctx.Interface, EventForwarder) {
				mockCtx := customctx.NewMockInterface(ctrl)
				mockEF := NewMockEventForwarder(ctrl)

				mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.New(), true)
				mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(uuid.Nil, false)

				return mockCtx, mockEF
			},
			expectedResult: nil,
			expectedError:  NewAuthError("Missing customer ID"),
		},
		{
			name: "Invalid request body",
			setupMocks: func(ctrl *gomock.Controller) (customctx.Interface, EventForwarder) {
				mockCtx := customctx.NewMockInterface(ctrl)
				mockEF := NewMockEventForwarder(ctrl)

				mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.New(), true)
				mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(uuid.New(), true)
				mockCtx.EXPECT().Bind(gomock.Any()).Return(errors.New("invalid JSON"))

				return mockCtx, mockEF
			},
			expectedResult: nil,
			expectedError:  NewProcessingError("Invalid request body"),
		},
		{
			name: "Forward event failure",
			setupMocks: func(ctrl *gomock.Controller) (customctx.Interface, EventForwarder) {
				mockCtx := customctx.NewMockInterface(ctrl)
				mockEF := NewMockEventForwarder(ctrl)

				tenantID := uuid.New()
				customerID := uuid.New()
				eventData := []byte(`{"key": "value"}`)

				mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(tenantID, true)
				mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(customerID, true)
				mockCtx.EXPECT().Bind(gomock.Any()).SetArg(0, eventData).Return(nil)
				mockCtx.EXPECT().Context().Return(&gofr.Context{})

				mockEF.EXPECT().ForwardEvent(gomock.Any(), tenantID, customerID, eventData).Return(errors.New("forwarding failed"))

				return mockCtx, mockEF
			},
			expectedResult: nil,
			expectedError:  NewProcessingError("Failed to forward event"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx, mockEF := tt.setupMocks(ctrl)

			server := NewHTTPServer(&gofr.App{}, mockEF)
			result, err := server.HandleEvent(&mockCtx)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
