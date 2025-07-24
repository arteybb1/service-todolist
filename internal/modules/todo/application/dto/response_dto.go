package dto

import (
	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoResponse struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title" bson:"title"`
	Status constants.Status   `json:"status" bson:"status"`
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

func TodoResponseMapper(todo schema.Todo) TodoResponse {
	return TodoResponse{
		ID:     todo.ID,
		Title:  todo.Title,
		Status: todo.Status,
		UserID: todo.UserID,
	}
}
