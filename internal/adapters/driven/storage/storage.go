package storage

import (
	"context"
	"fmt"
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
	host       string
	context    context.Context
	client     *mongo.Client
	collection *mongo.Collection
}

// Configure storage by chaning host or context.
type StorageOption func(*Storage)

func WithHost(uri string) StorageOption {

	return func(s *Storage) {

		s.host = uri
	}
}

// Database connection context.
func WithContext(ctx context.Context) StorageOption {
	return func(s *Storage) {
		s.context = ctx
	}
}

// Create new connection to the database instance.
func NewStorage(opts ...StorageOption) (*Storage, error) {

	ctx, _ := context.WithTimeout(context.Background(), dbOperationTimeout)
	// defer cancel()

	s := &Storage{
		host:    "mongodb://127.0.0.1:27017",
		context: ctx,
	}

	for _, opt := range opts {
		opt(s)
	}

	clientOptions := options.Client().ApplyURI(s.host)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	s.client = client

	s.collection = s.createCollection("user", "information")

	return s, nil
}

// Create collection in database
func (s *Storage) createCollection(database, collection string) *mongo.Collection {

	return s.client.Database(database).Collection(collection)
}

// Check if database is alive or not.
func (s *Storage) HealthCheck(ctx context.Context) error {
	err := s.client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Connection 200")
	return nil
}

// Create user in database with given type.
func (s *Storage) CreateUser(user *model.User) (*string, error) {
	//TODO: user id could be checked whether if exists or not. If exists generate another uuid to keep uniqueness.
	res, err := s.collection.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	insertionId := res.InsertedID.(string)

	return &insertionId, nil
}

// Update user in database with given type. Updates only one item.
func (s *Storage) UpdateUser(user *model.User) error {
	filterID := schemeIDFilter(user.ID)

	updateDocument := bson.M{
		"$set": toUserBson(user),
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filterID, updateDocument)

	if result.Err() != nil {
		return result.Err()
	}

	return nil

}

// Remove user in database with corresponding id.
func (s *Storage) RemoveUser(id string) error {

	filterId := schemeIDFilter(id)
	//TODO: user id could be checked whether if exists or not in the collection to inform client.
	_, err := s.collection.DeleteOne(context.Background(), filterId,

		&options.DeleteOptions{
			Comment: fmt.Sprintf("%s document is deleted from the collection", filterId),
		})

	if err != nil {
		return err
	}

	return nil

}

// Query users with a filter.
func (s *Storage) QueryUsers(filter *model.UserQuery) ([]model.User, error) {

	limit := int64(*filter.Size)
	// skip := int64(*filter.Page)*limit - limit

	queryFilter := filterBuilder(filter)

	fmt.Println("Query filter is", queryFilter)

	pipeline := createPipeline(queryFilter, limit)

	cur, err := s.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return fromBsonToUser(results)

}

// Disconnect from database.
func (s *Storage) GracefullShutdown() error {

	err := s.client.Disconnect(context.Background())

	if err != nil {
		return err
	}

	return nil
}

// Creates query pipeline for the aggregation.
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

// Builds the custom filter.
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

// Extract _id from filter.
func schemeIDFilter(filter any) *bson.D {

	return &bson.D{{Key: "_id", Value: filter}}
}

// Converts user struct to BSON Document pointer.
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

// Converts BSON Document to user slice.
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
