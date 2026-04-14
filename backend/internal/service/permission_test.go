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

func setupPermissionTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestPermissionService_CreatePermission(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	permission := &model.Permission{
		Name:        "user.list",
		Description: "List users permission",
		Resource:    "user",
		Action:      "list",
		Scope:       "system",
		Enabled:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `permissions`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreatePermission(context.Background(), permission)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_GetPermission(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true)

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	permission, err := service.GetPermission(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, permission)
	assert.Equal(t, "user.list", permission.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_GetPermission_NotFound(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	permission, err := service.GetPermission(context.Background(), 999)

	assert.NoError(t, err)
	assert.Nil(t, permission)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_ListPermissions(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `permissions`").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true).
		AddRow(2, "user.create", "Create users permission", "user", "create", "system", nil, true)

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WillReturnRows(rows)

	permissions, total, err := service.ListPermissions(context.Background(), nil, "", "", "", "", nil, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, permissions, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_UpdatePermission(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	permission := &model.Permission{
		ID:          1,
		Name:        "user.list.updated",
		Description: "Updated description",
		Resource:    "user",
		Action:      "list",
		Scope:       "system",
		Enabled:     false,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `permissions`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdatePermission(context.Background(), permission)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_DeletePermission(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `permissions`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DeletePermission(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_AssignPermissionToRole(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `role_permissions`").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `role_permissions`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.AssignPermissionToRole(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_RevokePermissionFromRole(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `role_permissions`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.RevokePermissionFromRole(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_GetPermissionsByRole(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true).
		AddRow(2, "user.create", "Create users permission", "user", "create", "system", nil, true)

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WillReturnRows(rows)

	permissions, err := service.GetPermissionsByRole(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPermissionService_GetUserPermissions(t *testing.T) {
	db, mock := setupPermissionTestDB(t)
	service := NewPermissionService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true).
		AddRow(2, "vm.create", "Create VMs permission", "vm", "create", "domain", nil, true)

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WillReturnRows(rows)

	permissions, err := service.GetUserPermissions(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}