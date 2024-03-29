package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	UserName     string `json:"username" bson:"username" binding:"required"`
	FirstName    string `json:"firstname" bson:"firstname" binding:"required"`
	LastName     string `json:"lastname" bson:"lastname" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	ProfileImage string `json:"profile,omitempty" bson:"profile,omitempty"`
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserName     string             `json:"username" bson:"username"`
	FirstName    string             `json:"firstname" bson:"firstname"`
	LastName     string             `json:"lastname" bson:"lastname"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Role         string             `json:"role" bson:"role"`
	ProfileImage string             `json:"profile,omitempty" bson:"profile,omitempty"`
}

type RequestUpdateUser struct {
	UserName     string `json:"username" bson:"username"`
	FirstName    string `json:"firstname" bson:"firstname"`
	LastName     string `json:"lastname" bson:"lastname"`
	Email        string `json:"email" bson:"email"`
	ProfileImage string `json:"profile,omitempty" bson:"profile,omitempty"`
}

// jwtCustomClaims represents the custom claims of a JWT token
type JwtCustomClaims struct {
	Role         string `json:"role"`
	UserName     string `json:"username" bson:"username"`
	FirstName    string `json:"firstname" bson:"firstname"`
	LastName     string `json:"lastname" bson:"lastname"`
	ProfileImage string `json:"profile,omitempty" bson:"profile,omitempty"`
	jwt.StandardClaims
}
