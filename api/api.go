package api

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	//"database/sql"
	"pastebin/db"
	"context"
	//"go.mongodb.org/mongo-driver/mongo"
	//"os"
)

var ConnectorPostgresDB  *db.PostgresDB
var ConnectorMongoDB  *db.MongoDB

func StartApiServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/check", ValidateJWTToken(ChekerHandler)).Methods("GET")
	r.HandleFunc("/api/checkandparse", ChekerHandlerParseToken).Methods("GET")
	r.HandleFunc("/api/createPaste", CreatePaste).Methods("POST")
	r.HandleFunc("/api/getPaste/{pasteKey}", GetPaste).Methods("GET")
	r.HandleFunc("/api/deletePaste", DeletePaste).Methods("POST")


	
	log.Println("Server started on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}


func StartApiServerAndPrepareDbConnection(){
	postgresClient, err := db.ConnectToPostgresDb("mydb", "postgres", "pass1234")
	if err != nil {
		log.Println(err)
		return
	}
	// ovde treba videti gde pozvati ovo za diskonektovanje sa baze
	defer db.DisconnectFromPostgresDb(postgresClient)
	ConnectorPostgresDB = db.NewPostgresDB(postgresClient)

	

	mongoClient, errM :=  db.ConnectToMongoDb(context.Background())
	if errM != nil {
		log.Println(err)
		return
	}
	defer db.DisconnectFromMongoDb(context.Background(), mongoClient)

	ConnectorMongoDB = db.NewMongoDB(mongoClient, context.Background(), "pastes", "messages")

	// err = ConnectorMongoDB.db.Drop(context.Background())
	// if err != nil {
	// 	log.Println(err)
	// }


	StartApiServer();
}

