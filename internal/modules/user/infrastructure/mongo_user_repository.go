package infrastructure

import (
	"context"

	"github.com/arteybb/service-todolist/internal/modules/user/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &userRepository{userCollection: collection}
}

func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schema.User, error) {
	var user schema.User

	err := r.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*schema.User, error) {
	var user schema.User

	err := r.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUser(ctx context.Context) ([]*schema.User, error) {
	cursor, err := r.userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*schema.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *schema.User) error {
	_, err := r.userCollection.InsertOne(ctx, user)
	return err
}
