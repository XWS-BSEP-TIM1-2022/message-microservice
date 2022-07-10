package startup

import (
	"context"
	"fmt"
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"io"
	"log"
	"message-microservice/application"
	"message-microservice/infrastructure/api"
	"message-microservice/infrastructure/persistance"
	"message-microservice/model"
	"message-microservice/startup/config"
	"net"
)

type Server struct {
	config      *config.Config
	tracer      otgo.Tracer
	closer      io.Closer
	mongoClient *mongo.Client
}

func NewServer(config *config.Config) *Server {
	tracer, closer := tracer.Init(config.MessageServiceName)
	otgo.SetGlobalTracer(tracer)
	return &Server{
		config: config,
		tracer: tracer,
		closer: closer,
	}
}

func (server *Server) GetTracer() otgo.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) Start() {
	server.mongoClient = server.initMongoClient()
	messageStore := server.initStoreStore(server.mongoClient)
	messageService := server.initMessageService(messageStore, server.config)
	notificationStore := server.initNotificationStoreStore(server.mongoClient)
	notificationService := server.initnotificationService(notificationStore, server.config)
	chatStore := server.initChatStoreStore(server.mongoClient)
	chatService := server.initChatService(chatStore, server.config)
	messagesStore := server.initMessagesStoreStore(server.mongoClient)
	messagesService := server.initMessagesService(messagesStore, server.config)

	messageHandler := server.initMessageHandler(messageService, notificationService, chatService, messagesService)

	server.startGrpcServer(messageHandler)
}

func (server *Server) Stop() {
	log.Println("stopping server")
	server.mongoClient.Disconnect(context.TODO())
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistance.GetClient(server.config.MessageDBHost, server.config.MessageDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) startGrpcServer(messageHandler *api.MessageHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	log.Println(fmt.Sprintf("started grpc server on localhost:%s", server.config.Port))
	messageService.RegisterMessageServiceServer(grpcServer, messageHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initStoreStore(client *mongo.Client) model.MessageStore {
	store := persistance.NewMessageMongoDBStore(client)
	return store
}

func (server *Server) initMessageService(store model.MessageStore, config *config.Config) *application.MessageService {
	return application.NewMessageService(store, config)
}

func (server *Server) initMessageHandler(service *application.MessageService, notificationService *application.NotificationService, chatService *application.ChatService, messageService *application.MessageService) *api.MessageHandler {
	return api.NewMessageHandler(service, notificationService, chatService, messageService)
}

func (server *Server) initNotificationStoreStore(client *mongo.Client) model.NotificationStore {
	store := persistance.NewNotificationMongoDBStore(client)
	return store
}

func (server *Server) initnotificationService(store model.NotificationStore, config *config.Config) *application.NotificationService {
	return application.NewNotificationService(store, config)
}

func (server *Server) initChatStoreStore(client *mongo.Client) model.ChatStore {
	store := persistance.NewChatMongoDBStore(client)
	return store
}

func (server *Server) initChatService(store model.ChatStore, config *config.Config) *application.ChatService {
	return application.NewChatService(store, config)
}

func (server *Server) initMessagesStoreStore(client *mongo.Client) model.MessageStore {
	store := persistance.NewMessageMongoDBStore(client)
	return store
}

func (server *Server) initMessagesService(store model.MessageStore, config *config.Config) *application.MessageService {
	return application.NewMessageService(store, config)
}
