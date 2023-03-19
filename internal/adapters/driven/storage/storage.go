package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/berkantay/user-management-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbOperationTimeout = 10 * time.Second
)

type UserRepository interface {
	AddUser(T any) error
	UpdateUser(filter string, update any) error
	RemoveUserById(filter any) error
	GetUserByFilter(T any) error
	HealthCheck(ctx context.Context) error
	GracefullShutdown() error
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

func NewStorage(opts ...StorageOption) (*Storage, error) {

	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	// defer cancel()

	s := &Storage{
		Host:    "mongodb://127.0.0.1",
		Port:    27017,
		Context: ctx,
	}

	for _, opt := range opts {
		opt(s)
	}

	uri := s.Host + ":" + strconv.Itoa(s.Port)

	fmt.Println("URI is:", uri)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected db instance..")

	s.Client = client

	s.Collection = s.Client.Database("user").Collection("information")

	return s, nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {

	err := s.Client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Connection 200")
	return nil
}

func (s *Storage) AddUser(T any) error {
	//TODO: user id could be checked whether if exists or not. If exists generate another uuid to keep uniqueness.
	_, err := s.Collection.InsertOne(context.Background(), T)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(filter string, update any) error {

	filterID := schemeIDFilter(filter)

	bsonUpdate, err := toUserBson(update)

	if err != nil {
		return err
	}

	updateDocument := bson.M{
		"$set": bsonUpdate,
	}

	result := s.Collection.FindOneAndUpdate(context.Background(), filterID, updateDocument)

	if result.Err() != nil {
		return result.Err()
	}

	return nil

}

func (s *Storage) RemoveUserById(filter any) error {

	if filter == nil {

		return errors.New("nil filter")
	}

	filterId := schemeIDFilter(filter)
	//TODO: user id could be checked whether if exists or not in the collection to inform client.
	_, err := s.Collection.DeleteOne(context.Background(), filterId,

		&options.DeleteOptions{
			Comment: fmt.Sprintf("%s document is deleted from the collection", filterId),
		})

	if err != nil {
		return err
	}

	return nil

}

func (s *Storage) GetUserByFilter(filter any) error {

	var info model.UserInfo

	data := s.Collection.FindOne(context.Background(), filter)

	return data.Decode(&info)
}

func (s *Storage) GracefullShutdown() error {

	err := s.Client.Disconnect(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func schemeIDFilter(filter any) *bson.D {

	return &bson.D{{Key: "_id", Value: filter}}
}

func toUserBson(update any) (*bson.D, error) {

	updateDoc := &bson.D{}

	data, err := bson.Marshal(update)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(data, &updateDoc)

	if err != nil {
		return nil, err
	}

	return updateDoc, nil

}
