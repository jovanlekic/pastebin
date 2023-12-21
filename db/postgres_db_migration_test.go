package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareMigrationTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Migration table if it exists
	dropScript := `
		DROP TABLE IF EXISTS Migration;
	`

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Migration table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE Migration (
			id UUID DEFAULT uuid_generate_v4(),
			migration VARCHAR(50) NOT NULL,
			PRIMARY KEY (id)
		);
	`

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Migration table created successfully!")
}

func TestInsertMigration(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareMigrationTable(t, testDB)

	migrationName := "TestMigration"

	err = testDB.InsertMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")

	// Check if the migration exists
	exists, err := testDB.CheckMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected migration to exist in the database")
}

func TestDeleteMigration(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareMigrationTable(t, testDB)

	migrationName := "TestMigration"

	// Insert a migration record
	err = testDB.InsertMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")

	// Delete the migration record
	err = testDB.DeleteMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")

	// Check if the migration exists after deletion
	exists, err := testDB.CheckMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")
	assert.False(t, exists, "Expected migration to be deleted from the database")
}

func TestCheckMigration(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareMigrationTable(t, testDB)

	migrationName := "TestMigration"

	// Insert a migration record
	err = testDB.InsertMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")

	// Check if the migration exists after deletion
	exists, err := testDB.CheckMigration(context.Background(), migrationName)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected migration to be found in the database")
}
