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

// CloudAccessGroupHandler 云访问组 Handler
type CloudAccessGroupHandler struct {
	service *service.CloudAccessGroupService
	logger  *zap.Logger
}

// NewCloudAccessGroupHandler 创建云访问组 Handler
func NewCloudAccessGroupHandler(db *gorm.DB, logger *zap.Logger) *CloudAccessGroupHandler {
	return &CloudAccessGroupHandler{
		service: service.NewCloudAccessGroupService(db),
		logger:  logger,
	}
}

// CreateCloudAccessGroupRequest 创建云访问组请求
type CreateCloudAccessGroupRequest struct {
	Name         string `json:"name" binding:"required"`
	Permissions  string `json:"permissions"`
	Platform     string `json:"platform"`
	SharedScope  string `json:"shared_scope"`
	DomainID     uint   `json:"domain_id" binding:"required"`
}

// UpdateCloudAccessGroupRequest 更新云访问组请求
type UpdateCloudAccessGroupRequest struct {
	Name         string `json:"name"`
	Permissions  string `json:"permissions"`
	Platform     string `json:"platform"`
	SharedScope  string `json:"shared_scope"`
}

// CloudAccessGroupBatchDeleteRequest 批量删除云访问组请求
type CloudAccessGroupBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// ProxySettingBatchDeleteRequest 批量删除代理设置请求
type ProxySettingBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// Create 创建云访问组
func (h *CloudAccessGroupHandler) Create(c *gin.Context) {
	var req CreateCloudAccessGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := &model.CloudAccessGroup{
		Name:        req.Name,
		Status:      "active",
		Permissions: req.Permissions,
		Platform:    req.Platform,
		SharedScope: req.SharedScope,
		DomainID:    req.DomainID,
	}

	if err := h.service.CreateCloudAccessGroup(c.Request.Context(), group); err != nil {
		h.logger.Error("failed to create cloud access group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// List 列出云访问组
func (h *CloudAccessGroupHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	limit := pageSize
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if platform := c.Query("platform"); platform != "" {
		filters["platform"] = platform
	}
	if domainID := c.Query("domain_id"); domainID != "" {
		id, _ := strconv.ParseUint(domainID, 10, 32)
		if id > 0 {
			filters["domain_id"] = uint(id)
		}
	}

	groups, total, err := h.service.ListCloudAccessGroups(c.Request.Context(), limit, offset, filters)
	if err != nil {
		h.logger.Error("failed to list cloud access groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"total": total,
	})
}

// Get 获取云访问组详情
func (h *CloudAccessGroupHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.service.GetCloudAccessGroup(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud access group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud access group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Update 更新云访问组
func (h *CloudAccessGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCloudAccessGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Permissions != "" {
		updates["permissions"] = req.Permissions
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.SharedScope != "" {
		updates["shared_scope"] = req.SharedScope
	}

	if err := h.service.UpdateCloudAccessGroup(c.Request.Context(), uint(id), updates); err != nil {
		h.logger.Error("failed to update cloud access group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除云访问组
func (h *CloudAccessGroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteCloudAccessGroup(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete cloud access group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// BatchDelete 批量删除云访问组
func (h *CloudAccessGroupHandler) BatchDelete(c *gin.Context) {
	var req CloudAccessGroupBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BatchDeleteCloudAccessGroups(c.Request.Context(), req.IDs); err != nil {
		h.logger.Error("failed to batch delete cloud access groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "batch deleted successfully"})
}

// ProxySettingHandler 代理设置 Handler
type ProxySettingHandler struct {
	service *service.ProxySettingService
	logger  *zap.Logger
}

// NewProxySettingHandler 创建代理设置 Handler
func NewProxySettingHandler(db *gorm.DB, logger *zap.Logger) *ProxySettingHandler {
	return &ProxySettingHandler{
		service: service.NewProxySettingService(db),
		logger:  logger,
	}
}

// CreateProxySettingRequest 创建代理设置请求
type CreateProxySettingRequest struct {
	Name        string `json:"name" binding:"required"`
	HttpsProxy  string `json:"https_proxy"`
	HttpProxy   string `json:"http_proxy"`
	NoProxy     string `json:"no_proxy"`
	SharedScope string `json:"shared_scope"`
	DomainID    uint   `json:"domain_id" binding:"required"`
}

// UpdateProxySettingRequest 更新代理设置请求
type UpdateProxySettingRequest struct {
	Name        string `json:"name"`
	HttpsProxy  string `json:"https_proxy"`
	HttpProxy   string `json:"http_proxy"`
	NoProxy     string `json:"no_proxy"`
	SharedScope string `json:"shared_scope"`
}

// SetSharingRequest 设置共享请求
type SetSharingRequest struct {
	SharedScope string `json:"shared_scope" binding:"required"`
}

// Create 创建代理设置
func (h *ProxySettingHandler) Create(c *gin.Context) {
	var req CreateProxySettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proxy := &model.ProxySetting{
		Name:        req.Name,
		HttpsProxy:  req.HttpsProxy,
		HttpProxy:   req.HttpProxy,
		NoProxy:     req.NoProxy,
		SharedScope: req.SharedScope,
		DomainID:    req.DomainID,
	}

	if err := h.service.CreateProxySetting(c.Request.Context(), proxy); err != nil {
		h.logger.Error("failed to create proxy setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, proxy)
}

// List 列出代理设置
func (h *ProxySettingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	limit := pageSize
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if domainID := c.Query("domain_id"); domainID != "" {
		id, _ := strconv.ParseUint(domainID, 10, 32)
		if id > 0 {
			filters["domain_id"] = uint(id)
		}
	}

	proxies, total, err := h.service.ListProxySettings(c.Request.Context(), limit, offset, filters)
	if err != nil {
		h.logger.Error("failed to list proxy settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": proxies,
		"total": total,
	})
}

// Get 获取代理设置详情
func (h *ProxySettingHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	proxy, err := h.service.GetProxySetting(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get proxy setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if proxy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "proxy setting not found"})
		return
	}

	c.JSON(http.StatusOK, proxy)
}

// Update 更新代理设置
func (h *ProxySettingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateProxySettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.HttpsProxy != "" {
		updates["https_proxy"] = req.HttpsProxy
	}
	if req.HttpProxy != "" {
		updates["http_proxy"] = req.HttpProxy
	}
	if req.NoProxy != "" {
		updates["no_proxy"] = req.NoProxy
	}
	if req.SharedScope != "" {
		updates["shared_scope"] = req.SharedScope
	}

	if err := h.service.UpdateProxySetting(c.Request.Context(), uint(id), updates); err != nil {
		h.logger.Error("failed to update proxy setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除代理设置
func (h *ProxySettingHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteProxySetting(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete proxy setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// BatchDelete 批量删除代理设置
func (h *ProxySettingHandler) BatchDelete(c *gin.Context) {
	var req ProxySettingBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BatchDeleteProxySettings(c.Request.Context(), req.IDs); err != nil {
		h.logger.Error("failed to batch delete proxy settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "batch deleted successfully"})
}

// SetSharing 设置共享
func (h *ProxySettingHandler) SetSharing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req SetSharingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetProxySettingSharing(c.Request.Context(), uint(id), req.SharedScope); err != nil {
		h.logger.Error("failed to set proxy sharing", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sharing set successfully"})
}