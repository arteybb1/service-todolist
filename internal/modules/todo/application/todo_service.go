package application

import (
	"context"
	"fmt"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/todo/application/dto"
	"github.com/arteybb/service-todolist/internal/modules/todo/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"github.com/arteybb/service-todolist/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
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
	err = s.todoRepo.Create(ctx, todo)
	if err != nil {
		return err
	}
	return nil
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

func (s *TodoService) UpdateTodoById(ctx context.Context, todoID primitive.ObjectID, userID primitive.ObjectID, status bool) error {
	var statusStr constants.Status
	if status {
		statusStr = constants.DONE
	} else {
		statusStr = constants.PENDING
	}
	return s.todoRepo.UpdateTodoById(ctx, todoID, userID, statusStr)
}

func (s *TodoService) GetTodosWithPendingCount(ctx context.Context, userIDHex string) ([]dto.TodoResponse, int, error) {
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	fmt.Println("userIDHex =", userIDHex)
	if err != nil {
		return nil, 0, err
	}

	var (
		todos        []schema.Todo
		pendingCount int
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		todos, err = s.todoRepo.GetTodosByUserID(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		pendingCount, err = s.todoRepo.CountTodosByStatus(ctx, userID, constants.PENDING)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	res := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		res[i] = dto.TodoResponseMapper(todo)
	}

	return res, pendingCount, nil
}
