package service

import (
	"testing"

	"github.com/yourorg/openCMP/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	err = db.AutoMigrate(&model.Domain{})
	if err != nil {
		panic("failed to migrate test database")
	}

	return db
}

func TestDomainService_CreateDomain(t *testing.T) {
	db := setupTestDB()
	service := NewDomainService(db)

	t.Run("Successful domain creation", func(t *testing.T) {
		domain := &model.Domain{
			Name:        "test-domain",
			Description: "A test domain",
		}

		err := service.CreateDomain(domain)
		assert.NoError(t, err)
		assert.NotZero(t, domain.ID)
		assert.Equal(t, "test-domain", domain.Name)
		assert.Equal(t, "A test domain", domain.Description)
		assert.Equal(t, "active", domain.Status)
	})

	t.Run("Duplicate domain name", func(t *testing.T) {
		domain1 := &model.Domain{
			Name: "duplicate-domain",
		}
		domain2 := &model.Domain{
			Name: "duplicate-domain",
		}

		err1 := service.CreateDomain(domain1)
		assert.NoError(t, err1)

		err2 := service.CreateDomain(domain2)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "already exists")
	})
}

func TestDomainService_GetDomainByID(t *testing.T) {
	db := setupTestDB()
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "get-test-domain",
		Description: "A test domain for getting",
	}
	err := service.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Get existing domain", func(t *testing.T) {
		result, err := service.GetDomainByID(domain.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.ID, result.ID)
		assert.Equal(t, domain.Name, result.Name)
		assert.Equal(t, domain.Description, result.Description)
	})

	t.Run("Get non-existing domain", func(t *testing.T) {
		result, err := service.GetDomainByID(999)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDomainService_GetDomainByName(t *testing.T) {
	db := setupTestDB()
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "get-by-name-test",
		Description: "A test domain for getting by name",
	}
	err := service.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Get existing domain by name", func(t *testing.T) {
		result, err := service.GetDomainByName("get-by-name-test")
		assert.NoError(t, err)
		assert.Equal(t, domain.ID, result.ID)
		assert.Equal(t, domain.Name, result.Name)
		assert.Equal(t, domain.Description, result.Description)
	})

	t.Run("Get non-existing domain by name", func(t *testing.T) {
		result, err := service.GetDomainByName("non-existent-domain")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDomainService_ListDomains(t *testing.T) {
	db := setupTestDB()
	service := NewDomainService(db)

	// Create multiple domains
	domains := []*model.Domain{
		{Name: "list-test-1", Description: "First test domain"},
		{Name: "list-test-2", Description: "Second test domain"},
		{Name: "list-test-3", Description: "Third test domain"},
	}

	for _, domain := range domains {
		err := service.CreateDomain(domain)
		assert.NoError(t, err)
	}

	t.Run("List all domains", func(t *testing.T) {
		results, total, err := service.ListDomains(0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 3)

		// Check that all created domains are in the result
		names := make(map[string]bool)
		for _, d := range results {
			names[d.Name] = true
		}
		assert.True(t, names["list-test-1"])
		assert.True(t, names["list-test-2"])
		assert.True(t, names["list-test-3"])
	})

	t.Run("List domains with pagination", func(t *testing.T) {
		results, total, err := service.ListDomains(0, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 2)

		results2, total2, err := service.ListDomains(2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total2)
		assert.Len(t, results2, 1)
	})
}

func TestDomainService_UpdateDomain(t *testing.T) {
	db := setupTestDB()
	service := NewDomainService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "update-test",
		Description: "Original description",
	}
	err := service.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Update domain successfully", func(t *testing.T) {
		updates := map[string]interface{}{
			"Description": "Updated description",
			"Status":      "inactive",
		}

		result, err := service.UpdateDomain(domain.ID, updates)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", result.Description)
		assert.Equal(t, "inactive", result.Status)
	})

	t.Run("Update domain name to existing name", func(t *testing.T) {
		// Create another domain
		domain2 := &model.Domain{
			Name: "other-domain",
		}
		err := service.CreateDomain(domain2)
		assert.NoError(t, err)

		// Try to update domain2's name to domain's name
		updates := map[string]interface{}{
			"Name": "update-test",
		}

		result, err := service.UpdateDomain(domain2.ID, updates)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
	})
}