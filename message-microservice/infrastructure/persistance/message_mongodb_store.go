package persistance

import (
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (store *MessageMongoDBStore) GetAllByChatId(ctx context.Context, chatId primitive.ObjectID) ([]*model.Message, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByChatId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"chatId": chatId.Hex()}
	return store.filter(ctx, filter)
}

func (store *MessageMongoDBStore) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.messages.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}
	message.Id = result.InsertedID.(primitive.ObjectID)
	return message, nil
}

func (store *MessageMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*model.Message, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.messages.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decodeMessage(ctx, cursor)
}

func (store *MessageMongoDBStore) filterOne(ctx context.Context, filter interface{}) (message *model.Message, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.messages.FindOne(ctx, filter)
	err = result.Decode(&message)
	return
}

func decodeMessage(ctx context.Context, cursor *mongo.Cursor) (messages []*model.Message, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	for cursor.Next(ctx) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil {
			return
		}
		messages = append(messages, &message)
	}
	err = cursor.Err()
	return
}
