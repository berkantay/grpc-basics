package database

import (
	"context"
	"log"
	"time"

	"github.com/berkantay/user-management-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Storage struct {
	host       string
	context    context.Context
	client     *mongo.Client
	collection *mongo.Collection
	logger     *log.Logger
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

func WithLogger(logger *log.Logger) StorageOption {
	return func(s *Storage) {
		s.logger = logger
	}
}

// Create new connection to the database instance.
func NewStorage(opts ...StorageOption) (*Storage, error) {
	s := &Storage{
		host:    "mongodb://127.0.0.1:27017",
		context: context.Background(),
	}

	for _, opt := range opts {
		opt(s)
	}

	clientOptions := options.Client().ApplyURI(s.host)
	s.logger.Printf("INFO:MongoDB|Connecting..")
	client, err := mongo.Connect(s.context, clientOptions)
	s.logger.Printf("INFO:MongoDB|Connected..")
	if err != nil {
		s.logger.Printf("ERROR:MongoDB|Could not connected [%s]", err)
		return nil, err
	}
	s.client = client
	s.logger.Printf("INFO:MongoDB|Creating collection..")
	s.collection = s.createCollection("user", "information")
	s.logger.Printf("INFO:MongoDB|Created collection..")
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
		s.logger.Printf("ERROR:MongoDB|Database is NOT alive..")
		return err
	}
	s.logger.Printf("INFO:MongoDB|Database alive..")
	return nil
}

// Create user in database with given type.
func (s *Storage) CreateUser(ctx context.Context, user *model.User) (*string, error) {
	//TODO: user id could be checked whether if exists or not. If exists generate another uuid to keep uniqueness.
	s.logger.Printf("INFO:MongoDB|Creating user.")
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		s.logger.Printf("ERROR:MongoDB|Could not create user. [%s]", err)
		return nil, err
	}
	insertionId := res.InsertedID.(string)
	s.logger.Printf("INFO:MongoDB|Creation successful.[%s]", insertionId)
	return &insertionId, nil
}

// Update user in database with given type. Updates only one item.
func (s *Storage) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	s.logger.Printf("INFO:MongoDB|Updating user")
	filterID := schemeIDFilter(user.ID)
	s.logger.Printf("INFO:MongoDB|Filtering on = [%s]", filterID)
	updateDocument := bson.M{
		"$set": toUserBson(user),
	}
	result := s.collection.FindOneAndUpdate(ctx, filterID, updateDocument)
	if result.Err() != nil {
		s.logger.Printf("ERROR:MongoDB|Update error is  [%s]", result.Err())
		return nil, result.Err()
	}
	update := model.User{}
	result.Decode(&update)
	s.logger.Printf("INFO:MongoDB|Update successful user. [%s]", update.ID)
	return user, nil
}

// Delete user in database with corresponding id.
func (s *Storage) DeleteUser(ctx context.Context, id string) (*string, error) {
	s.logger.Printf("INFO:MongoDB|Deleting user with id:[%s]", id)
	filterId := schemeIDFilter(id)
	res := s.collection.FindOneAndDelete(ctx, filterId)
	if res.Err() != nil {
		s.logger.Printf("ERROR:MongoDB|Could not delete [%s] error is: [%s]", id, res.Err())
		return nil, res.Err()
	}
	s.logger.Printf("INFO:MongoDB|Delete successful.")
	return &id, nil
}

// Query users with a filter.
func (s *Storage) QueryUsers(ctx context.Context, filter *model.UserQuery) ([]model.User, error) {
	s.logger.Printf("INFO:MongoDB|Querying user with given filter")
	var results []bson.M

	limit := int64(*filter.Size)
	skip := int64(*filter.Page)*limit - limit

	pipeline := createPipeline(filterBuilder(filter), limit, skip)
	cur, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Printf("ERROR:MongoDB|Aggregation error [%s]", err)
		return nil, err
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		s.logger.Printf("ERROR:MongoDB|Cursor error [%s]", err)
		return nil, err
	}

	res, err := fromBsonToUser(results)

	if err != nil {
		s.logger.Printf("ERROR:MongoDB|Could not convert BSON to User model. Error is: [%s]", err)
		return nil, err
	}
	s.logger.Printf("INFO:MongoDB|Query result")
	return res, nil
}

// Disconnect from database.
func (s *Storage) GracefullShutdown(ctx context.Context) error {
	s.logger.Printf("INFO:MongoDB|Shutting down..")
	err := s.client.Disconnect(ctx)
	s.logger.Printf("INFO:MongoDB|Closed.")
	if err != nil {
		return err
	}

	return nil
}

// Creates query pipeline for the aggregation.
func createPipeline(filter *bson.D, limit, skip int64) []bson.M {

	pipeline := []bson.M{
		{"$match": *filter},
		{"$facet": bson.M{
			"metadata": []bson.M{
				{"$count": "total"},
				{"$limit": limit},
				{"$skip": skip},
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
		titleCased := cases.Title(language.English, cases.Compact).String(*filter.FirstName)
		f = append(f, bson.E{Key: "first_name", Value: titleCased})
	}
	if filter.LastName != nil {
		titleCased := cases.Title(language.English, cases.Compact).String(*filter.LastName)
		f = append(f, bson.E{Key: "last_name", Value: titleCased})
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
		bson.E{Key: "updated_at", Value: time.Now()},
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
