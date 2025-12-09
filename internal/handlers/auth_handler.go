package handlers

import (
	"net/http"

	"github.com/bhaskar/todo-api/internal/services"
	"github.com/bhaskar/todo-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Registration data"
// @Success 201 {object} utils.APIResponse{data=services.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.ConflictError(c, err.Error())
			return
		}
		utils.InternalError(c, "Failed to register user")
		return
	}

	utils.Created(c, "User registered successfully", response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=services.AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedError(c, err.Error())
		return
	}

	utils.OK(c, "Login successful", response)
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=models.UserResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.InternalError(c, "Failed to fetch profile")
		return
	}
	if user == nil {
		utils.NotFoundError(c, "User")
		return
	}

	utils.OK(c, "Profile retrieved", user.ToResponse())
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Todo API is running",
	})
}
