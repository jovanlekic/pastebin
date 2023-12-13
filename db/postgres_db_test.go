package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"pastebin/db/models"
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

func prepareUsedTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Used table if it exists
	dropScript := `
		DROP TABLE IF EXISTS Used;
	`

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Used table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE Used (
			id UUID DEFAULT uuid_generate_v4(),
			key VARCHAR(20) NOT NULL,
			PRIMARY KEY (id)
		);
	`

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Used table created successfully!")
}

func prepareUnusedTable(t *testing.T, testDB *PostgresDB) {
	// Drop the Unused table if it exists
	dropScript := `
		DROP TABLE IF EXISTS Unused;
	`

	_, err := testDB.db.ExecContext(context.Background(), dropScript)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Unused table with your specified schema
	createScript := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE Unused (
			id UUID DEFAULT uuid_generate_v4(),
			key VARCHAR(20) NOT NULL,
			PRIMARY KEY (id)
		);
	`

	_, err = testDB.db.ExecContext(context.Background(), createScript)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Unused table created successfully!")
}

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

func TestFillUnusedTable(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	// Ensure Unused and Migration tables are created
	prepareUnusedTable(t, testDB)

	err = testDB.FillUnusedTable(context.Background())
	assert.NoError(t, err, "Expected no error")

	// Check if the migration exists
	exists, err := testDB.CheckMigration(context.Background(), "FillUnusedTable")
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected migration to exist in the database")

	// Check if Unused table is filled
	count, err := countRowsInTable(t, testDB.db, "Unused")
	assert.NoError(t, err, "Expected no error")
	assert.Greater(t, count, 0, "Expected Unused table to be filled")
}

func countRowsInTable(t *testing.T, db *sql.DB, tableName string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int
	err := db.QueryRowContext(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func TestInsertWordIntoUnused(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUnusedTable(t, testDB)

	word := "testword"

	err = testDB.InsertWordIntoUnused(context.Background(), word)
	assert.NoError(t, err, "Expected no error")

	// Check if the word exists in the Unused table
	exists, err := testDB.IsKeyInUsedTable(context.Background(), word)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected word to exist in the Unused table")
}

func TestGetFirstUnusedKey(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUnusedTable(t, testDB)

	// Insert a word into the Unused table
	word := "testword"
	err = testDB.InsertWordIntoUnused(context.Background(), word)
	assert.NoError(t, err, "Expected no error")

	// Get the first unused key
	firstKey, err := testDB.GetFirstUnusedKey(context.Background())
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, word, firstKey, "Expected first key to match the inserted word")
}

func TestDeleteKeyFromUnused(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUnusedTable(t, testDB)

	// Insert a word into the Unused table
	word := "testword"
	err = testDB.InsertWordIntoUnused(context.Background(), word)
	assert.NoError(t, err, "Expected no error")

	// Delete the word from the Unused table
	err = testDB.DeleteKeyFromUnused(context.Background(), word)
	assert.NoError(t, err, "Expected no error")

	// Check if the word exists in the Unused table after deletion
	exists, err := testDB.IsKeyInUsedTable(context.Background(), word)
	assert.NoError(t, err, "Expected no error")
	assert.False(t, exists, "Expected word to be deleted from the Unused table")
}

func TestInsertKeyIntoUsed(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUsedTable(t, testDB)

	key := "testkey"

	err = testDB.InsertKeyIntoUsed(context.Background(), key)
	assert.NoError(t, err, "Expected no error")

	// Check if the key exists in the Used table
	exists, err := testDB.IsKeyInUsedTable(context.Background(), key)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected key to exist in the Used table")
}

func TestMoveKeyFromUnusedToUsed(t *testing.T) {
	postgresClient, err := ConnectToPostgresDb("test_db", "postgres", "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromPostgresDb(postgresClient)
	testDB := NewPostgresDB(postgresClient)

	prepareUsedTable(t, testDB)

	// Insert a word into the Unused table
	word := "testword"
	err = testDB.InsertWordIntoUnused(context.Background(), word)
	assert.NoError(t, err, "Expected no error")

	// Move the key from Unused to Used
	movedKey, err := testDB.MoveKeyFromUnusedToUsed(context.Background())
	assert.NoError(t, err, "Expected no error")

	// Check if the moved key exists in the Used table
	exists, err := testDB.IsKeyInUsedTable(context.Background(), movedKey)
	assert.NoError(t, err, "Expected no error")
	assert.True(t, exists, "Expected moved key to exist in the Used table")

	// Check if the moved key is deleted from the Unused table
	exists, err = testDB.IsKeyInUsedTable(context.Background(), movedKey)
	assert.NoError(t, err, "Expected no error")
	assert.False(t, exists, "Expected moved key to be deleted from the Unused table")
}
