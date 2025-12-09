package handlers

import (
	"strconv"

	"github.com/bhaskar/todo-api/internal/middleware"
	"github.com/bhaskar/todo-api/internal/models"
	"github.com/bhaskar/todo-api/internal/services"
	"github.com/bhaskar/todo-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles todo endpoints
type TodoHandler struct {
	todoService *services.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateTodoRequest true "Todo data"
// @Success 201 {object} utils.APIResponse{data=models.TodoResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /api/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	todo, err := h.todoService.Create(userID, &req)
	if err != nil {
		utils.InternalError(c, "Failed to create todo")
		return
	}

	utils.Created(c, "Todo created successfully", todo)
}

// List godoc
// @Summary List todos
// @Description Get paginated list of todos for the authenticated user
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param completed query bool false "Filter by completed status"
// @Success 200 {object} utils.APIResponse{data=models.TodoListResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/todos [get]
func (h *TodoHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	var completed *bool
	if c.Query("completed") != "" {
		val := c.Query("completed") == "true"
		completed = &val
	}

	todos, err := h.todoService.List(userID, page, perPage, completed)
	if err != nil {
		utils.InternalError(c, "Failed to fetch todos")
		return
	}

	utils.OK(c, "Todos retrieved", todos)
}

// GetByID godoc
// @Summary Get a todo by ID
// @Description Get a specific todo item by ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} utils.APIResponse{data=models.TodoResponse}
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetByID(uint(todoID), userID)
	if err != nil {
		utils.NotFoundError(c, "Todo")
		return
	}

	utils.OK(c, "Todo retrieved", todo)
}

// Update godoc
// @Summary Update a todo
// @Description Update a specific todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param request body models.UpdateTodoRequest true "Update data"
// @Success 200 {object} utils.APIResponse{data=models.TodoResponse}
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	todo, err := h.todoService.Update(uint(todoID), userID, &req)
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFoundError(c, "Todo")
			return
		}
		utils.InternalError(c, "Failed to update todo")
		return
	}

	utils.OK(c, "Todo updated successfully", todo)
}

// Delete godoc
// @Summary Delete a todo
// @Description Delete a specific todo item
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 401 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestError(c, "Invalid todo ID")
		return
	}

	err = h.todoService.Delete(uint(todoID), userID)
	if err != nil {
		if err.Error() == "todo not found" {
			utils.NotFoundError(c, "Todo")
			return
		}
		utils.InternalError(c, "Failed to delete todo")
		return
	}

	utils.NoContent(c)
}

// GetStats godoc
// @Summary Get todo statistics
// @Description Get todo statistics for the authenticated user
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=map[string]int64}
// @Failure 401 {object} utils.APIResponse
// @Router /api/todos/stats [get]
func (h *TodoHandler) GetStats(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedError(c, "")
		return
	}

	stats, err := h.todoService.GetStats(userID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch statistics")
		return
	}

	utils.OK(c, "Statistics retrieved", stats)
}
