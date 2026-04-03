package service

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/utils"
)

// 自定义驱动值匹配器，用于匹配任意字符串
type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func setupUserTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestUserService_CreateUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	password := "password123"

	user := &model.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: password,
		DomainID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("testuser", "", "test@example.com", "", AnyString{}, 1, true, false, "", nil, "", nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password) // 密码应该被哈希
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "domain_id", "enabled"}).
		AddRow(1, "testuser", "test@example.com", 1, true)

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	user, err := service.GetUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUserByName(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "display_name", "email", "phone", "password", "domain_id", "enabled", "mfa_enabled", "mfa_secret", "last_login_at", "last_login_ip", "password_expire", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "testuser", "", "test@example.com", "", "$2a$10$Rwcb7tYbGB8Bz68NXyYefesOMG8cAwmSPZx4wO.SL3fSsKk3t5jEK", 1, true, false, "", nil, "", nil, "2026-04-02T14:45:16.719Z", "2026-04-02T14:45:16.719Z", nil)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE name = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	user, err := service.GetUserByName(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ListUsers(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users`").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "domain_id", "enabled"}).
		AddRow(1, "user1", "user1@example.com", 1, true).
		AddRow(2, "user2", "user2@example.com", 1, true)

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(rows)

	var domainID uint = 1
	users, total, err := service.ListUsers(context.Background(), &domainID, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_UpdateUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	user := &model.User{
		ID:          1,
		Name:        "updateduser",
		DisplayName: "Updated User",
		Email:       "updated@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET").
		WithArgs("Updated User", "updated@example.com", "", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_DeleteUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `users` WHERE").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_EnableUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `enabled`=? WHERE").
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.EnableUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_DisableUser(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `enabled`=? WHERE").
		WithArgs(false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DisableUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ResetUserPassword(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	newPassword := "newpassword123"
	expectedHashedPassword, _ := utils.HashPassword(newPassword)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `password`=? WHERE").
		WithArgs(expectedHashedPassword, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.ResetUserPassword(context.Background(), 1, newPassword)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUserRoles(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "admin", "Administrator role").
		AddRow(2, "user", "Regular user role")

	mock.ExpectQuery("SELECT \\* FROM `roles`").
		WithArgs(1, 1, 1).
		WillReturnRows(rows)

	roles, err := service.GetUserRoles(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, "admin", roles[0].Name)
	assert.Equal(t, "user", roles[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUserGroups(t *testing.T) {
	db, mock := setupUserTestDB(t)
	service := NewUserService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "group1", "First group").
		AddRow(2, "group2", "Second group")

	mock.ExpectQuery("SELECT \\* FROM `groups`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	groups, err := service.GetUserGroups(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "group1", groups[0].Name)
	assert.Equal(t, "group2", groups[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}