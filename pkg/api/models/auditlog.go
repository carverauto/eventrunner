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
