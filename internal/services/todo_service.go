package services

import (
	"errors"

	"github.com/bhaskar/todo-api/internal/models"
	"github.com/bhaskar/todo-api/internal/repository"
)

// TodoService handles todo business logic
type TodoService struct {
	todoRepo *repository.TodoRepository
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

// Create creates a new todo for a user
func (s *TodoService) Create(userID uint, req *models.CreateTodoRequest) (*models.TodoResponse, error) {
	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		DueDate:     req.DueDate,
		UserID:      userID,
		Completed:   false,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	response := todo.ToResponse()
	return &response, nil
}

// GetByID retrieves a todo by ID, with ownership validation
func (s *TodoService) GetByID(todoID, userID uint) (*models.TodoResponse, error) {
	todo, err := s.todoRepo.FindByIDAndUserID(todoID, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}

	response := todo.ToResponse()
	return &response, nil
}

// List retrieves paginated todos for a user
func (s *TodoService) List(userID uint, page, perPage int, completed *bool) (*models.TodoListResponse, error) {
	// Apply defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return s.todoRepo.ListByUserID(userID, page, perPage, completed)
}

// Update updates a todo
func (s *TodoService) Update(todoID, userID uint, req *models.UpdateTodoRequest) (*models.TodoResponse, error) {
	// Find todo with ownership check
	todo, err := s.todoRepo.FindByIDAndUserID(todoID, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}

	// Apply updates
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	response := todo.ToResponse()
	return &response, nil
}

// Delete removes a todo
func (s *TodoService) Delete(todoID, userID uint) error {
	// Verify ownership before delete
	todo, err := s.todoRepo.FindByIDAndUserID(todoID, userID)
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("todo not found")
	}

	return s.todoRepo.Delete(todoID)
}

// GetStats returns todo statistics for a user
func (s *TodoService) GetStats(userID uint) (map[string]int64, error) {
	total, err := s.todoRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	completed, err := s.todoRepo.CountCompletedByUserID(userID)
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"total":     total,
		"completed": completed,
		"pending":   total - completed,
	}, nil
}
