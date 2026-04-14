package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

func setupProjectTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&model.Domain{}, &model.Project{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestProjectService_CreateProject(t *testing.T) {
	db := setupProjectTestDB(t)
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)
	ctx := context.Background()

	// Create a domain first
	domain := &model.Domain{
		Name:        "test-domain",
		Description: "A test domain",
	}
	err := domainService.CreateDomain(ctx, domain)
	assert.NoError(t, err)

	t.Run("Successful project creation", func(t *testing.T) {
		project := &model.Project{
			Name:        "test-project",
			Description: "A test project",
			DomainID:    domain.ID,
		}

		err := projectService.CreateProject(ctx, project)
		assert.NoError(t, err)
		assert.NotZero(t, project.ID)
		assert.Equal(t, "test-project", project.Name)
		assert.Equal(t, "A test project", project.Description)
		assert.Equal(t, domain.ID, project.DomainID)
	})
}

func TestProjectService_GetProject(t *testing.T) {
	db := setupProjectTestDB(t)
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)
	ctx := context.Background()

	// Create a domain and project first
	domain := &model.Domain{
		Name:        "get-test-domain",
		Description: "A test domain",
	}
	err := domainService.CreateDomain(ctx, domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "get-test-project",
		Description: "A test project",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(ctx, project)
	assert.NoError(t, err)

	t.Run("Get existing project", func(t *testing.T) {
		result, err := projectService.GetProject(ctx, project.ID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, project.ID, result.ID)
		assert.Equal(t, project.Name, result.Name)
		assert.Equal(t, project.Description, result.Description)
		assert.Equal(t, project.DomainID, result.DomainID)
	})

	t.Run("Get non-existing project", func(t *testing.T) {
		result, err := projectService.GetProject(ctx, 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestProjectService_ListProjects(t *testing.T) {
	db := setupProjectTestDB(t)
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)
	ctx := context.Background()

	// Create domains and projects
	domain1 := &model.Domain{
		Name: "list-test-domain-1",
	}
	domain2 := &model.Domain{
		Name: "list-test-domain-2",
	}
	err := domainService.CreateDomain(ctx, domain1)
	assert.NoError(t, err)
	err = domainService.CreateDomain(ctx, domain2)
	assert.NoError(t, err)

	projects := []*model.Project{
		{Name: "list-test-1", Description: "First test project", DomainID: domain1.ID},
		{Name: "list-test-2", Description: "Second test project", DomainID: domain1.ID},
		{Name: "list-test-3", Description: "Third test project", DomainID: domain2.ID},
	}

	for _, project := range projects {
		err := projectService.CreateProject(ctx, project)
		assert.NoError(t, err)
	}

	t.Run("List all projects", func(t *testing.T) {
		results, total, err := projectService.ListProjects(ctx, nil, "", nil, 10, 0)
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
		results, total, err := projectService.ListProjects(ctx, &domain1.ID, "", nil, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)

		// Check that only projects from domain1 are in the result
		for _, p := range results {
			assert.Equal(t, domain1.ID, p.DomainID)
		}
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	db := setupProjectTestDB(t)
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)
	ctx := context.Background()

	// Create domain and project
	domain := &model.Domain{
		Name: "update-test-domain",
	}
	err := domainService.CreateDomain(ctx, domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "update-test-project",
		Description: "Original description",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(ctx, project)
	assert.NoError(t, err)

	t.Run("Update project successfully", func(t *testing.T) {
		project.Description = "Updated description"
		err := projectService.UpdateProject(ctx, project)
		assert.NoError(t, err)

		result, err := projectService.GetProject(ctx, project.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", result.Description)
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	db := setupProjectTestDB(t)
	domainService := NewDomainService(db)
	projectService := NewProjectService(db)
	ctx := context.Background()

	// Create domain and project
	domain := &model.Domain{
		Name: "delete-test-domain",
	}
	err := domainService.CreateDomain(ctx, domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "delete-test-project",
		Description: "A test project",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(ctx, project)
	assert.NoError(t, err)

	// Delete project
	err = projectService.DeleteProject(ctx, project.ID)
	assert.NoError(t, err)

	// Verify project is deleted
	result, err := projectService.GetProject(ctx, project.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}