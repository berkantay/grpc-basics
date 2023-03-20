package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/berkantay/user-management-service/internal/adapters/driving/proto"
	"github.com/berkantay/user-management-service/internal/model"
	"github.com/berkantay/user-management-service/pkg/utility"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"google.golang.org/grpc"
)

type UserService interface {
	CreateUser(user *model.User) (*string, error)
	UpdateUser(user *model.User) error
	RemoveUser(userId string) error
	QueryUsers(query *model.UserQuery) ([]model.User, error)
}

type Server struct {
	service UserService
	pb.UnimplementedUserApiServer
}

func NewServer(service UserService) *Server {

	return &Server{
		service: service,
	}
}

// Run the gRPC server.
func (s *Server) Run() {

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen on port 8080: %v", err)
	}

	userManagementService := s

	grpcServer := grpc.NewServer()

	pb.RegisterUserApiServer(grpcServer, userManagementService)

	grpcServer.Serve(listen)
	defer grpcServer.Stop()

}

// Implements CreateUser function according to proto definition.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	wrappedMessage := createUserRequestToUser(req)

	isValidEmail := utility.CheckIsValidMail(req.Email)

	if !isValidEmail {

		return &pb.CreateUserResponse{
			Status: &pb.Status{
				Code:    "INVALID_ARGUMENT",
				Message: "Invalid email.",
			},
		}, errors.New("invalid email")
	}

	insertionId, err := s.service.CreateUser(wrappedMessage)

	if err != nil {
		return &pb.CreateUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not create user.",
			},
		}, err
	}

	return &pb.CreateUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User created.",
		},
		Payload: &pb.UserPayload{
			Id: *insertionId,
		}, //TODO Fill user info from db
	}, nil
}

// Implements DeleteUser function according to proto definition.
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	fmt.Println("Deleting user", req.Id)

	err := s.service.RemoveUser(req.Id)

	if err != nil {
		return &pb.DeleteUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not delete user.",
			},
		}, err
	}

	return &pb.DeleteUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User deleted.",
		}, //TODO Fill user info from db
	}, nil
}

// Implements UpdateUser function according to proto definition.
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Println("Updating user", req.Id)

	err := s.service.UpdateUser(updateUserRequestToUser(req))

	if err != nil {
		return &pb.UpdateUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not update user.",
			},
		}, err
	}

	return &pb.UpdateUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User updated.",
		}, //TODO Fill user info from db
	}, nil

}

// Implements QueryUsers function according to proto definition.
func (s *Server) QueryUsers(ctx context.Context, req *pb.QueryUsersRequest) (*pb.QueryUsersResponse, error) {

	userQuery := toUserQuery(req)

	user, err := s.service.QueryUsers(userQuery)

	if err != nil {
		return &pb.QueryUsersResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Internal error occured",
			},
		}, err
	}

	if user == nil {
		return &pb.QueryUsersResponse{
			Status: &pb.Status{
				Code:    "NOT_FOUND",
				Message: "Could not found any user.",
			},
		}, err
	}

	return toPbQueryResponse(user, req), nil

}

// Convert protobuf CREATE request structure to User model.
func createUserRequestToUser(req *pb.CreateUserRequest) *model.User { //TODO:move this wrapping layer from server logic

	return &model.User{
		FirstName: cases.Title(language.English, cases.Compact).String(req.FirstName),
		LastName:  cases.Title(language.English, cases.Compact).String(req.LastName),
		NickName:  req.NickName,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}
}

// Convert protobuf UPDATE request structure to User model.
func updateUserRequestToUser(req *pb.UpdateUserRequest) *model.User { //TODO:move this wrapping layer from server logic
	return &model.User{
		ID:        req.Id,
		FirstName: cases.Title(language.English, cases.Compact).String(req.FirstName),
		LastName:  cases.Title(language.English, cases.Compact).String(req.LastName),
		NickName:  req.NickName,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}
}

// Convert protobuf QUERY request structure to User model.
func toUserQuery(req *pb.QueryUsersRequest) *model.UserQuery {
	defaultPage := int64(1)
	defaultSize := int64(10)

	if req.Page == nil {
		req.Page = &defaultPage
	}

	if req.Size == nil {
		req.Size = &defaultSize
	}

	return &model.UserQuery{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		NickName:  req.NickName,
		Email:     req.Email,
		Country:   req.Country,
		Page:      req.Page,
		Size:      req.Size,
	}
}

// Convert User model to protobuf QueryUsersResponse.
func toPbQueryResponse(users []model.User, req *pb.QueryUsersRequest) *pb.QueryUsersResponse {
	payload := make([]*pb.UserPayload, 0)
	for _, u := range users {
		payload = append(payload, &pb.UserPayload{
			Id:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			NickName:  u.NickName,
			Password:  u.Password,
			Email:     u.Email,
			Country:   u.Country,
		})
	}

	return &pb.QueryUsersResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "Users queried",
		},
		Payload: payload,
		Meta: &pb.Meta{
			Page: *req.Page,
			Size: *req.Size,
		},
	}
}
