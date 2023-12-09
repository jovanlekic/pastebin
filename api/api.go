package api

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)



func StartApiServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/logout", LogoutHandler).Methods("POST")

	log.Println("Server started on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

