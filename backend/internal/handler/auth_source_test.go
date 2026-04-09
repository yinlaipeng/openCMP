package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

func setupAuthSourceHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.AuthSource{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func newTestAuthSourceHandler(db *gorm.DB) *AuthSourceHandler {
	logger, _ := zap.NewDevelopment()
	return &AuthSourceHandler{
		service: service.NewAuthSourceService(db),
		logger:  logger,
	}
}

// createTestAuthSource 在测试 DB 中创建一条认证源，返回其 ID
func createTestAuthSource(t *testing.T, db *gorm.DB, name, srcType string) *model.AuthSource {
	t.Helper()
	cfgJSON := []byte(`{"url":"ldap://ldap.example.com:389"}`)
	src := &model.AuthSource{
		Name:       name,
		Type:       srcType,
		Scope:      "system",
		Enabled:    true,
		AutoCreate: false,
		Config:     cfgJSON,
	}
	if err := db.Create(src).Error; err != nil {
		t.Fatalf("failed to create test auth source: %v", err)
	}
	return src
}

// TestAuthSourceHandler_List - GET /auth-sources 返回列表
func TestAuthSourceHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	createTestAuthSource(t, db, "ldap-1", "ldap")
	createTestAuthSource(t, db, "ldap-2", "ldap")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/auth-sources", nil)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(2), resp["total"])
	items, ok := resp["items"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, items, 2)
}

// TestAuthSourceHandler_Get - GET /auth-sources/:id 返回单条
func TestAuthSourceHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	src := createTestAuthSource(t, db, "test-get", "ldap")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/auth-sources/"+strconv.Itoa(int(src.ID)), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got model.AuthSource
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, src.ID, got.ID)
	assert.Equal(t, "test-get", got.Name)
}

// TestAuthSourceHandler_Get_NotFound - 不存在的 ID 返回 404
func TestAuthSourceHandler_Get_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/auth-sources/9999", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "9999"}}

	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestAuthSourceHandler_Create - POST 创建 LDAP 认证源
func TestAuthSourceHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	body := map[string]interface{}{
		"name":        "new-ldap",
		"type":        "ldap",
		"scope":       "system",
		"enabled":     true,
		"auto_create": false,
		"config": map[string]interface{}{
			"url":     "ldap://ldap.example.com:389",
			"base_dn": "dc=example,dc=com",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/auth-sources", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got model.AuthSource
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "new-ldap", got.Name)
	assert.Equal(t, "ldap", got.Type)
	assert.Equal(t, "system", got.Scope)
}

// TestAuthSourceHandler_Create_MissingName - 缺少 name 返回 400
func TestAuthSourceHandler_Create_MissingName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	body := map[string]interface{}{
		"type":  "ldap",
		"scope": "system",
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/auth-sources", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAuthSourceHandler_Update - PUT 更新认证源
func TestAuthSourceHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	src := createTestAuthSource(t, db, "update-me", "ldap")

	body := map[string]interface{}{
		"name":        "updated-ldap",
		"type":        "ldap",
		"scope":       "system",
		"enabled":     true,
		"auto_create": true,
		"config": map[string]interface{}{
			"url": "ldap://new.example.com:389",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPut, "/auth-sources/"+strconv.Itoa(int(src.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got model.AuthSource
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "updated-ldap", got.Name)
	assert.True(t, got.AutoCreate)
}

// TestAuthSourceHandler_Delete - DELETE 删除认证源
func TestAuthSourceHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	src := createTestAuthSource(t, db, "delete-me", "ldap")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/auth-sources/"+strconv.Itoa(int(src.ID)), nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证已被软删除（GetAuthSource 返回 nil）
	svc := service.NewAuthSourceService(db)
	deleted, err := svc.GetAuthSource(c.Request.Context(), src.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

// TestAuthSourceHandler_Enable - POST /enable
func TestAuthSourceHandler_Enable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	// 先创建一个禁用的认证源
	cfgJSON := []byte(`{}`)
	src := &model.AuthSource{
		Name:    "enable-me",
		Type:    "ldap",
		Scope:   "system",
		Enabled: false,
		Config:  cfgJSON,
	}
	db.Create(src)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/auth-sources/"+strconv.Itoa(int(src.ID))+"/enable", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Enable(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "enabled", resp["message"])
}

// TestAuthSourceHandler_Disable - POST /disable
func TestAuthSourceHandler_Disable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	src := createTestAuthSource(t, db, "disable-me", "ldap")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/auth-sources/"+strconv.Itoa(int(src.ID))+"/disable", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Disable(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "disabled", resp["message"])
}

// TestAuthSourceHandler_Test - POST /test（连接测试）
// local 类型默认返回 valid=true（无配置即为真实连接通过的 stub）
func TestAuthSourceHandler_Test(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthSourceHandlerTestDB(t)
	h := newTestAuthSourceHandler(db)

	// 使用 local 类型（TestAuthSource 对 local 直接返回 true）
	cfgJSON := []byte(`{}`)
	src := &model.AuthSource{
		Name:    "local-test",
		Type:    "local",
		Scope:   "system",
		Enabled: true,
		Config:  cfgJSON,
	}
	db.Create(src)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/auth-sources/"+strconv.Itoa(int(src.ID))+"/test", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(src.ID))}}

	h.Test(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["valid"])
}
