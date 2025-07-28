package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

// Task status enum
type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusInProgress TaskStatus = "InProgress"
	StatusCompleted  TaskStatus = "Completed"
	StatusCancelled  TaskStatus = "Cancelled"
)

type Task struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	Title       string     `json:"title" gorm:"not null" validate:"required"`
	Description string     `json:"description" gorm:"type:text"`
	Status      TaskStatus `json:"status" gorm:"not null;default:'Pending'" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
}

type TaskRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
}

type TaskResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
