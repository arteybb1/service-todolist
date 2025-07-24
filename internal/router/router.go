package router

import (
	"time"

	"github.com/arteybb/service-todolist/internal/config"

	"github.com/arteybb/service-todolist/internal/modules/todo/application"
	todoInfrastructure "github.com/arteybb/service-todolist/internal/modules/todo/infrastructure"
	userApplication "github.com/arteybb/service-todolist/internal/modules/user/application"
	userInfrastructure "github.com/arteybb/service-todolist/internal/modules/user/infrastructure"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{config.AppConfig.BaseURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	routerGroup := r.Group("/api")

	userCol := config.MongoDB.Collection("users")
	userRepo := userInfrastructure.NewUserRepository(userCol)
	userService := userApplication.NewUserService(userRepo)

	todoCol := config.MongoDB.Collection("todos")
	todoRepo := todoInfrastructure.NewTodoRepository(todoCol)

	todoService := application.NewTodoService(todoRepo)

	TodoRoute(routerGroup, todoService)
	UserRoute(routerGroup, userService)
	AuthRoute(routerGroup, userService)

	return r
}
