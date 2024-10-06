//go:build !migrations

package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func createIndexes() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			indexes := []struct {
				collection string
				keys       string
				unique     bool
			}{
				{collection: "tenants", keys: "name", unique: true},
				{collection: "users", keys: "email", unique: true},
				{collection: "users", keys: "tenant_id", unique: false},
				{collection: "api_keys", keys: "key", unique: true},
			}

			for _, idx := range indexes {
				uniqueStr := "false"
				if idx.unique {
					uniqueStr = "true"
				}
				_, err := d.SQL.Exec("db." + idx.collection + ".createIndex({" + idx.keys + ": 1}, {unique: " + uniqueStr + "})")
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
