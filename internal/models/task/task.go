package task

import (
	"context"

	"task-management/internal/entities"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Task interface defines methods for task data operations
type Task interface {
	Create(ctx context.Context, task *entities.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Task, error)
	GetAll(ctx context.Context, page, pageSize int, status *entities.TaskStatus) ([]entities.Task, error)
	Update(ctx context.Context, task *entities.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskModel struct {
	db *gorm.DB
}

// New creates a new instance of Task
func New(db *gorm.DB) Task {
	return &taskModel{db: db}
}

// Create adds a new task to the database
func (m *taskModel) Create(ctx context.Context, task *entities.Task) error {
	// Generate a new UUID if not provided
	if task.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		task.ID = id
	}

	return m.db.WithContext(ctx).Create(task).Error
}

// GetByID retrieves a task by its ID
func (m *taskModel) GetByID(ctx context.Context, id uuid.UUID) (*entities.Task, error) {
	var task entities.Task
	result := m.db.WithContext(ctx).First(&task, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

// GetAll retrieves all tasks with pagination and optional status filtering
func (m *taskModel) GetAll(ctx context.Context, page, pageSize int, status *entities.TaskStatus) ([]entities.Task, error) {
	var tasks []entities.Task
	query := m.db.WithContext(ctx)

	// Apply status filter if provided
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

// Update updates an existing task
func (m *taskModel) Update(ctx context.Context, task *entities.Task) error {
	return m.db.WithContext(ctx).Save(task).Error
}

// Delete removes a task by its ID
func (m *taskModel) Delete(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Delete(&entities.Task{}, "id = ?", id).Error
}
