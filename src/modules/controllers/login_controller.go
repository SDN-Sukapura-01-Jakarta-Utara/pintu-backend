package controllers

import (
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/services"

	"github.com/gin-gonic/gin"
)

// LoginController handles HTTP requests for authentication
type LoginController struct {
	service services.LoginService
}

// NewLoginController creates a new Login controller
func NewLoginController(service services.LoginService) *LoginController {
	return &LoginController{service: service}
}

// Login handles user login
func (c *LoginController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// Logout handles user logout
func (c *LoginController) Logout(ctx *gin.Context) {
	// In JWT stateless approach, logout is client-side
	// Just return success message, client should delete token from environment

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logout successful, please delete your token",
		"user_id": userID,
	})
}
