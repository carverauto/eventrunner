package context

import (
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type CustomContext struct {
	gofrContext *gofr.Context
	claims      map[string]interface{}
}

// NewCustomContext creates a new Context.
func NewCustomContext(c *gofr.Context) *CustomContext {
	return &CustomContext{
		gofrContext: c,
		claims:      make(map[string]interface{}),
	}
}

// ContextAdapter is a test helper that implements customctx.Context
type ContextAdapter struct {
	MockContext *MockContext
	GofrContext *gofr.Context
}

func (a *ContextAdapter) SetClaim(key string, value interface{}) {
	a.MockContext.SetClaim(key, value)
}

func (a *ContextAdapter) GetClaim(key string) (interface{}, bool) {
	return a.MockContext.GetClaim(key)
}

func (a *ContextAdapter) GetStringClaim(key string) (string, bool) {
	return a.MockContext.GetStringClaim(key)
}

func (a *ContextAdapter) GetUUIDClaim(key string) (uuid.UUID, bool) {
	return a.MockContext.GetUUIDClaim(key)
}

func (a *ContextAdapter) GetAPIKey() (string, bool) {
	return a.MockContext.GetAPIKey()
}

func (a *ContextAdapter) Bind(v interface{}) error {
	return a.MockContext.Bind(v)
}

func (a *ContextAdapter) Context() *gofr.Context {
	return a.GofrContext
}

// SetClaim sets the claim value.
func (c *CustomContext) SetClaim(key string, value interface{}) {
	c.claims[key] = value
}

// GetClaim returns the claim value.
func (c *CustomContext) GetClaim(key string) (interface{}, bool) {
	value, ok := c.claims[key]

	return value, ok
}

// GetStringClaim returns the claim value as a string.
func (c *CustomContext) GetStringClaim(key string) (string, bool) {
	value, ok := c.claims[key]
	if !ok {
		return "", false
	}

	strValue, ok := value.(string)

	return strValue, ok
}

// GetUUIDClaim returns the claim value as a UUID.
func (c *CustomContext) GetUUIDClaim(key string) (uuid.UUID, bool) {
	value, ok := c.claims[key]
	if !ok {
		return uuid.UUID{}, false
	}

	strValue, ok := value.(uuid.UUID)

	return strValue, ok
}

// GetAPIKey returns the API key from the context.
func (c *CustomContext) GetAPIKey() (string, bool) {
	if apiKey, ok := c.gofrContext.Request.Context().Value("APIKey").(string); ok {
		return apiKey, true
	}

	return "", false
}

// Bind binds the input to the given value.
func (c *CustomContext) Bind(v interface{}) error {
	return c.gofrContext.Bind(v)
}

// Context returns the underlying gofr.Context.
func (c *CustomContext) Context() *gofr.Context {
	return c.gofrContext
}
