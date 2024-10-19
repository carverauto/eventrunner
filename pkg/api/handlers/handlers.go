package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/middleware"
)

type TenantHandler struct{}

// HandlerFunc is an adapter to allow the use of ordinary functions as Handlers.
type HandlerFunc func(*gofr.Context) (interface{}, error)

// Handle calls f(c).
func (f HandlerFunc) Handle(c *gofr.Context) (interface{}, error) {
	return f(c)
}

// Middleware defines the standard middleware signature.
type Middleware func(Handler) Handler

func (*TenantHandler) Create(c *gofr.Context) (interface{}, error) {
	var tenant models.Tenant
	if err := c.Bind(&tenant); err != nil {
		return nil, err
	}

	tenant.TenantID = uuid.New()
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	_, err := c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (*TenantHandler) GetAll(c *gofr.Context) (interface{}, error) {
	var tenants []models.Tenant
	err := c.Mongo.Find(c, "tenants", bson.M{}, &tenants)

	return tenants, err
}

type UserHandler struct{}

func (h *UserHandler) Create(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, err
	}

	// Retrieve JWT claims
	claimData := c.Context.Value(middleware.JWTClaim)
	claims, ok := claimData.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claim data type")
	}

	// Use claims to set user data
	if tenantID, ok := claims["tenant_id"].(string); ok {
		user.TenantID, _ = uuid.Parse(tenantID)
	}
	if customerID, ok := claims["customer_id"].(string); ok {
		user.CustomerID, _ = uuid.Parse(customerID)
	}

	user.UserID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (*UserHandler) GetAll(c *gofr.Context) (interface{}, error) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = c.Mongo.Find(c, "users", bson.M{"tenant_id": tenantID}, &users)

	return users, err
}

func (h *UserHandler) CreateSuperUser(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, err
	}

	// Create user in Ory Hydra (pseudo-code)
	hydraUser, err := h.hydraClient.CreateUser(user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	// Create user in MongoDB
	user.UserID = hydraUser.ID
	user.Roles = []string{"superuser"}
	result, err := c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) CreateUser(c *gofr.Context) (interface{}, error) {
	// Check if the requester is a superuser
	if !isSuperUser(c) {
		return nil, errors.New("unauthorized")
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, err
	}

	// Create user in Ory Hydra (pseudo-code)
	hydraUser, err := h.hydraClient.CreateUser(user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	// Create user in MongoDB
	user.UserID = hydraUser.ID
	result, err := c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
