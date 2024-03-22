package model

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	UserName     string `json:"username" bson:"username"`
	Password     string `json:"password" binding:"required,min=8"`
	Role         string `json:"role" binding:"required,oneof=admin exhibitor"`
	ProfileImage string `json:"profile,omitempty" bson:"profile,omitempty"`
}

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	Email        string `json:"email" bson:"email"`
	Password     string `json:"password" bson:"password"`
	Role         string `json:"role" bson:"role"`
	ProfileImage string `json:"profile,omitempty" bson:"profile,omitempty"`
}
