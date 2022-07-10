package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStore interface {
	GetAllByChatId(ctx context.Context, id primitive.ObjectID) ([]*Message, error)
	Create(ctx context.Context, message *Message) (*Message, error)
}
