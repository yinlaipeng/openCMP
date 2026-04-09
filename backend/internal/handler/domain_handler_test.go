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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTestDB() *gorm.DB {
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

func TestDomainHandler_CreateDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	domainHandler := NewDomainHandler(domainService)

	t.Run("Successful domain creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		domainData := model.Domain{
			Name:        "test-create-domain",
			Description: "A test domain for creation",
		}
		jsonData, _ := json.Marshal(domainData)

		req, _ := http.NewRequest("POST", "/domains", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		domainHandler.CreateDomain(ginCtx)

		assert.Equal(t, http.StatusCreated, w.Code)

		var responseDomain model.Domain
		err := json.Unmarshal(w.Body.Bytes(), &responseDomain)
		assert.NoError(t, err)
		assert.Equal(t, "test-create-domain", responseDomain.Name)
		assert.Equal(t, "A test domain for creation", responseDomain.Description)
		assert.Equal(t, "active", responseDomain.Status)
	})

	t.Run("Create domain with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("POST", "/domains", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		domainHandler.CreateDomain(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDomainHandler_GetDomainByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	domainHandler := NewDomainHandler(domainService)

	// Create a domain first
	domain := &model.Domain{
		Name:        "get-test-domain",
		Description: "A test domain for getting",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Get existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains/"+strconv.Itoa(int(domain.ID)), nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(domain.ID))}}

		domainHandler.GetDomainByID(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseDomain model.Domain
		err := json.Unmarshal(w.Body.Bytes(), &responseDomain)
		assert.NoError(t, err)
		assert.Equal(t, domain.ID, responseDomain.ID)
		assert.Equal(t, domain.Name, responseDomain.Name)
	})

	t.Run("Get non-existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		domainHandler.GetDomainByID(ginCtx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get domain with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		domainHandler.GetDomainByID(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDomainHandler_ListDomains(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	domainHandler := NewDomainHandler(domainService)

	// Create multiple domains
	domains := []*model.Domain{
		{Name: "list-test-1", Description: "First test domain"},
		{Name: "list-test-2", Description: "Second test domain"},
		{Name: "list-test-3", Description: "Third test domain"},
	}

	for _, domain := range domains {
		err := domainService.CreateDomain(domain)
		assert.NoError(t, err)
	}

	t.Run("List all domains", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains", nil)
		ginCtx.Request = req

		domainHandler.ListDomains(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].([]interface{})
		total := int(response["total"].(float64))

		assert.Equal(t, 3, total)
		assert.Len(t, data, 3)
	})

	t.Run("List domains with pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains?offset=0&limit=2", nil)
		ginCtx.Request = req
		q := req.URL.Query()
		q.Add("offset", "0")
		q.Add("limit", "2")
		req.URL.RawQuery = q.Encode()
		ginCtx.Request = req

		domainHandler.ListDomains(ginCtx)

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
}

func TestDomainHandler_UpdateDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	domainHandler := NewDomainHandler(domainService)

	// Create a domain first
	domain := &model.Domain{
		Name:        "update-test-domain",
		Description: "Original description",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Update domain successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		updates := map[string]interface{}{
			"Description": "Updated description",
			"Status":      "inactive",
		}
		jsonData, _ := json.Marshal(updates)

		req, _ := http.NewRequest("PUT", "/domains/"+strconv.Itoa(int(domain.ID)), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(domain.ID))}}

		domainHandler.UpdateDomain(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseDomain model.Domain
		err := json.Unmarshal(w.Body.Bytes(), &responseDomain)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", responseDomain.Description)
		assert.Equal(t, "inactive", responseDomain.Status)
	})

	t.Run("Update domain with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		updates := map[string]interface{}{
			"Description": "Updated description",
		}
		jsonData, _ := json.Marshal(updates)

		req, _ := http.NewRequest("PUT", "/domains/invalid", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		domainHandler.UpdateDomain(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update domain with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("PUT", "/domains/"+strconv.Itoa(int(domain.ID)), bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(domain.ID))}}

		domainHandler.UpdateDomain(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDomainHandler_DeleteDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	domainService := service.NewDomainService(db)
	domainHandler := NewDomainHandler(domainService)

	// Create a domain first
	domain := &model.Domain{
		Name:        "delete-test-domain",
		Description: "A test domain for deletion",
	}
	err := domainService.CreateDomain(domain)
	assert.NoError(t, err)

	t.Run("Delete existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/domains/"+strconv.Itoa(int(domain.ID)), nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(domain.ID))}}

		domainHandler.DeleteDomain(ginCtx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Delete non-existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/domains/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		domainHandler.DeleteDomain(ginCtx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Delete domain with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/domains/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		domainHandler.DeleteDomain(ginCtx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}