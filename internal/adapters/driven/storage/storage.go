package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (s *Storage) CreateUser(user *model.User) error {
	//TODO: user id could be checked whether if exists or not. If exists generate another uuid to keep uniqueness.
	_, err := s.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(user *model.User) error {
	filterID := schemeIDFilter(user.ID)

	updateDocument := bson.M{
		"$set": toUserBson(user),
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

func (s *Storage) QueryUsers(filter *model.UserQuery, numberOfEntry, pageNumber int) ([]model.User, error) {

	var users []model.User

	// limit := int64(numberOfEntry)
	// skip := int64(pageNumber)*limit - limit

	queryFilter := filterBuilder(filter)

	fmt.Println("Query filter is", queryFilter)

	cur, err := s.Collection.Find(context.Background(), queryFilter)

	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		var user model.User

		if err := cur.Decode(&user); err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) GracefullShutdown() error {

	err := s.Client.Disconnect(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func filterBuilder(filter *model.UserQuery) *bson.M {

	f := bson.M{}

	if filter.FirstName != "" {
		f["first_name"] = filter.FirstName
	}
	if filter.LastName != "" {
		f["last_name"] = filter.LastName
	}
	if filter.NickName != "" {
		f["nickname"] = filter.NickName
	}
	if filter.Country != "" {
		f["country"] = filter.Country
	}

	// return &bson.D{{Key: "first_name", Value: filter.FirstName}}

	return &f
}

func schemeIDFilter(filter any) *bson.D {

	return &bson.D{{Key: "_id", Value: filter}}
}

func toUserBson(user *model.User) *bson.D {

	return &bson.D{
		bson.E{Key: "first_name", Value: user.FirstName},
		bson.E{Key: "last_name", Value: user.LastName},
		bson.E{Key: "nickname", Value: user.NickName},
		bson.E{Key: "password", Value: user.Password},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "country", Value: user.Country},
	}
}
