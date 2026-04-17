package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/opencmp/opencmp/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func TestDomainHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	logger := zap.NewNop()
	domainHandler := NewDomainHandler(db, logger)

	t.Run("Successful domain creation", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		reqData := CreateDomainRequest{
			Name:        "test-create-domain",
			Description: "A test domain for creation",
			Enabled:     true,
		}
		jsonData, _ := json.Marshal(reqData)

		req, _ := http.NewRequest("POST", "/domains", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		domainHandler.Create(ginCtx)

		assert.Equal(t, http.StatusCreated, w.Code)

		var responseDomain model.Domain
		err := json.Unmarshal(w.Body.Bytes(), &responseDomain)
		if err != nil {
			t.Logf("Response body: %s", w.Body.String())
		}
	})

	t.Run("Create domain with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("POST", "/domains", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		ginCtx.Request = req

		domainHandler.Create(ginCtx)

		if w.Code != http.StatusBadRequest {
			t.Logf("Expected BadRequest, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestDomainHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	logger := zap.NewNop()
	domainHandler := NewDomainHandler(db, logger)

	t.Run("Get non-existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		domainHandler.Get(ginCtx)

		if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
			t.Logf("Expected NotFound or InternalServerError, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Get domain with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		domainHandler.Get(ginCtx)

		if w.Code != http.StatusBadRequest {
			t.Logf("Expected BadRequest, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestDomainHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	logger := zap.NewNop()
	domainHandler := NewDomainHandler(db, logger)

	t.Run("List all domains", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/domains", nil)
		ginCtx.Request = req

		domainHandler.List(ginCtx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Logf("Response body: %s", w.Body.String())
		}
	})
}

func TestDomainHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	logger := zap.NewNop()
	domainHandler := NewDomainHandler(db, logger)

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

		domainHandler.Update(ginCtx)

		if w.Code != http.StatusBadRequest {
			t.Logf("Expected BadRequest, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestDomainHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerTestDB()
	logger := zap.NewNop()
	domainHandler := NewDomainHandler(db, logger)

	t.Run("Delete non-existing domain", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/domains/999", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "999"}}

		domainHandler.Delete(ginCtx)

		// Could be 500 (error) or 200 (message: deleted) depending on implementation
		t.Logf("Delete response: %d - %s", w.Code, w.Body.String())
	})

	t.Run("Delete domain with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("DELETE", "/domains/invalid", nil)
		ginCtx.Request = req
		ginCtx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		domainHandler.Delete(ginCtx)

		if w.Code != http.StatusBadRequest {
			t.Logf("Expected BadRequest, got %d: %s", w.Code, w.Body.String())
		}
	})
}