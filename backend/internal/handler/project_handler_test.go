package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProjectHandler_CreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	projectService := service.NewProjectService(db)
	projectHandler := NewProjectHandler(projectService)

	// Create a domain first
	domain := &model.Domain{
		Name:        "test-domain-for-project",
		Description: "A test domain for project creation",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Successful project creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		projectData := model.Project{
			Name:        "test-create-project",
			Description: "A test project for creation",
			DomainID:    domain.ID,
		}
		jsonData, _ := json.Marshal(projectData)

		req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		projectHandler.CreateProject(ginCtx)

		assert.Equal(t, http.StatusCreated, w.Code)

		var responseProject model.Project
		err := json.Unmarshal(w.Body.Bytes(), &responseProject)
		assert.NoError(t, err)
		assert.Equal(t, "test-create-project", responseProject.Name)
		assert.Equal(t, "A test project for creation", responseProject.Description)
		assert.Equal(t, domain.ID, responseProject.DomainID)
		assert.Equal(t, "active", responseProject.Status)
	})

	t.Run("Create project with non-existent domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		projectData := model.Project{
			Name:        "test-invalid-domain-project",
			Description: "A test project with invalid domain",
			DomainID:    999, // Non-existent domain
		}
		jsonData, _ := json.Marshal(projectData)

		req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		projectHandler.CreateProject(ginCtx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Create project with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("POST", "/projects", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		projectHandler.CreateProject(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_GetProjectByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	projectService := service.NewProjectService(db)
	projectHandler := NewProjectHandler(projectService)

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
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects/"+strconv.Itoa(int(project.ID)), nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(project.ID))}}

		projectHandler.GetProjectByID(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseProject model.Project
		err := json.Unmarshal(w.Body.Bytes(), &responseProject)
		assert.NoError(t, err)
		assert.Equal(t, project.ID, responseProject.ID)
		assert.Equal(t, project.Name, responseProject.Name)
		assert.Equal(t, domain.Name, responseProject.Domain.Name) // Check if domain is preloaded
	})

	t.Run("Get non-existing project", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		projectHandler.GetProjectByID(ginCtx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get project with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		projectHandler.GetProjectByID(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_ListProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	projectService := service.NewProjectService(db)
	projectHandler := NewProjectHandler(projectService)

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
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects", nil)
		ginCtx.Request = req

		projectHandler.ListProjects(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].([]interface{})
		total := int(response["total"].(float64))

		assert.Equal(t, 3, total)
		assert.Len(t, data, 3)
	})

	t.Run("List projects with domain filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects?domain_id="+strconv.Itoa(int(domain1.ID)), nil)
		ginCtx.Request = req
		q := req.URL.Query()
		q.Add("domain_id", strconv.Itoa(int(domain1.ID)))
		req.URL.RawQuery = q.Encode()
		ginCtx.Request = req

		projectHandler.ListProjects(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].([]interface{})
		total := int(response["total"].(float64))

		assert.Equal(t, 2, total)
		assert.Len(t, data, 2)
	})

	t.Run("List projects with pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects?offset=0&limit=2", nil)
		ginCtx.Request = req
		q := req.URL.Query()
		q.Add("offset", "0")
		q.Add("limit", "2")
		req.URL.RawQuery = q.Encode()
		ginCtx.Request = req

		projectHandler.ListProjects(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].([]interface{})
		total := int(response["total"].(float64))
		limit := int(response["limit"].(float64))

		assert.Equal(t, 3, total)
		assert.Len(t, data, 2)
		assert.Equal(t, 2, limit)
	})

	t.Run("List projects with invalid domain ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/projects?domain_id=invalid", nil)
		ginCtx.Request = req
		q := req.URL.Query()
		q.Add("domain_id", "invalid")
		req.URL.RawQuery = q.Encode()
		ginCtx.Request = req

		projectHandler.ListProjects(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	projectService := service.NewProjectService(db)
	projectHandler := NewProjectHandler(projectService)

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
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		updates := map[string]interface{}{
			"Description": "Updated description",
			"Status":      "inactive",
		}
		jsonData, _ := json.Marshal(updates)

		req, _ := http.NewRequest("PUT", "/projects/"+strconv.Itoa(int(project.ID)), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(project.ID))}}

		projectHandler.UpdateProject(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseProject model.Project
		err := json.Unmarshal(w.Body.Bytes(), &responseProject)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", responseProject.Description)
		assert.Equal(t, "inactive", responseProject.Status)
	})

	t.Run("Update project with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		updates := map[string]interface{}{
			"Description": "Updated description",
		}
		jsonData, _ := json.Marshal(updates)

		req, _ := http.NewRequest("PUT", "/projects/invalid", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		projectHandler.UpdateProject(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update project with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("PUT", "/projects/"+strconv.Itoa(int(project.ID)), bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(project.ID))}}

		projectHandler.UpdateProject(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	projectService := service.NewProjectService(db)
	projectHandler := NewProjectHandler(projectService)

	// Create a domain and project first
	domain := &model.Domain{
		Name: "delete-test-domain",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	project := &model.Project{
		Name:        "delete-test-project",
		Description: "A test project for deletion",
		DomainID:    domain.ID,
	}
	err = projectService.CreateProject(project)
	assert.NoError(t, err)

	t.Run("Delete existing project", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/projects/"+strconv.Itoa(int(project.ID)), nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(project.ID))}}

		projectHandler.DeleteProject(ginCtx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Delete non-existing project", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/projects/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		projectHandler.DeleteProject(ginCtx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Delete project with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/projects/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		projectHandler.DeleteProject(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}