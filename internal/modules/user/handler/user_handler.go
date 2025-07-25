package handler

import (
	"net/http"

	"github.com/arteybb/service-todolist/internal/constants"
	"github.com/arteybb/service-todolist/internal/modules/user/application"
	"github.com/arteybb/service-todolist/internal/modules/user/application/dto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.userService.GetAllUser(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var userDTO dto.UserCreate

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.CreateUser(ctx, userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": constants.CREATE_SUCCESS,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.USER_NOT_FOUND})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.INVALID_ID_FORMAT})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
