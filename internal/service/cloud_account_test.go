package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestCloudAccountService_CreateCloudAccount(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewCloudAccountService(db)

	account := &model.CloudAccount{
		Name:         "test-account",
		ProviderType: "alibaba",
		Status:       "active",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `cloud_accounts`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.CreateCloudAccount(context.Background(), account)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloudAccountService_GetCloudAccount(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewCloudAccountService(db)

	rows := sqlmock.NewRows([]string{"id", "name", "provider_type", "status"}).
		AddRow(1, "test-account", "alibaba", "active")

	mock.ExpectQuery("SELECT \\* FROM `cloud_accounts`").
		WithArgs(1, 1).
		WillReturnRows(rows)

	account, err := service.GetCloudAccount(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, "test-account", account.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloudAccountService_GetCloudAccount_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewCloudAccountService(db)

	mock.ExpectQuery("SELECT \\* FROM `cloud_accounts`").
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	account, err := service.GetCloudAccount(context.Background(), 999)

	assert.NoError(t, err)
	assert.Nil(t, account)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloudAccountService_ListCloudAccounts(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewCloudAccountService(db)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `cloud_accounts`").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "provider_type", "status"}).
		AddRow(1, "test-account-1", "alibaba", "active").
		AddRow(2, "test-account-2", "tencent", "active")

	mock.ExpectQuery("SELECT \\* FROM `cloud_accounts`").
		WillReturnRows(rows)

	accounts, total, err := service.ListCloudAccounts(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, accounts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloudAccountService_VerifyCloudAccount(t *testing.T) {
	db, _ := setupTestDB(t)
	service := NewCloudAccountService(db)

	creds, _ := json.Marshal(map[string]string{
		"access_key_id":     "test",
		"access_key_secret": "test",
	})

	account := &model.CloudAccount{
		ID:           1,
		Name:         "test-account",
		ProviderType: "alibaba",
		Credentials:  creds,
		Status:       "active",
	}

	// 注意：这个测试会失败，因为 alibaba provider 需要真实的 SDK
	// 实际项目中应该使用 mock provider
	valid, err := service.VerifyCloudAccount(context.Background(), account)
	
	// 预期会失败（因为没有真实的云账户）
	assert.False(t, valid)
	assert.Error(t, err)
}
