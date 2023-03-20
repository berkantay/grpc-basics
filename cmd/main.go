package main

import (
	"context"
	"log"

	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
	"github.com/berkantay/user-management-service/internal/adapters/driving/grpcserver"
	"github.com/berkantay/user-management-service/internal/user"
)

func main() {

	database, err := storage.NewStorage()

	if err != nil {
		log.Fatal(err)
	}

	err = database.HealthCheck(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	defer database.GracefullShutdown()

	application := user.NewService(database)

	server := grpcserver.NewServer(application)

	server.Run()
}
