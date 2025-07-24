package infrastructure

import (
	"context"
	"errors"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/todo/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(col *mongo.Collection) domain.TodoRepository {
	return &todoRepository{collection: col}
}

func (r *todoRepository) GetAll(ctx context.Context) ([]schema.Todo, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []schema.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *schema.Todo) error {
	_, err := r.collection.InsertOne(ctx, todo)
	return err
}

func (r *todoRepository) DeleteTodoById(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *todoRepository) GetTodoById(ctx context.Context, id string) (*schema.Todo, error) {
	var todo schema.Todo

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) GetTodosByUserID(ctx context.Context, userID primitive.ObjectID) ([]schema.Todo, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []schema.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) UpdateTodoById(ctx context.Context, todoID primitive.ObjectID, userID primitive.ObjectID, status constants.Status) error {
	filter := bson.M{
		"_id":     todoID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{"status": status},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
