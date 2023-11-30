package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"pastebin/db/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	ctx       context.Context
	dbName    string
	tableName string
	db        *mongo.Collection
}

func ConnectToMongoDb() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	fmt.Print("Successfully connected to Mongo!\n")
	return client
}

func DisconnectFromMongoDb(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	fmt.Print("Disconnected from Mongo!\n")
}

func NewMongoDB(client *mongo.Client, dbName string, tableName string) (dbObj *MongoDB) {
	dbObj = new(MongoDB) // mora da se rezervise mem za obj
	dbObj.ctx = context.TODO()
	dbObj.dbName = dbName
	dbObj.tableName = tableName
	dbObj.db = client.Database(dbName).Collection(tableName)
	return
}

func (dbObj *MongoDB) CreateMessage(newMessage models.Message) (string, error) {
	result, err := dbObj.db.InsertOne(dbObj.ctx, newMessage)
	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID.Hex(), err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID.Hex(), fmt.Errorf("Failed to get ObjectId")
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
		bson.D{{"$set", updatedMessage}},
	)
	return err
}

func (dbObj *MongoDB) DeleteMessage(id primitive.ObjectID) error {
	_, err := dbObj.db.DeleteOne(dbObj.ctx, bson.M{"_id": id})
	return err
}
