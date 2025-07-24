package router

import (
	"github.com/arteybb/service-todolist/internal/middleware"
	"github.com/arteybb/service-todolist/internal/modules/todo/application"
	todoHandler "github.com/arteybb/service-todolist/internal/modules/todo/handler"

	"github.com/gin-gonic/gin"
)

func TodoRoute(r *gin.RouterGroup, todoService *application.TodoService) {
	handler := todoHandler.NewTodoHandler(todoService)
	group := r.Group("/todos")
	group.Use(middleware.JWTAuthMiddleware())
	group.GET("", handler.GetTodosByUserID)
	group.GET(":id", handler.GetTodoById)
	group.POST("create", handler.CreateTodo)
	group.DELETE("delete/:id", handler.DeleteTodoById)
}
