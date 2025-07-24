package handler

import (
	"net/http"

	"github.com/arteybb/service-todolist/internal/modules/todo/application"
	"github.com/arteybb/service-todolist/internal/modules/todo/application/dto"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service *application.TodoService
}

func NewTodoHandler(service *application.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	ctx := c.Request.Context()

	todos, err := h.service.GetAllTodos(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todoDTO dto.TodoCreateDTO
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&todoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateTodo(ctx, todoDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created successfully",
	})
}

func (h *TodoHandler) GetTodoById(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	todo, err := h.service.GetTodoById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodoById(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	err := h.service.DeleteTodoById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

func (h *TodoHandler) GetTodosByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	todos, err := h.service.GetTodosByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}
