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
	RemoveUser(ctx context.Context, id string) (*string, error)
	QueryUsers(ctx context.Context, filter *model.UserQuery) ([]model.User, error)
	HealthCheck(ctx context.Context) error
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
	service.logger.Printf("Create called with[%s]", *user)
	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	service.logger.Printf("Wrapped to[%s]", *user)

	hashed, err := hashPassword(user.Password)
	if err != nil {
		service.logger.Printf("Could not hash password[%s]", err)
		return nil, err
	}
	user.Password = hashed

	insertionId, err := service.db.CreateUser(ctx, user)
	if err != nil {
		service.logger.Printf("Could not create user [%s]", err)
		return nil, err
	}
	service.logger.Printf("User created with id[%s]", *insertionId)

	return insertionId, nil

}

// Updates user.
func (service *Service) Update(ctx context.Context, user *model.User) (*model.User, error) {
	service.logger.Printf("Update called with[%s]", *user)
	update, err := service.db.UpdateUser(ctx, user)
	if err != nil {
		service.logger.Printf("Could not update user[%s]", err)
		return nil, err
	}
	service.logger.Printf("User updated with [%s]", update)

	return update, nil
}

// Remove user by given id.
func (service *Service) Remove(ctx context.Context, userId string) (*string, error) {
	service.logger.Printf("Remove called with[%s]", userId)
	id, err := service.db.RemoveUser(ctx, userId)
	service.logger.Printf("User removed with id[%s]", userId)
	if err != nil {
		service.logger.Printf("Could not remove user[%s]", err)
		return nil, err
	}
	return id, nil
}

// Query user for the given UserQuery, return list of users after query operation and error if exists.
func (service *Service) Query(ctx context.Context, query *model.UserQuery) ([]model.User, error) {
	service.logger.Printf("Query called for")
	users, err := service.db.QueryUsers(ctx, query)
	if err != nil {
		service.logger.Printf("User could not queried[%v]", &query)
		return nil, err
	}
	return users, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
