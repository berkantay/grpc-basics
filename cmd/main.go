package main

import (
	"context"
	"log"

	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
	"github.com/berkantay/user-management-service/internal/adapters/driving/grpcserver"
	"github.com/berkantay/user-management-service/internal/application"
)

func main() {

	database, err := storage.NewStorage()

	if err != nil {
		log.Fatal(err)
	}

	database.HealthCheck(context.Background())

	defer database.GracefullShutdown()

	application := application.NewApplication(database)

	server := grpcserver.NewServer(application)

	server.Run()
}
