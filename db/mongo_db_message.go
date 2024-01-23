package db

import (
	"fmt"
	"log"
	"pastebin/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (dbObj *MongoDB) CreateMessage(messageBody string) (string, error) {
	result, err := dbObj.db.InsertOne(dbObj.ctx, models.Message{MessageBody: messageBody})
	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID.Hex(), err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID.Hex(), fmt.Errorf("failed to get ObjectId")
	}

	return insertedID.Hex(), nil
}

func (dbObj *MongoDB) ReadMessage(id primitive.ObjectID) (*models.Message, error) {
	var message models.Message
	err := dbObj.db.FindOne(dbObj.ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (dbObj *MongoDB) UpdateMessage(id primitive.ObjectID, updatedMessage models.Message) error {
	_, err := dbObj.db.UpdateOne(
		dbObj.ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: updatedMessage}},
	)
	return err
}

func (dbObj *MongoDB) DeleteMessage(id primitive.ObjectID) error {
	_, err := dbObj.db.DeleteOne(dbObj.ctx, bson.M{"_id": id})
	return err
}
