package api

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	//"database/sql"
	"pastebin/db"
)

var ConnectorPostresDB  *db.PostgresDB

func StartApiServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")

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
	ConnectorPostresDB = db.NewPostgresDB(postgresClient)


	StartApiServer();
}

