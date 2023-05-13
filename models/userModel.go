package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	UserID         string             `bson:"user_id"`
	FirstName      string             `bson:"first_name" validate:"required,min=2,max=100"`
	LastName       string             `bson:"last_name" validate:"required,min=2,max=100"`
	Email          string             `bson:"email" validate:"email,required"`
	PhoneNumber    string             `bson:"phone_number" validate:"required,min=10"`
	Type           string             `bson:"type" validate:"required,eq=ADMIN|eq=USER"`
	HashedPassword string             `bson:"hashed_password" validate:"required"`
	AccessToken    string             `bson:"access_token"`
	RefreshToken   string             `bson:"refresh_token"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

type PublishableUserInfo struct {
	UserID         string    `json:"user_id"`
	DisplayingName string    `json:"displaying_name"`
	Email          string    `json:"email" validate:"email,required"`
	PhoneNumber    string    `json:"phone_number" validate:"required,min=10"`
	Type           string    `json:"type" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisteringUser struct {
	Email           string `json:"email" validate:"email,required"`
	FirstName       string `json:"first_name" validate:"required,min=2,max=100"`
	LastName        string `json:"last_name" validate:"required,min=2,max=100"`
	PhoneNumber     string `json:"phone_number" validate:"required,min=10"`
	Password        string `json:"password" validate:"required,min=8,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
}
