package application

import (
	"context"
	"errors"
	"time"

	"github.com/arteybb/service-todolist/internal/config"
	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/auth/application/dto"
	"github.com/arteybb/service-todolist/internal/modules/user/application"
	"github.com/arteybb/service-todolist/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *application.UserService
}

func NewAuthService(userService *application.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Login(ctx context.Context, loginDto dto.LoginDto) (*dto.TokenPair, error) {
	user, err := s.userService.GetByUsername(ctx, loginDto.Username)
	if err != nil {
		return nil, errors.New(string(constants.USER_NOT_FOUND))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
		return nil, errors.New(string(constants.INVALID_CREDENTIAL))
	}

	accessToken, err := s.generateToken(user, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(user *schema.User, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"_id": user.ID.Hex(),
		"exp": time.Now().Add(expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTSecret())
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New(string(constants.INVALID_REFRESH_TOKEN))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(string(constants.INVALID_CLAIMS))
	}

	userIDStr, ok := claims["_id"].(string)
	if !ok {
		return nil, errors.New(string(constants.INVALID_ID))
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New(string(constants.INVALID_ID_FORMAT))
	}

	user, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New(string(constants.USER_NOT_FOUND))
	}

	newAccessToken, err := s.generateToken(user, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
