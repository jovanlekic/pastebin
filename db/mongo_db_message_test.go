package db

import (
	"context"
	"pastebin/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestReadMessages(t *testing.T) {
	// Set up MongoDB client for testing
	client, err := ConnectToMongoDb(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromMongoDb(context.Background(), client)

	// Create a test database and collection
	testDB := NewMongoDB(client, context.Background(), "test_db", "messages")

	// Drop the existing collection to start with a clean slate
	err = testDB.db.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Insert some test messages into the collection
	message1 := models.Message{MessageBody: "Test message 1"}
	message2 := models.Message{MessageBody: "Test message 2"}

	insertResult, err := testDB.db.InsertMany(context.Background(), []interface{}{message1, message2})
	if err != nil {
		t.Fatal(err)
	}

	// Convert the inserted IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range insertResult.InsertedIDs {
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			t.Fatal("Failed to convert ID to ObjectID")
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Call the ReadMessages function with the test IDs
	messages, err := testDB.ReadMessages(objectIDs)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the retrieved messages match the inserted ones
	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Compare message bodies
	if messages[0].MessageBody != "Test message 1" || messages[1].MessageBody != "Test message 2" {
		t.Fatal("Retrieved messages do not match expected values")
	}
}

func TestCreateMessage(t *testing.T) {
	client, err := ConnectToMongoDb(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromMongoDb(context.Background(), client)

	testDB := NewMongoDB(client, context.Background(), "test_db", "messages")

	err = testDB.db.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	messageBody := "testing"
	objectID, err := testDB.CreateMessage(messageBody)

	assert.NoError(t, err, "Expected no error")
	assert.NotEmpty(t, objectID, "Expected non-empty ObjectID")

	var storedMessage models.Message
	err = testDB.db.FindOne(testDB.ctx, primitive.M{"message_body": messageBody}).Decode(&storedMessage)
	assert.NoError(t, err, "Error fetching stored message from the database")
	assert.Equal(t, messageBody, storedMessage.MessageBody, "Stored message body does not match")
	assert.Equal(t, objectID, storedMessage.ID.Hex(), "Stored ObjectID does not match")
}

func TestReadMessage(t *testing.T) {
	client, err := ConnectToMongoDb(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromMongoDb(context.Background(), client)

	testDB := NewMongoDB(client, context.Background(), "test_db", "messages")

	err = testDB.db.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	messageBody := "testing"
	insertedID, err := testDB.CreateMessage(messageBody)
	if err != nil {
		t.Fatal(err)
	}

	objectID, err := primitive.ObjectIDFromHex(insertedID)
	if err != nil {
		t.Fatal(err)
	}

	readMessage, err := testDB.ReadMessage(objectID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, readMessage, "Expected a non-nil message")
	assert.Equal(t, objectID, readMessage.ID, "ObjectID does not match")
	assert.Equal(t, messageBody, readMessage.MessageBody, "Message body does not match")
}

func TestUpdateMessage(t *testing.T) {
	client, err := ConnectToMongoDb(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromMongoDb(context.Background(), client)

	testDB := NewMongoDB(client, context.Background(), "test_db", "messages")

	err = testDB.db.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	messageBody := "testing"
	insertedID, err := testDB.CreateMessage(messageBody)
	if err != nil {
		t.Fatal(err)
	}

	updatedMessage := models.Message{
		MessageBody: "Updated message!",
	}

	objectID, err := primitive.ObjectIDFromHex(insertedID)
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.UpdateMessage(objectID, updatedMessage)
	if err != nil {
		t.Fatal(err)
	}

	readMessage, err := testDB.ReadMessage(objectID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, readMessage, "Expected a non-nil message")
	assert.Equal(t, objectID, readMessage.ID, "ObjectID does not match")
	assert.Equal(t, updatedMessage.MessageBody, readMessage.MessageBody, "Message body does not match")
}

func TestDeleteMessage(t *testing.T) {
	client, err := ConnectToMongoDb(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer DisconnectFromMongoDb(context.Background(), client)

	testDB := NewMongoDB(client, context.Background(), "test_db", "messages")

	err = testDB.db.Drop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	messageBody := "testing"
	insertedID, err := testDB.CreateMessage(messageBody)
	if err != nil {
		t.Fatal(err)
	}

	objectID, err := primitive.ObjectIDFromHex(insertedID)
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.DeleteMessage(objectID)
	if err != nil {
		t.Fatal(err)
	}

	deletedMessage, err := testDB.ReadMessage(objectID)

	assert.Error(t, err, "Expected an error as the message should be deleted")
	assert.Nil(t, deletedMessage, "Expected a nil message as it should be deleted")
}
