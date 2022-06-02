package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
