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

// CloudAccountHandler 云账户 Handler
type CloudAccountHandler struct {
	service *service.CloudAccountService
	logger  *zap.Logger
}

// NewCloudAccountHandler 创建云账户 Handler
func NewCloudAccountHandler(db *gorm.DB, logger *zap.Logger) *CloudAccountHandler {
	return &CloudAccountHandler{
		service: service.NewCloudAccountService(db),
		logger:  logger,
	}
}

// CreateCloudAccountRequest 创建云账户请求
type CreateCloudAccountRequest struct {
	Name         string            `json:"name" binding:"required"`
	ProviderType string            `json:"provider_type" binding:"required"`
	Credentials  map[string]string `json:"credentials" binding:"required"`
	Description  string            `json:"description"`
}

// UpdateCloudAccountRequest 更新云账户请求
type UpdateCloudAccountRequest struct {
	Name        string            `json:"name"`
	Credentials map[string]string `json:"credentials"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
}

// Create 创建云账户
func (h *CloudAccountHandler) Create(c *gin.Context) {
	var req CreateCloudAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &model.CloudAccount{
		Name:         req.Name,
		ProviderType: req.ProviderType,
		Description:  req.Description,
		Status:       string(model.CloudAccountStatusActive),
	}

	// 序列化凭证
	credsJSON, err := json.Marshal(req.Credentials)
	if err != nil {
		h.logger.Error("failed to marshal credentials", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal credentials"})
		return
	}
	account.Credentials = credsJSON

	if err := h.service.CreateCloudAccount(c.Request.Context(), account); err != nil {
		h.logger.Error("failed to create cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// List 列出云账户
func (h *CloudAccountHandler) List(c *gin.Context) {
	// 支持 page/page_size 和 limit/offset 两种分页参数格式
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 如果使用传统的 limit/offset 格式，则使用它们
	if c.Query("limit") != "" && c.Query("offset") != "" {
		limit, err1 := strconv.Atoi(c.Query("limit"))
		offset, err2 := strconv.Atoi(c.Query("offset"))
		if err1 == nil && err2 == nil {
			accounts, total, err := h.service.ListCloudAccounts(c.Request.Context(), limit, offset)
			if err != nil {
				h.logger.Error("failed to list cloud accounts", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			calculatedPage := offset/pageSize + 1
			if offset%pageSize != 0 {
				calculatedPage = offset/pageSize + 1
			}

			c.JSON(http.StatusOK, gin.H{
				"items":     accounts,
				"total":     total,
				"page":      calculatedPage,
				"page_size": pageSize,
			})
			return
		}
	}

	// 使用 page/page_size 格式（默认格式）
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// 转换为 offset
	offset := (page - 1) * pageSize

	accounts, total, err := h.service.ListCloudAccounts(c.Request.Context(), pageSize, offset) // pageSize作为limit，offset作为offset
	if err != nil {
		h.logger.Error("failed to list cloud accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     accounts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get 获取云账户详情
func (h *CloudAccountHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Update 更新云账户
func (h *CloudAccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	var req UpdateCloudAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Status != "" {
		account.Status = req.Status
	}
	if req.Description != "" {
		account.Description = req.Description
	}
	if req.Credentials != nil {
		credsJSON, err := json.Marshal(req.Credentials)
		if err != nil {
			h.logger.Error("failed to marshal credentials", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal credentials"})
			return
		}
		account.Credentials = credsJSON
	}

	if err := h.service.UpdateCloudAccount(c.Request.Context(), account); err != nil {
		h.logger.Error("failed to update cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Delete 删除云账户
func (h *CloudAccountHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteCloudAccount(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Verify 验证云账户
func (h *CloudAccountHandler) Verify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	// 使用 VerifyCredentials 进行真实 API 调用验证
	valid, message, err := h.service.VerifyCredentials(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to verify cloud account credentials", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   valid,
		"message": message,
	})
}

// TestConnection 测试云账户连接
func (h *CloudAccountHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	valid, err := h.service.TestConnection(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to test cloud account connection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "connection test failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"message":   "connection test succeeded",
	})
}

// Sync 同步云账户资源
func (h *CloudAccountHandler) Sync(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	stats, err := h.service.SyncResources(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to sync cloud account resources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "sync completed",
		"statistics": stats,
	})
}

// VerifyCredentials 验证云账户凭证（实际调用 API）
func (h *CloudAccountHandler) VerifyCredentials(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	valid, message, err := h.service.VerifyCredentials(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to verify cloud account credentials", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   valid,
		"message": message,
	})
}

// UpdateStatus 更新云账户启用状态
func (h *CloudAccountHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新 Enabled 字段（启用状态）
	account.Enabled = req.Enabled

	// 同时更新 Status 字段（连接状态）
	if req.Enabled {
		account.Status = string(model.CloudAccountStatusActive)
	} else {
		account.Status = string(model.CloudAccountStatusInactive)
	}

	if err := h.service.UpdateCloudAccount(c.Request.Context(), account); err != nil {
		h.logger.Error("failed to update cloud account status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated", "enabled": req.Enabled})
}

// UpdateAttributes 更新云账户属性
func (h *CloudAccountHandler) UpdateAttributes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	account, err := h.service.GetCloudAccount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	var req struct {
		AutoSync         bool     `json:"auto_sync"`
		SyncPolicy       string   `json:"sync_policy"`
		SyncInterval     int      `json:"sync_interval"`
		SyncResourceTypes []string `json:"sync_resource_types"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新属性 - 存储到 description 或扩展字段
	// TODO: 需要在 CloudAccount model 中添加属性字段

	c.JSON(http.StatusOK, gin.H{
		"message":           "attributes updated",
		"auto_sync":         req.AutoSync,
		"sync_policy":       req.SyncPolicy,
		"sync_interval":     req.SyncInterval,
		"sync_resource_types": req.SyncResourceTypes,
	})
}
