package context

import (
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type Context struct {
	*gofr.Context
	claims map[string]interface{}
}

// NewCustomContext creates a new Context.
func NewCustomContext(c *gofr.Context) *Context {
	return &Context{
		Context: c,
		claims:  make(map[string]interface{}),
	}
}

// SetClaim sets the claim value.
func (c *Context) SetClaim(key string, value interface{}) {
	c.claims[key] = value
}

// GetClaim returns the claim value.
func (c *Context) GetClaim(key string) (interface{}, bool) {
	value, ok := c.claims[key]

	return value, ok
}

// GetStringClaim returns the claim value as a string.
func (c *Context) GetStringClaim(key string) (string, bool) {
	value, ok := c.claims[key]
	if !ok {
		return "", false
	}

	strValue, ok := value.(string)

	return strValue, ok
}

// GetUUIDClaim returns the claim value as a UUID.
func (c *Context) GetUUIDClaim(key string) (uuid.UUID, bool) {
	value, ok := c.claims[key]
	if !ok {
		return uuid.UUID{}, false
	}

	strValue, ok := value.(uuid.UUID)

	return strValue, ok
}

// GetAPIKey returns the API key from the context.
func (c *Context) GetAPIKey() (string, bool) {
	if apiKey, ok := c.Context.Request.Context().Value("APIKey").(string); ok {
		return apiKey, true
	}

	return "", false
}
