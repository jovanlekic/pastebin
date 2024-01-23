package api

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	//"database/sql"
	"pastebin/db"
	"pastebin/kgs"
	"context"
	//"go.mongodb.org/mongo-driver/mongo"
	//"os"
)

var ConnectorPostgresDB  *db.PostgresDB
var ConnectorMongoDB  *db.MongoDB
var KgsPasteKeys  kgs.KGS
var KgsDevKeys kgs.KGS

func StartApiServer() {
	r := mux.NewRouter()

	//corsOpts := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // Set your frontend origin here

	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/check", ValidateJWTToken(ChekerHandler)).Methods("GET")
	r.HandleFunc("/api/checkandparse", ChekerHandlerParseToken).Methods("GET")
	r.HandleFunc("/api/createPaste", CreatePaste).Methods("POST")
	r.HandleFunc("/api/getPaste/{pasteKey}", GetPaste).Methods("GET")
	r.HandleFunc("/api/deletePaste", DeletePaste).Methods("POST")


	
	log.Println("Server started on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
    	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
    	handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type", "Authorization"}),
		)(r))
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


	// add KGS for pastekeys
	postgresClientKgsPasteKey, errP := db.ConnectToPostgresDb("pastekeys", "postgres", "pass1234")
	if errP != nil {
		log.Println(errP)
		return
	}
	// ovde treba videti gde pozvati ovo za diskonektovanje sa baze
	defer db.DisconnectFromPostgresDb(postgresClientKgsPasteKey)
	KgsPasteKeys = kgs.GetInstance(postgresClientKgsPasteKey)


	// add KGS for devkeys
	postgresClientKgsDevKey, errD := db.ConnectToPostgresDb("devkeys", "postgres", "pass1234")
	if errD != nil {
		log.Println(errD)
		return
	}
	// ovde treba videti gde pozvati ovo za diskonektovanje sa baze
	defer db.DisconnectFromPostgresDb(postgresClientKgsDevKey)
	KgsDevKeys = kgs.GetInstance(postgresClientKgsDevKey)

	log.Println("Uspesna konekcija ostvarena na svim bazama!")

	StartApiServer();
}

