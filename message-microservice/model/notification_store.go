package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationStore interface {
	GetAllByUserId(ctx context.Context, id primitive.ObjectID) ([]*Notification, error)
	Create(ctx context.Context, notification *Notification) (*Notification, error)
}
