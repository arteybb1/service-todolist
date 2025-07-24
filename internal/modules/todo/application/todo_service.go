package application

import (
	"context"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/todo/application/dto"
	"github.com/arteybb/service-todolist/internal/modules/todo/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"github.com/arteybb/service-todolist/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService struct {
	todoRepo domain.TodoRepository
}

func NewTodoService(todoRepo domain.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]dto.TodoResponse, error) {
	todos, err := s.todoRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return utils.MapSlice(todos, dto.TodoResponseMapper), nil
}

func (s *TodoService) CreateTodo(ctx context.Context, todoDTO dto.TodoCreateDTO) error {
	userID, err := primitive.ObjectIDFromHex(todoDTO.UserID)
	if err != nil {
		return err
	}
	todo := &schema.Todo{
		ID:     primitive.NewObjectID(),
		Title:  todoDTO.Title,
		Status: constants.PENDING,
		UserID: userID,
	}
	return s.todoRepo.Create(ctx, todo)
}

func (s *TodoService) GetTodoById(ctx context.Context, id string) (*dto.TodoResponse, error) {
	todo, err := s.todoRepo.GetTodoById(ctx, id)
	if err != nil {
		return nil, err
	}
	res := dto.TodoResponseMapper(*todo)
	return &res, nil
}

func (s *TodoService) DeleteTodoById(ctx context.Context, id string) error {
	return s.todoRepo.DeleteTodoById(ctx, id)
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userIDHex string) ([]dto.TodoResponse, error) {
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}

	todos, err := s.todoRepo.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		res[i] = dto.TodoResponseMapper(todo)
	}

	return res, nil
}
