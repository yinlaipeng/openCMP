package service

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// 自定义驱动值匹配器，用于匹配JSON参数
type AnyJSON struct{}

func (a AnyJSON) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func setupAuthSourceTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestAuthSourceService_CreateAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"host":     "ldap.example.com",
		"port":     389,
		"base_dn":  "dc=example,dc=com",
		"use_ssl":  false,
		"username": "admin",
		"password": "secret",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		Name:        "test-ldap",
		Type:        "ldap",
		Description: "Test LDAP auth source",
		Enabled:     true,
		Config:      configBytes,
		AutoCreate:  true,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `auth_sources`").
		WithArgs(sqlmock.AnyArg(), "test-ldap", "Test LDAP auth source", "ldap", "system", nil, true, AnyJSON{}, true, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateAuthSource(context.Background(), source)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_GetAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"host":     "ldap.example.com",
		"port":     389,
		"base_dn":  "dc=example,dc=com",
		"use_ssl":  false,
		"username": "admin",
		"password": "secret",
	}
	configBytes, _ := json.Marshal(config)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "description", "enabled", "config", "auto_create"}).
		AddRow(1, "test-ldap", "ldap", "Test LDAP auth source", true, configBytes, true)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	source, err := service.GetAuthSource(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "test-ldap", source.Name)
	assert.Equal(t, "ldap", source.Type)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_ListAuthSources(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `auth_sources`").
		WillReturnRows(countRows)

	config1 := map[string]interface{}{
		"host":     "ldap1.example.com",
		"port":     389,
		"base_dn":  "dc=example,dc=com",
		"use_ssl":  false,
	}
	config1Bytes, _ := json.Marshal(config1)

	config2 := map[string]interface{}{
		"client_id":     "test-client",
		"client_secret": "test-secret",
		"issuer":        "https://oidc.example.com",
	}
	config2Bytes, _ := json.Marshal(config2)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "description", "enabled", "config", "auto_create"}).
		AddRow(1, "ldap-source", "ldap", "LDAP auth source", true, config1Bytes, true).
		AddRow(2, "oidc-source", "oidc", "OIDC auth source", true, config2Bytes, false)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(rows)

	sources, total, err := service.ListAuthSources(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, sources, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_UpdateAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"host":     "updated-ldap.example.com",
		"port":     636,
		"base_dn":  "dc=example,dc=com",
		"use_ssl":  true,
		"username": "admin",
		"password": "newsecret",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		ID:          1,
		Name:        "updated-ldap",
		Type:        "ldap",
		Description: "Updated LDAP auth source",
		Config:      configBytes,
		Enabled:     true,
		AutoCreate:  false,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auth_sources` SET").
		WithArgs("updated-ldap", "Updated LDAP auth source", "ldap", true, AnyJSON{}, false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdateAuthSource(context.Background(), source)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_DeleteAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `auth_sources` WHERE").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DeleteAuthSource(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_EnableAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auth_sources` SET `enabled`=? WHERE").
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.EnableAuthSource(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_DisableAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auth_sources` SET `enabled`=? WHERE").
		WithArgs(false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.DisableAuthSource(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestAuthSourceService_GetAuthSourcesByScope_System(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "scope", "enabled"}).
		AddRow(1, "system-ldap", "ldap", "system", true)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(rows)

	sources, err := svc.GetAuthSourcesByScope(context.Background(), "system", nil)

	assert.NoError(t, err)
	assert.Len(t, sources, 1)
	assert.Equal(t, "system", sources[0].Scope)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_GetAuthSourcesByScope_Domain(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	domainID := uint(42)
	rows := sqlmock.NewRows([]string{"id", "name", "type", "scope", "domain_id", "enabled"}).
		AddRow(2, "domain-ldap", "ldap", "domain", domainID, true)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(rows)

	sources, err := svc.GetAuthSourcesByScope(context.Background(), "domain", &domainID)

	assert.NoError(t, err)
	assert.Len(t, sources, 1)
	assert.Equal(t, "domain", sources[0].Scope)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_AuthenticateUser_LocalSuccess(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	// bcrypt hash for "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 4)

	rows := sqlmock.NewRows([]string{"id", "name", "password", "enabled", "domain_id"}).
		AddRow(1, "testuser", string(hash), true, 0)

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	user, authSource, err := svc.AuthenticateUser(context.Background(), "testuser", "password123", nil)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Nil(t, authSource) // local auth, no auth source
	assert.Equal(t, "testuser", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_AuthenticateUser_WrongPassword(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), 4)

	userRows := sqlmock.NewRows([]string{"id", "name", "password", "enabled", "domain_id"}).
		AddRow(1, "testuser", string(hash), true, 0)
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs("testuser", 1).
		WillReturnRows(userRows)

	// After local failure, query LDAP sources
	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(sqlmock.NewRows([]string{"id"})) // empty: no LDAP sources

	_, _, err := svc.AuthenticateUser(context.Background(), "testuser", "wrongpassword", nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_AuthenticateUser_DisabledUser(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 4)

	rows := sqlmock.NewRows([]string{"id", "name", "password", "enabled", "domain_id"}).
		AddRow(1, "disableduser", string(hash), false, 0)
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs("disableduser", 1).
		WillReturnRows(rows)

	_, _, err := svc.AuthenticateUser(context.Background(), "disableduser", "password123", nil)

	assert.Error(t, err)
	assert.Equal(t, "user is disabled", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_AuthenticateUser_UserNotFound(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	// User not found
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs("nobody", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	// Query LDAP sources - none available
	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := svc.AuthenticateUser(context.Background(), "nobody", "pass", nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_TestAuthSource_Local(t *testing.T) {
	db, _ := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	// local type always returns true
	source := &model.AuthSource{Type: "local"}
	valid, err := svc.TestAuthSource(context.Background(), source)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestAuthSourceService_TestAuthSource_LDAPNoConfig(t *testing.T) {
	db, _ := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	source := &model.AuthSource{Type: "ldap", Config: nil}
	valid, err := svc.TestAuthSource(context.Background(), source)
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "ldap config is empty")
}

func TestAuthSourceService_SyncUsers_UnsupportedType(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	svc := NewAuthSourceService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "enabled"}).
		AddRow(1, "oidc-src", "oidc", true)
	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	_, err := svc.SyncUsers(context.Background(), 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sync not supported")
}
