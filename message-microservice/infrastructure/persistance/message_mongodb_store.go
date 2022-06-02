package persistance

import (
	"go.mongodb.org/mongo-driver/mongo"
	"message-microservice/model"
)

const (
	DATABASE   = "messagesDB"
	COLLECTION = "messages"
)

type MessageMongoDBStore struct {
	messages *mongo.Collection
}

func NewMessageMongoDBStore(client *mongo.Client) model.MessageStore {
	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDBStore{
		messages: messages,
	}
}
