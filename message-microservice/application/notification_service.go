package application

import (
	"context"
	"errors"
	"fmt"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message-microservice/model"
	"message-microservice/startup/config"
	"time"
)

type NotificationService struct {
	store      model.NotificationStore
	config     *config.Config
	userClient userService.UserServiceClient
}

func NewNotificationService(store model.NotificationStore, c *config.Config) *NotificationService {
	return &NotificationService{
		store:      store,
		config:     c,
		userClient: services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (service *NotificationService) GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]*model.Notification, error) {
	Log.Info("Get all notifications of user with id: " + userId.Hex())

	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllByUserId(ctx, userId)
}

func (service *NotificationService) Create(ctx context.Context, notification *model.Notification, notificationType int32) (*model.Notification, error) {
	Log.Info("Create notification")

	span := tracer.StartSpanFromContextMetadata(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.userClient.GetRequest(ctx, &userService.UserIdRequest{UserId: notification.UserID})
	if err != nil {
		Log.Error("Create notification. Error: " + err.Error())
		return nil, err
	}

	fromUser, err := service.userClient.GetRequest(ctx, &userService.UserIdRequest{UserId: notification.FromUserID})
	if err != nil {
		Log.Error("Create notification. Error: " + err.Error())
		return nil, err
	}

	notification.Message = "Dear " + user.User.Name + ", "

	if notificationType == 1 {
		notification.Message += "you have a new message from user " + fromUser.User.Name + " " + fromUser.User.Surname + "."
	} else if notificationType == 2 {
		notification.Message += "user " + fromUser.User.Name + " " + fromUser.User.Surname + ", created new post."
	} else if notificationType == 3 {
		notification.Message += "user " + fromUser.User.Name + " " + fromUser.User.Surname + ", created new comment."
	} else {
		return nil, errors.New("not recognize notification type")
	}

	notification.Date = time.Now()

	return service.store.Create(ctx, notification)
}
