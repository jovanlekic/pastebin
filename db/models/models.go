package models

import "github.com/google/uuid"

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
	MessageID   string `bson:"message_id"`
	MessageBody string `bson:"message_body"`
}
