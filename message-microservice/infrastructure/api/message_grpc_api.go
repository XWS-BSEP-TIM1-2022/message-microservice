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
	chatService         *application.ChatService
	messageService      *application.MessageService
}

func NewMessageHandler(service *application.MessageService, notificationService *application.NotificationService, chatService *application.ChatService, messageService *application.MessageService) *MessageHandler {
	return &MessageHandler{
		service:             service,
		notificationService: notificationService,
		chatService:         chatService,
		messageService:      messageService,
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

func (handler *MessageHandler) GetAllChatsForUser(ctx context.Context, in *messageService.UserIdRequest) (*messageService.GetAllChatsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllChats")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := in.UserId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	chats, err := handler.chatService.GetAllByUserId(ctx, objectId)

	if err != nil {
		return nil, err
	}
	response := &messageService.GetAllChatsResponse{
		Chat: []*messageService.Chat{},
	}
	for _, chat := range chats {
		current := mapChat(chat)
		response.Chat = append(response.Chat, current)
	}
	return response, nil
}

func (handler *MessageHandler) CreateChat(ctx context.Context, in *messageService.NewChatRequest) (*messageService.GetChatResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "CreateChat")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	chatFromRequest := mapChatPb(in.Chat)
	chat, err := handler.chatService.Create(ctx, chatFromRequest)
	if err != nil {
		return nil, err
	}
	chatPb := mapChat(chat)
	response := &messageService.GetChatResponse{
		Chat: chatPb,
	}
	return response, nil
}

func (handler *MessageHandler) GetAllMessagesForUser(ctx context.Context, in *messageService.ChatIdRequest) (*messageService.GetAllMessagesResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllMessagesForUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := in.ChatId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	messages, err := handler.messageService.GetAllByChatId(ctx, objectId)

	if err != nil {
		return nil, err
	}
	response := &messageService.GetAllMessagesResponse{
		Messages: []*messageService.Message{},
	}
	for _, message := range messages {
		current := mapMessage(message)
		response.Messages = append(response.Messages, current)
	}
	return response, nil
}

func (handler *MessageHandler) CreateMessage(ctx context.Context, in *messageService.NewMessageRequest) (*messageService.GetMessageResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "CreateMessage")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	messageFromRequest := mapMessagePb(in.Message)
	message, err := handler.messageService.Create(ctx, messageFromRequest)
	if err != nil {
		return nil, err
	}
	messagePb := mapMessage(message)
	response := &messageService.GetMessageResponse{
		Message: messagePb,
	}

	objectId, err := primitive.ObjectIDFromHex(message.ChatID)
	if err != nil {
		return nil, err
	}
	chat, _ := handler.chatService.GetChatById(ctx, objectId)
	if chat.Username == message.Username {
		notification := mapMessageToNotification(message, chat.FromUserID, chat.UserID)
		handler.notificationService.Create(ctx, mapNotificationPb(notification), 1)
	}

	if chat.FromUsername == message.Username {
		notification := mapMessageToNotification(message, chat.UserID, chat.FromUserID)
		handler.notificationService.Create(ctx, mapNotificationPb(notification), 1)
	}

	return response, nil
}
