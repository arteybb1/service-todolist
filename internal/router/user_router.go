package router

import (
	"github.com/arteybb/service-todolist/internal/middleware"
	"github.com/arteybb/service-todolist/internal/modules/user/application"
	"github.com/arteybb/service-todolist/internal/modules/user/handler"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.RouterGroup, userService *application.UserService) {
	userHandler := handler.NewUserHandler(userService)

	group := r.Group("/user")
	group.POST("create", userHandler.CreateUser)

	protected := group.Group("")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("", userHandler.GetAllUser)
	protected.GET("/profile", userHandler.GetProfile)
}
