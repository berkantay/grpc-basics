package grpcserver

import (
	"testing"

	pb "github.com/berkantay/user-management-service/internal/adapters/driving/proto"
	"github.com/berkantay/user-management-service/internal/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestCreateUserRequestToUser(t *testing.T) {
	req := &pb.CreateUserRequest{
		FirstName: "john",
		LastName:  "doe",
		NickName:  "johndoe",
		Password:  "password123",
		Email:     "john.doe@example.com",
		Country:   "US",
	}

	user := createUserRequestToUser(req)

	assert.Equal(t, "John", user.FirstName, "Unexpected first name")
	assert.Equal(t, "Doe", user.LastName, "Unexpected last name")
	assert.Equal(t, "johndoe", user.NickName, "Unexpected nickname")
	assert.Equal(t, "password123", user.Password, "Unexpected password")
	assert.Equal(t, "john.doe@example.com", user.Email, "Unexpected email")
	assert.Equal(t, "US", user.Country, "Unexpected country")
}

func TestUpdateUserRequestToUser(t *testing.T) {
	// create a sample request
	req := &pb.UpdateUserRequest{
		Id:        "usertest",
		FirstName: "john",
		LastName:  "doe",
		NickName:  "johndoe",
		Password:  "password",
		Email:     "john.doe@example.com",
		Country:   "USA",
	}

	// convert the request to a user model
	user := updateUserRequestToUser(req)

	// verify the user model has the expected values
	if user.ID != "usertest" {
		t.Errorf("unexpected ID, got %s, want %s", user.ID, "usertest")
	}
	if user.FirstName != "John" {
		t.Errorf("unexpected first name, got %s, want %s", user.FirstName, "John")
	}
	if user.LastName != "Doe" {
		t.Errorf("unexpected last name, got %s, want %s", user.LastName, "Doe")
	}
	if user.NickName != "johndoe" {
		t.Errorf("unexpected nickname, got %s, want %s", user.NickName, "johndoe")
	}
	if user.Password != "password" {
		t.Errorf("unexpected password, got %s, want %s", user.Password, "password")
	}
	if user.Email != "john.doe@example.com" {
		t.Errorf("unexpected email, got %s, want %s", user.Email, "john.doe@example.com")
	}
	if user.Country != "USA" {
		t.Errorf("unexpected country, got %s, want %s", user.Country, "USA")
	}
}

func TestToPbQueryResponse(t *testing.T) {
	// Define sample input data
	users := []model.User{
		{
			ID:        "123",
			FirstName: "John",
			LastName:  "Doe",
			NickName:  "jdoe",
			Password:  "password123",
			Email:     "jdoe@example.com",
			Country:   "USA",
		},
		{
			ID:        "456",
			FirstName: "Jane",
			LastName:  "Doe",
			NickName:  "jadoe",
			Password:  "password456",
			Email:     "jadoe@example.com",
			Country:   "Canada",
		},
	}

	id := "123"

	req := &pb.QueryUsersRequest{
		Id:   &id,
		Page: proto.Int64(1),
		Size: proto.Int64(10),
	}

	// Call the function being tested
	resp := toPbQueryResponse(users, req)

	// Assert the result
	if resp.Status == nil || resp.Status.Code != "OK" || resp.Status.Message != "Users queried" {
		t.Errorf("Unexpected response status: %v", resp.Status)
	}

	if len(resp.Payload) != 2 {
		t.Errorf("Unexpected number of users in response: %v", len(resp.Payload))
	}

	// Check the first user in the response
	u1 := resp.Payload[0]
	if u1.Id != "123" || u1.FirstName != "John" || u1.LastName != "Doe" ||
		u1.NickName != "jdoe" || u1.Password != "password123" || u1.Email != "jdoe@example.com" ||
		u1.Country != "USA" {
		t.Errorf("Unexpected user in response: %v", u1)
	}

	// Check the second user in the response
	u2 := resp.Payload[1]
	if u2.Id != "456" || u2.FirstName != "Jane" || u2.LastName != "Doe" ||
		u2.NickName != "jadoe" || u2.Password != "password456" || u2.Email != "jadoe@example.com" ||
		u2.Country != "Canada" {
		t.Errorf("Unexpected user in response: %v", u2)
	}
}
