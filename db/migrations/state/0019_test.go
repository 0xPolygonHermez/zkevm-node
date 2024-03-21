package migrations_test

import (
	"database/sql"
	"testing"
)

type migrationTest0019 struct{}

func (m migrationTest0019) InsertData(db *sql.DB) error {
	//TODO: Add insert data
	return nil
}

func (m migrationTest0019) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	//TODO: Add checks
}

func (m migrationTest0019) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	//TODO: Add checks
}
func TestMigration0019(t *testing.T) {
	runMigrationTest(t, 19, migrationTest0019{})
}
