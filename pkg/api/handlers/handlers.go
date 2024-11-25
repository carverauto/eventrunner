package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/models"
	customctx "github.com/carverauto/eventrunner/pkg/context"
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

// getCustomContext is a helper function to extract the custom context
func getCustomContext(c *gofr.Context) (customctx.Context, error) {
	customCtxVal := c.Request.Context().Value("customCtx")
	if customCtxVal == nil {
		return nil, errors.NewAppError(500, "custom context not found")
	}

	customCtx, ok := customCtxVal.(customctx.Context)
	if !ok {
		return nil, errors.NewAppError(500, "invalid custom context type")
	}

	return customCtx, nil
}

type UserInfo struct {
	UserID   uuid.UUID
	Roles    []string
	TenantID uuid.UUID
	Email    string
}

func getUserInfo(c *gofr.Context) (*UserInfo, error) {
	customCtx, err := getCustomContext(c)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := customCtx.GetStringClaim("X-User-Id")
	if !ok {
		return nil, errors.NewMissingParamError("user ID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.NewInvalidParamError("user ID")
	}

	role, ok := customCtx.GetStringClaim("X-User-Role")
	if !ok {
		return nil, errors.NewMissingParamError("user role")
	}

	tenantIDStr, ok := customCtx.GetStringClaim("X-Tenant-Id")
	if !ok {
		return nil, errors.NewMissingParamError("tenant ID")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, errors.NewInvalidParamError("tenant ID")
	}

	email, _ := customCtx.GetStringClaim("X-User-Email") // Optional

	return &UserInfo{
		UserID:   userID,
		Roles:    []string{role},
		TenantID: tenantID,
		Email:    email,
	}, nil
}

/*
func UpdateRoles(c *gofr.Context) (interface{}, error) {
	var req struct {
		UserID string   `json:"user_id"`
		Roles  []string `json:"roles"`
	}

	if err := c.Bind(&req); err != nil {
		return nil, err
	}

	// Get the admin user's email from the session
	adminEmail := c.Get("user_email")

	// Check if admin is from threadr.ai domain
	if !strings.HasSuffix(adminEmail, "@threadr.ai") {
		return nil, errors.New("unauthorized")
	}

	// Update roles using Kratos Admin API
	// Implementation depends on your Kratos client setup
	return kratosClient.UpdateIdentity(req.UserID, req.Roles)
}

*/

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

	userInfo, err := getUserInfo(c)
	if err != nil {
		return nil, err
	}

	// look through userInfo.Roles to see if they are the superUser
	isSuperUser := containsRole(userInfo.Roles, "superuser")
	if !isSuperUser {
		return nil, errors.NewForbiddenError("Only superuser can create tenants")
	}

	tenant.TenantID = uuid.New()
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	_, err = c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert tenant into database")
	}

	err = c.Mongo.UpdateOne(c, "users",
		bson.M{"user_id": userInfo.UserID},
		bson.M{"$addToSet": bson.M{"tenant_ids": tenant.TenantID}})
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to update user's tenants")
	}

	return tenant, nil
}

// containsRole checks if a role is present in a list of roles from the UserInfo struct.
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
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

	// Get custom context
	customCtx, err := getCustomContext(c)
	if err != nil {
		return nil, err
	}

	// Get requester's role
	requesterRole, ok := customCtx.GetStringClaim("X-User-Role")
	if !ok {
		return nil, errors.NewMissingParamError("X-User-Role header")
	}

	// Check if the requester has the right to create users
	isSuperuser := requesterRole == "superuser"
	isTenantAdmin := requesterRole == "tenant_admin"

	if !isSuperuser && !isTenantAdmin {
		return nil, errors.NewForbiddenError("Insufficient permissions to create user")
	}

	// Get the tenant ID from the request parameters
	tenantIDStr := c.Param("tenant_id")
	if tenantIDStr == "" {
		return nil, errors.NewMissingParamError("tenant_id")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, errors.NewInvalidParamError("tenant_id")
	}

	// Verify that the requester has access to this tenant
	requesterTenantID, ok := customCtx.GetStringClaim("X-Tenant-ID")
	if !ok {
		return nil, errors.NewMissingParamError("X-Tenant-ID header")
	}

	// If not superuser, verify tenant access
	if !isSuperuser {
		requesterTenantUUID, err := uuid.Parse(requesterTenantID)
		if err != nil {
			return nil, errors.NewInvalidParamError("requester tenant ID")
		}

		if requesterTenantUUID != tenantID {
			return nil, errors.NewForbiddenError("You don't have access to this tenant")
		}
	}

	// Create identity in Ory
	identity := ory.CreateIdentityBody{
		SchemaId: os.Getenv("ORY_SCHEMA_ID"),
		Traits: map[string]interface{}{
			"email":     user.Email,
			"roles":     user.Roles,
			"tenant_id": tenantID.String(),
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

	// If there's a customer ID in the request, validate and assign it
	if customerIDStr, ok := customCtx.GetStringClaim("X-Customer-ID"); ok {
		customerID, err := uuid.Parse(customerIDStr)
		if err != nil {
			return nil, errors.NewInvalidParamError("customer ID")
		}
		user.CustomerID = customerID
	}

	_, err = c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "Failed to insert user into database")
	}

	// Get the creator's email
	creatorEmail, ok := customCtx.GetStringClaim("X-User-Email")
	if !ok {
		creatorEmail = "unknown" // fallback value if email is not available
	}

	// Create the audit log
	auditLog := models.AuditLog{
		ID:          uuid.New(),
		TenantID:    tenantID,
		UserID:      user.UserID,
		Action:      "create_user",
		Description: fmt.Sprintf("User %s created by %s", user.Email, creatorEmail),
		CreatedAt:   time.Now(),
	}

	_, err = c.Mongo.InsertOne(c, "audit_logs", auditLog)
	if err != nil {
		c.Logger.Errorf("Failed to create audit log: %v", err)
		// Don't return the error as this is not critical
	}

	return user, nil
}

type AuditAction string

const (
	AuditActionCreateUser AuditAction = "create_user"
	AuditActionUpdateUser AuditAction = "update_user"
	AuditActionDeleteUser AuditAction = "delete_user"
)

func createAuditLog(c *gofr.Context, action AuditAction, tenantID uuid.UUID, userID uuid.UUID, description string) error {
	auditLog := models.AuditLog{
		ID:          uuid.New(),
		TenantID:    tenantID,
		UserID:      userID,
		Action:      string(action),
		Description: description,
		CreatedAt:   time.Now(),
	}

	_, err := c.Mongo.InsertOne(c, "audit_logs", auditLog)
	return err
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
