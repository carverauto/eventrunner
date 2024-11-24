package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID `bson:"_id" json:"id"`
	TenantID    uuid.UUID `bson:"tenant_id" json:"tenant_id"`
	UserID      uuid.UUID `bson:"user_id" json:"user_id"`
	Action      string    `bson:"action" json:"action"`
	Description string    `bson:"description" json:"description"`
	Metadata    any       `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
