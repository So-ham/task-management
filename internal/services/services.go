package services

import (
	"context"

	"task-management/internal/entities"
	"task-management/internal/models"

	"github.com/gofrs/uuid"
)

// Service represents the service layer having
// all the services from all service packages
type service struct {
	model models.Model
}

// New creates a new instance of Service
func New(model *models.Model) Service {
	m := &service{model: *model}
	return m
}

type Service interface {

	// Task services
	CreateTask(ctx context.Context, req *entities.TaskRequest) error
	GetTaskByID(ctx context.Context, id uuid.UUID) (*entities.TaskResponse, error)
	GetAllTasks(ctx context.Context, page, pageSize int, status *entities.TaskStatus) ([]*entities.TaskResponse, error)
	UpdateTask(ctx context.Context, id uuid.UUID, req *entities.TaskRequest) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
