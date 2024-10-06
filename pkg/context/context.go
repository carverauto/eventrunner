package context

import (
	"net/http"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type CustomContext struct {
	gofrContext *gofr.Context
	claims      map[string]interface{}
	headers     http.Header
}

// NewCustomContext creates a new Context.
func NewCustomContext(c *gofr.Context) *CustomContext {
	return &CustomContext{
		gofrContext: c,
		claims:      make(map[string]interface{}),
		headers:     http.Header{},
	}
}

// Adapter is a test helper that implements customctx.Context.
type Adapter struct {
	MockContext *MockContext
	GofrContext *gofr.Context
}

func (a *Adapter) SetClaim(key string, value interface{}) {
	a.MockContext.SetClaim(key, value)
}

func (a *Adapter) GetClaim(key string) (interface{}, bool) {
	return a.MockContext.GetClaim(key)
}

func (a *Adapter) GetStringClaim(key string) (string, bool) {
	return a.MockContext.GetStringClaim(key)
}

func (a *Adapter) GetUUIDClaim(key string) (uuid.UUID, bool) {
	return a.MockContext.GetUUIDClaim(key)
}

func (a *Adapter) GetAPIKey() (string, bool) {
	return a.MockContext.GetAPIKey()
}

func (a *Adapter) Bind(v interface{}) error {
	return a.MockContext.Bind(v)
}

func (a *Adapter) Context() *gofr.Context {
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

// SetHeader sets an HTTP header.
func (c *CustomContext) SetHeader(key, value string) {
	c.headers.Set(key, value)
}

// GetHeader retrieves an HTTP header value.
func (c *CustomContext) GetHeader(key string) string {
	return c.headers.Get(key)
}

// Headers returns all HTTP headers.
func (c *CustomContext) Headers() http.Header {
	return c.headers
}
