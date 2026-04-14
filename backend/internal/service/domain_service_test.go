package service

import (
	"context"
	"testing"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDomainTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&model.Domain{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestDomainService_CreateDomain(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	domain := &model.Domain{
		Name:        "test-domain",
		Description: "A test domain",
		Enabled:     true,
	}

	err := service.CreateDomain(context.Background(), domain)
	assert.NoError(t, err)
	assert.NotZero(t, domain.ID)
	assert.Equal(t, "test-domain", domain.Name)
}

func TestDomainService_GetDomain(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "get-test-domain",
		Description: "A test domain for getting",
		Enabled:     true,
	}
	err := service.CreateDomain(context.Background(), domain)
	assert.NoError(t, err)

	t.Run("Get existing domain", func(t *testing.T) {
		result, err := service.GetDomain(context.Background(), domain.ID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.ID, result.ID)
		assert.Equal(t, domain.Name, result.Name)
	})

	t.Run("Get non-existing domain", func(t *testing.T) {
		result, err := service.GetDomain(context.Background(), 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestDomainService_ListDomains(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	// Create multiple domains
	domains := []*model.Domain{
		{Name: "list-test-1", Description: "First test domain", Enabled: true},
		{Name: "list-test-2", Description: "Second test domain", Enabled: true},
		{Name: "list-test-3", Description: "Third test domain", Enabled: true},
	}

	for _, d := range domains {
		err := service.CreateDomain(context.Background(), d)
		assert.NoError(t, err)
	}

	results, total, err := service.ListDomains(context.Background(), "", nil, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, results, 3)
}

func TestDomainService_UpdateDomain(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "update-test",
		Description: "Original description",
		Enabled:     true,
	}
	err := service.CreateDomain(context.Background(), domain)
	assert.NoError(t, err)

	domain.Description = "Updated description"
	err = service.UpdateDomain(context.Background(), domain)
	assert.NoError(t, err)

	result, err := service.GetDomain(context.Background(), domain.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated description", result.Description)
}

func TestDomainService_DeleteDomain(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "delete-test",
		Description: "A test domain for deleting",
		Enabled:     true,
	}
	err := service.CreateDomain(context.Background(), domain)
	assert.NoError(t, err)

	err = service.DeleteDomain(context.Background(), domain.ID)
	assert.NoError(t, err)

	result, err := service.GetDomain(context.Background(), domain.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDomainService_EnableDisableDomain(t *testing.T) {
	db := setupDomainTestDB(t)
	service := NewDomainService(db)

	// Create a disabled domain
	domain := &model.Domain{
		Name:    "enable-test",
		Enabled: false,
	}
	err := service.CreateDomain(context.Background(), domain)
	assert.NoError(t, err)

	err = service.EnableDomain(context.Background(), domain.ID)
	assert.NoError(t, err)

	result, err := service.GetDomain(context.Background(), domain.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Enabled)

	err = service.DisableDomain(context.Background(), domain.ID)
	assert.NoError(t, err)

	result, err = service.GetDomain(context.Background(), domain.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Enabled)
}