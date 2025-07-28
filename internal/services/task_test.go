package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"task-management/internal/entities"
	"task-management/internal/models"
	taskMock "task-management/internal/models/task/mocks"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

func Test_service_CreateTask(t *testing.T) {
	taskID, _ := uuid.NewV4()
	testTime := time.Now()
	pendingStatus := entities.StatusPending

	req := &entities.TaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      pendingStatus,
		DueDate:     &testTime,
	}

	successMock := taskMock.Task{}
	successMock.On("Create", mock.Anything, mock.MatchedBy(func(t *entities.Task) bool {
		return t.Title == req.Title &&
			t.Description == req.Description &&
			t.Status == req.Status &&
			t.DueDate.Equal(*req.DueDate)
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(1).(*entities.Task)
		task.ID = taskID
	})

	errorMock := taskMock.Task{}
	errorMock.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	tests := []struct {
		name    string
		s       *service
		req     *entities.TaskRequest
		wantErr bool
	}{
		{
			name: "successful creation",
			s: &service{
				model: models.Model{
					Task: &successMock,
				},
			},
			req:     req,
			wantErr: false,
		},
		{
			name: "database error",
			s: &service{
				model: models.Model{
					Task: &errorMock,
				},
			},
			req:     req,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.CreateTask(context.Background(), tt.req); (err != nil) != tt.wantErr {
				t.Errorf("service.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetTaskByID(t *testing.T) {
	taskID, _ := uuid.NewV4()
	invalidID, _ := uuid.NewV4()
	testTime := time.Now()
	pendingStatus := entities.StatusPending

	task := &entities.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      pendingStatus,
		DueDate:     &testTime,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
	}

	expected := &entities.TaskResponse{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      pendingStatus,
		DueDate:     &testTime,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
	}

	successMock := taskMock.Task{}
	successMock.On("GetByID", mock.Anything, taskID).Return(task, nil)

	errorMock := taskMock.Task{}
	errorMock.On("GetByID", mock.Anything, invalidID).Return(nil, errors.New("not found"))

	tests := []struct {
		name    string
		s       *service
		id      uuid.UUID
		want    *entities.TaskResponse
		wantErr bool
	}{
		{
			name: "successful retrieval",
			s: &service{
				model: models.Model{
					Task: &successMock,
				},
			},
			id:      taskID,
			want:    expected,
			wantErr: false,
		},
		{
			name: "not found",
			s: &service{
				model: models.Model{
					Task: &errorMock,
				},
			},
			id:      invalidID,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetTaskByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTaskByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTaskByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
