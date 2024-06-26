package handlers

import (
	"atommuse/backend/authentication-service/pkg/model"
	services "atommuse/backend/authentication-service/pkg/service"
	"fmt"
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

//	@Summary		Edit User
//	@Description	Edit User
//	@Tags			User
//	@Security		BearerAuth
//	@ID				UpdateUserByID
//	@Produce		json
//	@Param			id					path	string					true	"User ID"
//	@Param			RequestUpdateUser	body	model.RequestUpdateUser	true	"User data to edit"
//	@Success		200
//
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api-users/users/{id} [put]
func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	// Get user ID from URL path
	userID := c.Param("id")

	// Parse request body to get the updated user data
	var updateUser model.RequestUpdateUser
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Update user and get the token
	token, err := h.userService.UpdateUserByID(userID, &updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "User updated successfully"})
}

//	@Summary		Change Password
//	@Description	Change Password
//	@Tags			User
//	@Security		BearerAuth
//	@ID				ChangePassword
//	@Produce		json
//	@Param			RequestUpdateUserPassword	body	model.RequestUpdateUserPassword	true	"User password to change password"
//	@Success		200
//	@Failure		401
//	@Failure		500
//	@Router			/api-users/users/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var changePasswordRequest model.RequestUpdateUserPassword
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from the context or session, assuming it's stored there after login
	userID, _ := c.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert user ID to string
	userIDString := fmt.Sprint(userID)

	// Call the ChangePassword method in the authService
	err := h.userService.ChangePassword(userIDString, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

//	@Summary		BanUser
//	@Description	BanUser
//	@Tags			Admin
//	@ID				BanUser
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200
//
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api-users/users/{id}/ban [post]
func (h *UserHandler) BanUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.userService.BanUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

//	@Summary		GetAllUsers
//	@Description	GetAllUsers
//	@Tags			User
//	@Security		BearerAuth
//	@ID				GetAllUsers
//	@Produce		json
//	@Success		200
//
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/api-users/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
