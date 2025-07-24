package domain

import (
	"context"

	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]*schema.User, error)
	CreateUser(ctx context.Context, user *schema.User) error
	FindByUsername(ctx context.Context, username string) (*schema.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*schema.User, error)
}
