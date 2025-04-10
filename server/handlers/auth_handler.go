// handlers/auth_handler.go
package handlers

import (
	"bustracking/middleware"
	"bustracking/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	AuthService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}

	// If user is a driver, include bus information
	if user.Role == "driver" {
		bus, err := h.AuthService.GetAssignedBus(user.ID)
		if err == nil {
			response["bus"] = gin.H{
				"id":        bus.ID,
				"busNumber": bus.BusNumber,
				"routeId":   bus.RouteID,
				"status":    bus.Status,
				"latitude":  bus.Latitude,
				"longitude": bus.Longitude,
			}
		} else {
			response["bus"] = nil
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.AuthService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	response := gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}

	// If user is a driver, include bus information
	if user.Role == "driver" {
		bus, err := h.AuthService.GetAssignedBus(user.ID)
		if err == nil {
			response["bus"] = gin.H{
				"id":        bus.ID,
				"busNumber": bus.BusNumber,
				"routeId":   bus.RouteID,
				"status":    bus.Status,
				"latitude":  bus.Latitude,
				"longitude": bus.Longitude,
			}
		} else {
			response["bus"] = nil
		}
	}

	c.JSON(http.StatusOK, response)
}
