package schema

import (
	"github.com/arteybb/service-todolist/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title" bson:"title"`
	Status constants.Status   `json:"status" bson:"status"`
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}
