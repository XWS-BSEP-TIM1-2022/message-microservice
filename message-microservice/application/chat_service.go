package application

import (
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message-microservice/model"
	"message-microservice/startup/config"
)

type ChatService struct {
	store  model.ChatStore
	config *config.Config
}

func NewChatService(store model.ChatStore, config *config.Config) *ChatService {
	return &ChatService{
		store:  store,
		config: config,
	}
}

func (service *ChatService) GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]*model.Chat, error) {
	Log.Info("Get all chats of user with id: " + userId.Hex())

	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllByUserId(ctx, userId)
}

func (service *ChatService) Create(ctx context.Context, chat *model.Chat) (*model.Chat, error) {
	Log.Info("Create chat")

	span := tracer.StartSpanFromContextMetadata(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	userId, err := primitive.ObjectIDFromHex(chat.UserID)
	if err != nil {
		return nil, err
	}

	forUserId, err := primitive.ObjectIDFromHex(chat.FromUserID)
	if err != nil {
		return nil, err
	}

	return service.store.Create(ctx, chat, userId, forUserId)
}
