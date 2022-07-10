package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatStore interface {
	GetAllByUserId(ctx context.Context, id primitive.ObjectID) ([]*Chat, error)
	GetChatById(ctx context.Context, id primitive.ObjectID) (*Chat, error)
	Create(ctx context.Context, chat *Chat, userId primitive.ObjectID, fromUser primitive.ObjectID) (*Chat, error)
}
