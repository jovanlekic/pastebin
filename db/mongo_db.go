package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	ctx       context.Context
	dbName    string
	tableName string
	db        *mongo.Collection
}

const wrkDir = "app"

func ConnectToMongoDb(ctx context.Context) (*mongo.Client, error) {
	re := regexp.MustCompile(`^(.*` + wrkDir + `)`)
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
