package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/middleware"
	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/carverauto/eventrunner/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
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

/*
	TODO: Finish this function
func (h *Handlers) CreateAPICredentials(c *gofr.Context) (interface{}, error) {
	// Get user info from session headers
	userID := c.Header("X-User")
	tenantID := c.Header("X-Tenant-ID")

	// Create OAuth2 client in Ory Hydra
	client := ory.CreateOAuth2ClientBody{
		Scope:                   "api:access",
		GrantTypes:              []string{"client_credentials"},
		TokenEndpointAuthMethod: "client_secret_post",
		// Add other required fields
	}

	createdClient, _, err := h.OryClient.OAuth2API.CreateOAuth2Client(c.Context).
		CreateOAuth2ClientBody(client).Execute()
	if err != nil {
		return nil, err
	}

	// Store association in your database
	apiCreds := models.APICredentials{
		UserID:   userID,
		TenantID: tenantID,
		ClientID: createdClient.ClientId,
		// Store other metadata
	}

	_, err = c.Mongo.InsertOne(c, "api_credentials", apiCreds)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"client_id":     createdClient.ClientId,
		"client_secret": createdClient.ClientSecret,
	}, nil
}
*/

func (h *Handlers) CreateSuperUser(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, errors.NewInvalidParamError("user data")
	}

	// Check if this is the first user
	var existingUsers []models.User
	err := c.Mongo.Find(c, "users", bson.M{}, &existingUsers)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to query users")
	}

	if len(existingUsers) > 0 {
		return nil, errors.NewAppError(409, "Superuser already exists")
	}

	// Generate a new UUID for the default tenant
	defaultTenantID := uuid.New()

	// Create default tenant
	defaultTenant := models.Tenant{
		ID:        uuid.New(),
		TenantID:  defaultTenantID,
		Name:      "Default Tenant",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = c.Mongo.InsertOne(c, "tenants", defaultTenant)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert default tenant")
	}

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: os.Getenv("ORY_SCHEMA_ID"),
		Traits: map[string]interface{}{
			"email":     user.Email,
			"roles":     []string{"superuser"},
			"tenant_id": defaultTenantID.String(),
		},
	}

	createdIdentity, _, err := h.OryClient.IdentityAPI.CreateIdentity(c.Context).CreateIdentityBody(identity).Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create identity in Ory: %v", err))
	}

	// Create user in MongoDB
	user.UserID = uuid.New()
	user.OryID = createdIdentity.Id
	user.TenantIDs = []uuid.UUID{defaultTenantID}
	user.Roles = []string{"superuser"}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert user into database")
	}

	return map[string]interface{}{
		"user":   user,
		"tenant": defaultTenant,
	}, nil
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

	roles, ok := claims["roles"].([]interface{})
	if !ok || !containsRole(roles, "superuser") {
		return nil, errors.NewForbiddenError("Only superuser can create tenants")
	}

	tenant.TenantID = uuid.New()
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert tenant into database")
	}

	// Update the superuser's tenants
	userID, _ := uuid.Parse(claims["sub"].(string))
	err = c.Mongo.UpdateOne(c, "users",
		bson.M{"user_id": userID},
		bson.M{"$addToSet": bson.M{"tenant_ids": tenant.TenantID}})
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to update user's tenants")
	}

	return tenant, nil
}

func containsRole(roles []interface{}, role string) bool {
	for _, r := range roles {
		if r.(string) == role {
			return true
		}
	}
	return false
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
	requesterRoles, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, errors.NewForbiddenError("Invalid role data")
	}

	isSuperuser := containsRole(requesterRoles, "superuser")
	isTenantAdmin := containsRole(requesterRoles, "tenant_admin")

	if !isSuperuser && !isTenantAdmin {
		return nil, errors.NewForbiddenError("Insufficient permissions to create user")
	}

	// Get the tenant ID from the request
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return nil, errors.NewInvalidParamError("tenant_id")
	}

	// Verify that the requester has access to this tenant
	requesterTenantIDs, ok := claims["tenant_ids"].([]interface{})
	if !ok || !containsTenantID(requesterTenantIDs, tenantID) {
		return nil, errors.NewForbiddenError("You don't have access to this tenant")
	}

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: os.Getenv("ORY_SCHEMA_ID"),
		Traits: map[string]interface{}{
			"email":     user.Email,
			"roles":     user.Roles,
			"tenant_id": tenantID.String(), // We still need to include a primary tenant_id for Ory
		},
	}

	createdIdentity, _, err := h.OryClient.IdentityAPI.CreateIdentity(c.Context).CreateIdentityBody(identity).Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create identity in Ory: %v", err))
	}

	// Create user in MongoDB
	user.UserID = uuid.New()
	user.OryID = createdIdentity.Id
	user.TenantIDs = []uuid.UUID{tenantID}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert user into database")
	}

	return user, nil
}

func containsTenantID(tenantIDs []interface{}, tenantID uuid.UUID) bool {
	for _, id := range tenantIDs {
		if id.(string) == tenantID.String() {
			return true
		}
	}
	return false
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

func (h *Handlers) Login(c *gofr.Context) (interface{}, error) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return nil, errors.NewInvalidParamError("login data")
	}

	// Create a login flow
	flow, _, err := h.OryClient.FrontendAPI.CreateNativeLoginFlow(c.Context).Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create login flow: %v", err))
	}

	// Prepare the update login flow body
	updateLoginFlowBody := ory.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &ory.UpdateLoginFlowWithPasswordMethod{
			Password:   loginReq.Password,
			Identifier: loginReq.Email,
			Method:     "password",
		},
	}

	// Update the login flow
	resp, _, err := h.OryClient.FrontendAPI.UpdateLoginFlow(c.Context).
		Flow(flow.Id).
		UpdateLoginFlowBody(updateLoginFlowBody).
		Execute()

	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials")
	}

	// Check if the login was successful
	if resp.Session.Identity == nil {
		return nil, errors.NewUnauthorizedError("Login failed")
	}

	// Extract user information from Ory identity traits
	traits, ok := resp.Session.Identity.Traits.(map[string]interface{})
	if !ok {
		return nil, errors.NewAppError(500, "Failed to parse identity traits")
	}

	email, _ := traits["email"].(string)
	roles, _ := traits["roles"].([]interface{})
	tenantIDs, _ := traits["tenant_ids"].([]interface{})
	customerID, _ := traits["customer_id"].(string)

	// Convert roles to []string
	stringRoles := make([]string, len(roles))
	for i, role := range roles {
		stringRoles[i], _ = role.(string)
	}

	// Convert tenant_ids to []string
	stringTenantIDs := make([]string, len(tenantIDs))
	for i, id := range tenantIDs {
		stringTenantIDs[i], _ = id.(string)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         resp.Session.Identity.Id,
		"email":       email,
		"roles":       stringRoles,
		"tenant_ids":  stringTenantIDs,
		"customer_id": customerID,
		"exp":         time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days expiration
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to generate token")
	}

	response := map[string]interface{}{
		"token": tokenString,
	}

	// Only include session_token if it's not nil
	if resp.SessionToken != nil {
		response["session_token"] = *resp.SessionToken
	}

	return response, nil
}
