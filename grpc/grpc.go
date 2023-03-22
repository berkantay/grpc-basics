package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/mail"

	pb "github.com/berkantay/user-management-service/grpc/proto"
	"github.com/berkantay/user-management-service/model"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"google.golang.org/grpc"
)

const (
	userCreated = "user_created"
	userDeleted = "user_deleted"
	userUpdated = "user_updated"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) (*string, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, userId string) (*string, error)
	Query(ctx context.Context, query *model.UserQuery) ([]model.User, error)
}

type EventPublisher interface {
	Publish(topic string, payload []byte) error
}

type Server struct {
	user UserService
	pb.UnimplementedUserAPIServer
	logger    *log.Logger
	publisher EventPublisher
}

func NewServer(service UserService, publisher EventPublisher, logger *log.Logger) *Server {
	return &Server{
		user:      service,
		logger:    logger,
		publisher: publisher,
	}
}

// Run the gRPC server.
func (s *Server) Run() {
	s.logger.Printf("INFO:gRPC|Connecting tcp socket..")
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		s.logger.Printf("ERROR:gRPC|failed to listen on port 8080 | [%v]", err)
	}

	userManagementService := s

	grpcServer := grpc.NewServer()
	s.logger.Printf("INFO:gRPC|Registering to User API")
	pb.RegisterUserAPIServer(grpcServer, userManagementService)
	s.logger.Printf("INFO:gRPC|Registered to User API")

	s.logger.Printf("INFO:gRPC|Binding TCP Socket")
	grpcServer.Serve(listen)
	s.logger.Printf("INFO:gRPC|Binded")
	defer grpcServer.Stop()

}

// Implements CreateUser function according to proto definition.
func (s *Server) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.logger.Printf("INFO:gRPC|Create called")
	wrappedMessage := createUserRequestToUser(req)
	s.logger.Printf("INFO:gRPC|Request converted to user model")
	s.logger.Printf("INFO:gRPC|Checking if email is valid")
	isValidEmail := checkIsValidMail(req.Email)
	if !isValidEmail {
		s.logger.Printf("WARNING:gRPC|Invalid email")
		return &pb.CreateUserResponse{
			Status: &pb.Status{
				Code:    "INVALID_ARGUMENT",
				Message: "Invalid email.",
			},
		}, errors.New("invalid email")
	}
	insertionId, err := s.user.Create(ctx, wrappedMessage)
	if err != nil {
		s.logger.Printf("ERROR:gRPC|Could not create user. [%s]", err)
		return &pb.CreateUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not create user.",
			},
		}, err
	}

	t := toEventMessage(userCreated, &pb.UserPayload{
		Id: *insertionId,
	})

	go func() {
		userIdByte, err := json.Marshal(t)
		if err != nil {
			s.logger.Printf("ERROR:gRPC|Could not marshal user create. [%s]", err)
		}
		s.publisher.Publish("user", userIdByte)
	}()
	s.logger.Printf("INFO:gRPC|User created.")
	return &pb.CreateUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User created.",
		},
		Payload: &pb.UserPayload{
			Id: *insertionId,
		},
	}, nil
}

// Implements DeleteUser function according to proto definition.
func (s *Server) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.Printf("INFO:gRPC|Delete user called")
	id, err := s.user.Delete(ctx, req.Id)

	if err != nil {
		s.logger.Printf("ERROR:gRPC|Could not delete user. [%s]", err)
		return &pb.DeleteUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not delete user. User not found",
			},
		}, nil //TODO: Check if err should return or not?
	}

	t := toEventMessage(userDeleted, &pb.UserPayload{
		Id: *id,
	})
	go func() {
		userIdByte, err := json.Marshal(t)
		if err != nil {
			s.logger.Printf("ERROR:gRPC|Could not marshal delete. [%s]", err)
		}
		s.publisher.Publish("user", userIdByte)
	}()
	s.logger.Printf("INFO:gRPC|User deleted.")
	return &pb.DeleteUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User deleted.",
		},
		UserIdResponse: &pb.UserIdResponse{
			Id: *id,
		},
	}, nil
}

// Implements UpdateUser function according to proto definition.
func (s *Server) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.logger.Printf("INFO:gRPC|Updat called")
	update, err := s.user.Update(ctx, updateUserRequestToUser(req))
	s.logger.Printf("INFO:gRPC|Checking if email is valid")
	isValidEmail := checkIsValidMail(req.Email)
	if !isValidEmail {
		s.logger.Printf("WARNING:gRPC|Invalid email")
		return &pb.UpdateUserResponse{
			Status: &pb.Status{
				Code:    "INVALID_ARGUMENT",
				Message: "Invalid email.",
			},
		}, errors.New("invalid email")
	}

	if err != nil {
		s.logger.Printf("ERROR:gRPC|Could not update user. [%s]", err)
		return &pb.UpdateUserResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not update user.",
			},
		}, err
	}

	t := toEventMessage(userUpdated, &pb.UserPayload{
		Id:        update.ID,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		NickName:  update.NickName,
		Password:  update.Password,
		Email:     update.Email,
		Country:   update.Country,
	})

	go func() {
		userIdByte, err := json.Marshal(t)
		if err != nil {
			s.logger.Printf("ERROR:gRPC|Could not marshal user update. [%s]", err)
		}
		s.publisher.Publish("user", userIdByte)
	}()
	s.logger.Println("INFO:gRPC|User updated")
	return &pb.UpdateUserResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User updated.",
		}, Payload: toUserUpdatePayload(update),
	}, nil

}

// Implements QueryUsers function according to proto definition.
func (s *Server) Query(ctx context.Context, req *pb.QueryUsersRequest) (*pb.QueryUsersResponse, error) {
	s.logger.Printf("INFO:gRPC|Query called.")
	user, err := s.user.Query(ctx, toUserQuery(req))

	if err != nil {
		s.logger.Printf("ERROR:gRPC|Query error. [%s]", err)
		return &pb.QueryUsersResponse{
			Status: &pb.Status{
				Code:    "INTERNAL",
				Message: "Could not query user",
			},
		}, err
	}

	if user == nil {
		s.logger.Printf("WARNING:gRPC|Could not found user.")
		return &pb.QueryUsersResponse{
			Status: &pb.Status{
				Code:    "NOT_FOUND",
				Message: "Could not found any user.",
			},
		}, err
	}
	s.logger.Printf("INFO:gRPC|Query done.")
	return toPbQueryResponse(user, req), nil
}

// Check if the service is alive or not.
func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthcheckRequest) (*pb.HealthcheckResponse, error) {
	return &pb.HealthcheckResponse{
		Status: &pb.Status{
			Code:    "OK",
			Message: "User service alive.",
		},
	}, nil
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

func toUserUpdatePayload(update *model.User) *pb.UserPayload {
	return &pb.UserPayload{
		Id:        update.ID,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		NickName:  update.NickName,
		Email:     update.Email,
		Country:   update.Country,
	}
}

func checkIsValidMail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func toEventMessage(eventName string, payload any) *model.UserEvent {
	return &model.UserEvent{
		EventName: eventName,
		Payload:   payload,
	}
}
