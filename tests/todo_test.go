package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhaskar/todo-api/internal/config"
	"github.com/bhaskar/todo-api/internal/handlers"
	"github.com/bhaskar/todo-api/internal/middleware"
	"github.com/bhaskar/todo-api/internal/models"
	"github.com/bhaskar/todo-api/internal/repository"
	"github.com/bhaskar/todo-api/internal/services"
	"github.com/bhaskar/todo-api/pkg/database"
	"github.com/bhaskar/todo-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TodoTestSuite is the test suite for todo endpoints
type TodoTestSuite struct {
	suite.Suite
	router      *gin.Engine
	todoHandler *handlers.TodoHandler
	authHandler *handlers.AuthHandler
	jwtManager  *utils.JWTManager
	authToken   string
}

// SetupSuite runs before all tests
func (s *TodoTestSuite) SetupSuite() {
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
	todoRepo := repository.NewTodoRepository(db)
	authService := services.NewAuthService(userRepo, s.jwtManager)
	todoService := services.NewTodoService(todoRepo)

	s.authHandler = handlers.NewAuthHandler(authService)
	s.todoHandler = handlers.NewTodoHandler(todoService)

	// Setup router
	s.router = gin.New()
	
	// Auth routes
	s.router.POST("/api/auth/register", s.authHandler.Register)
	s.router.POST("/api/auth/login", s.authHandler.Login)

	// Protected todo routes
	protected := s.router.Group("/api/todos")
	protected.Use(middleware.AuthMiddleware(s.jwtManager))
	{
		protected.POST("", s.todoHandler.Create)
		protected.GET("", s.todoHandler.List)
		protected.GET("/stats", s.todoHandler.GetStats)
		protected.GET("/:id", s.todoHandler.GetByID)
		protected.PUT("/:id", s.todoHandler.Update)
		protected.DELETE("/:id", s.todoHandler.Delete)
	}

	// Register and login to get auth token
	s.setupTestUser()
}

// setupTestUser creates a test user and gets auth token
func (s *TodoTestSuite) setupTestUser() {
	body := map[string]string{
		"email":    "todotest@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.authToken = response.Data.Token
}

// TestCreateTodo tests creating a new todo
func (s *TodoTestSuite) TestCreateTodo() {
	body := models.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "This is a test todo",
		Priority:    "high",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Success)
}

// TestCreateTodoWithoutAuth tests creating todo without authentication
func (s *TodoTestSuite) TestCreateTodoWithoutAuth() {
	body := models.CreateTodoRequest{
		Title: "Unauthorized Todo",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// TestListTodos tests listing todos
func (s *TodoTestSuite) TestListTodos() {
	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Success)
}

// TestListTodosWithPagination tests listing with pagination params
func (s *TodoTestSuite) TestListTodosWithPagination() {
	req := httptest.NewRequest(http.MethodGet, "/api/todos?page=1&per_page=5", nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestGetTodoByID tests getting a specific todo
func (s *TodoTestSuite) TestGetTodoByID() {
	// First create a todo
	body := models.CreateTodoRequest{
		Title: "Get By ID Test",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var createResponse struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	// Get the todo
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/todos/%d", createResponse.Data.ID), nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestUpdateTodo tests updating a todo
func (s *TodoTestSuite) TestUpdateTodo() {
	// Create a todo first
	createBody := models.CreateTodoRequest{
		Title: "Update Test",
	}
	jsonBody, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var createResponse struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	// Update the todo
	completed := true
	updateBody := models.UpdateTodoRequest{
		Completed: &completed,
	}
	jsonBody, _ = json.Marshal(updateBody)

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/todos/%d", createResponse.Data.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestDeleteTodo tests deleting a todo
func (s *TodoTestSuite) TestDeleteTodo() {
	// Create a todo first
	createBody := models.CreateTodoRequest{
		Title: "Delete Test",
	}
	jsonBody, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var createResponse struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	// Delete the todo
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/todos/%d", createResponse.Data.ID), nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w = httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
}

// TestGetTodoStats tests getting todo statistics
func (s *TodoTestSuite) TestGetTodoStats() {
	req := httptest.NewRequest(http.MethodGet, "/api/todos/stats", nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestGetNonExistentTodo tests getting a todo that doesn't exist
func (s *TodoTestSuite) TestGetNonExistentTodo() {
	req := httptest.NewRequest(http.MethodGet, "/api/todos/99999", nil)
	req.Header.Set("Authorization", "Bearer "+s.authToken)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNotFound, w.Code)
}

// TestTodoTestSuite runs the test suite
func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
