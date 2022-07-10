package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string             `json:"userId" bson:"userId"`
	Username     string             `json:"username"`
	FromUserID   string             `json:"fromUserId" bson:"fromUserID"`
	FromUsername string             `json:"fromUsername"`
}
