package persistance

import (
	"go.mongodb.org/mongo-driver/mongo"
	"message-microservice/model"
)

const (
	COLLECTION_NOTIF = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func NewNotificationMongoDBStore(client *mongo.Client) model.NotificationStore {
	notifications := client.Database(DATABASE).Collection(COLLECTION_NOTIF)
	return &NotificationMongoDBStore{
		notifications: notifications,
	}
}
