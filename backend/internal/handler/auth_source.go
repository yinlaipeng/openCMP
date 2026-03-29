package handler

import (
	"encoding/json"
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
	AutoCreate  bool                   `json:"auto_create"`
}

// List 列出认证源
func (h *AuthSourceHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	sources, total, err := h.service.ListAuthSources(c.Request.Context(), limit, offset)
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
		AutoCreate:  req.AutoCreate,
	}

	c.JSON(http.StatusCreated, source)
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
	source.AutoCreate = req.AutoCreate
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
