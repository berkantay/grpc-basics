package main

import (
	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
	"github.com/berkantay/user-management-service/internal/adapters/driving/grpcserver"
	"github.com/berkantay/user-management-service/internal/application"
)

func main() {

	database := storage.NewStorage()

	server := grpcserver.NewServer("tcp", "localhost:8080")

	server.RegisterApiServer()

	application := application.NewApplication(database, server)

	application.Listen()

}
