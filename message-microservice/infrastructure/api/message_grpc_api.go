package api

import (
	"context"
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message-microservice/application"
	"sort"
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

func (handler *MessageHandler) GetAllNotifications(ctx context.Context, in *messageService.UserIdRequest) (*messageService.GetAllResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllNotifications")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := in.UserId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	notifications, err := handler.notificationService.GetAllByUserId(ctx, objectId)
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[j].Date.Before(notifications[i].Date)
	})

	if err != nil {
		return nil, err
	}
	response := &messageService.GetAllResponse{
		Notifications: []*messageService.Notification{},
	}
	for _, notification := range notifications {
		current := mapNotification(notification)
		response.Notifications = append(response.Notifications, current)
	}
	return response, nil
}

func (handler *MessageHandler) CreateNotification(ctx context.Context, in *messageService.NewNotificationRequest) (*messageService.GetResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "CreateNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	notificationFromRequest := mapNotificationPb(in.Notification)
	notification, err := handler.notificationService.Create(ctx, notificationFromRequest, in.NotificationType)
	if err != nil {
		return nil, err
	}
	notificationPb := mapNotification(notification)
	response := &messageService.GetResponse{
		Notification: notificationPb,
	}
	return response, nil
}
