package main


import "pastebin/api"

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations() {
	db, err := sql.Open("postgres", "postgres://postgres:pass1234@localhost:5432/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success migration!")
}

func main() {
	// mongoClient, err := db.ConnectToMongoDb(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.DisconnectFromMongoDb(context.Background(), mongoClient)

	// dbObj := db.NewMongoDB(mongoClient, context.Background(), "sample_joca", "novi_messages")

	//runMigrations()
	api.StartApiServerAndPrepareDbConnection();

	// postgresClient, err := db.ConnectToPostgresDb("", "", "")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.DisconnectFromPostgresDb(postgresClient)
	//IGRANJE SA DB

	// dbObj := db.NewPostgresDB(postgresClient)

	// novi := db.User{
	// 	Name:     "OpetNovi",
	// 	Password: "1a1dc91c907325c69271ddf0c944bc72",
	// 	PasteNum: 0,
	// 	DevKey:   "e77989ed21758e78331b20e477fc5582",
	// 	Email:    "dev1@gmail.com",
	// }
	// newId, err := dbObj.CreateUser(&novi)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print("Succes -> ", newId, "\n")
	// parsedUUID, err := uuid.Parse("6316cbc8-f91e-4903-b07d-f4251b3c48f3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// usr, err := dbObj.ReadUser(parsedUUID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = dbObj.DeleteUser(parsedUUID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print(usr, "\n")
	// rows, err := postgresClient.Query("SELECT * FROM Users;")
	// if err != nil {
	// 	panic(err)
	// }
	// for rows.Next() {
	// 	var d int
	// 	var id []uint8
	// 	var a, b, c, e string
	// 	err := rows.Scan(&id, &a, &b, &d, &c, &e)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Print(id, " ", a, " ", b, " ", c, " ", d, " ", e, "\n")
	// }

	// dbObj := db.NewMongoDB(client, "sample_mflix", "movies")

	// //neki random filter
	// title := "Back to the Future"
	// filter := bson.D{{Key: "title", Value: title}}

	// // FindOne - vraca prvi document kome je title = "Back to the Future"
	// result, err := dbObj.FindOne(filter) // kako bi kod koji koristimo bio sto prostiji napravicu wrapper oko najcescih db poziva
	// if err != nil {
	// 	panic(err)
	// }

	// dbObj := db.NewDB(client, "sample_message", "messages")

	// filter := bson.D{{Key: "message_id", Value: "1"}}

	// result, err := dbObj.FindOne(filter)
	// if err != nil {
	// 	panic(err)
	// }

	// newMessage := db.Message{
	// 	MessageID:   "7",
	// 	MessageBody: "New message content",
	// }

	// dbObj.InsertOne(newMessage)

	// filter := bson.D{{Key: "message_id", Value: "7"}}

	// err := dbObj.DeleteMany(filter)
	// if err != nil {
	// 	panic(err)
	// }

	// results, err := dbObj.FindAll()
	// if err != nil {
	// 	panic(err)
	// }

	// filter := bson.M{"message_id": bson.M{"$in": []string{"1", "2"}}}
	// results, err := dbObj.Find(filter, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, result := range results {
	// 	fmt.Printf("%s\n", result) // ispis u json-u
	// }

	// jsonData, err := json.MarshalIndent(result, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", jsonData) // ispis u json-u

}
