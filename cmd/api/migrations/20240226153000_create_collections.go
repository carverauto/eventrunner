//go:build !migrations

package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createCollections() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			collections := []string{"tenants", "users", "api_keys"}
			for _, coll := range collections {
				_, err := d.SQL.Exec("db.createCollection(\"" + coll + "\")")
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
