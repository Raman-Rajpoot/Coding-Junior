package models

import "time"

// User defines the structure for the user model
type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"userName" bson:"userName"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	FullName  string    `json:"fullName" bson:"fullName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

// LoginRequest defines the structure for login data
type LoginRequest struct {
	UserName string `json:"userName" bson:"userName"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
