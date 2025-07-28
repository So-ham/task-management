package task

import (
	"context"
	"regexp"
	"task-management/internal/entities"
	"testing"

	"errors"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening gorm stub database connection", err)
	}

	return gormDB, mock
}

func Test_taskModel_Create(t *testing.T) {
	gormDB, mock := NewMock()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	task := entities.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
	}

	stmt := regexp.QuoteMeta(
		gormDB.Session(&gorm.Session{DryRun: true}).Create(&task).Statement.SQL.String())

	type args struct {
		ctx  context.Context
		task entities.Task
	}

	tests := []struct {
		name    string
		m       *taskModel
		args    args
		wantErr bool
	}{
		{
			name: "Successful create",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx:  context.Background(),
				task: task,
			},
			wantErr: false,
		},
		{
			name: "Failed create",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx:  context.Background(),
				task: task,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectQuery(stmt).WillReturnError(errors.New("DB Closed"))
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectQuery(stmt).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			}

			err := tt.m.Create(tt.args.ctx, &tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskModel.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskModel_GetByID(t *testing.T) {
	gormDB, mock := NewMock()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	taskID, _ := uuid.NewV4()
	task := entities.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
	}

	stmt := regexp.QuoteMeta(
		"SELECT * FROM \"tasks\" WHERE \"tasks\".\"id\" = $1 AND \"tasks\".\"deleted_at\" IS NULL ORDER BY \"tasks\".\"id\" LIMIT 1")

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	tests := []struct {
		name    string
		m       *taskModel
		args    args
		want    *entities.Task
		wantErr bool
	}{
		{
			name: "Successful get by ID",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx: context.Background(),
				id:  taskID,
			},
			want:    &task,
			wantErr: false,
		},
		{
			name: "Failed get by ID",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx: context.Background(),
				id:  taskID,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectQuery(stmt).WillReturnError(errors.New("not found"))
			} else {
				mock.ExpectQuery(stmt).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status"}).AddRow(task.ID, task.Title, task.Description, task.Status))
			}

			got, err := tt.m.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskModel.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("taskModel.GetByID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_taskModel_Update(t *testing.T) {
	gormDB, mock := NewMock()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	task := entities.Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "InProgress",
	}

	stmt := regexp.QuoteMeta(
		gormDB.Session(&gorm.Session{DryRun: true}).Save(&task).Statement.SQL.String())

	type args struct {
		ctx  context.Context
		task entities.Task
	}

	tests := []struct {
		name    string
		m       *taskModel
		args    args
		wantErr bool
	}{
		{
			name: "Successful update",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx:  context.Background(),
				task: task,
			},
			wantErr: false,
		},
		{
			name: "Failed update",
			m:    &taskModel{db: gormDB},
			args: args{
				ctx:  context.Background(),
				task: task,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectExec(stmt).WillReturnError(errors.New("DB Closed"))
			} else {
				mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := tt.m.Update(tt.args.ctx, &tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskModel.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func Test_taskModel_Delete(t *testing.T) {
// 	gormDB, mock := NewMock()

// 	defer func() {
// 		db, _ := gormDB.DB()
// 		db.Close()
// 	}()

// 	taskID, _ := uuid.NewV4()

// 	stmt := regexp.QuoteMeta(
// 		"DELETE FROM \"tasks\" WHERE \"tasks\".\"id\" = $1")

// 	type args struct {
// 		ctx context.Context
// 		id  uuid.UUID
// 	}

// 	tests := []struct {
// 		name    string
// 		m       *taskModel
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successful delete",
// 			m:    &taskModel{db: gormDB},
// 			args: args{
// 				ctx: context.Background(),
// 				id:  taskID,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Failed delete",
// 			m:    &taskModel{db: gormDB},
// 			args: args{
// 				ctx: context.Background(),
// 				id:  taskID,
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if (err := tt.m.Delete(tt.args.ctx, tt.args.id)) != nil {
// 				if !tt.wantErr {
// 					t.Errorf("taskModel.Delete() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 			} else if tt.wantErr {
// 				t.Errorf("taskModel.Delete() error = nil, wantErr %v", tt.wantErr)
// 			}
// 		})
// 	}
// }
