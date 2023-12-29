package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pastebin/models"
	"context"
)

func makeDevKey() string{
	return "devkey1";
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUserReg models.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&newUserReg); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// checking if all data is sent
	if newUserReg.Username == "" || newUserReg.Password == "" || newUserReg.Email == "" {
		log.Println("Bad request for registration: insufficient number of fields")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// ovde bih mozda uradio md5 ili neku hes funkciju na pasvordu ali to mozemo i posle
	
	newUser := models.User{
		Name: 		newUserReg.Username,
		Password:	newUserReg.Password,
		PasteNum: 	0,
		DevKey: 	makeDevKey(),
		Email: 		newUserReg.Email,
	}


	if _, err := ConnectorPostgresDB.CreateUser(context.Background(), &newUser); err!=nil{
		log.Println(err)
		http.Error(w,"Impossible to register", http.StatusBadRequest)
		return
	} 

	w.WriteHeader(http.StatusCreated)
}



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// check if all data is sent
	if loginRequest.Username == "" || loginRequest.Password == "" {
		log.Println("Bad request for login: insufficient number of fields")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := ConnectorPostgresDB.ReadUserByUsername(context.Background(), loginRequest.Username);
	if err!=nil{
		log.Println(err)
		http.Error(w,"Cannot login with these credentials", http.StatusBadRequest)
		return
	} 
	
	// check if user password is ok
	if user.Password != loginRequest.Password {
		log.Println("Error: " + loginRequest.Username + " tried to login with bad password")
		http.Error(w,"Bad credentials", http.StatusBadRequest)
		return
	}


	// make jwt token and send back to user
	newToken, err := CreateNewToken(user.Name, user.DevKey);
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return;
	}


	w.WriteHeader(http.StatusAccepted)
	data,_ := json.Marshal(map[string]interface{}{"Token": newToken})
	w.Write(data)

}
