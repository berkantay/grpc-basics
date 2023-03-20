package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/berkantay/user-management-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Host:    "mongodb://mongodb",
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

func (s *Storage) QueryUsers(filter *model.UserQuery) ([]model.User, error) {

	limit := int64(*filter.Size)
	// skip := int64(*filter.Page)*limit - limit

	queryFilter := filterBuilder(filter)

	fmt.Println("Query filter is", queryFilter)

	pipeline := createPipeline(queryFilter, limit)

	cur, err := s.Collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return fromBsonToUser(results)

}

func (s *Storage) GracefullShutdown() error {

	err := s.Client.Disconnect(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func createPipeline(filter *bson.D, limit int64) []bson.M {
	pipeline := []bson.M{
		{"$match": *filter},
		{"$facet": bson.M{
			"metadata": []bson.M{
				{"$count": "total"},
			},
			"data": []bson.M{
				{"$limit": limit},
			},
		}},
	}
	return pipeline
}

func filterBuilder(filter *model.UserQuery) *bson.D {

	f := bson.D{}

	if filter.FirstName != nil {
		f = append(f, bson.E{Key: "first_name", Value: *filter.FirstName})
	}
	if filter.LastName != nil {
		f = append(f, bson.E{Key: "last_name", Value: *filter.LastName})
	}
	if filter.NickName != nil {
		f = append(f, bson.E{Key: "nickname", Value: *filter.NickName})
	}
	if filter.Country != nil {
		f = append(f, bson.E{Key: "country", Value: *filter.Country})
	}

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

func fromBsonToUser(queriedUser []primitive.M) ([]model.User, error) {

	var result []model.User

	for _, q := range queriedUser {
		for _, d := range q["data"].(primitive.A) {

			result = append(result, model.User{
				ID:        d.(primitive.M)["_id"].(string),
				FirstName: d.(primitive.M)["first_name"].(string),
				LastName:  d.(primitive.M)["last_name"].(string),
				NickName:  d.(primitive.M)["nickname"].(string),
				Password:  d.(primitive.M)["password"].(string),
				Email:     d.(primitive.M)["email"].(string),
				Country:   d.(primitive.M)["country"].(string),
				CreatedAt: d.(primitive.M)["created_at"].(primitive.DateTime).Time(),
				UpdatedAt: d.(primitive.M)["updated_at"].(primitive.DateTime).Time(),
			})

		}
	}

	return result, nil
}
