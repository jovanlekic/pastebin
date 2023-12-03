package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"pastebin/db/models"
	"regexp"

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

const projName = "pastebin"

func ConnectToMongoDb(ctx context.Context) (*mongo.Client, error) {
	re := regexp.MustCompile(`^(.*` + projName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Println("No .env file found")
		return nil, err
	}
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		return nil, errors.New("Must set MongoURI!")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	fmt.Print("Successfully connected to Mongo!\n")
	return client, err
}

func DisconnectFromMongoDb(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
	fmt.Print("Disconnected from Mongo!\n")
}

func NewMongoDB(client *mongo.Client, ctx context.Context, dbName string, tableName string) (dbObj *MongoDB) {
	dbObj = new(MongoDB) // mora da se rezervise mem za obj
	dbObj.ctx = ctx
	dbObj.dbName = dbName
	dbObj.tableName = tableName
	dbObj.db = client.Database(dbName).Collection(tableName)
	return
}

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
