package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     string             `json:"userId" bson:"userId"`
	FromUserID string             `json:"fromUserId"`
	Message    string             `json:"message"`
	Date       time.Time          `json:"date"`
}
