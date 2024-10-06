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
				if err := d.Mongo.CreateCollection(d.Context, coll); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
