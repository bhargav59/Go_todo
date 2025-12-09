package repository

import (
	"errors"
	"math"

	"github.com/bhaskar/todo-api/internal/models"
	"gorm.io/gorm"
)

// TodoRepository handles todo data operations
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create inserts a new todo into the database
func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

// FindByID retrieves a todo by ID
func (r *TodoRepository) FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &todo, err
}

// FindByIDAndUserID retrieves a todo by ID and user ID (ownership check)
func (r *TodoRepository) FindByIDAndUserID(id, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &todo, err
}

// ListByUserID retrieves paginated todos for a user
func (r *TodoRepository) ListByUserID(userID uint, page, perPage int, completed *bool) (*models.TodoListResponse, error) {
	var todos []models.Todo
	var total int64

	query := r.db.Model(&models.Todo{}).Where("user_id = ?", userID)

	// Filter by completed status if provided
	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get paginated results
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, err
	}

	// Convert to response
	todoResponses := make([]models.TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponses[i] = todo.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.TodoListResponse{
		Todos:      todoResponses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// Update updates a todo record
func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

// Delete soft-deletes a todo
func (r *TodoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

// DeleteByIDAndUserID deletes a todo by ID only if owned by user
func (r *TodoRepository) DeleteByIDAndUserID(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// CountByUserID counts todos for a user
func (r *TodoRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountCompletedByUserID counts completed todos for a user
func (r *TodoRepository) CountCompletedByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userID, true).Count(&count).Error
	return count, err
}
