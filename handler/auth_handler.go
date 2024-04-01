package handlers

import (
	"atommuse/backend/authentication-service/pkg/model"
	services "atommuse/backend/authentication-service/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

//	@Summary		Login
//	@Description	Login user
//	@Tags			Authentications
//	@ID				Login
//	@Produce		json
//	@Param			loginRequest	body	model.LoginRequest	true	"User data to login"
//	@Success		200
//
//	@Failure		400
//	@Failure		401
//
//	@Failure		500
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest model.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

//	@Summary		Register
//	@Description	Register user
//	@Tags			Authentications
//	@ID				Register
//	@Produce		json
//	@Param			registerRequest	body	model.RegisterRequest	true	"User data to create"
//	@Success		201
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest model.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(registerRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registered Successfully"})
}
