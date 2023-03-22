package main

import (
	"context"
	"log"
	"os"

	"github.com/berkantay/user-management-service/broker"
	"github.com/berkantay/user-management-service/database"
	"github.com/berkantay/user-management-service/grpc"
	"github.com/berkantay/user-management-service/user"
)

var Version = "development"

func main() {
	file, err := os.OpenFile("user-management-service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := log.New(file, "User Management Server Log | ", log.LstdFlags)
	logger.Printf("User Management Service [%s]", Version)
	database, err := database.NewStorage(
		database.WithHost(os.Getenv("MONGO_URL")),
		database.WithLogger(logger),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = database.HealthCheck(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer database.GracefullShutdown(context.Background())

	application := user.NewService(database, logger)

	publisher, err := broker.NewBrokerHandler(logger)

	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(application, publisher, logger)

	server.Run()
}
