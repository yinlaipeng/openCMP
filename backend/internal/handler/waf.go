package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// WAFHandler WAF策略 Handler
type WAFHandler struct {
	service *service.WAFService
	logger  *zap.Logger
}

// NewWAFHandler 创建 WAF Handler
func NewWAFHandler(db *gorm.DB, logger *zap.Logger) *WAFHandler {
	return &WAFHandler{
		service: service.NewWAFService(db),
		logger:  logger,
	}
}

// ListWAFInstances 列出 WAF 策略实例
// GET /api/v1/waf
func (h *WAFHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 搜索参数
	name := c.Query("name")
	status := c.Query("status")
	platform := c.Query("platform")
	cloudAccountID := c.Query("cloud_account_id")
	domainID := c.Query("domain_id")

	filter := service.WAFFilter{
		Name:          name,
		Status:        status,
		Platform:      platform,
		CloudAccountID: cloudAccountID,
		DomainID:      domainID,
	}

	items, total, err := h.service.List(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		h.logger.Error("failed to list waf instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetWAFInstance 获取单个 WAF 实例详情
// GET /api/v1/waf/:id
func (h *WAFHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	instance, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get waf instance", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "waf instance not found"})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// CreateWAFInstanceRequest 创建 WAF 实例请求
type CreateWAFInstanceRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Type           string                 `json:"type"`
	Platform       string                 `json:"platform"`
	CloudAccountID uint                   `json:"cloud_account_id"`
	DomainID       uint                   `json:"domain_id"`
	RegionID       string                 `json:"region_id"`
	Tags           map[string]string      `json:"tags"`
	Description    string                 `json:"description"`
}

// CreateWAFInstance 创建 WAF 实例
// POST /api/v1/waf
func (h *WAFHandler) Create(c *gin.Context) {
	var req CreateWAFInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance := &model.WAFInstance{
		Name:           req.Name,
		Type:           req.Type,
		Platform:       req.Platform,
		CloudAccountID: req.CloudAccountID,
		DomainID:       req.DomainID,
		RegionID:       req.RegionID,
		Description:    req.Description,
		Status:         "creating",
		Enabled:        true,
	}

	if req.Tags != nil {
		instance.Tags = model.MapToJSON(req.Tags)
	}

	err := h.service.Create(c.Request.Context(), instance)
	if err != nil {
		h.logger.Error("failed to create waf instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// UpdateWAFInstanceRequest 更新 WAF 实例请求
type UpdateWAFInstanceRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
	Enabled     *bool             `json:"enabled"`
}

// UpdateWAFInstance 更新 WAF 实例
// PUT /api/v1/waf/:id
func (h *WAFHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateWAFInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "waf instance not found"})
		return
	}

	if req.Name != "" {
		instance.Name = req.Name
	}
	if req.Description != "" {
		instance.Description = req.Description
	}
	if req.Tags != nil {
		instance.Tags = model.MapToJSON(req.Tags)
	}
	if req.Enabled != nil {
		instance.Enabled = *req.Enabled
	}

	err = h.service.Update(c.Request.Context(), instance)
	if err != nil {
		h.logger.Error("failed to update waf instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// DeleteWAFInstance 删除 WAF 实例
// DELETE /api/v1/waf/:id
func (h *WAFHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to delete waf instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "waf instance deleted"})
}

// BatchDeleteWAFRequest 批量删除请求
type BatchDeleteWAFRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// BatchDelete 批量删除 WAF 实例
// POST /api/v1/waf/batch-delete
func (h *WAFHandler) BatchDelete(c *gin.Context) {
	var req BatchDeleteWAFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.BatchDelete(c.Request.Context(), req.IDs)
	if err != nil {
		h.logger.Error("failed to batch delete waf instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "waf instances deleted", "count": len(req.IDs)})
}

// SyncStatus 同步状态
// POST /api/v1/waf/:id/sync
func (h *WAFHandler) SyncStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.SyncStatus(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to sync waf status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync completed"})
}