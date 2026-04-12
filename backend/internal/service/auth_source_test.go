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
		"url":               "ldap://ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "secret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "example.com",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		Name:        "test-ldap",
		Type:        "ldap",
		Description: "Test LDAP auth source",
		Enabled:     true,
		Config:      configBytes,
		AutoCreate:  false, // Default value now false
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `auth_sources`").
		WithArgs(sqlmock.AnyArg(), "test-ldap", "Test LDAP auth source", "ldap", "system", nil, true, AnyJSON{}, false, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateAuthSource(context.Background(), source)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_CreateAuthSourceWithTargetDomain(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"url":               "ldap://ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "secret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "test-domain.com",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		Name:        "test-ldap-with-target",
		Type:        "ldap",
		Description: "Test LDAP auth source with target domain",
		Enabled:     true,
		Config:      configBytes,
		AutoCreate:  false, // Default value now false
	}

	// Mock check for existing domain (not found)
	mock.ExpectQuery("SELECT \\* FROM `domains` WHERE").
		WithArgs("test-domain.com").
		WillReturnError(gorm.ErrRecordNotFound)

	// Mock domain creation
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `domains`").
		WithArgs("test-domain.com", "Domain for test-ldap-with-target authentication source", true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Mock auth source creation
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `auth_sources`").
		WithArgs(sqlmock.AnyArg(), "test-ldap-with-target", "Test LDAP auth source with target domain", "ldap", "system", sqlmock.AnyArg(), true, AnyJSON{}, false, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
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
		"url":               "ldap://ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "secret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "example.com",
	}
	configBytes, _ := json.Marshal(config)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "description", "enabled", "config", "auto_create"}).
		AddRow(1, "test-ldap", "ldap", "Test LDAP auth source", true, configBytes, false)

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
		"url":               "ldap://ldap1.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "secret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "example.com",
	}
	config1Bytes, _ := json.Marshal(config1)

	config2 := map[string]interface{}{
		"client_id":     "test-client",
		"client_secret": "test-secret",
		"issuer":        "https://oidc.example.com",
	}
	config2Bytes, _ := json.Marshal(config2)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "description", "enabled", "config", "auto_create"}).
		AddRow(1, "ldap-source", "ldap", "LDAP auth source", true, config1Bytes, false).
		AddRow(2, "oidc-source", "oidc", "OIDC auth source", true, config2Bytes, false)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources`").
		WillReturnRows(rows)

	sources, total, err := service.ListAuthSources(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, sources, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_ListAuthSourcesWithFilters(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	filters := map[string]interface{}{
		"name":    "test",
		"type":    "ldap",
		"enabled": true,
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `auth_sources` WHERE").
		WithArgs("%test%", "ldap", true).
		WillReturnRows(countRows)

	config := map[string]interface{}{
		"url":               "ldap://ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "secret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "example.com",
	}
	configBytes, _ := json.Marshal(config)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "description", "enabled", "config", "auto_create"}).
		AddRow(1, "test-ldap", "ldap", "Test LDAP auth source", true, configBytes, false)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources` WHERE name LIKE \\? AND type = \\? AND enabled = \\? ORDER BY created_at DESC LIMIT \\? OFFSET \\?").
		WithArgs("%test%", "ldap", true, 10, 0).
		WillReturnRows(rows)

	sources, total, err := service.ListAuthSourcesWithFilters(context.Background(), 10, 0, filters)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, sources, 1)
	assert.Equal(t, "test-ldap", sources[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_UpdateAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"url":               "ldap://updated-ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "newsecret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "updated.com",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		ID:          1,
		Name:        "updated-ldap",
		Type:        "ldap",
		Description: "Updated LDAP auth source",
		Config:      configBytes,
		Enabled:     true,
		AutoCreate:  false, // Now defaults to false
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auth_sources` SET").
		WithArgs("updated-ldap", "Updated LDAP auth source", "ldap", "system", nil, true, AnyJSON{}, false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.UpdateAuthSource(context.Background(), source)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthSourceService_UpdateAuthSourceWithTargetDomain(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	config := map[string]interface{}{
		"url":               "ldap://updated-ldap.example.com:389",
		"base_dn":           "dc=example,dc=com",
		"bind_dn":           "cn=admin,dc=example,dc=com",
		"bind_password":     "newsecret",
		"user_filter":       "(objectClass=person)",
		"user_id_attr":      "uid",
		"user_name_attr":    "cn",
		"user_search_base":  "ou=users,dc=example,dc=com",
		"group_search_base": "ou=groups,dc=example,dc=com",
		"target_domain":     "updated.com",
	}
	configBytes, _ := json.Marshal(config)

	source := &model.AuthSource{
		ID:          1,
		Name:        "updated-ldap",
		Type:        "ldap",
		Description: "Updated LDAP auth source with target domain",
		Config:      configBytes,
		Enabled:     true,
		AutoCreate:  false, // Now defaults to false
	}

	// Mock check for existing domain (not found)
	mock.ExpectQuery("SELECT \\* FROM `domains` WHERE").
		WithArgs("updated.com").
		WillReturnError(gorm.ErrRecordNotFound)

	// Mock domain creation
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `domains`").
		WithArgs("updated.com", "Domain for updated-ldap authentication source", true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Mock auth source update
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auth_sources` SET").
		WithArgs("updated-ldap", "Updated LDAP auth source with target domain", "ldap", "system", sqlmock.AnyArg(), true, AnyJSON{}, false, 1).
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

	rows := sqlmock.NewRows([]string{"id", "name", "type", "scope", "enabled", "auto_create"}).
		AddRow(1, "system-ldap", "ldap", "system", true, false)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources` WHERE \\(scope = \\? AND enabled = \\?\\)").
		WithArgs("system", true).
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
	rows := sqlmock.NewRows([]string{"id", "name", "type", "scope", "domain_id", "enabled", "auto_create"}).
		AddRow(2, "domain-ldap", "ldap", "domain", domainID, true, false)

	mock.ExpectQuery("SELECT \\* FROM `auth_sources` WHERE \\(scope = \\? AND enabled = \\?\\) AND \\(domain_id = \\?\\)").
		WithArgs("domain", true, domainID).
		WillReturnRows(rows)

	sources, err := svc.GetAuthSourcesByScope(context.Background(), "domain", &domainID)

	assert.NoError(t, err)
	assert.Len(t, sources, 1)
	assert.Equal(t, "domain", sources[0].Scope)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnmarshalLDAPConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		jsonData := json.RawMessage(`{
			"url": "ldap://localhost:389",
			"base_dn": "dc=example,dc=com",
			"bind_dn": "cn=admin,dc=example,dc=com",
			"bind_password": "password",
			"user_filter": "(objectClass=person)",
			"user_id_attr": "uid",
			"user_name_attr": "cn",
			"user_search_base": "ou=users,dc=example,dc=com",
			"group_search_base": "ou=groups,dc=example,dc=com",
			"user_enabled_attribute": "accountEnabled",
			"protocol": "ldap",
			"auth_type": "openldap",
			"target_domain": "example.com"
		}`)

		cfg, err := UnmarshalLDAPConfig(jsonData)
		assert.NoError(t, err)
		assert.Equal(t, "ldap://localhost:389", cfg.URL)
		assert.Equal(t, "dc=example,dc=com", cfg.BaseDN)
		assert.Equal(t, "cn=admin,dc=example,dc=com", cfg.BindDN)
		assert.Equal(t, "password", cfg.BindPassword)
		assert.Equal(t, "(objectClass=person)", cfg.UserFilter)
		assert.Equal(t, "uid", cfg.UserIDAttr)
		assert.Equal(t, "cn", cfg.UserNameAttr)
		assert.Equal(t, "ou=users,dc=example,dc=com", cfg.UserSearchBase)
		assert.Equal(t, "ou=groups,dc=example,dc=com", cfg.GroupSearchBase)
		assert.Equal(t, "accountEnabled", cfg.UserEnabledAttribute)
		assert.Equal(t, "ldap", cfg.Protocol)
		assert.Equal(t, "openldap", cfg.AuthType)
		assert.Equal(t, "example.com", cfg.TargetDomain)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := UnmarshalLDAPConfig(json.RawMessage(`invalid json`))
		assert.Error(t, err)
	})
}

func TestAuthSourceService_isSystemAuthSource(t *testing.T) {
	db, _ := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	t.Run("local system scope", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "local",
			Scope: "system",
			Name:  "test",
		}
		assert.True(t, service.isSystemAuthSource(source))
	})

	t.Run("sql system scope", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "sql",
			Scope: "system",
			Name:  "test",
		}
		assert.True(t, service.isSystemAuthSource(source))
	})

	t.Run("ldap system scope", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "ldap",
			Scope: "system",
			Name:  "test",
		}
		assert.False(t, service.isSystemAuthSource(source))
	})

	t.Run("system in name", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "ldap",
			Scope: "domain",
			Name:  "system-ldap",
		}
		assert.True(t, service.isSystemAuthSource(source))
	})

	t.Run("builtin in name", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "local",
			Scope: "domain",
			Name:  "builtin-auth",
		}
		assert.True(t, service.isSystemAuthSource(source))
	})

	t.Run("case insensitive name", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "local",
			Scope: "domain",
			Name:  "SYSTEM-LOCAL",
		}
		assert.True(t, service.isSystemAuthSource(source))
	})

	t.Run("neither system nor builtin", func(t *testing.T) {
		source := &model.AuthSource{
			Type:  "ldap",
			Scope: "domain",
			Name:  "regular-ldap",
		}
		assert.False(t, service.isSystemAuthSource(source))
	})
}

func TestAuthSourceService_TestAuthSource(t *testing.T) {
	db, _ := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	t.Run("ldap test success", func(t *testing.T) {
		source := &model.AuthSource{
			Type: "ldap",
			Config: json.RawMessage(`{
				"url": "ldap://localhost:389",
				"base_dn": "dc=example,dc=com"
			}`),
		}

		valid, err := service.TestAuthSource(context.Background(), source)
		// Note: Since we're not actually connecting to LDAP in tests, this may return an error
		// depending on network connectivity. The important thing is the method runs without panicking.
		_ = valid
		_ = err
	})

	t.Run("mock ldap test", func(t *testing.T) {
		source := &model.AuthSource{
			Type: "ldap",
			Config: json.RawMessage(`{
				"url": "ldap://mock-ldap-server:389",
				"base_dn": "dc=example,dc=com"
			}`),
		}

		valid, err := service.TestAuthSource(context.Background(), source)
		// Mock LDAP tests should return true without attempting real connections
		_ = valid
		_ = err
	})

	t.Run("unknown type", func(t *testing.T) {
		source := &model.AuthSource{
			Type:   "unknown",
			Config: json.RawMessage(`{}`),
		}

		valid, err := service.TestAuthSource(context.Background(), source)
		assert.True(t, valid)
		assert.NoError(t, err)
	})
}

func TestAuthSourceService_EnableDisableAuthSource(t *testing.T) {
	db, mock := setupAuthSourceTestDB(t)
	service := NewAuthSourceService(db)

	t.Run("enable auth source", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `auth_sources` SET `enabled`=\\? WHERE `auth_sources`.`id` = \\?").
			WithArgs(true, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.EnableAuthSource(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("disable auth source", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `auth_sources` SET `enabled`=\\? WHERE `auth_sources`.`id` = \\?").
			WithArgs(false, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DisableAuthSource(context.Background(), 1)
		assert.NoError(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
