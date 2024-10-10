package handlers

import (
	"context"

	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/google/uuid"
	ory "github.com/ory/client-go"
	"go.mongodb.org/mongo-driver/bson"
	"gofr.dev/pkg/gofr"
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

func (*TenantHandler) Create(c *gofr.Context) (models.Tenant, error) {
	var tenant models.Tenant
	if err := c.Bind(&tenant); err != nil {
		return models.Tenant{}, err
	}

	result, err := c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return models.Tenant{}, err
	}

	tenant.ID = result.(uuid.UUID)

	return tenant, nil
}

func (*TenantHandler) GetAll(c *gofr.Context) (interface{}, error) {
	var tenants []models.Tenant
	err := c.Mongo.Find(c, "tenants", bson.M{}, &tenants)

	return tenants, err
}

type UserHandler struct {
	OryClient *ory.APIClient
}

func (h *UserHandler) Login(c *gofr.Context) (interface{}, error) {
	flow, _, err := h.OryClient.FrontendAPI.CreateNativeLoginFlow(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	// Return the login flow to the client
	return flow, nil
}

func (h *UserHandler) SubmitLogin(c *gofr.Context) (interface{}, error) {
	var loginData ory.UpdateLoginFlowBody
	if err := c.Bind(&loginData); err != nil {
		return nil, err
	}

	response, _, err := h.OryClient.FrontendAPI.UpdateLoginFlow(context.Background()).
		Flow(c.Param("flow")).
		UpdateLoginFlowBody(loginData).
		Execute()
	if err != nil {
		return nil, err
	}

	// Return the session token or redirect URL
	return response, nil
}

/*
func (h *UserHandler) GetUserInfo(c *gofr.Context) (interface{}, error) {
	session, _, err := h.OryClient.FrontendAPI.ToSession(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	traits := session.Identity.Traits.(map[string]interface{})
	tenantID, _ := uuid.Parse(traits["tenant_id"].(string))
	customerID, _ := uuid.Parse(traits["customer_id"].(string))

	// Use these IDs for authorization or to fetch additional user data from MongoDB
	// ...

	return session.Identity, nil
}
*/

func (h *UserHandler) Create(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, err
	}

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: "default",
		Traits: map[string]interface{}{
			"email":       user.Email,
			"tenant_id":   user.TenantID.String(),
			"customer_id": user.CustomerID.String(),
		},
	}

	createdIdentity, _, err := h.OryClient.IdentityAPI.CreateIdentity(context.Background()).CreateIdentityBody(identity).Execute()
	if err != nil {
		return nil, err
	}

	// Store additional user data in MongoDB if needed
	user.OryID = createdIdentity.Id
	_, err = c.Mongo.InsertOne(c, "users", user)
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
