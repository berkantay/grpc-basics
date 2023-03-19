package model

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	NickName  string    `bson:"nickname" json:"nickname"`
	Password  string    `bson:"password" json:"password"`
	Email     string    `bson:"email" json:"email"`
	Country   string    `bson:"country" json:"country"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserQuery struct {
	ID        string `bson:"_id,omitempty"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	NickName  string `bson:"nickname" json:"nickname"`
	Email     string `bson:"email" json:"email"`
	Country   string `bson:"country" json:"country"`
}

type UserPage struct {
	NumberOfItem int
	Limit        int
}
