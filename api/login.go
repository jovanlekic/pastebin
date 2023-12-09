package api

import (
	"encoding/json"
	"log"
	"net/http"
)


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// send request to db


	w.WriteHeader(http.StatusNoContent)
}



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}


	// send request to base to check if the user actually exist
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	
}