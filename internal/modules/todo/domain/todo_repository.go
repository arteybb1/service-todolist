package domain

import (
	"context"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoRepository interface {
	GetAll(ctx context.Context) ([]schema.Todo, error)
	GetTodoById(ctx context.Context, id string) (*schema.Todo, error)
	Create(ctx context.Context, todo *schema.Todo) error
	DeleteTodoById(ctx context.Context, id string) error
	GetTodosByUserID(ctx context.Context, userID primitive.ObjectID) ([]schema.Todo, error)
	UpdateTodoById(ctx context.Context, todoID primitive.ObjectID, userID primitive.ObjectID, status constants.Status) error
	CountTodosByStatus(ctx context.Context, userID primitive.ObjectID, status constants.Status) (int, error)
}
