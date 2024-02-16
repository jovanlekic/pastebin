package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// look what is the meaning of `json:"username" binding:"required"`

type UserLogin struct {  // for communication between frontend and servers
	Username 	string `json:"username"`
	Password 	string `json:"password"`
}

type UserRegistration struct { // for communication between frontend and servers
	Username 	string `json:"username" binding:"required"`
	Password 	string `json:"password" binding:"required"`
	Email 	 	string `json:"email" binding:"required"`
}

type Paste struct{
	DevKey	 	string `json:"devkey"`
	PasteKey	string `json:"pastekey"`
	Message 	string `json:"message"`
}

type DeleteRequest struct{
	DevKey  string	`json:"devkey"`
	PasteKey    string	`json:"pastekey"`
}


// communication with relational PostgreSQL database
type User struct { // for communication between api servers and database
	UserID   uuid.UUID
	Name     string
	Password string
	PasteNum int
	DevKey   string
	Email    string
}

// communication with relational PostgreSQL database
type Object struct { // for communication between api servers and database
	PasteKey  string
	DevKey    string
	MessageID string
}

// communication with non-relational Mongo database
type Message struct { // for communication between api servers and database
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MessageBody string             `bson:"message_body"`
}

