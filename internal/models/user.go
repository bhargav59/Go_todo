package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // Never expose password in JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Todos     []Todo         `gorm:"foreignKey:UserID" json:"todos,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserResponse is the safe representation of user data for API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
