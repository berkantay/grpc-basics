package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/berkantay/user-management-service/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbOperationTimeout = 10 * time.Second
)

type UserRepository interface {
	Connect(ctx context.Context) error
	AddUser(userInfo *model.UserInfo) error
	UpdateUser(filter, update any) error
	RemoveUser(filter any) error
	GetUserByFilter(T any) error
	HealthCheck(ctx context.Context) error
}

type Storage struct {
	Host       string
	Port       int
	Context    context.Context
	Client     *mongo.Client
	Collection *mongo.Collection
}

type StorageOption func(*Storage)

func WithHost(uri string) StorageOption {

	return func(s *Storage) {

		s.Host = uri
	}
}

func WithPort(port int) StorageOption {
	return func(s *Storage) {
		s.Port = port
	}
}

func WithContext(ctx context.Context) StorageOption {
	return func(s *Storage) {
		s.Context = ctx
	}
}

func NewStorage(opts ...StorageOption) *Storage {

	ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
	defer cancel()

	s := &Storage{
		Host:    "mongodb://localhost",
		Port:    27017,
		Context: ctx,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Storage) Connect(ctx context.Context) error {

	uri := s.Host + ":" + strconv.Itoa(s.Port)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	s.Client = client

	collection := s.Client.Database("user").Collection("information")

	s.Collection = collection

	return nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	err := s.Client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddUser(user *model.UserInfo) error {

	_, err := s.Collection.InsertOne(s.Context, user)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(filter, update any) error {

	_, err := s.Collection.UpdateOne(s.Context, filter, update)

	if err != nil {
		return err
	}

	return nil

}

func (s *Storage) RemoveUser(filter any) error {

	_, err := s.Collection.DeleteOne(s.Context, filter)

	if err != nil {
		return err
	}

	return nil

}

func (s *Storage) GetUserByFilter(filter any) error {

	var info model.UserInfo

	data := s.Collection.FindOne(s.Context, filter)

	return data.Decode(&info)
}
