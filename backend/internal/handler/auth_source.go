package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// AuthSourceHandler 认证源 Handler
type AuthSourceHandler struct {
	service *service.AuthSourceService
	logger  *zap.Logger
}

// NewAuthSourceHandler 创建认证源 Handler
func NewAuthSourceHandler(db *gorm.DB, logger *zap.Logger) *AuthSourceHandler {
	return &AuthSourceHandler{
		service: service.NewAuthSourceService(db),
		logger:  logger,
	}
}

// CreateAuthSourceRequest 创建认证源请求
type CreateAuthSourceRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Type        string                 `json:"type" binding:"required"` // ldap/oidc/saml
	Scope       string                 `json:"scope"`                   // system/domain
	DomainID    *uint                  `json:"domain_id"`
	Config      map[string]interface{} `json:"config"`
	Enabled     bool                   `json:"enabled"`
}

// List 列出认证源
func (h *AuthSourceHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get filter parameters
	filters := make(map[string]interface{})

	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if sourceType := c.Query("type"); sourceType != "" {
		filters["type"] = sourceType
	}
	if scope := c.Query("scope"); scope != "" {
		filters["scope"] = scope
	}
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			filters["enabled"] = enabled
		}
	}

	sources, total, err := h.service.ListAuthSourcesWithFilters(c.Request.Context(), limit, offset, filters)
	if err != nil {
		h.logger.Error("failed to list auth sources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sources,
		"total": total,
	})
}

// Get 获取认证源详情
func (h *AuthSourceHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	source, err := h.service.GetAuthSource(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auth source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

// Create 创建认证源
func (h *AuthSourceHandler) Create(c *gin.Context) {
	var req CreateAuthSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认 scope
	if req.Scope == "" {
		req.Scope = "system"
	}

	// 如果是域范围，验证 domain_id
	if req.Scope == "domain" && req.DomainID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain_id is required for domain scope"})
		return
	}

	source := &model.AuthSource{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Scope:       req.Scope,
		DomainID:    req.DomainID,
		Enabled:     req.Enabled,
	}

	// 解析并存储 config JSON
	if req.Config != nil {
		configData, err := json.Marshal(req.Config)
		if err != nil {
			h.logger.Error("failed to marshal config", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config"})
			return
		}
		source.Config = configData
	}

	if err := h.service.CreateAuthSource(c.Request.Context(), source); err != nil {
		h.logger.Error("failed to create auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, source)
}

// Sync 同步认证源用户
func (h *AuthSourceHandler) Sync(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.service.SyncUsers(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to sync auth source users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Test 测试认证源连接
func (h *AuthSourceHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	source, err := h.service.GetAuthSource(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auth source not found"})
		return
	}

	// Allow mock LDAP tests without actual connection
	valid, err := h.service.TestAuthSource(c.Request.Context(), source)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "connection failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "connection successful"})
}

// Update 更新认证源
func (h *AuthSourceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	source, err := h.service.GetAuthSource(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auth source not found"})
		return
	}

	var req CreateAuthSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source.Name = req.Name
	source.Description = req.Description
	source.Scope = req.Scope
	source.DomainID = req.DomainID
	source.Enabled = req.Enabled
	if req.Config != nil {
		configData, err := json.Marshal(req.Config)
		if err != nil {
			h.logger.Error("failed to marshal config", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config"})
			return
		}
		source.Config = configData
	}

	if err := h.service.UpdateAuthSource(c.Request.Context(), source); err != nil {
		h.logger.Error("failed to update auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

// Delete 删除认证源
func (h *AuthSourceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteAuthSource(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用认证源
func (h *AuthSourceHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnableAuthSource(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用认证源
func (h *AuthSourceHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisableAuthSource(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable auth source", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// TestLDAPUsers 测试 LDAP 用户查询
func (h *AuthSourceHandler) TestLDAPUsers(c *gin.Context) {
	var req struct {
		URL                  string `json:"url"`
		BaseDN               string `json:"base_dn"`
		BindDN               string `json:"bind_dn"`
		BindPassword         string `json:"bind_password"`
		UserFilter           string `json:"user_filter"`
		UserIDAttr           string `json:"user_id_attr"`
		UserNameAttr         string `json:"user_name_attr"`
		UserSearchBase       string `json:"user_search_base"`
		GroupSearchBase      string `json:"group_search_base"`
		UserEnabledAttribute string `json:"user_enabled_attribute"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构造临时的LDAP配置对象
	configObj := map[string]string{
		"url":                    req.URL,
		"base_dn":                req.BaseDN,
		"bind_dn":                req.BindDN,
		"bind_password":          req.BindPassword,
		"user_filter":            req.UserFilter,
		"user_id_attr":           req.UserIDAttr,
		"user_name_attr":         req.UserNameAttr,
		"user_search_base":       req.UserSearchBase,
		"group_search_base":      req.GroupSearchBase,
		"user_enabled_attribute": req.UserEnabledAttribute,
	}

	configJSON, err := json.Marshal(configObj)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to marshal config"})
		return
	}

	tempConfig := model.AuthSource{
		Config: configJSON,
	}

	users, err := h.service.TestLDAPUsers(c.Request.Context(), &tempConfig)
	if err != nil {
		h.logger.Error("failed to test ldap users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
		"message": fmt.Sprintf("成功查询到 %d 个用户", len(users)),
	})
}
