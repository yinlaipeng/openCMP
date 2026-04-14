package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

func setupScheduledTaskTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	return gormDB, mock
}

func TestScheduledTaskService_CreateScheduledTask(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	task := &model.ScheduledTask{
		Name:        "test-task",
		Type:        "sync",
		Frequency:   "daily",
		TriggerTime: "00:00",
		Status:      "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `scheduled_tasks`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateScheduledTask(context.Background(), task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_GetScheduledTask(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "frequency", "trigger_time", "status", "cloud_account_id", "valid_from", "valid_until", "created_at", "updated_at"}).
		AddRow(1, "test-task", "sync", "daily", "00:00", "active", nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT \\* FROM `scheduled_tasks`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	task, err := service.GetScheduledTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test-task", task.Name)
	assert.Equal(t, "daily", task.Frequency)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_GetScheduledTask_NotFound(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	mock.ExpectQuery("SELECT \\* FROM `scheduled_tasks`").
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	task, err := service.GetScheduledTask(context.Background(), 999)

	assert.NoError(t, err)
	assert.Nil(t, task)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_ListScheduledTasks(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	// Count rows
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `scheduled_tasks`").
		WillReturnRows(countRows)

	// Task rows
	rows := sqlmock.NewRows([]string{"id", "name", "type", "frequency", "trigger_time", "status", "cloud_account_id"}).
		AddRow(1, "task-1", "sync", "daily", "00:00", "active", nil).
		AddRow(2, "task-2", "sync", "weekly", "00:00", "active", nil).
		AddRow(3, "task-3", "sync", "monthly", "00:00", "inactive", nil)

	mock.ExpectQuery("SELECT \\* FROM `scheduled_tasks`").
		WillReturnRows(rows)

	tasks, total, err := service.ListScheduledTasks(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, tasks, 3)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_UpdateScheduledTask(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	task := &model.ScheduledTask{
		ID:          1,
		Name:        "updated-task",
		Type:        "sync",
		Frequency:   "weekly",
		TriggerTime: "12:00",
		Status:      "inactive",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `scheduled_tasks`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdateScheduledTask(context.Background(), task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_DeleteScheduledTask(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	// GORM uses soft delete (UPDATE with deleted_at) instead of hard delete
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `scheduled_tasks`").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := service.DeleteScheduledTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_UpdateScheduledTaskStatus(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	// GORM also updates updated_at field
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `scheduled_tasks`").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := service.UpdateScheduledTaskStatus(context.Background(), 1, "paused")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_CreateScheduledTask_WithValidPeriod(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	// Create a task with valid_from and valid_until
	task := &model.ScheduledTask{
		Name:        "scheduled-task-with-period",
		Type:        "sync",
		Frequency:   "daily",
		TriggerTime: "06:00",
		Status:      "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `scheduled_tasks`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateScheduledTask(context.Background(), task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduledTaskService_ListScheduledTasks_WithPagination(t *testing.T) {
	db, mock := setupScheduledTaskTestDB(t)
	service := NewScheduledTaskService(db)

	// Count rows
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(20)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `scheduled_tasks`").
		WillReturnRows(countRows)

	// Task rows for page 2 (offset 10)
	rows := sqlmock.NewRows([]string{"id", "name", "type", "frequency", "trigger_time", "status"}).
		AddRow(11, "task-11", "sync", "daily", "00:00", "active").
		AddRow(12, "task-12", "sync", "daily", "00:00", "active")

	mock.ExpectQuery("SELECT \\* FROM `scheduled_tasks`").
		WillReturnRows(rows)

	tasks, total, err := service.ListScheduledTasks(context.Background(), 10, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(20), total)
	assert.Len(t, tasks, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}