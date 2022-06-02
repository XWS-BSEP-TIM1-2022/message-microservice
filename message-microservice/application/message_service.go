package application

import (
	"message-microservice/model"
	"message-microservice/startup/config"
)

type MessageService struct {
	store  model.MessageStore
	config *config.Config
}

func NewMessageService(store model.MessageStore, config *config.Config) *MessageService {
	return &MessageService{
		store:  store,
		config: config,
	}
}
