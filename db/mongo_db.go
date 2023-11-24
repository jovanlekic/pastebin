package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	ctx       context.Context
	dbName    string
	tableName string
	db        *mongo.Collection
}

type Message struct {
	MessageID   string `bson:"message_id"`
	MessageBody string `bson:"message_body"`
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
	return client
}

func DisconnectFromMongoDb(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	fmt.Print("Disconnected!\n")
}

func NewDB(client *mongo.Client, dbName string, tableName string) (dbObj *DB) {
	dbObj = new(DB) // mora da se rezervise mem za obj
	dbObj.ctx = context.TODO()
	dbObj.dbName = dbName
	dbObj.tableName = tableName
	dbObj.db = client.Database(dbName).Collection(tableName)
	return
}

func (dbObj *DB) FindOne(filter interface{}) (result Message, err error) {
	err = dbObj.db.FindOne(dbObj.ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Print("No document was found with this filter!\n", err)
		return
	}
	if err != nil {
		panic(err)
	}
	return
}

func (dbObj *DB) Find(filter interface{}) (results []Message, err error) {
	cursor, err := dbObj.db.Find(dbObj.ctx, filter)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return
}

func (dbObj *DB) FindAll() (results []Message, err error) {
	cursor, err := dbObj.db.Find(dbObj.ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return
}

func (dbObj *DB) InsertOne(newMessage Message) (err error) {
	_, err = dbObj.db.InsertOne(dbObj.ctx, newMessage, &options.InsertOneOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (dbObj *DB) InsertMany(messages []interface{}) (err error) {
	_, err = dbObj.db.InsertMany(dbObj.ctx, messages)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (dbObj *DB) DeleteOne(filter interface{}) (err error) {
	_, err = dbObj.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return
}

func (dbObj *DB) DeleteMany(filter interface{}) (err error) {
	_, err = dbObj.db.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return
}
