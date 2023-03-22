package user

import (
	"context"
	"log"
	"time"

	"github.com/berkantay/user-management-service/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*string, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) (*string, error)
	QueryUsers(ctx context.Context, filter *model.UserQuery) ([]model.User, error)
}

type Service struct {
	db     UserRepository
	logger *log.Logger
}

// Create new user service.
func NewService(db UserRepository, logger *log.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

// Fills necessary informations and creates user.
func (service *Service) Create(ctx context.Context, user *model.User) (*string, error) {
	service.logger.Printf("INFO:Create operation started.")
	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashed, err := hashPassword(user.Password)
	if err != nil {
		service.logger.Printf("ERROR:Could not hash password[%s]", err)
		return nil, err
	}
	user.Password = hashed

	insertionId, err := service.db.CreateUser(ctx, user)
	if err != nil {
		service.logger.Printf("ERROR:Create operation failed [%s]", err)
		return nil, err
	}
	service.logger.Printf("INFO:Create operation done.")
	service.logger.Printf("INFO:User created with id[%s]", *insertionId)
	return insertionId, nil

}

// Updates user.
func (service *Service) Update(ctx context.Context, user *model.User) (*model.User, error) {
	service.logger.Printf("INFO:Update operation started.")
	update, err := service.db.UpdateUser(ctx, user)
	if err != nil {
		service.logger.Printf("ERROR:Could not update user[%s]", err)
		return nil, err
	}
	service.logger.Printf("INFO:Update operation done.")
	service.logger.Printf("INFO:User updated with [%s]", update)
	return update, nil
}

// Remove user by given id.
func (service *Service) Delete(ctx context.Context, userId string) (*string, error) {
	service.logger.Printf("INFO:Delete operation started.")
	id, err := service.db.DeleteUser(ctx, userId)
	if err != nil {
		service.logger.Printf("ERROR:Could not delete user[%s]", err)
		return nil, err
	}
	service.logger.Printf("INFO:Delete operation done.")
	service.logger.Printf("INFO:User deleted with id[%s]", userId)
	return id, nil
}

// Query user for the given UserQuery, return list of users after query operation and error if exists.
func (service *Service) Query(ctx context.Context, query *model.UserQuery) ([]model.User, error) {
	service.logger.Printf("INFO:Query operation started.")
	users, err := service.db.QueryUsers(ctx, query)
	if err != nil {
		service.logger.Printf("ERROR:User could not queried[%s]", err)
		return nil, err
	}
	service.logger.Printf("INFO:Query operation done.")
	return users, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
