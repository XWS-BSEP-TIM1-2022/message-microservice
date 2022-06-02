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
