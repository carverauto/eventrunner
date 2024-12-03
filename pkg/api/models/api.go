/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID        uuid.UUID   `bson:"_id,omitempty" json:"-"`
	TenantID  uuid.UUID   `bson:"tenant_id" json:"id"`
	Name      string      `bson:"name" json:"name"`
	Customers []uuid.UUID `bson:"customers" json:"customers"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
	Active    bool        `bson:"active" json:"active"`
}

type Customer struct {
	ID         uuid.UUID `bson:"_id,omitempty" json:"-"`
	CustomerID uuid.UUID `bson:"customer_id" json:"id"`
	TenantID   uuid.UUID `bson:"tenant_id" json:"tenant_id"`
	Name       string    `bson:"name" json:"name"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	Active     bool      `bson:"active" json:"active"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OryID      string             `bson:"ory_id" json:"ory_id"`
	UserID     uuid.UUID          `bson:"user_id" json:"user_id"`
	TenantIDs  []uuid.UUID        `bson:"tenant_ids" json:"tenant_ids"`
	CustomerID uuid.UUID          `bson:"customer_id" json:"customer_id,omitempty"`
	Username   string             `bson:"username" json:"username"`
	Email      string             `bson:"email" json:"email"`
	Roles      []string           `bson:"roles" json:"roles"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin  time.Time          `bson:"last_login" json:"last_login"`
	Active     bool               `bson:"active" json:"active"`
}

type APICredentials struct {
	ID         uuid.UUID `bson:"_id,omitempty" json:"-"`
	UserID     uuid.UUID `bson:"user_id" json:"user_id"`
	KeyID      uuid.UUID `bson:"key_id" json:"id"`
	TenantID   uuid.UUID `bson:"tenant_id" json:"tenant_id"`
	CustomerID uuid.UUID `bson:"customer_id" json:"customer_id"`
	ClientID   string    `bson:"client_id" json:"client_id"`
	Key        string    `bson:"key" json:"key"`
	Active     bool      `bson:"active" json:"active"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt  time.Time `bson:"expires_at" json:"expires_at"`
	LastUsed   time.Time `bson:"last_used" json:"last_used"`
}
