package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/middleware"
	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/carverauto/eventrunner/pkg/errors"
	"github.com/google/uuid"
	ory "github.com/ory/client-go"
	"go.mongodb.org/mongo-driver/bson"
	"gofr.dev/pkg/gofr"
)

type Handlers struct {
	OryClient *ory.APIClient
}

func NewHandlers(oryClient *ory.APIClient) *Handlers {
	return &Handlers{OryClient: oryClient}
}

func (h *Handlers) CreateSuperUser(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, errors.NewInvalidParamError("user data")
	}

	// Generate a new UUID for the tenant_id
	tenantID := uuid.New()

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: os.Getenv("ORY_SCHEMA_ID"),
		Traits: map[string]interface{}{
			"email":     user.Email,
			"roles":     []string{"superuser"},
			"tenant_id": tenantID.String(),
		},
	}

	// Debug logging
	payloadBytes, _ := json.MarshalIndent(identity, "", "  ")
	fmt.Printf("Payload being sent to Ory:\n%s\n", string(payloadBytes))

	createdIdentity, _, err := h.OryClient.IdentityAPI.CreateIdentity(c.Context).CreateIdentityBody(identity).Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create identity in Ory: %v", err))
	}

	// Create user in MongoDB
	user.UserID = uuid.New()
	user.OryID = createdIdentity.Id
	user.TenantID = tenantID
	user.Roles = []string{"superuser"}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert user into database")
	}

	return user, nil
}

func (h *Handlers) CreateTenant(c *gofr.Context) (interface{}, error) {
	var tenant models.Tenant
	if err := c.Bind(&tenant); err != nil {
		return nil, errors.NewInvalidParamError("tenant data")
	}

	claims, err := middleware.GetJWTClaims(c)
	if err != nil {
		return nil, err
	}

	if claims["role"] != "superuser" {
		return nil, errors.NewForbiddenError("Only superuser can create tenants")
	}

	tenant.TenantID = uuid.New()
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert tenant into database")
	}

	return tenant, nil
}

func (h *Handlers) CreateUser(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, errors.NewInvalidParamError("user data")
	}

	claims, err := middleware.GetJWTClaims(c)
	if err != nil {
		return nil, err
	}

	// Check if the requester has the right to create users
	requesterRole := claims["role"].(string)
	if requesterRole != "superuser" && requesterRole != "tenant_admin" {
		return nil, errors.NewForbiddenError("Insufficient permissions to create user")
	}

	// Use claims to set user data
	tenantID, ok := claims["tenantId"].(string)
	if !ok {
		return nil, errors.NewMissingParamError("tenantId")
	}
	user.TenantID, _ = uuid.Parse(tenantID)

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: os.Getenv("ORY_SCHEMA_ID"),
		Traits: map[string]interface{}{
			"email":    user.Email,
			"role":     user.Roles[0], // Assuming the first role is the primary role
			"tenantId": user.TenantID.String(),
		},
	}

	createdIdentity, _, err := h.OryClient.IdentityAPI.CreateIdentity(c.Context).CreateIdentityBody(identity).Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create identity in Ory: %v", err))
	}

	// Create user in MongoDB
	user.UserID = uuid.New()
	user.OryID = createdIdentity.Id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert user into database")
	}

	return user, nil
}

func (h *Handlers) GetAllUsers(c *gofr.Context) (interface{}, error) {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return nil, errors.NewInvalidParamError("tenant_id")
	}

	var users []models.User
	err = c.Mongo.Find(c, "users", bson.M{"tenant_id": tenantID}, &users)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to fetch users from database")
	}

	return users, nil
}
