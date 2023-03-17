package application

import (
	"context"
	"fmt"

	"github.com/berkantay/user-management-service/internal/adapters/driven/storage"
)

type ApplicationRepository interface {
	AddUser(T any) error
	UpdateUser(filter, update any) error
	RemoveUser(filter any) error
	GetUserByFilter(T any) error
	DatabaseHealthCheck(ctx context.Context) error
	Echo(ctx context.Context) error
	GracefullShutdown() error
}

type Application struct {
	db storage.UserRepository
}

func NewApplication(db storage.UserRepository) *Application {

	return &Application{
		db: db,
	}
}

func (app *Application) AddUser(T any) error {

	err := app.db.AddUser(T)

	if err != nil {
		return err
	}

	return nil

}

func (app *Application) UpdateUser(filter, update any) error {

	err := app.db.UpdateUser(filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (app *Application) RemoveUser(filter any) error {

	err := app.db.RemoveUser(filter)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) GetUserByFilter(T any) error {

	err := app.db.GetUserByFilter(T)

	if err != nil {
		return err
	}

	return nil
}

func (app *Application) DatabaseHealthCheck(ctx context.Context) error {

	err := app.db.HealthCheck(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (app *Application) Echo(ctx context.Context) error {

	fmt.Println("Echo back")

	return nil
}

func (app *Application) GracefullShutdown() error {

	err := app.db.GracefullShutdown()

	if err != nil {
		return err
	}

	return nil

}
