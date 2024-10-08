//go:build !migrations

package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		20240226153000: createCollections(),
		20240226153100: createIndexes(),
	}
}
