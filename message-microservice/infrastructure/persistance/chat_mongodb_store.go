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
	COLLECTION_CHAT = "chat"
)

type ChatMongoDBStore struct {
	chats *mongo.Collection
}

func NewChatMongoDBStore(client *mongo.Client) model.ChatStore {
	chats := client.Database(DATABASE).Collection(COLLECTION_CHAT)
	return &ChatMongoDBStore{
		chats: chats,
	}
}

func (store *ChatMongoDBStore) GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]*model.Chat, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"userId": userId.Hex()}
	users, err1 := store.filter(ctx, filter)
	if err1 != nil {
		return nil, err1
	}
	forFilter := bson.M{"fromUserID": userId.Hex()}
	forUser, err2 := store.filter(ctx, forFilter)
	if err2 != nil {
		return nil, err2
	}
	return append(users, forUser...), nil
}

func (store *ChatMongoDBStore) GetChatById(ctx context.Context, chatId primitive.ObjectID) (*model.Chat, error) {
	span := tracer.StartSpanFromContext(ctx, "GetChatById")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": chatId}
	users, err1 := store.filterOne(ctx, filter)
	if err1 != nil {
		return nil, err1
	}
	return users, nil
}

func (store *ChatMongoDBStore) Create(ctx context.Context, chat *model.Chat, userId primitive.ObjectID, fromUser primitive.ObjectID) (*model.Chat, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"userId": userId.Hex(), "fromUserID": fromUser.Hex()}
	chat1, err1 := store.filterOne(ctx, filter)
	if err1 != nil {
		forFilter := bson.M{"userId": fromUser.Hex(), "fromUserID": userId.Hex()}
		chat2, err := store.filterOne(ctx, forFilter)
		if err != nil {
			result, err := store.chats.InsertOne(ctx, chat)
			if err != nil {
				return nil, err
			}
			chat.Id = result.InsertedID.(primitive.ObjectID)
			return chat, nil
		}
		return chat2, nil
	}
	return chat1, nil
}

func (store *ChatMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*model.Chat, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.chats.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decodeChat(ctx, cursor)
}

func (store *ChatMongoDBStore) filterOne(ctx context.Context, filter interface{}) (chat *model.Chat, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.chats.FindOne(ctx, filter)
	err = result.Decode(&chat)
	return
}

func decodeChat(ctx context.Context, cursor *mongo.Cursor) (chats []*model.Chat, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	for cursor.Next(ctx) {
		var chat model.Chat
		err = cursor.Decode(&chat)
		if err != nil {
			return
		}
		chats = append(chats, &chat)
	}
	err = cursor.Err()
	return
}
