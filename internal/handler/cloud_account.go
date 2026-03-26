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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	accounts, total, err := h.service.ListCloudAccounts(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list cloud accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": accounts,
		"total": total,
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

	valid, err := h.service.VerifyCloudAccount(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to verify cloud account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "verification succeeded"})
}
