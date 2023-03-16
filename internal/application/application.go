package application

import "github.com/berkantay/user-management-service/internal/adapters/driven/storage"

type Application struct {
	db storage.UserRepository
}

func NewApplication(db storage.UserRepository) *Application {

	return &Application{
		db: db,
	}
}
