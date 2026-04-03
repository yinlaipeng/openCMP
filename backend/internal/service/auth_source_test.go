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