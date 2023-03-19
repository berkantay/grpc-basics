package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/berkantay/user-management-service/internal/adapters/driving/proto"
	"github.com/berkantay/user-management-service/internal/application"
	"github.com/berkantay/user-management-service/internal/model"
	"github.com/google/uuid"

	"google.golang.org/grpc"
)

type ServerRepository interface {
}

type Server struct {
	app application.ApplicationRepository
	pb.UnimplementedUserApiServer
}

func NewServer(app application.ApplicationRepository) *Server {

	return &Server{
		app: app,
	}
}

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

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

	wrappedMessage := toUser(req)

	err := s.app.AddUser(wrappedMessage)

	if err != nil {
		return &pb.UserResponse{
			ReturnCode: 404,
		}, err
	}

	return &pb.UserResponse{
		ReturnCode: 200,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.UserResponse, error) {

	fmt.Println("Deleting user", req.Id)

	err := s.app.RemoveUser(req.Id)

	if err != nil {
		return &pb.UserResponse{
			ReturnCode: 404,
		}, err
	}

	return &pb.UserResponse{
		ReturnCode: 200,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	fmt.Println("Updating user", req.Id)

	err := s.app.UpdateUser(req.Id, req.Update)

	if err != nil {
		return &pb.UserResponse{
			ReturnCode: 404,
		}, err
	}

	return &pb.UserResponse{
		ReturnCode: 200,
	}, nil

}

func toUser(req *pb.CreateUserRequest) *model.UserInfo { //TODO:move this wrapping layer from server logic

	return &model.UserInfo{
		UUID:      uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		NickName:  req.NickName,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
