package db

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"pastebin/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func prepareUserTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Users table if it exists
	dropScript := `
		DROP TABLE IF EXISTS Users;
	`

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Users table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE Users (
			user_id UUID DEFAULT uuid_generate_v4(),
			name VARCHAR(20) NOT NULL,
			password VARCHAR(32) NOT NULL,
			pasteNum INT NOT NULL,
			dev_key VARCHAR(32) NOT NULL,
			email VARCHAR(32) NOT NULL,
			PRIMARY KEY (user_id)
		);
	`

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Users table created successfully!")
}

func TestCreateUser(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUserTable(t, testDB)

	hasher := md5.New()
	hasher.Write([]byte("test"))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	testUser := models.User{
		Name:     "test",
		Password: hashString,
		PasteNum: 0,
		DevKey:   hashString,
		Email:    "test@gmail.com",
	}

	uuid, err := testDB.CreateUser(context.Background(), &testUser)

	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, uuid, "Expected non-empty uuid")

}

func TestReadUserById(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUserTable(t, testDB)

	hasher := md5.New()
	hasher.Write([]byte("test"))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// Create a test user to be inserted into the database
	testUser := models.User{
		Name:     "test",
		Password: hashString,
		PasteNum: 0,
		DevKey:   hashString,
		Email:    "test@example.com",
	}

	// Insert the test user into the database
	tuuid, err := testDB.CreateUser(context.Background(), &testUser)
	if err != nil {
		t.Fatal(err)
	}

	testUser.UserID = tuuid

	// Read the user by ID from the database
	resultUser, err := testDB.ReadUserById(context.Background(), tuuid)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, testUser, resultUser, "Expected the retrieved user to match the inserted user")

	// Test case where the user ID does not exist in the database
	nonExistentUserID := uuid.New()
	nonExistentUser, err := testDB.ReadUserById(context.Background(), nonExistentUserID)
	assert.Error(t, err, "Expected an error for a non-existent user")
	assert.Equal(t, models.User{}, nonExistentUser, "Expected an empty user for a non-existent user ID")
}

func TestReadUserByUsername(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUserTable(t, testDB)

	hasher := md5.New()
	hasher.Write([]byte("test"))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// Create a test user to be inserted into the database
	testUser := models.User{
		Name:     "test",
		Password: hashString,
		PasteNum: 0,
		DevKey:   hashString,
		Email:    "test@example.com",
	}

	// Insert the test user into the database
	testUser.UserID, err = testDB.CreateUser(context.Background(), &testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Read the user by username from the database
	resultUser, err := testDB.ReadUserByUsername(context.Background(), testUser.Name)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, testUser, resultUser, "Expected the retrieved user to match the inserted user")

	// Test case where the username does not exist in the database
	nonExistentUsername := "non_existent_username"
	nonExistentUser, err := testDB.ReadUserByUsername(context.Background(), nonExistentUsername)
	assert.Error(t, err, "Expected an error for a non-existent username")
	assert.Equal(t, models.User{}, nonExistentUser, "Expected an empty user for a non-existent username")
}

func TestUpdateUser(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUserTable(t, testDB)

	// Create a test user to be inserted into the database
	testUser := models.User{
		Name:     "test_username",
		Password: "test_password",
		PasteNum: 0,
		DevKey:   "test_dev_key",
		Email:    "test@example.com",
	}

	// Insert the test user into the database
	testUser.UserID, err = testDB.CreateUser(context.Background(), &testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Update the user in the database
	updatedUser := models.User{
		UserID:   testUser.UserID,
		Name:     "updated_username",
		Password: "updated_password",
		PasteNum: 1,
		DevKey:   "updated_dev_key",
		Email:    "updated@example.com",
	}

	err = testDB.UpdateUser(context.Background(), testUser.UserID, updatedUser)
	assert.NoError(t, err, "Expected no error")

	// Read the user from the database to check if it was updated successfully
	resultUser, err := testDB.ReadUserById(context.Background(), testUser.UserID)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, updatedUser, resultUser, "Expected the retrieved user to match the updated user")
}

func TestDeleteUser(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUserTable(t, testDB)

	// Create a test user to be inserted into the database
	testUser := models.User{
		Name:     "test_username",
		Password: "test_password",
		PasteNum: 0,
		DevKey:   "test_dev_key",
		Email:    "test@example.com",
	}

	// Insert the test user into the database
	_, err = testDB.CreateUser(context.Background(), &testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the user from the database
	err = testDB.DeleteUser(context.Background(), testUser.UserID)
	assert.NoError(t, err, "Expected no error")

	// Try to read the user from the database to verify it was deleted
	deletedUser, err := testDB.ReadUserById(context.Background(), testUser.UserID)
	assert.Error(t, err, "Expected an error as the user should be deleted")
	assert.Equal(t, models.User{}, deletedUser, "Expected an empty user for a deleted user")
}
