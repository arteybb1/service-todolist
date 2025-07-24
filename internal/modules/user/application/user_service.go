package application

import (
	"context"

	"github.com/arteybb/service-todolist/internal/modules/user/application/dto"
	"github.com/arteybb/service-todolist/internal/modules/user/domain"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetAllUser(ctx context.Context) ([]*schema.User, error) {
	users, err := s.userRepository.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}

	if users == nil {
		users = []*schema.User{}
	}

	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, userDto dto.UserCreate) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &schema.User{
		ID:       primitive.NewObjectID(),
		Username: userDto.Username,
		Password: string(hashedPassword),
	}

	return s.userRepository.CreateUser(ctx, user)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*schema.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}

func (s *UserService) GetByID(ctx context.Context, id primitive.ObjectID) (*schema.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*schema.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, objID)
}
