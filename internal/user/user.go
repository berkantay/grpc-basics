package user

import (
	"context"
	"fmt"
	"time"

	"github.com/berkantay/user-management-service/internal/model"
	"github.com/berkantay/user-management-service/pkg/encryption"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	RemoveUserById(filter any) error
	QueryUsers(filter *model.UserQuery) ([]model.User, error)
	HealthCheck(ctx context.Context) error
	GracefullShutdown() error
}

type Service struct {
	db UserRepository
}

func NewService(db UserRepository) *Service {

	return &Service{
		db: db,
	}
}

func (app *Service) CreateUser(user *model.User) error {

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashed, err := encryption.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	err = app.db.CreateUser(user)

	if err != nil {
		return err
	}

	return nil

}

func (app *Service) UpdateUser(user *model.User) error {

	err := app.db.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (app *Service) RemoveUser(userId string) error {

	err := app.db.RemoveUserById(userId)
	if err != nil {
		return err
	}

	return nil
}

func (app *Service) QueryUsers(query *model.UserQuery) ([]model.User, error) {

	users, err := app.db.QueryUsers(query)

	if err != nil {
		return nil, err
	}

	fmt.Println(users)

	return users, nil
}

func (app *Service) DatabaseHealthCheck(ctx context.Context) error {

	err := app.db.HealthCheck(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (app *Service) Echo(ctx context.Context) error {

	fmt.Println("Echo back")

	return nil
}

func (app *Service) GracefullShutdown() error {

	err := app.db.GracefullShutdown()

	if err != nil {
		return err
	}

	return nil

}
