package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID   uuid.UUID
	Name     string
	Password string
	PasteNum int
	DevKey   string
	Email    string
}

type Object struct {
	PasteKey  string
	DevKey    string
	MessageID string
}

type Message struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MessageBody string             `bson:"message_body"`
}
