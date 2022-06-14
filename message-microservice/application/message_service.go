package application

import (
	"github.com/sirupsen/logrus"
	"message-microservice/model"
	"message-microservice/startup/config"
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
