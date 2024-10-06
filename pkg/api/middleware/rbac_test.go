package middleware

import (
	"testing"

	"github.com/carverauto/eventrunner/pkg/eventingest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthenticateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContext := NewMockCustomContext(ctrl)

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

				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)
				mockContext.EXPECT().FindAPIKey(apiKey).Return(tenantID, customerID, nil)
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
			name: "Invalid API Key",
			setupMocks: func() {
				apiKey := "invalid_key"
				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)
				mockContext.EXPECT().FindAPIKey(apiKey).Return(uuid.Nil, uuid.Nil, eventingest.NewAuthError("Invalid API Key"))
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Invalid API Key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			middleware := AuthenticateAPIKey(func(cc CustomContext) (interface{}, error) {
				return "success", nil
			})

			result, err := middleware(mockContext)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
