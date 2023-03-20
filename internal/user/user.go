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
	CreateUser(user *model.User) (*string, error)
	UpdateUser(user *model.User) error
	RemoveUser(id string) error
	QueryUsers(filter *model.UserQuery) ([]model.User, error)
	HealthCheck(ctx context.Context) error
	GracefullShutdown() error
}

type Service struct {
	db UserRepository
}

// Create new user service.
func NewService(db UserRepository) *Service {

	return &Service{
		db: db,
	}
}

// Fills necessary informations and creates user.
func (app *Service) CreateUser(user *model.User) (*string, error) {

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashed, err := encryption.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed

	insertionId, err := app.db.CreateUser(user)

	if err != nil {
		return nil, err
	}

	return insertionId, nil

}

// Updates user.
func (app *Service) UpdateUser(user *model.User) error {

	err := app.db.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

// Remove user by given id.
func (app *Service) RemoveUser(userId string) error {

	err := app.db.RemoveUser(userId)
	if err != nil {
		return err
	}

	return nil
}

// Query user for the given UserQuery, return list of users after query operation and error if exists.
func (app *Service) QueryUsers(query *model.UserQuery) ([]model.User, error) {
	users, err := app.db.QueryUsers(query)

	if err != nil {
		return nil, err
	}

	fmt.Println(users)

	return users, nil
}

func (app *Service) Echo(ctx context.Context) error {

	fmt.Println("Echo back")

	return nil
}
