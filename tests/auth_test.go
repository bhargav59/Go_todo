package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhaskar/todo-api/internal/config"
	"github.com/bhaskar/todo-api/internal/handlers"
	"github.com/bhaskar/todo-api/internal/middleware"
	"github.com/bhaskar/todo-api/internal/repository"
	"github.com/bhaskar/todo-api/internal/services"
	"github.com/bhaskar/todo-api/pkg/database"
	"github.com/bhaskar/todo-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthTestSuite is the test suite for authentication endpoints
type AuthTestSuite struct {
	suite.Suite
	router      *gin.Engine
	authHandler *handlers.AuthHandler
	jwtManager  *utils.JWTManager
}

// SetupSuite runs before all tests
func (s *AuthTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Use in-memory SQLite for testing - fresh database each run
	cfg := &config.DatabaseConfig{
		Host:   "sqlite",
		DBName: ":memory:",
	}

	db, err := database.Connect(cfg)
	s.Require().NoError(err)
	s.Require().NoError(database.Migrate(db))

	// Setup JWT manager
	s.jwtManager = utils.NewJWTManager("test-secret", time.Hour, "test")

	// Setup repositories and services
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, s.jwtManager)
	s.authHandler = handlers.NewAuthHandler(authService)

	// Setup router
	s.router = gin.New()
	s.router.POST("/api/auth/register", s.authHandler.Register)
	s.router.POST("/api/auth/login", s.authHandler.Login)
	
	// Protected route
	protected := s.router.Group("")
	protected.Use(middleware.AuthMiddleware(s.jwtManager))
	protected.GET("/api/auth/profile", s.authHandler.GetProfile)
}

// TestRegister tests user registration
func (s *AuthTestSuite) TestRegister() {
	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Success)
}

// TestRegisterDuplicateEmail tests duplicate email registration
func (s *AuthTestSuite) TestRegisterDuplicateEmail() {
	body := map[string]string{
		"email":    "duplicate@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	// First registration should succeed
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusCreated, w.Code)

	// Second registration should fail
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	assert.Equal(s.T(), http.StatusConflict, w.Code)
}

// TestRegisterInvalidEmail tests registration with invalid email
func (s *AuthTestSuite) TestRegisterInvalidEmail() {
	body := map[string]string{
		"email":    "invalid-email",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// TestLogin tests user login
func (s *AuthTestSuite) TestLogin() {
	// Register user first
	registerBody := map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Login
	loginBody := map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Success)
}

// TestLoginInvalidPassword tests login with wrong password
func (s *AuthTestSuite) TestLoginInvalidPassword() {
	body := map[string]string{
		"email":    "login@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// TestGetProfileUnauthorized tests accessing profile without token
func (s *AuthTestSuite) TestGetProfileUnauthorized() {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// TestAuthTestSuite runs the test suite
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
