package user

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/berkantay/user-management-service/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *model.User) (*string, error) {
	id := uuid.NewString()
	return &id, nil
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return user, nil
}

func (m *mockUserRepository) RemoveUser(ctx context.Context, id string) (*string, error) {
	return &id, nil
}

func (m *mockUserRepository) QueryUsers(ctx context.Context, filter *model.UserQuery) ([]model.User, error) {
	users := []model.User{
		{
			ID:        "123",
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Password:  "password",
		},
	}
	return users, nil
}

func (m *mockUserRepository) HealthCheck(ctx context.Context) error {
	return nil
}

func TestUserServiceCreate(t *testing.T) {
	userService := NewService(&mockUserRepository{}, log.Default())
	testUser := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password",
	}

	ctx := context.Background()
	insertionId, err := userService.Create(ctx, testUser)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}
	if insertionId == nil {
		t.Error("Create returned nil insertion ID")
	}
}

func TestUserServiceUpdate(t *testing.T) {
	userService := NewService(&mockUserRepository{}, log.Default())
	testUser := &model.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password",
	}

	ctx := context.Background()
	updatedUser, err := userService.Update(ctx, testUser)
	if err != nil {
		t.Fatalf("Update returned unexpected error: %v", err)
	}
	if diff := cmp.Diff(testUser, updatedUser); diff != "" {
		t.Errorf("Update returned unexpected result (-want +got):\n%s", diff)
	}
}

func TestUserServiceRemove(t *testing.T) {
	userService := NewService(&mockUserRepository{}, log.Default())

	ctx := context.Background()
	removedId, err := userService.Remove(ctx, "123")
	if err != nil {
		t.Fatalf("Remove returned unexpected error: %v", err)
	}
	if *removedId != "123" {
		t.Errorf("Remove returned unexpected ID: %s", *removedId)
	}
}

func TestHashPassword(t *testing.T) {
	tests := []string{
		"password",
		"this_is_a_longer_password",
	}

	for _, test := range tests {
		t.Run("check hash", func(t *testing.T) {
			hashed, err := hashPassword(test)
			if err != nil {
				t.Errorf("hashPassword(%s) error: %s", test, err)
			}
			// Check that the password can be successfully verified
			err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(test))
			if err != nil {
				t.Errorf("hashPassword(%s) produced invalid hash: %s", test, err)
			}
		})
	}
}
