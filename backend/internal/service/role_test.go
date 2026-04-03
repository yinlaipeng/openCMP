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

func setupRoleTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestRoleService_CreateRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	role := &model.Role{
		Name:        "test-role",
		DisplayName: "Test Role",
		Description: "A test role",
		DomainID:    1,
		Type:        "custom",
		Enabled:     true,
		IsPublic:    false,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `roles`").
		WithArgs(sqlmock.AnyArg(), "test-role", "Test Role", "A test role", 1, "custom", true, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateRole(context.Background(), role)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_GetRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "display_name", "description", "domain_id", "type", "enabled", "is_public"}).
		AddRow(1, "admin", "Administrator", "Administrator role", 1, "system", true, false)

	mock.ExpectQuery("SELECT \\* FROM `roles`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	role, err := service.GetRole(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "admin", role.Name)
	assert.Equal(t, "Administrator", role.DisplayName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_ListRoles(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `roles`").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "display_name", "description", "domain_id", "type", "enabled", "is_public"}).
		AddRow(1, "admin", "Administrator", "Administrator role", 1, "system", true, false).
		AddRow(2, "user", "User", "Regular user role", 1, "system", true, false)

	mock.ExpectQuery("SELECT \\* FROM `roles`").
		WillReturnRows(rows)

	var domainID uint = 1
	roles, total, err := service.ListRoles(context.Background(), &domainID, "", "", nil, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, roles, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_UpdateRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	role := &model.Role{
		ID:          1,
		Name:        "updated-role",
		DisplayName: "Updated Role",
		Description: "An updated role",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `roles` SET").
		WithArgs("updated-role", "Updated Role", "An updated role", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdateRole(context.Background(), role)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_DeleteRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `roles` WHERE").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DeleteRole(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_EnableRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `roles` SET `enabled`=? WHERE").
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.EnableRole(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_DisableRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `roles` SET `enabled`=? WHERE").
		WithArgs(false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DisableRole(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_GetRolePermissions(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "display_name", "resource", "action", "type"}).
		AddRow(1, "user_list", "List Users", "user", "list", "system").
		AddRow(2, "user_create", "Create User", "user", "create", "system")

	mock.ExpectQuery("SELECT \\* FROM `permissions`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	permissions, err := service.GetRolePermissions(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, permissions, 2)
	assert.Equal(t, "user_list", permissions[0].Name)
	assert.Equal(t, "user_create", permissions[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_AssignPermissionToRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `role_permissions`").
		WithArgs(sqlmock.AnyArg(), 1, 1, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.AssignPermissionToRole(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleService_RevokePermissionFromRole(t *testing.T) {
	db, mock := setupRoleTestDB(t)
	service := NewRoleService(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `role_permissions` WHERE").
		WithArgs(1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.RevokePermissionFromRole(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}