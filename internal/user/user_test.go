package user

import (
	"context"
	"testing"

	"github.com/berkantay/user-management-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m MockRepository) CreateUser(user *model.User) (*string, error) {

	mockId := uuid.NewString()

	return &mockId, nil
}
func (m MockRepository) UpdateUser(user *model.User) error {
	return nil
}
func (m MockRepository) RemoveUser(user string) error {
	return nil
}
func (m MockRepository) QueryUsers(filter *model.UserQuery) ([]model.User, error) {

	users := make([]model.User, 0)

	users = append(users, model.User{
		ID:        "testid",
		FirstName: *filter.FirstName,
		Password:  "hashed&passwd",
	})

	return users, nil
}
func (m MockRepository) HealthCheck(ctx context.Context) error {
	return nil
}
func (m MockRepository) GracefullShutdown() error {
	return nil
}

func TestNewService(t *testing.T) {
	mockRepo := MockRepository{}
	userService := NewService(mockRepo)
	assert.NotNil(t, userService, "Service object is nil")
}

func TestServiceCreateUser(t *testing.T) {
	mockRepo := MockRepository{}
	userService := NewService(mockRepo)
	t.Run("given client wants create user", func(t *testing.T) {
		mockUser := &model.User{
			FirstName: "John",
			LastName:  "Doe",
			NickName:  "johndoe",
			Password:  "passwd",
			Email:     "johndoe@gmail.com",
			Country:   "Turkey",
		}
		t.Run("when user created", func(t *testing.T) {
			insertionId, err := userService.CreateUser(mockUser)
			assert.Nil(t, err)
			t.Run("then it should create user return id should not be nil", func(t *testing.T) {
				assert.NotNil(t, insertionId)
			})
		})
	})
}

func TestServiceRemoveUser(t *testing.T) {
	mockRepo := MockRepository{}
	userService := NewService(mockRepo)
	t.Run("given client wants remove user", func(t *testing.T) {
		id := "testid"
		t.Run("then user must be deleted", func(t *testing.T) {
			err := userService.RemoveUser(id)
			assert.Nil(t, err)
		})
	})

}

func TestServiceUpdateUser(t *testing.T) {
	mockRepo := MockRepository{}
	userService := NewService(mockRepo)
	t.Run("given client wants to updateuser", func(t *testing.T) {
		mockUser := &model.User{
			ID:        "testid",
			FirstName: "John",
			LastName:  "Doe",
			NickName:  "johndoe",
			Password:  "passwd",
			Email:     "johndoe@gmail.com",
			Country:   "Turkey",
		}
		t.Run("then user must be updated", func(t *testing.T) {
			err := userService.UpdateUser(mockUser)
			assert.Nil(t, err)
		})
	})
}

func TestServiceQueryUser(t *testing.T) {
	mockRepo := MockRepository{}
	userService := NewService(mockRepo)
	t.Run("given client wants to updateuser", func(t *testing.T) {
		id := "testid"
		filterName := "John"
		mockUser := &model.UserQuery{
			ID:        &id,
			FirstName: &filterName,
		}
		t.Run("then user must be updated", func(t *testing.T) {
			users, err := userService.QueryUsers(mockUser)
			assert.Nil(t, err)
			assert.Equal(t, users[0].FirstName, filterName)
		})
	})

}
