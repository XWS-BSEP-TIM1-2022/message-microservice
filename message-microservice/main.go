package main

import (
	"message-microservice/startup"
	cfg "message-microservice/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	defer server.Stop()
}
