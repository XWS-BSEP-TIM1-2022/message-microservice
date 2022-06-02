package config

import (
	"os"
)

type Config struct {
	Port               string
	MessageDBHost      string
	MessageDBPort      string
	MessageServiceName string
	UserServiceHost    string
	UserServicePort    string
}

func NewConfig() *Config {
	return &Config{
		Port:               getEnv("MESSAGE_SERVICE_PORT", "8089"),
		MessageDBHost:      getEnv("MESSAGE_DB_HOST", "dislinkt:WiYf6BvFmSpJS2Ob@xws.cjx50.mongodb.net/messagesDB"),
		MessageDBPort:      getEnv("MESSAGE_DB_PORT", ""),
		MessageServiceName: getEnv("MESSAGE_SERVICE_NAME", "message_service"),
		UserServiceHost:    getEnv("USER_SERVICE_HOST", "localhost"),
		UserServicePort:    getEnv("USER_SERVICE_PORT", "8085"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
