//go:build !migrations

package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createIndexes() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// Create unique index on tenant name
			if err := d.Mongo.CreateIndex(d.Context, "tenants", map[string]interface{}{"name": 1}, true); err != nil {
				return err
			}

			// Create unique index on user email
			if err := d.Mongo.CreateIndex(d.Context, "users", map[string]interface{}{"email": 1}, true); err != nil {
				return err
			}

			// Create index on user's tenant_id for faster queries
			if err := d.Mongo.CreateIndex(d.Context, "users", map[string]interface{}{"tenant_id": 1}, false); err != nil {
				return err
			}

			// Create unique index on API key
			if err := d.Mongo.CreateIndex(d.Context, "api_keys", map[string]interface{}{"key": 1}, true); err != nil {
				return err
			}

			return nil
		},
	}
}
