package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pastebin/db"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDb() *mongo.Client {
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

func DisconnectFromDb(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	fmt.Print("Disconnected!\n")
}

func main() {
	client := ConnectToDb()
	defer DisconnectFromDb(client)

	//IGRANJE SA DB

	dbObj := db.NewDB(client, "sample_mflix", "movies")

	//neki random filter
	title := "Back to the Future"
	filter := bson.D{{Key: "title", Value: title}}

	// FindOne - vraca prvi document kome je title = "Back to the Future"
	result, err := dbObj.FindOne(filter, nil) // kako bi kod koji koristimo bio sto prostiji napravicu wrapper oko najcescih db poziva
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s\n", result) // ispis u bson.M

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData) // ispis u json-u
}
