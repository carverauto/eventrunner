//go:build !migrations

package migrations

import (
	"time"

	"gofr.dev/pkg/gofr/migration"
)

const (
	// Table creation statement for Cassandra
	createTableCassandra = `CREATE TABLE IF NOT EXISTS events (
		id TEXT,
		source TEXT,
		type TEXT,
		subject TEXT,
		time TIMESTAMP,
		data_contenttype TEXT,
		data TEXT,
		specversion TEXT,
		PRIMARY KEY ((type), time, id)
	) WITH CLUSTERING ORDER BY (time DESC);`

	// Example batch insert into the events table
	addCassandraRecords = `BEGIN BATCH
		INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion) VALUES ('1', 'sourceA', 'typeA', 'subjectA', toTimestamp(now()), 'application/json', '{ "key": "value" }', '1.0');
		INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion) VALUES ('2', 'sourceB', 'typeB', 'subjectB', toTimestamp(now()), 'application/xml', '<key>value</key>', '1.0');
	APPLY BATCH;`

	// Template query for inserting data into events
	eventDataCassandra = `INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
)

func createTableEventsCassandra() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// Create the events table
			if err := d.Cassandra.Exec(createTableCassandra); err != nil {
				return err
			}

			// Execute a batch operation to add records
			if err := d.Cassandra.Exec(addCassandraRecords); err != nil {
				return err
			}

			// Create a new batch for further operations
			batchName := "eventsBatch"
			if err := d.Cassandra.NewBatch(batchName, 0); err != nil { // 0 indicates a LoggedBatch for atomicity
				return err
			}

			now := time.Now()

			// Add records to the batch
			if err := d.Cassandra.BatchQuery(batchName, eventDataCassandra, "3", "sourceC", "typeC", "subjectC", now, "application/json", "{ \"key\": \"valueC\" }", "1.0"); err != nil {
				return err
			}

			if err := d.Cassandra.BatchQuery(batchName, eventDataCassandra, "4", "sourceD", "typeD", "subjectD", now, "application/xml", "<key>valueD</key>", "1.0"); err != nil {
				return err
			}

			// Execute the batch
			if err := d.Cassandra.ExecuteBatch(batchName); err != nil {
				return err
			}

			return nil
		},
	}
}
