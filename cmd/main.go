package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
	"github.com/berkantay/user-management-service/internal/adapters/driving/grpcserver"
	"github.com/berkantay/user-management-service/internal/user"
)

func main() {

	fmt.Println("url", os.Getenv("MONGO_URL"))

	database, err := storage.NewStorage(
		storage.WithHost(os.Getenv("MONGO_URL")),
	)

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
