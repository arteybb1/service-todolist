package router

import (
	auth "github.com/arteybb/service-todolist/internal/modules/auth/application"
	"github.com/arteybb/service-todolist/internal/modules/auth/handler"
	user "github.com/arteybb/service-todolist/internal/modules/user/application"
	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.RouterGroup, userService *user.UserService) {
	authService := auth.NewAuthService(userService)
	authHandler := handler.NewAuthHandler(authService)

	group := r.Group("/auth")
	group.POST("", authHandler.Login)
	group.POST("/refresh_token", authHandler.RefreshToken)
}
