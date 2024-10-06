package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/carverauto/eventrunner/pkg/config"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
)

// MockRequest for use in testing, similar to what you used in another test file.
type MockRequest struct {
	ctx    context.Context
	params map[string][]string // Adjusted to hold slices of strings
	body   []byte
	header http.Header
}

func (r *MockRequest) Context() context.Context {
	return r.ctx
}

func (r *MockRequest) Param(key string) string {
	if vals, ok := r.params[key]; ok && len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (r *MockRequest) Params(key string) []string {
	return r.params[key]
}

func (r *MockRequest) PathParam(key string) string {
	if vals, ok := r.params[key]; ok && len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (r *MockRequest) Bind(i interface{}) error {
	return json.Unmarshal(r.body, i)
}

func (*MockRequest) HostName() string {
	return "localhost"
}

func (r *MockRequest) Header() http.Header {
	return r.header
}

const authorizationKey contextKey = "Authorization"

func TestJWTMiddleware_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVerifier := NewMockIDTokenVerifier(ctrl)
	mockToken := NewMockToken(ctrl)

	jwtMiddleware := &JWTMiddleware{
		verifier: mockVerifier,
		config:   &config.OAuthConfig{},
	}

	tests := []struct {
		name           string
		setupRequest   func(*MockRequest)
		setupMocks     func()
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "Valid Token",
			setupRequest: func(req *MockRequest) {
				req.ctx = context.WithValue(context.Background(), authorizationKey, "Bearer valid_token")
			},
			setupMocks: func() {
				mockVerifier.EXPECT().Verify(gomock.Any(), "valid_token").Return(mockToken, nil).Times(1)
				mockToken.EXPECT().Claims(gomock.Any()).DoAndReturn(func(v interface{}) error {
					claims := struct {
						TenantID   string `json:"tenant_id"`
						CustomerID string `json:"customer_id"`
					}{
						TenantID:   "test-tenant-id",
						CustomerID: "test-customer-id",
					}
					claimsBytes, _ := json.Marshal(claims)
					return json.Unmarshal(claimsBytes, v)
				}).Times(1)
			},
			expectedResult: "success",
			expectedError:  nil,
		},
		{
			name:           "Missing Authorization Header",
			setupRequest:   func(*MockRequest) {}, // No Authorization set in context
			setupMocks:     func() {},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Missing or invalid authorization header"),
		},
		{
			name: "Invalid Authorization Header",
			setupRequest: func(req *MockRequest) {
				req.ctx = context.WithValue(context.Background(), authorizationKey, "InvalidHeader")
			},
			setupMocks:     func() {},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Invalid authorization header format"),
		},
		{
			name: "Invalid Token",
			setupRequest: func(req *MockRequest) {
				req.ctx = context.WithValue(context.Background(), authorizationKey, "Bearer invalid_token")
			},
			setupMocks: func() {
				mockVerifier.EXPECT().Verify(gomock.Any(), "invalid_token").Return(nil, eventingest.NewAuthError("Invalid token")).Times(1)
			},
			expectedResult: nil,
			expectedError:  eventingest.NewAuthError("Invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := &MockRequest{
				ctx:    context.Background(),
				params: make(map[string][]string),
				header: http.Header{},
			}
			tt.setupRequest(mockRequest)

			// Create a new GoFr container
			c, _ := container.NewMockContainer(t)

			// Create the GoFr context
			gofrCtx := &gofr.Context{
				Context:   mockRequest.Context(),
				Request:   mockRequest,
				Container: c,
			}

			tt.setupMocks()

			handler := jwtMiddleware.Validate(func(customctx.Context) (interface{}, error) {
				return "success", nil
			})

			result, err := handler(gofrCtx)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
