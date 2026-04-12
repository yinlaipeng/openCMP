package service

import (
	"encoding/json"
	"testing"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProjectTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	err = db.AutoMigrate(&model.Domain{}, &model.Project{})
	if err != nil {
		panic("failed to migrate test database")
	}

	return db
}

func TestProjectService_CreateProject(t *testing.T) {
	db := setupProjectTestDB()
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)

	// Create a domain first
	domain := &model.Domain{
		Name:        "test-domain",
		Description: "A test domain",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Successful project creation", func(t *testing.T) {
		project := &model.Project{
			Name:        "test-project",
			Description: "A test project",
			DomainID:    domain.ID,
		}

		err := projectService.CreateProject(project)
		assert.NoError(t, err)
		assert.NotZero(t, project.ID)
		assert.Equal(t, "test-project", project.Name)
		assert.Equal(t, "A test project", project.Description)
		assert.Equal(t, domain.ID, project.DomainID)
		assert.Equal(t, "active", project.Status)
	})

	t.Run("Create project with non-existent domain", func(t *testing.T) {
		project := &model.Project{
			Name:        "invalid-project",
			Description: "A project with invalid domain",
			DomainID:    999, // Non-existent domain ID
		}

		err := projectService.CreateProject(project)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Create duplicate project name in same domain", func(t *testing.T) {
		project1 := &model.Project{
			Name:        "duplicate-project",
			Description: "First project",
			DomainID:    domain.ID,
		}
		project2 := &model.Project{
			Name:        "duplicate-project",
			Description: "Second project",
			DomainID:    domain.ID,
		}

		err1 := projectService.CreateProject(project1)
		assert.NoError(t, err1)

		err2 := projectService.CreateProject(project2)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "already exists")
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	db := setupProjectTestDB()
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)

	// Create a domain and project first
	domain := &model.Domain{
		Name:        "get-test-domain",
		Description: "A test domain",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "get-test-project",
		Description: "A test project",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(project)
	assert.NoError(t, err)

	t.Run("Get existing project", func(t *testing.T) {
		result, err := projectService.GetProjectByID(project.ID)
		assert.NoError(t, err)
		assert.Equal(t, project.ID, result.ID)
		assert.Equal(t, project.Name, result.Name)
		assert.Equal(t, project.Description, result.Description)
		assert.Equal(t, project.DomainID, result.DomainID)
		assert.Equal(t, domain.Name, result.Domain.Name) // Check if domain is preloaded
	})

	t.Run("Get non-existing project", func(t *testing.T) {
		result, err := projectService.GetProjectByID(999)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProjectService_GetProjectByName(t *testing.T) {
	db := setupProjectTestDB()
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)

	// Create a domain and project first
	domain := &model.Domain{
		Name:        "get-by-name-domain",
		Description: "A test domain",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "get-by-name-project",
		Description: "A test project",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(project)
	assert.NoError(t, err)

	t.Run("Get existing project by name", func(t *testing.T) {
		result, err := projectService.GetProjectByName(domain.ID, "get-by-name-project")
		assert.NoError(t, err)
		assert.Equal(t, project.ID, result.ID)
		assert.Equal(t, project.Name, result.Name)
		assert.Equal(t, project.Description, result.Description)
		assert.Equal(t, project.DomainID, result.DomainID)
	})

	t.Run("Get non-existing project by name", func(t *testing.T) {
		result, err := projectService.GetProjectByName(domain.ID, "non-existent-project")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Get project by name in wrong domain", func(t *testing.T) {
		// Create another domain
		domain2 := &model.Domain{
			Name:        "another-domain",
			Description: "Another test domain",
		}
		err := domainService.CreateDomain(domain2)
		assert.NoError(t, err)

		result, err := projectService.GetProjectByName(domain2.ID, "get-by-name-project")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProjectService_ListProjects(t *testing.T) {
	db := setupProjectTestDB()
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)

	// Create domains and projects
	domain1 := &model.Domain{
		Name: "list-test-domain-1",
	}
	domain2 := &model.Domain{
		Name: "list-test-domain-2",
	}
	err := domainService.CreateDomain(domain1)
	assert.NoError(t, err)
	err = domainService.CreateDomain(domain2)
	assert.NoError(t, err)

	projects := []*model.Project{
		{Name: "list-test-1", Description: "First test project", DomainID: domain1.ID},
		{Name: "list-test-2", Description: "Second test project", DomainID: domain1.ID},
		{Name: "list-test-3", Description: "Third test project", DomainID: domain2.ID},
	}

	for _, project := range projects {
		err := projectService.CreateProject(project)
		assert.NoError(t, err)
	}

	t.Run("List all projects", func(t *testing.T) {
		results, total, err := projectService.ListProjects(0, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 3)

		// Check that all created projects are in the result
		names := make(map[string]bool)
		for _, p := range results {
			names[p.Name] = true
		}
		assert.True(t, names["list-test-1"])
		assert.True(t, names["list-test-2"])
		assert.True(t, names["list-test-3"])
	})

	t.Run("List projects with domain filter", func(t *testing.T) {
		results, total, err := projectService.ListProjects(0, 10, &domain1.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)

		// Check that only projects from domain1 are in the result
		names := make(map[string]bool)
		for _, p := range results {
			assert.Equal(t, domain1.ID, p.DomainID)
			names[p.Name] = true
		}
		assert.True(t, names["list-test-1"])
		assert.True(t, names["list-test-2"])
		assert.False(t, names["list-test-3"])
	})

	t.Run("List projects with pagination", func(t *testing.T) {
		results, total, err := projectService.ListProjects(0, 2, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 2)

		results2, total2, err := projectService.ListProjects(2, 2, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total2)
		assert.Len(t, results2, 1)
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	db := setupProjectTestDB()
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)

	// Create domains and project
	domain1 := &model.Domain{
		Name: "update-test-domain-1",
	}
	domain2 := &model.Domain{
		Name: "update-test-domain-2",
	}
	err := domainService.CreateDomain(domain1)
	assert.NoError(t, err)
	err = domainService.CreateDomain(domain2)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "update-test-project",
		Description: "Original description",
		DomainID:    domain1.ID,
	}
	err = projectService.CreateProject(project)
	assert.NoError(t, err)

	t.Run("Update project successfully", func(t *testing.T) {
		updates := map[string]interface{}{
			"Description": "Updated description",
			"Status":      "inactive",
		}

		result, err := projectService.UpdateProject(project.ID, updates)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", result.Description)
		assert.Equal(t, "inactive", result.Status)
	})

	t.Run("Update project to different domain", func(t *testing.T) {
		updates := map[string]interface{}{
			"DomainID": domain2.ID,
		}

		result, err := projectService.UpdateProject(project.ID, updates)
		assert.NoError(t, err)
		assert.Equal(t, domain2.ID, result.DomainID)
	})

	t.Run("Update project to duplicate name in target domain", func(t *testing.T) {
		// Create another project in domain2
		project2 := &model.Project{
			Name:        "other-project",
			Description: "Another project",
			DomainID:    domain2.ID,
		}
		err := projectService.CreateProject(project2)
		assert.NoError(t, err)

		// Try to update project to have same name as project2 in domain2
		updates := map[string]interface{}{
			"Name":     "other-project", // Same name as project2
			"DomainID": domain2.ID,      // Move to same domain as project2
		}

		result, err := projectService.UpdateProject(project.ID, updates)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Update project quota with valid JSON", func(t *testing.T) {
		quota := `{"vm_count": 10, "storage_gb": 100}`
		updates := map[string]interface{}{
			"Quota": quota,
		}

		result, err := projectService.UpdateProject(project.ID, updates)
		assert.NoError(t, err)
		assert.Equal(t, quota, result.Quota)
	})

	t.Run("Update project quota with invalid JSON", func(t *testing.T) {
		updates := map[string]interface{}{
			"Quota": "{invalid_json}",
		}

		result, err := projectService.UpdateProject(project.ID, updates)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid quota JSON format")
	})
}
