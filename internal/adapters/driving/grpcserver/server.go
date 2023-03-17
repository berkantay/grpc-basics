package grpcserver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/berkantay/user-management-service/internal/adapters/driving/proto"

	"google.golang.org/grpc"
)

type ServerRepository interface {
	Listen() error
}

type Server struct {
	pb.UnimplementedUserApiServer
	listener *net.Listener
	instance *grpc.Server
}

func NewServer(connType, uri string) *Server {

	listener, err := net.Listen(connType, uri)

	if err != nil {
		fmt.Println(err)
	}

	grpcServer := grpc.NewServer()

	return &Server{
		listener: &listener,
		instance: grpcServer,
	}
}

func (s *Server) RegisterApiServer() {

	pb.RegisterUserApiServer(s.instance, &Server{})

	fmt.Print("User API registered to gRPC Server...")
}

func (s *Server) Listen() error {

	err := s.instance.Serve(*s.listener)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Echo(ctx context.Context, req *pb.ResponseRequest) (*pb.ResponseRequest, error) {

	return req, nil
}
