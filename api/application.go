package api

import (
	"log"
	"net/http"
	"encoding/json"
	"pastebin/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
)


// tester function
func ChekerHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Authorization passed")
	w.WriteHeader(http.StatusOK)
}

// tester function
func ChekerHandlerParseToken(w http.ResponseWriter, r *http.Request){
	mapClaims, error := ParseAccesToken(r)
	if error != nil {
		http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
		log.Println("Unauthorized access: Try to access " + r.URL.String())
		return
	}
	
	username := mapClaims["username"].(string)
	devkey := mapClaims["devkey"].(string)


	log.Println("Authorization passed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"TokenInfo" : map[string]interface{}{
			"username": username,
			"devkey":	devkey,
		},
	})
}


func CreatePaste(w http.ResponseWriter, r *http.Request){
	mapClaims, error := ParseAccesToken(r)
	if error != nil {
		http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
		log.Println("Unauthorized access: Try to access " + r.URL.String())
		return
	}
	//usernameToken := mapClaims["username"].(string)
	devKeyToken := mapClaims["devkey"].(string)
	

	var requestData models.Paste

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// pasteKey is not mandatory
	if requestData.Message == "" || requestData.DevKey == "" {
		log.Println("Bad request for deleting paste: insufficient number of fields")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if devKeyToken != requestData.DevKey {
		log.Println(devKeyToken)
			log.Println(requestData.DevKey)
		http.Error(w,"Not valid data!", http.StatusUnauthorized)
		log.Println("Unauthorized function: Try to access " + r.URL.String() + " with different devkey")
		return
	}

	pastekey, errKey := KgsPasteKeys.Check(requestData.PasteKey)
	if errKey != nil {
		http.Error(w,"Error: Cannot create paste", http.StatusInternalServerError)
		log.Println("Error: Cannot create key for paste: "+ requestData.PasteKey + ": " + errKey.Error())
		return 
	}

	// first create message
	messageId, errMsg := ConnectorMongoDB.CreateMessage(requestData.Message)
	if errMsg != nil {
		http.Error(w,"Error: Cannot create paste", http.StatusInternalServerError)
		log.Println("Error: Cannot create message for paste: "+ requestData.PasteKey + "!")
		return 
	}

	

	newObject := models.Object{
		PasteKey: 	 pastekey,
		DevKey: 	requestData.DevKey,
		MessageID: 	messageId,
	}
	

	errObj := ConnectorPostgresDB.CreateObject(context.Background(), &newObject)
	if errObj != nil {
		idM , _ := primitive.ObjectIDFromHex(messageId)
		ConnectorMongoDB.DeleteMessage(idM)
		http.Error(w,"Error: Cannot create paste", http.StatusInternalServerError)
		log.Println("Error: Cannot create object for paste: "+ requestData.PasteKey + "!")
		return 
	}

	w.WriteHeader(http.StatusCreated)
	data,_ := json.Marshal(map[string]interface{}{"PasteKey": newObject.PasteKey})
	w.Write(data)

}


func GetPaste(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pasteKey := vars["pasteKey"]

	object, errObj := ConnectorPostgresDB.ReadObjectWithoutDevKey(context.Background(), pasteKey)
	if errObj != nil {
		http.Error(w,"Paste not found", http.StatusNotFound)
		log.Println("Error: paste"+ pasteKey + " not found!")
		return 
	}

	messageId, errMes := primitive.ObjectIDFromHex(object.MessageID)
	if errMes != nil {
		http.Error(w,"Error", http.StatusInternalServerError)
		log.Println("Error: Cannot convert from string to primitive.ObjectId")
		return 
	}

	message, errMsg := ConnectorMongoDB.ReadMessage(messageId);
	if errMsg!= nil {
		http.Error(w,"Error: Cannot retrieve paste", http.StatusInternalServerError)
		log.Println("Error: Cannot retrieve paste: "+ pasteKey + "!")
		return 
	}
	
	w.WriteHeader(http.StatusOK)
	data,_ := json.Marshal(map[string]interface{}{"Message": message.MessageBody})
	w.Write(data)
}


func DeletePaste(w http.ResponseWriter, r *http.Request){
	mapClaims, error := ParseAccesToken(r)
	if error != nil {
		http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
		log.Println("Unauthorized access: Try to access " + r.URL.String())
		return
	}

	//usernameToken := mapClaims["username"].(string)
	devKeyToken := mapClaims["devkey"].(string)

	var requestData models.DeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if requestData.PasteKey == "" || requestData.DevKey == "" {
		log.Println("Bad request for deleting paste: insufficient number of fields")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if devKeyToken != requestData.DevKey {
		log.Println(devKeyToken)
		log.Println(requestData.DevKey)

		http.Error(w,"Not valid data!", http.StatusUnauthorized)
		log.Println("Unauthorized function: Try to access " + r.URL.String() + " with different devkey")
		return
	}


	// now call function to get Object 
	object, errObj := ConnectorPostgresDB.ReadObject(context.Background(), requestData.PasteKey, requestData.DevKey)
	if errObj != nil {
		http.Error(w,"Not valid data!", http.StatusBadRequest)
		log.Println("Error: User devkey: " + requestData.DevKey + " tried to delete paste: " + requestData.PasteKey + " but paste doesnt exist or he is not authorized!")
		return 
	}

	// delete message from MongoDb
	messageId, errMes := primitive.ObjectIDFromHex(object.MessageID)
	if errMes != nil {
		http.Error(w,"Error", http.StatusInternalServerError)
		log.Println("Error: Cannot convert from string to primitive.ObjectId")
		return 
	}

	if deletedError := ConnectorMongoDB.DeleteMessage(messageId); deletedError != nil{
		http.Error(w,"Error: Cannot delete paste", http.StatusBadRequest)
		log.Println("Error: Cannot delete message: "+ object.MessageID + "!")
		return 
	}

	// delete object from PostgresDb
	errObj = ConnectorPostgresDB.DeleteObject(context.Background(), requestData.PasteKey, requestData.DevKey)
	if errObj != nil {
		http.Error(w,"Error: Cannot delete paste", http.StatusInternalServerError)
		log.Println("Error: Cannot delete paste: "+ requestData.PasteKey + "!")
		return 
	}

	
	w.WriteHeader(http.StatusAccepted)
}

func getUserInfo(w http.ResponseWriter, r *http.Request){
	mapClaims, error := ParseAccesToken(r)
	if error != nil {
		http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
		log.Println("Unauthorized access: Try to access " + r.URL.String())
		return
	}

	username := mapClaims["username"].(string)
	//devkey := mapClaims["devkey"].(string)

	user, err := ConnectorPostgresDB.ReadUserByUsername(context.Background(), username);
	if err!=nil{
		log.Println(err)
		http.Error(w,"Error: user doesn't exist", http.StatusNotFound)
		return
	} 

	w.WriteHeader(http.StatusOK)
	data,_ := json.Marshal(map[string]interface{}{
		"username": user.Name,
		"email": user.Email,
		"pastenum": user.PasteNum,
		"devkey": user.DevKey,
	})
	w.Write(data)
}