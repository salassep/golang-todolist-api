package models

import (
	"time"

	"github.com/google/uuid"
)

type SubTodo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Todo      uuid.UUID `gorm:"not null" json:"todo,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateSubTodoRequest struct {
	Content   string    `json:"content" binding:"required"`
	Todo      uuid.UUID `json:"todo,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdateSubTodo struct {
	Content   string    `json:"content"`
	Todo      uuid.UUID `json:"todo,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
