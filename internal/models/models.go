package models

import (
	"task-management/internal/models/task"

	"gorm.io/gorm"
)

type Model struct {
	Task task.Task
}

// New creates a new instance of Model
func New(gdb *gorm.DB) *Model {
	return &Model{
		Task: task.New(gdb),
	}
}
