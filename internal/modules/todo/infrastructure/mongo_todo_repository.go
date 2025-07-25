package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/todo/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoRepository struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

func NewTodoRepository(col *mongo.Collection, redisClient *redis.Client) domain.TodoRepository {
	return &todoRepository{
		collection:  col,
		redisClient: redisClient,
	}
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
	if err == nil {
		_ = r.redisClient.Del(ctx, "todos:user:"+todo.UserID.Hex()).Err()
	}
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

	_ = r.redisClient.Del(ctx, "todo:"+id).Err()

	return nil
}

func (r *todoRepository) GetTodoById(ctx context.Context, id string) (*schema.Todo, error) {
	cacheKey := "todo:" + id
	var todo schema.Todo

	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &todo); err == nil {
			return &todo, nil
		}
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(todo)
	if err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, bytes, 5*time.Minute).Err()
	}

	return &todo, nil
}

func (r *todoRepository) GetTodosByUserID(ctx context.Context, userID primitive.ObjectID) ([]schema.Todo, error) {
	cacheKey := "todos:user:" + userID.Hex()
	var todos []schema.Todo
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &todos); err == nil {
			return todos, nil
		}
	}

	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(todos)
	if err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, bytes, 5*time.Minute).Err()
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

	err = r.redisClient.Del(ctx, "todo:"+todoID.Hex()).Err()
	if err != nil {
		log.Println("[Redis DEL Error] key:", "todo:"+todoID.Hex(), "error:", err)
	} else {
		log.Println("[Cache Invalidate] Redis key deleted:", "todo:"+todoID.Hex())
	}

	err = r.redisClient.Del(ctx, "todos:user:"+userID.Hex()).Err()
	if err != nil {
		log.Println("[Redis DEL Error] key:", "todos:user:"+userID.Hex(), "error:", err)
	} else {
		log.Println("[Cache Invalidate] Redis key deleted:", "todos:user:"+userID.Hex())
	}

	return nil
}
