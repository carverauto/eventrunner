package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name"`
	Customers []primitive.ObjectID `bson:"customers" json:"customers"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
	Active    bool                 `bson:"active" json:"active"`
}

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TenantID  primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Active    bool               `bson:"active" json:"active"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TenantID   primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Username   string             `bson:"username" json:"username"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"-"` // Don't expose password in JSON
	Role       string             `bson:"role" json:"role"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin  time.Time          `bson:"last_login" json:"last_login"`
	Active     bool               `bson:"active" json:"active"`
}

type APIKey struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TenantID   primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Key        string             `bson:"key" json:"key"`
	Active     bool               `bson:"active" json:"active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt  time.Time          `bson:"expires_at" json:"expires_at"`
	LastUsed   time.Time          `bson:"last_used" json:"last_used"`
}
