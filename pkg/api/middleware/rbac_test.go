// File: middleware/rbac_test.go
package middleware

import (
	"testing"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

func TestAuthenticateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMocks     func(*customctx.MockContext)
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "Valid API Key",
			setupMocks: func(mockContext *customctx.MockContext) {
				apiKey := "valid_key"
				tenantID := uuid.MustParse("26cd72f6-7ef6-4665-a39d-daaf3463760d")
				customerID := uuid.MustParse("75060cd8-b8b7-455b-868f-e300b979cc43")

				mockContext.EXPECT().GetAPIKey().Return(apiKey, true)
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(tenantID, true)
				mockContext.EXPECT().GetUUIDClaim("customer_id").Return(customerID, true)
				mockContext.EXPECT().SetClaim("api_key", apiKey)
				mockContext.EXPECT().SetClaim("tenant_id", tenantID)
				mockContext.EXPECT().SetClaim("customer_id", customerID)
				mockContext.EXPECT().Context().Return(&gofr.Context{})
			},
			expectedResult: "success",
			expectedError:  nil,
		},
		{
			name: "Missing API Key",
			setupMocks: func(mockContext *customctx.MockContext) {
				mockContext.EXPECT().GetAPIKey().Return("", false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing API Key"),
		},
		{
			name: "Missing Tenant ID",
			setupMocks: func(mockContext *customctx.MockContext) {
				mockContext.EXPECT().GetAPIKey().Return("valid_key", true)
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.Nil, false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing tenant ID"),
		},
		{
			name: "Missing Customer ID",
			setupMocks: func(mockContext *customctx.MockContext) {
				mockContext.EXPECT().GetAPIKey().Return("valid_key", true)
				mockContext.EXPECT().GetUUIDClaim("tenant_id").Return(uuid.New(), true)
				mockContext.EXPECT().GetUUIDClaim("customer_id").Return(uuid.Nil, false)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing customer ID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContext := customctx.NewMockContext(ctrl)
			tt.setupMocks(mockContext)

			nextHandler := handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
				return "success", nil
			})

			result, err := authenticateAPIKey(mockContext, nextHandler)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
