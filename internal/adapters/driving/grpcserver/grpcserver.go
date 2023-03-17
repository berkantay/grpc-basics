package grpcserver

import (
	"context"
	"log"
	"net"

	pb "github.com/berkantay/user-management-service/internal/adapters/driving/proto"
	"github.com/berkantay/user-management-service/internal/application"

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

}

// func (s *Server) Close() error {

// }
func (s *Server) Echo(ctx context.Context, req *pb.ResponseRequest) (*pb.ResponseRequest, error) {

	return req, nil
}

// func NewServer(connType, uri string) *Server {

// 	listener, err := net.Listen(connType, uri)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	grpcServer := grpc.NewServer()

// 	userApi := &User{}

// 	pb.RegisterUserApiServer(grpcServer, userApi)

// 	err = grpcServer.Serve(listener)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &Server{
// 		serverInstance: grpcServer,
// 		userApi:        userApi,
// 	}

// }
