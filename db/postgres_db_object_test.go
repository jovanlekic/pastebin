package db

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"pastebin/db/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareObjectTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Object table if it exists
	dropScript := `
		DROP TABLE IF EXISTS Object;
	`

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Object table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE Object (
			dev_key        varchar(32) NOT NULL,
			paste_key      varchar(20) NOT NULL,
			message_id     varchar(32),
			PRIMARY KEY (dev_key, paste_key)
		);
	`

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Object table created successfully!")
}

func TestCreateObject(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareObjectTable(t, testDB)

	hasher := md5.New()
	hasher.Write([]byte("test"))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	testObject := models.Object{
		PasteKey:  "12345678901234567890",
		DevKey:    hashString,
		MessageID: hashString,
	}
	err = testDB.CreateObject(context.Background(), &testObject)

	assert.NoError(t, err, "Expected no error")
}

func TestReadObject(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareObjectTable(t, testDB)

	// Create a test object to be inserted into the database
	testObject := models.Object{
		PasteKey:  "test_paste_key",
		DevKey:    "test_dev_key",
		MessageID: "test_message_id",
	}

	// Insert the test object into the database
	err = testDB.CreateObject(context.Background(), &testObject)
	if err != nil {
		t.Fatal(err)
	}

	// Read the object from the database
	resultObject, err := testDB.ReadObject(context.Background(), testObject.PasteKey, testObject.DevKey)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, &testObject, resultObject, "Expected the retrieved object to match the inserted object")

	// Test case where the object does not exist in the database
	nonExistentPasteKey := "non_existent_paste_key"
	nonExistentDevKey := "non_existent_dev_key"
	nonExistentObject, err := testDB.ReadObject(context.Background(), nonExistentPasteKey, nonExistentDevKey)
	assert.Error(t, err, "Expected an error for a non-existent object")
	assert.Nil(t, nonExistentObject, "Expected a nil object for a non-existent pasteKey and devKey")
}

func TestUpdateObject(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareObjectTable(t, testDB)

	// Create a test object to be inserted into the database
	testObject := models.Object{
		PasteKey:  "test_paste_key",
		DevKey:    "test_dev_key",
		MessageID: "test_message_id",
	}

	// Insert the test object into the database
	err = testDB.CreateObject(context.Background(), &testObject)
	if err != nil {
		t.Fatal(err)
	}

	// Update the object in the database
	updatedObject := models.Object{
		PasteKey:  "test_paste_key",
		DevKey:    "test_dev_key",
		MessageID: "updated_message_id",
	}

	err = testDB.UpdateObject(context.Background(), &updatedObject)
	assert.NoError(t, err, "Expected no error")

	// Read the object from the database to check if it was updated successfully
	resultObject, err := testDB.ReadObject(context.Background(), updatedObject.PasteKey, updatedObject.DevKey)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, updatedObject.MessageID, resultObject.MessageID, "Expected the retrieved object to have the updated message ID")
}

func TestDeleteObject(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareObjectTable(t, testDB)

	// Create a test object to be inserted into the database
	testObject := models.Object{
		PasteKey:  "test_paste_key",
		DevKey:    "test_dev_key",
		MessageID: "test_message_id",
	}

	// Insert the test object into the database
	err = testDB.CreateObject(context.Background(), &testObject)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the object from the database
	err = testDB.DeleteObject(context.Background(), testObject.PasteKey, testObject.DevKey)
	assert.NoError(t, err, "Expected no error")

	// Try to read the object from the database to verify it was deleted
	deletedObject, err := testDB.ReadObject(context.Background(), testObject.PasteKey, testObject.DevKey)
	assert.Error(t, err, "Expected an error as the object should be deleted")
	assert.Nil(t, deletedObject, "Expected a nil object for a deleted object")
}
