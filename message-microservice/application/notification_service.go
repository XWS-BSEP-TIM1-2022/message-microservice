package application

import (
	"message-microservice/model"
	"message-microservice/startup/config"
)

type NotificationService struct {
	store  model.NotificationStore
	config *config.Config
}

func NewNotificationService(store model.MessageStore, config *config.Config) *NotificationService {
	return &NotificationService{
		store:  store,
		config: config,
	}
}
