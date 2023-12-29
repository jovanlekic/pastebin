package db

import (
	"context"
	"pastebin/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
