package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserResponse struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
}

type UserCreate struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
