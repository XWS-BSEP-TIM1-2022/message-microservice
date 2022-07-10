package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChatID   string             `json:"chatId" bson:"chatId"`
	Message  string             `json:"message"`
	Date     time.Time          `json:"date"`
	Username string             `json:"username"`
}
