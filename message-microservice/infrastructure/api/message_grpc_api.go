package api

import (
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"message-microservice/application"
)

type MessageHandler struct {
	messageService.UnimplementedMessageServiceServer
	service             *application.MessageService
	notificationService *application.NotificationService
}

func NewMessageHandler(service *application.MessageService, notificationService *application.NotificationService) *MessageHandler {
	return &MessageHandler{
		service:             service,
		notificationService: notificationService,
	}
}
