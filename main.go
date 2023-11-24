package main

import (
	"pastebin/db"
)

func main() {
	mongoClient := db.ConnectToMongoDb()
	defer db.DisconnectFromMongoDb(mongoClient)

	postgresClient := db.ConnectToPostgresDb()
	defer db.DisconnectFromPostgresDb(postgresClient)
	//IGRANJE SA DB

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

	// dbObj := db.NewDB(client, "sample_mflix", "movies")

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
