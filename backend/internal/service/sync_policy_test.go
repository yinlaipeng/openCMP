package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupSyncPolicyTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestSyncPolicyService_GetSyncPolicy(t *testing.T) {
	db, mock := setupSyncPolicyTestDB(t)
	service := NewSyncPolicyService(db)

	// Policy row
	policyRows := sqlmock.NewRows([]string{"id", "name", "remarks", "status", "enabled", "scope", "domain_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "test-policy", "test remarks", "active", true, "domain", 1, nil, nil, nil)

	mock.ExpectQuery("SELECT \\* FROM `sync_policies`").
		WithArgs(1, 1).
		WillReturnRows(policyRows)

	// Rules rows (using sync_policy_rules table)
	rulesRows := sqlmock.NewRows([]string{"id", "sync_policy_id", "condition_type", "resource_mapping", "target_project_id", "target_project_name", "deleted_at"})
	mock.ExpectQuery("SELECT \\* FROM `sync_policy_rules`").
		WillReturnRows(rulesRows)

	policy, err := service.GetSyncPolicy(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, policy)
	assert.Equal(t, "test-policy", policy.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncPolicyService_GetSyncPolicy_NotFound(t *testing.T) {
	db, mock := setupSyncPolicyTestDB(t)
	service := NewSyncPolicyService(db)

	mock.ExpectQuery("SELECT \\* FROM `sync_policies`").
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	policy, err := service.GetSyncPolicy(context.Background(), 999)

	assert.NoError(t, err)
	assert.Nil(t, policy)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncPolicyService_ListSyncPolicies(t *testing.T) {
	db, mock := setupSyncPolicyTestDB(t)
	service := NewSyncPolicyService(db)

	// Count rows
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `sync_policies`").
		WillReturnRows(countRows)

	// Policy rows
	policyRows := sqlmock.NewRows([]string{"id", "name", "remarks", "status", "enabled", "scope", "domain_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "policy-1", "remarks 1", "active", true, "domain", 1, nil, nil, nil).
		AddRow(2, "policy-2", "remarks 2", "active", true, "domain", 1, nil, nil, nil)

	mock.ExpectQuery("SELECT \\* FROM `sync_policies`").
		WillReturnRows(policyRows)

	// Preload rules for policy 1
	mock.ExpectQuery("SELECT \\* FROM `sync_policy_rules`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "sync_policy_id"}))

	policies, total, err := service.ListSyncPolicies(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, policies, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncPolicyService_ListSyncPolicies_Empty(t *testing.T) {
	db, mock := setupSyncPolicyTestDB(t)
	service := NewSyncPolicyService(db)

	// Count rows
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `sync_policies`").
		WillReturnRows(countRows)

	// Empty policy rows
	policyRows := sqlmock.NewRows([]string{"id", "name", "remarks", "status", "enabled", "scope", "domain_id"})
	mock.ExpectQuery("SELECT \\* FROM `sync_policies`").
		WillReturnRows(policyRows)

	policies, total, err := service.ListSyncPolicies(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, policies, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}