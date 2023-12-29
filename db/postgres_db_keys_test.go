package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareKeysTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Keys table if it exists
	dropScript := `
        DROP TABLE IF EXISTS Keys;
    `

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Keys table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS Keys (
            id UUID DEFAULT uuid_generate_v4(),
            key VARCHAR(32) NOT NULL,
            used BOOLEAN NOT NULL,
            PRIMARY KEY (id)
        );
    `

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Keys table created successfully!")
}

func TestMarkKeyAsUsed(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Insert a key into the Keys table
	testKey := "test_key"
	err = testDB.InsertKeyIntoKeys(context.Background(), testKey)
	if err != nil {
		t.Fatal(err)
	}

	// Mark the key as used
	err = testDB.MarkKeyAsUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")

	// Check if the key is marked as used
	isUsed, err := testDB.IsKeyUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, isUsed, "Expected the key to be marked as used")
}

func TestIsKeyUsed(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Insert a key into the Keys table
	testKey := "test_key"
	err = testDB.InsertKeyIntoKeys(context.Background(), testKey)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the key is initially not marked as used
	isUsed, err := testDB.IsKeyUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")
	assert.False(t, isUsed, "Expected the key to not be marked as used")

	// Mark the key as used
	err = testDB.MarkKeyAsUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")

	// Check if the key is now marked as used
	isUsed, err = testDB.IsKeyUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, isUsed, "Expected the key to be marked as used")
}

func TestGetFirstUnusedKey(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Insert keys into the Keys table
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err := testDB.InsertKeyIntoKeys(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get the first unused key
	firstKey, err := testDB.GetFirstUnusedKey(context.Background())
	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, firstKey, "Expected a non-empty key")
}

func TestGetAndMarkFirstUnusedKey(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Insert keys into the Keys table
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err := testDB.InsertKeyIntoKeys(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get and mark the first unused key
	firstKey, err := testDB.GetAndMarkFirstUnusedKey(context.Background())
	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, firstKey, "Expected a non-empty key")

	// Check if the key is marked as used
	isUsed, err := testDB.IsKeyUsed(context.Background(), firstKey)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, isUsed, "Expected the key to be marked as used")
}

func TestInsertKeyIntoKeys(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Insert a key into the Keys table
	testKey := "test_key"
	err = testDB.InsertKeyIntoKeys(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")

	// Check if the key is in the Keys table
	isUsed, err := testDB.IsKeyUsed(context.Background(), testKey)
	assert.NoError(t, err, "Expected no error")
	assert.False(t, isUsed, "Expected the key to not be marked as used")
}

func TestFillKeysTable(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareKeysTable(t, testDB)

	// Fill the Keys table with 6-character words
	err = testDB.FillKeysTable(context.Background())
	assert.NoError(t, err, "Expected no error")

	// Check if all keys are in the Keys table
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				for l := 0; l < 3; l++ {
					for m := 0; m < 3; m++ {
						for n := 0; n < 3; n++ {
							key := string(rune('a'+i)) + string(rune('a'+j)) + string(rune('a'+k)) +
								string(rune('a'+l)) + string(rune('a'+m)) + string(rune('a'+n))
							isUsed, err := testDB.IsKeyUsed(context.Background(), key)
							assert.NoError(t, err, "Expected no error")
							assert.False(t, isUsed, "Expected the key to not be marked as used")
						}
					}
				}
			}
		}
	}
}
