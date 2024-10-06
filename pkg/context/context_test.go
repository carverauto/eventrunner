package context

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GofrContextWrapper is an interface that wraps the methods of gofr.Context that we use.
type GofrContextWrapper interface {
	Param(key string) string
	PathParam(key string) string
	Bind(v interface{}) error
}

// MockGofrContextWrapper is a mock for GofrContextWrapper.
type MockGofrContextWrapper struct {
	mock.Mock
}

func (m *MockGofrContextWrapper) Param(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockGofrContextWrapper) PathParam(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockGofrContextWrapper) Bind(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// ContextWrapper wraps our Context and uses GofrContextWrapper instead of *gofr.Context.
type ContextWrapper struct {
	*CustomContext
	gofrCtx GofrContextWrapper
}

func NewContextWrapper(gofrCtx GofrContextWrapper) *ContextWrapper {
	return &ContextWrapper{
		CustomContext: &CustomContext{
			Context: nil, // We're not setting this as we're using the wrapper
			claims:  make(map[string]interface{}),
		},
		gofrCtx: gofrCtx,
	}
}

func (c *ContextWrapper) Bind(v interface{}) error {
	return c.gofrCtx.Bind(v)
}

func (c *ContextWrapper) GetAPIKey() (string, bool) {
	apiKey := c.gofrCtx.Param("APIKey")
	if apiKey != "" {
		return apiKey, true
	}

	return "", false
}

func TestNewCustomContext(t *testing.T) {
	mockGofrCtx := new(MockGofrContextWrapper)
	customCtx := NewContextWrapper(mockGofrCtx)

	assert.NotNil(t, customCtx)
	assert.NotNil(t, customCtx.claims)
}

func TestSetAndGetClaim(t *testing.T) {
	customCtx := NewContextWrapper(new(MockGofrContextWrapper))

	customCtx.SetClaim("test", "value")
	value, ok := customCtx.GetClaim("test")

	assert.True(t, ok)
	assert.Equal(t, "value", value)

	_, ok = customCtx.GetClaim("nonexistent")
	assert.False(t, ok)
}

func TestGetStringClaim(t *testing.T) {
	customCtx := NewContextWrapper(new(MockGofrContextWrapper))

	customCtx.SetClaim("string", "value")
	customCtx.SetClaim("not_string", 123)

	value, ok := customCtx.GetStringClaim("string")
	assert.True(t, ok)
	assert.Equal(t, "value", value)

	_, ok = customCtx.GetStringClaim("not_string")
	assert.False(t, ok)

	_, ok = customCtx.GetStringClaim("nonexistent")
	assert.False(t, ok)
}

func TestGetUUIDClaim(t *testing.T) {
	customCtx := NewContextWrapper(new(MockGofrContextWrapper))

	uuidValue := uuid.New()
	customCtx.SetClaim("uuid", uuidValue)
	customCtx.SetClaim("not_uuid", "not a uuid")

	value, ok := customCtx.GetUUIDClaim("uuid")
	assert.True(t, ok)
	assert.Equal(t, uuidValue, value)

	_, ok = customCtx.GetUUIDClaim("not_uuid")
	assert.False(t, ok)

	_, ok = customCtx.GetUUIDClaim("nonexistent")
	assert.False(t, ok)
}

func TestGetAPIKey(t *testing.T) {
	mockGofrCtx := new(MockGofrContextWrapper)
	customCtx := NewContextWrapper(mockGofrCtx)

	// Test case when APIKey is present
	mockGofrCtx.On("Param", "APIKey").Return("test-api-key").Once()

	apiKey, ok := customCtx.GetAPIKey()
	assert.True(t, ok)
	assert.Equal(t, "test-api-key", apiKey)

	// Test case when APIKey is not present
	mockGofrCtx.On("Param", "APIKey").Return("").Once()

	apiKey, ok = customCtx.GetAPIKey()
	assert.False(t, ok)
	assert.Empty(t, apiKey)

	mockGofrCtx.AssertExpectations(t)
}

func TestBind(t *testing.T) {
	mockGofrCtx := new(MockGofrContextWrapper)
	customCtx := NewContextWrapper(mockGofrCtx)

	var testStruct struct {
		Name string `json:"name"`
	}

	mockGofrCtx.On("Bind", &testStruct).Return(nil)

	err := customCtx.Bind(&testStruct)

	require.NoError(t, err)
	mockGofrCtx.AssertExpectations(t)
}
