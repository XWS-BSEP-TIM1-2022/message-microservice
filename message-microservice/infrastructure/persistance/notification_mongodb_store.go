package persistance

import (
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"message-microservice/model"
)

const (
	COLLECTION_NOTIF = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func NewNotificationMongoDBStore(client *mongo.Client) model.NotificationStore {
	notifications := client.Database(DATABASE).Collection(COLLECTION_NOTIF)
	return &NotificationMongoDBStore{
		notifications: notifications,
	}
}

func (store *NotificationMongoDBStore) GetAllByUserId(ctx context.Context, userId primitive.ObjectID) ([]*model.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"userId": userId.Hex()}
	return store.filter(ctx, filter)
}

func (store *NotificationMongoDBStore) Create(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "Create")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.notifications.InsertOne(ctx, notification)
	if err != nil {
		return nil, err
	}
	notification.Id = result.InsertedID.(primitive.ObjectID)
	return notification, nil
}

func (store *NotificationMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*model.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "filter")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.notifications.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decode(ctx, cursor)
}

func (store *NotificationMongoDBStore) filterOne(ctx context.Context, filter interface{}) (notification *model.Notification, err error) {
	span := tracer.StartSpanFromContext(ctx, "filterOne")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.notifications.FindOne(ctx, filter)
	err = result.Decode(&notification)
	return
}

func decode(ctx context.Context, cursor *mongo.Cursor) (notifications []*model.Notification, err error) {
	span := tracer.StartSpanFromContext(ctx, "decode")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	for cursor.Next(ctx) {
		var notification model.Notification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}
