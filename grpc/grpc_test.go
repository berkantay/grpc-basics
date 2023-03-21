package grpc

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"testing"

	pb "github.com/berkantay/user-management-service/grpc/proto"
	"github.com/berkantay/user-management-service/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type UserServiceMock struct{}

var lis *bufconn.Listener

func (s UserServiceMock) Create(ctx context.Context, user *model.User) (*string, error) {
	// Return a mock user ID.
	return stringPtr("123"), nil
}

func (s UserServiceMock) Update(ctx context.Context, user *model.User) (*model.User, error) {
	// Return the input user as is.
	if user.ID == "test-id" {
		return user, nil

	}

	user.ID = "wrong-test-id"

	return user, errors.New("mock mismatch id")
}

func (s UserServiceMock) Delete(ctx context.Context, userId string) (*string, error) {
	// Return the input user ID as is.
	if userId == "test-id" {
		return stringPtr("test-id"), nil
	}
	return nil, errors.New("mock mismatch id")
}

type EventPublisherMock struct{}

func (e EventPublisherMock) Publish(topic string, payload []byte) error {
	// Return the input user ID as is.
	return nil
}

func (s UserServiceMock) Query(ctx context.Context, query *model.UserQuery) ([]model.User, error) {
	// Return a mock list of users.
	return []model.User{
		{
			ID:        "123",
			FirstName: "User-1",
			LastName:  "Lastname-1",
			NickName:  "u1",
			Password:  "PASSWD123",
			Email:     "test@gmail.com",
			Country:   "Turkey",
		},
		{
			ID:        "4321",
			FirstName: "User-2",
			LastName:  "Lastname-2",
			NickName:  "u2",
			Password:  "123qwert",
			Email:     "mock@gmail.com",
			Country:   "UK",
		},
	}, nil
}

func stringPtr(s string) *string {
	return &s
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserAPIServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestHealthCheck(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserAPIClient(conn)
	resp, err := client.HealthCheck(ctx, &pb.HealthcheckRequest{})
	if err != nil {
		t.Fatalf("Healthcheck failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
}

func TestCreateUserValidMail(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Create(context.Background(), &pb.CreateUserRequest{
		FirstName: "Cagatay",
		LastName:  "Ay",
		NickName:  "excalibur",
		Password:  "test123",
		Email:     "berkantay.5@gmail.com",
		Country:   "zimbambwe",
	})

	if err != nil {
		t.Fatalf("User failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
}

func TestCreateUserInvalidMail(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Create(context.Background(), &pb.CreateUserRequest{
		FirstName: "Cagatay",
		LastName:  "Ay",
		NickName:  "excalibur",
		Password:  "test123",
		Email:     "berkantay.5mail.com",
		Country:   "zimbambwe",
	})

	assert.Equal(t, resp, &pb.CreateUserResponse{
		Status: &pb.Status{
			Code:    "INVALID_ARGUMENT",
			Message: "Invalid email.",
		},
	})

	assert.NotNil(t, err)

}

func TestDeleteUserValidID(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Delete(context.Background(), &pb.DeleteUserRequest{
		Id: "test-id",
	})

	assert.Nil(t, err)

	assert.Equal(t, resp, &pb.DeleteUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User deleted.",
		},
		UserIdResponse: &pb.UserIdResponse{
			Id: "test-id",
		},
	})
}

func TestDeleteUserInvalidID(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Delete(context.Background(), &pb.DeleteUserRequest{
		Id: "123",
	})

	assert.Nil(t, err)

	assert.Equal(t, resp, &pb.DeleteUserResponse{
		Status: &pb.Status{
			Code:    "INTERNAL",
			Message: "Could not delete user. User not found",
		},
	})
}

func TestUpdateUserValidID(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Update(context.Background(), &pb.UpdateUserRequest{
		Id:        "test-id",
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johndoe",
		Password:  "password123",
		Email:     "johndoe@example.com",
		Country:   "US",
	})

	assert.Nil(t, err)

	assert.Equal(t, resp.Status, &pb.Status{
		Code:    "OK",
		Message: "User updated.",
	})

}

func TestUpdateUserInvalidID(t *testing.T) {
	mockUserService := &UserServiceMock{}
	mockEventPublisher := &EventPublisherMock{}

	logger := log.New(nil, "User Management Server Log | ", log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	s := NewServer(mockUserService, mockEventPublisher, logger)

	resp, err := s.Update(context.Background(), &pb.UpdateUserRequest{
		Id:        "test-id123",
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johndoe",
		Password:  "password123",
		Email:     "johndoe@example.com",
		Country:   "US",
	})

	assert.NotNil(t, err)

	assert.Equal(t, resp, &pb.UpdateUserResponse{
		Status: &pb.Status{
			Code:    "INTERNAL",
			Message: "Could not update user.",
		},
	})
}

func TestCreateUserRequestToUser(t *testing.T) {
	req := &pb.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johndoe",
		Password:  "password123",
		Email:     "johndoe@example.com",
		Country:   "US",
	}

	expectedUser := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johndoe",
		Password:  "password123",
		Email:     "johndoe@example.com",
		Country:   "US",
	}

	user := createUserRequestToUser(req)

	assert.Equal(t, expectedUser, user)
}

func TestToUserQuery(t *testing.T) {
	req := &pb.QueryUsersRequest{
		Id:        stringPtr("1"),
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		NickName:  stringPtr("johndoe"),
		Email:     stringPtr("johndoe@example.com"),
		Country:   stringPtr("US"),
		Page:      nil,
		Size:      nil,
	}

	expectedUserQuery := &model.UserQuery{
		ID:        stringPtr("1"),
		FirstName: stringPtr("John"),
		LastName:  stringPtr("Doe"),
		NickName:  stringPtr("johndoe"),
		Email:     stringPtr("johndoe@example.com"),
		Country:   stringPtr("US"),
		Page:      new(int64),
		Size:      new(int64),
	}

	userQuery := toUserQuery(req)

	assert.Equal(t, expectedUserQuery.FirstName, userQuery.FirstName)
	assert.Equal(t, expectedUserQuery.LastName, userQuery.LastName)
	assert.Equal(t, expectedUserQuery.NickName, userQuery.NickName)
	assert.Equal(t, expectedUserQuery.Email, userQuery.Email)
	assert.Equal(t, expectedUserQuery.Country, userQuery.Country)
	assert.Equal(t, int64(1), *userQuery.Page)
	assert.Equal(t, int64(10), *userQuery.Size)
}
