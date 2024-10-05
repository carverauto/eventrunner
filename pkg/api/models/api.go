package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TenantID primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"` // Don't expose password in JSON
	Role     string             `bson:"role" json:"role"`
}

type APIKey struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Key    string             `bson:"key" json:"key"`
	Active bool               `bson:"active" json:"active"`
}
