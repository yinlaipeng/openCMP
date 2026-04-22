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

// WebappHandler 应用程序服务 Handler
type WebappHandler struct {
	service *service.WebappService
	logger  *zap.Logger
}

// NewWebappHandler 创建 Webapp Handler
func NewWebappHandler(db *gorm.DB, logger *zap.Logger) *WebappHandler {
	return &WebappHandler{
		service: service.NewWebappService(db),
		logger:  logger,
	}
}

// List 列出应用程序服务实例
// GET /api/v1/webapp
func (h *WebappHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 搜索参数
	name := c.Query("name")
	status := c.Query("status")
	platform := c.Query("platform")
	cloudAccountID := c.Query("cloud_account_id")
	projectID := c.Query("project_id")
	stack := c.Query("stack")

	filter := service.WebappFilter{
		Name:           name,
		Status:         status,
		Platform:       platform,
		CloudAccountID: cloudAccountID,
		ProjectID:      projectID,
		Stack:          stack,
	}

	items, total, err := h.service.List(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		h.logger.Error("failed to list webapp instances", zap.Error(err))
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

// Get 获取单个应用程序服务详情
// GET /api/v1/webapp/:id
func (h *WebappHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	instance, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get webapp instance", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "webapp instance not found"})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// CreateWebappRequest 创建请求
type CreateWebappRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Stack          string                 `json:"stack"`
	OsType         string                 `json:"os_type"`
	IpAddr         string                 `json:"ip_addr"`
	Domain         string                 `json:"domain"`
	ServerFarm     string                 `json:"server_farm"`
	Platform       string                 `json:"platform"`
	CloudAccountID uint                   `json:"cloud_account_id"`
	RegionID       string                 `json:"region_id"`
	ProjectID      uint                   `json:"project_id"`
	Tags           map[string]string      `json:"tags"`
	Description    string                 `json:"description"`
}

// Create 创建应用程序服务
// POST /api/v1/webapp
func (h *WebappHandler) Create(c *gin.Context) {
	var req CreateWebappRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance := &model.WebappInstance{
		Name:           req.Name,
		Stack:          req.Stack,
		OsType:         req.OsType,
		IpAddr:         req.IpAddr,
		Domain:         req.Domain,
		ServerFarm:     req.ServerFarm,
		Platform:       req.Platform,
		CloudAccountID: req.CloudAccountID,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		Description:    req.Description,
		Status:         "creating",
		Enabled:        true,
	}

	if req.Tags != nil {
		instance.Tags = model.MapToJSON(req.Tags)
	}

	err := h.service.Create(c.Request.Context(), instance)
	if err != nil {
		h.logger.Error("failed to create webapp instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// UpdateWebappRequest 更新请求
type UpdateWebappRequest struct {
	Name        string            `json:"name"`
	Stack       string            `json:"stack"`
	OsType      string            `json:"os_type"`
	IpAddr      string            `json:"ip_addr"`
	Domain      string            `json:"domain"`
	ServerFarm  string            `json:"server_farm"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
	Enabled     *bool             `json:"enabled"`
}

// Update 更新应用程序服务
// PUT /api/v1/webapp/:id
func (h *WebappHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateWebappRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webapp instance not found"})
		return
	}

	if req.Name != "" {
		instance.Name = req.Name
	}
	if req.Stack != "" {
		instance.Stack = req.Stack
	}
	if req.OsType != "" {
		instance.OsType = req.OsType
	}
	if req.IpAddr != "" {
		instance.IpAddr = req.IpAddr
	}
	if req.Domain != "" {
		instance.Domain = req.Domain
	}
	if req.ServerFarm != "" {
		instance.ServerFarm = req.ServerFarm
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
		h.logger.Error("failed to update webapp instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// Delete 删除应用程序服务
// DELETE /api/v1/webapp/:id
func (h *WebappHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to delete webapp instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webapp instance deleted"})
}

// WebappBatchDeleteRequest 批量删除请求
type WebappBatchDeleteRequest struct {
IDs []uint `json:"ids" binding:"required"`
}

// BatchDelete 批量删除
// POST /api/v1/webapp/batch-delete
func (h *WebappHandler) BatchDelete(c *gin.Context) {
var req WebappBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.BatchDelete(c.Request.Context(), req.IDs)
	if err != nil {
		h.logger.Error("failed to batch delete webapp instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webapp instances deleted", "count": len(req.IDs)})
}

// SyncStatus 同步状态
// POST /api/v1/webapp/:id/sync
func (h *WebappHandler) SyncStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.SyncStatus(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to sync webapp status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync completed"})
}