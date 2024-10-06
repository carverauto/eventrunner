package handlers

import (
	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gofr.dev/pkg/gofr"
)

type TenantHandler struct{}

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

type UserHandler struct{}

func (*UserHandler) Create(c *gofr.Context) (interface{}, error) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return nil, err
	}

	// TODO: Hash password before storing

	result, err := c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, err
	}

	user.ID = result.(uuid.UUID)

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
