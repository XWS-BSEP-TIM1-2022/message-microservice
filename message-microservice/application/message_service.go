package application

import (
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message-microservice/model"
	"message-microservice/startup/config"
	"time"
)

type MessageService struct {
	store  model.MessageStore
	config *config.Config
}

var Log = logrus.New()

func NewMessageService(store model.MessageStore, config *config.Config) *MessageService {
	return &MessageService{
		store:  store,
		config: config,
	}
}

func (service *MessageService) GetAllByChatId(ctx context.Context, chatId primitive.ObjectID) ([]*model.Message, error) {
	Log.Info("Get all messages of chat with id: " + chatId.Hex())

	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllByChatId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllByChatId(ctx, chatId)
}

func (service *MessageService) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	Log.Info("Create message")

	span := tracer.StartSpanFromContextMetadata(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	message.Date = time.Now()

	return service.store.Create(ctx, message)
}
