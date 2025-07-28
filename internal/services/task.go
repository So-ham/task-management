package services

import (
	"context"

	"task-management/internal/entities"

	"github.com/gofrs/uuid"
)

// CreateTask adds a new task
func (s *service) CreateTask(ctx context.Context, req *entities.TaskRequest) error {
	// Create task entity from request
	task := &entities.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueDate:     req.DueDate,
	}

	// Save to database
	err := s.model.Task.Create(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskByID retrieves a task by its ID
func (s *service) GetTaskByID(ctx context.Context, id uuid.UUID) (*entities.TaskResponse, error) {
	// Get task from database
	task, err := s.model.Task.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entities.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

// GetAllTasks retrieves all tasks with pagination and optional status filtering
func (s *service) GetAllTasks(ctx context.Context, page, pageSize int, status *entities.TaskStatus) ([]*entities.TaskResponse, error) {
	// Get tasks from database with pagination and filtering
	tasks, err := s.model.Task.GetAll(ctx, page, pageSize, status)
	if err != nil {
		return nil, err
	}

	// Convert to response objects
	response := make([]*entities.TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = &entities.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return response, nil
}

// UpdateTask updates an existing task
func (s *service) UpdateTask(ctx context.Context, id uuid.UUID, req *entities.TaskRequest) error {
	// Check if task exists
	existingTask, err := s.model.Task.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	existingTask.Title = req.Title
	existingTask.Description = req.Description
	existingTask.Status = req.Status
	existingTask.DueDate = req.DueDate

	// Save to database
	return s.model.Task.Update(ctx, existingTask)
}

// DeleteTask removes a task by its ID
func (s *service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	// Check if task exists
	_, err := s.model.Task.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from database
	return s.model.Task.Delete(ctx, id)
}
