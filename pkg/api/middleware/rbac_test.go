package middleware

import (
	"testing"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthenticateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := customctx.NewMockContext(ctrl)

	tests := []struct {
		name           string
		setupMocks     func()
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "Valid API Key",
			setupMocks: func() {
				apiKey := "valid_key"
				tenantID := uuid.New()
				customerID := uuid.New()

				// Set up mock expectations for API key retrieval
				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)

				// Set up mock expectations for claims
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(tenantID, true)
				mockContext.EXPECT().GetUUIDClaim("customer_id").Return(customerID, true)

				// Set up mock expectations for setting claims
				mockContext.EXPECT().SetClaim("api_key", apiKey)
				mockContext.EXPECT().SetClaim("tenant_id", tenantID)
				mockContext.EXPECT().SetClaim("customer_id", customerID)
			},
			expectedResult: "success",
			expectedError:  nil,
		},
		{
			name: "Missing API Key",
			setupMocks: func() {
				mockContext.EXPECT().GetAPIKey().Return("", false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing API Key"),
		},
		{
			name: "Missing Tenant or Customer ID Claims",
			setupMocks: func() {
				apiKey := "valid_key"

				// Set up mock expectations for API key retrieval
				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)

				// Mock missing tenant or customer ID claims
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.Nil, false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing tenant ID"),
		},
		{
			name: "Invalid API Key",
			setupMocks: func() {
				apiKey := "invalid_key"

				// Set up mock expectations for API key retrieval
				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)

				// Mock missing tenant and customer ID claims
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.Nil, false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing tenant ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			middleware := AuthenticateAPIKey(func(cc customctx.Context) (interface{}, error) {
				return "success", nil
			})

			result, err := middleware(mockContext)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
