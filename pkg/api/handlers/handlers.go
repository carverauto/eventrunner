package handlers

import (
	"github.com/carverauto/eventrunner/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gofr.dev/pkg/gofr"
)

type TenantHandler struct{}

func (h *TenantHandler) Create(c *gofr.Context) (interface{}, error) {
	var tenant models.Tenant
	if err := c.Bind(&tenant); err != nil {
		return nil, err
	}

	result, err := c.Mongo.InsertOne(c, "tenants", tenant)
	if err != nil {
		return nil, err
	}

	tenant.ID = result.(primitive.ObjectID)

	return tenant, nil
}

func (h *TenantHandler) GetAll(c *gofr.Context) (interface{}, error) {
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

	// TODO: Hash password before storing

	result, err := c.Mongo.InsertOne(c, "users", user)
	if err != nil {
		return nil, err
	}

	user.ID = result.(primitive.ObjectID)

	return user, nil
}

func (h *UserHandler) GetAll(c *gofr.Context) (interface{}, error) {
	tenantID, err := primitive.ObjectIDFromHex(c.Param("tenant_id"))
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = c.Mongo.Find(c, "users", bson.M{"tenant_id": tenantID}, &users)

	return users, err
}
