package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	UUID      primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	NickName  string             `bson:"nickname" json:"nickname"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	Country   string             `bson:"country" json:"country"`
}
