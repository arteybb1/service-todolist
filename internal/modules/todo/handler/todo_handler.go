package handler

import (
	"net/http"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/todo/application"
	"github.com/arteybb/service-todolist/internal/modules/todo/application/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		"message": constants.CREATE_SUCCESS,
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
	c.JSON(http.StatusOK, gin.H{"message": constants.DELETE_SUCCESS})
}

func (h *TodoHandler) GetTodosByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.USER_NOT_FOUND})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.INVALID_ID})
		return
	}

	todos, err := h.service.GetTodosByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (h *TodoHandler) UpdateTodoStatus(c *gin.Context) {
	todoIdParams := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIdParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.INVALID_ID})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.USER_NOT_FOUND})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.INVALID_ID})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.USER_NOT_FOUND})
		return
	}

	var req dto.UpdateTodoStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateTodoById(c.Request.Context(), todoID, userID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constants.UPDATE_SUCCESS})
}

func (h *TodoHandler) GetTodosWithPendingCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.UNAUTHORIZED})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.UNAUTHORIZED})
		return
	}

	todos, pendingCount, err := h.service.GetTodosWithPendingCount(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.INTERNAL_ERROR, "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos":        todos,
		"pendingCount": pendingCount,
	})
}
