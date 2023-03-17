package application

import (
	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
	"github.com/berkantay/user-management-service/internal/adapters/driving/grpcserver"
)

type Application struct {
	db     storage.UserRepository
	server grpcserver.ServerRepository
}

func NewApplication(db storage.UserRepository, server grpcserver.ServerRepository) *Application {

	return &Application{
		db:     db,
		server: server,
	}
}

func (app *Application) Listen() {

	app.server.Listen()

}
