package handlers

import (
	"atommuse/backend/authentication-service/pkg/model"
	services "atommuse/backend/authentication-service/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateUserByIDHandler handles HTTP requests to update a user by their ID.

// @Summary		Edit User
// @Description	Edit User
// @Tags			User
// @Security		BearerAuth
// @ID				UpdateUserByID
// @Produce		json
// @Param			id					path	string					true	"User ID"
// @Param			RequestUpdateUser	body	model.RequestUpdateUser	true	"User data to edit"
// @Success		200
// @Failure		500
// @Router			/api/user/{id} [put]
func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	// Get user ID from URL path
	userID := c.Param("id")

	// Parse request body to get the updated user data
	var updateUser model.RequestUpdateUser
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Update user
	err := h.userService.UpdateUserByID(userID, &updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
