package eventingest

import (
	"testing"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

type testCase struct {
	name           string
	setupMocks     func(*gomock.Controller) (customctx.Context, EventForwarder)
	expectedResult interface{}
	expectedError  error
}

func TestHandleEvent(t *testing.T) {
	tests := []testCase{
		{
			name:           "Success",
			setupMocks:     setupSuccessCase,
			expectedResult: map[string]string{"status": "received"},
			expectedError:  nil,
		},
		{
			name:           "Missing tenant ID",
			setupMocks:     setupMissingTenantIDCase,
			expectedResult: nil,
			expectedError:  NewAuthError("Missing tenant ID"),
		},
		{
			name:           "Missing customer ID",
			setupMocks:     setupMissingCustomerIDCase,
			expectedResult: nil,
			expectedError:  NewAuthError("Missing customer ID"),
		},
		{
			name:           "Invalid request body",
			setupMocks:     setupInvalidRequestBodyCase,
			expectedResult: nil,
			expectedError:  NewProcessingError("Invalid request body"),
		},
		{
			name:           "Forward event failure",
			setupMocks:     setupForwardEventFailureCase,
			expectedResult: nil,
			expectedError:  NewProcessingError("Failed to forward event"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}

func runTestCase(t *testing.T, tc testCase) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCtx, mockEF := tc.setupMocks(ctrl)

	server := NewHTTPServer(&gofr.App{}, mockEF)
	result, err := server.HandleEvent(mockCtx)

	assert.Equal(t, tc.expectedResult, result)
	assert.Equal(t, tc.expectedError, err)
}

func setupSuccessCase(ctrl *gomock.Controller) (customctx.Context, EventForwarder) {
	mockCtx := customctx.NewMockContext(ctrl)
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
}

func setupMissingTenantIDCase(ctrl *gomock.Controller) (customctx.Context, EventForwarder) {
	mockCtx := customctx.NewMockContext(ctrl)
	mockEF := NewMockEventForwarder(ctrl)

	mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.Nil, false)

	return mockCtx, mockEF
}

func setupMissingCustomerIDCase(ctrl *gomock.Controller) (customctx.Context, EventForwarder) {
	mockCtx := customctx.NewMockContext(ctrl)
	mockEF := NewMockEventForwarder(ctrl)

	mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.New(), true)
	mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(uuid.Nil, false)

	return mockCtx, mockEF
}

func setupInvalidRequestBodyCase(ctrl *gomock.Controller) (customctx.Context, EventForwarder) {
	mockCtx := customctx.NewMockContext(ctrl)
	mockEF := NewMockEventForwarder(ctrl)

	mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.New(), true)
	mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(uuid.New(), true)
	mockCtx.EXPECT().Bind(gomock.Any()).Return(errInvalidJSON)

	return mockCtx, mockEF
}

func setupForwardEventFailureCase(ctrl *gomock.Controller) (customctx.Context, EventForwarder) {
	mockCtx := customctx.NewMockContext(ctrl)
	mockEF := NewMockEventForwarder(ctrl)

	tenantID := uuid.New()
	customerID := uuid.New()
	eventData := []byte(`{"key": "value"}`)

	mockCtx.EXPECT().GetUUIDClaim("tenant_id").Return(tenantID, true)
	mockCtx.EXPECT().GetUUIDClaim("customer_id").Return(customerID, true)
	mockCtx.EXPECT().Bind(gomock.Any()).SetArg(0, eventData).Return(nil)
	mockCtx.EXPECT().Context().Return(&gofr.Context{})

	mockEF.EXPECT().ForwardEvent(gomock.Any(), tenantID, customerID, eventData).Return(errForwardFail)

	return mockCtx, mockEF
}
