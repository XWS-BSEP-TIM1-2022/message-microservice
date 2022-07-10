package api

import (
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message-microservice/model"
	"strconv"
	"strings"
	"time"
)

func mapNotification(notification *model.Notification) *messageService.Notification {
	notificationPb := &messageService.Notification{
		Id:         notification.Id.Hex(),
		UserId:     notification.UserID,
		FromUserId: notification.FromUserID,
		Message:    notification.Message,
		Date:       notification.Date.String(),
	}
	return notificationPb
}

func mapNotificationPb(notificationPb *messageService.Notification) *model.Notification {
	id, _ := primitive.ObjectIDFromHex(notificationPb.Id)
	t := time.Now()
	if notificationPb.Date != "" {
		dateString := strings.Split(notificationPb.Date, " ")
		date := strings.Split(dateString[0], "-")
		year, _ := strconv.Atoi(date[0])
		month, _ := strconv.Atoi(date[1])
		day, _ := strconv.Atoi(date[2])

		timeString := strings.Split(dateString[1], ":")
		hour, _ := strconv.Atoi(timeString[0])
		minutes, _ := strconv.Atoi(timeString[1])
		t = time.Date(year, time.Month(month), day, hour, minutes, 0, 0, time.UTC)
	}
	notification := &model.Notification{
		Id:         id,
		UserID:     notificationPb.UserId,
		FromUserID: notificationPb.FromUserId,
		Message:    notificationPb.Message,
		Date:       t,
	}
	return notification
}

func mapChat(chat *model.Chat) *messageService.Chat {
	chatPb := &messageService.Chat{
		Id:           chat.Id.Hex(),
		UserId:       chat.UserID,
		FromUserId:   chat.FromUserID,
		Username:     chat.Username,
		FromUsername: chat.FromUsername,
	}
	return chatPb
}

func mapChatPb(chatPb *messageService.Chat) *model.Chat {
	id, _ := primitive.ObjectIDFromHex(chatPb.Id)
	chat := &model.Chat{
		Id:           id,
		UserID:       chatPb.UserId,
		FromUserID:   chatPb.FromUserId,
		Username:     chatPb.Username,
		FromUsername: chatPb.FromUsername,
	}
	return chat
}

func mapMessage(message *model.Message) *messageService.Message {
	messagePb := &messageService.Message{
		Id:       message.Id.Hex(),
		ChatId:   message.ChatID,
		Message:  message.Message,
		Date:     message.Date.String(),
		Username: message.Username,
	}
	return messagePb
}

func mapMessagePb(messagePb *messageService.Message) *model.Message {
	id, _ := primitive.ObjectIDFromHex(messagePb.Id)
	t := time.Now()
	if messagePb.Date != "" {
		dateString := strings.Split(messagePb.Date, " ")
		date := strings.Split(dateString[0], "-")
		year, _ := strconv.Atoi(date[0])
		month, _ := strconv.Atoi(date[1])
		day, _ := strconv.Atoi(date[2])

		timeString := strings.Split(dateString[1], ":")
		hour, _ := strconv.Atoi(timeString[0])
		minutes, _ := strconv.Atoi(timeString[1])
		t = time.Date(year, time.Month(month), day, hour, minutes, 0, 0, time.UTC)
	}
	message := &model.Message{
		Id:       id,
		ChatID:   messagePb.ChatId,
		Message:  messagePb.Message,
		Date:     t,
		Username: messagePb.Username,
	}
	return message
}

func mapMessageToNotification(message *model.Message, id string, fromUser string) *messageService.Notification {
	messagePb := &messageService.Notification{
		Message:    "You have new message from " + message.Username,
		Date:       message.Date.String(),
		UserId:     id,
		FromUserId: fromUser,
	}
	return messagePb
}
