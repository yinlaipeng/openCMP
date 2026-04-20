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

	// 解析搜索参数
	searchParams := &service.CloudAccountSearchParams{
		ID:            c.Query("id"),
		Name:          c.Query("name"),
		Remarks:       c.Query("remarks"),
		ProviderType:  c.Query("provider_type"),
		Status:        c.Query("status"),
		HealthStatus:  c.Query("health_status"),
		AccountNumber: c.Query("account_number"),
	}

	// 解析 enabled 参数
	if c.Query("enabled") != "" {
		enabled, err := strconv.ParseBool(c.Query("enabled"))
		if err == nil {
			searchParams.Enabled = &enabled
		}
	}

	// 解析 domain_id 参数
	if c.Query("domain_id") != "" {
		domainIDVal, err := strconv.ParseUint(c.Query("domain_id"), 10, 32)
		if err == nil {
			domainID := uint(domainIDVal)
			searchParams.DomainID = &domainID
		}
	}

	// 如果使用传统的 limit/offset 格式，则使用它们
	limit := pageSize
	offset := (page - 1) * pageSize
	if c.Query("limit") != "" && c.Query("offset") != "" {
		limitVal, err1 := strconv.Atoi(c.Query("limit"))
		offsetVal, err2 := strconv.Atoi(c.Query("offset"))
		if err1 == nil && err2 == nil {
			limit = limitVal
			offset = offsetVal
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
	offset = (page - 1) * pageSize

	// 检查是否有搜索参数
	hasSearchParams := searchParams.ID != "" || searchParams.Name != "" || searchParams.Remarks != "" ||
		searchParams.ProviderType != "" || searchParams.Status != "" || searchParams.HealthStatus != "" ||
		searchParams.AccountNumber != "" || searchParams.Enabled != nil || searchParams.DomainID != nil

	var accounts []*model.CloudAccount
	var total int64
	var err error

	if hasSearchParams {
		// 使用搜索方法
		accounts, total, err = h.service.ListCloudAccountsWithSearch(c.Request.Context(), searchParams, limit, offset)
	} else {
		// 使用普通列表方法
		accounts, total, err = h.service.ListCloudAccounts(c.Request.Context(), limit, offset)
	}

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
		// 更新状态为连接断开
		h.service.RefreshAccountConnectionStatus(c.Request.Context(), uint(id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		// 更新状态为连接断开
		h.service.RefreshAccountConnectionStatus(c.Request.Context(), uint(id))
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "connection test failed",
		})
		return
	}

	// 更新状态为已连接
	h.service.RefreshAccountConnectionStatus(c.Request.Context(), uint(id))

	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"message":   "connection test succeeded",
	})
}

// SyncRequest 同步请求参数
type SyncRequest struct {
	Mode          string   `json:"mode"`           // full 或 incremental
	ResourceTypes []string `json:"resource_types"` // 可选的资源类型列表
}

// Sync 同步云账户资源
func (h *CloudAccountHandler) Sync(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// 解析请求参数
	var syncReq SyncRequest
	if err := c.ShouldBindJSON(&syncReq); err != nil {
		// 如果没有请求体，使用默认值（增量同步）
		syncReq.Mode = "incremental"
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

	// 根据请求参数选择同步模式
	syncMode := model.SyncModeIncremental
	if syncReq.Mode == "full" {
		syncMode = model.SyncModeFull
	}

	// 获取要同步的资源类型列表
	resourceTypes := syncReq.ResourceTypes
	// 如果选择了"all"，清空列表表示同步所有类型
	if len(resourceTypes) == 1 && resourceTypes[0] == "all" {
		resourceTypes = nil
	}

	stats, err := h.service.SyncResourcesWithMode(c.Request.Context(), account, syncMode, model.SyncTriggerManual, nil, resourceTypes)
	if err != nil {
		h.logger.Error("failed to sync cloud account resources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "sync completed",
		"statistics":     stats,
		"sync_mode":      string(syncMode),
		"resource_types": syncReq.ResourceTypes,
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

	// 刷新账号连接状态
	h.service.RefreshAccountConnectionStatus(c.Request.Context(), uint(id))

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

// GetSupportedResourceTypes 获取支持的资源类型列表
func (h *CloudAccountHandler) GetSupportedResourceTypes(c *gin.Context) {
	// 返回所有支持的资源类型
	resourceTypes := []map[string]string{
		{"id": "subnet", "name": "IP子网"},
		{"id": "vm", "name": "主机"},
		{"id": "load_balancer", "name": "负载均衡实例"},
		{"id": "oss", "name": "对象存储"},
		{"id": "rds", "name": "RDS"},
		{"id": "redis", "name": "缓存"},
		{"id": "nat_gateway", "name": "NAT网关"},
		{"id": "file_storage", "name": "文件存储"},
		{"id": "waf", "name": "WAF策略"},
		{"id": "mongodb", "name": "MongoDB实例"},
		{"id": "elasticsearch", "name": "Elasticsearch"},
		{"id": "kafka", "name": "Kafka"},
		{"id": "k8s", "name": "K8s集群"},
		{"id": "vpc_peering", "name": "VPC互联网络"},
		{"id": "cdn", "name": "内容分发网络"},
		{"id": "dns", "name": "DNS解析"},
		{"id": "image", "name": "镜像"},
		{"id": "disk", "name": "云硬盘"},
		{"id": "snapshot", "name": "快照"},
		{"id": "vpc", "name": "VPC"},
		{"id": "security_group", "name": "安全组"},
		{"id": "eip", "name": "弹性公网IP"},
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resourceTypes,
		"total": len(resourceTypes),
	})
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
		SyncPolicyID     *uint    `json:"sync_policy_id"`
		SyncInterval     int      `json:"sync_interval"`
		SyncResourceTypes []string `json:"sync_resource_types"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新云账户的同步策略绑定
	account.SyncPolicyID = req.SyncPolicyID
	if err := h.service.UpdateCloudAccount(c.Request.Context(), account); err != nil {
		h.logger.Error("failed to update cloud account sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "attributes updated",
		"auto_sync":           req.AutoSync,
		"sync_policy_id":      req.SyncPolicyID,
		"sync_interval":       req.SyncInterval,
		"sync_resource_types": req.SyncResourceTypes,
	})
}

// TestConnectionWithCredentials 使用新凭证测试连接
func (h *CloudAccountHandler) TestConnectionWithCredentials(c *gin.Context) {
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
		AccessKeyId     string `json:"access_key_id" binding:"required"`
		AccessKeySecret string `json:"access_key_secret" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credentials are required"})
		return
	}

	// 使用新凭证测试连接
	valid, message, regions, err := h.service.TestConnectionWithCredentials(c.Request.Context(), account, req.AccessKeyId, req.AccessKeySecret)
	if err != nil {
		h.logger.Error("failed to test connection with credentials", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "连接失败: " + err.Error(),
			"regions":   []string{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": valid,
		"message":   message,
		"regions":   regions,
	})
}

// GetRegions 获取云账号可同步区域列表
func (h *CloudAccountHandler) GetRegions(c *gin.Context) {
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

	// 获取云厂商支持的区域列表
	regions, err := h.service.GetAvailableRegions(c.Request.Context(), account)
	if err != nil {
		h.logger.Error("failed to get available regions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": regions,
		"total": len(regions),
	})
}

// BatchSyncRequest 批量同步请求
type BatchSyncRequest struct {
	AccountIds    []uint   `json:"account_ids" binding:"required"`
	Mode          string   `json:"mode"`           // full 或 incremental
	ResourceTypes []string `json:"resource_types"` // 可选的资源类型列表
}

// BatchSync 批量同步云账号
func (h *CloudAccountHandler) BatchSync(c *gin.Context) {
	var req BatchSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_ids are required"})
		return
	}

	if len(req.AccountIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no accounts selected"})
		return
	}

	// 同步模式
	syncMode := model.SyncModeIncremental
	if req.Mode == "full" {
		syncMode = model.SyncModeFull
	}

	// 资源类型
	resourceTypes := req.ResourceTypes
	if len(resourceTypes) == 1 && resourceTypes[0] == "all" {
		resourceTypes = nil
	}

	// 执行批量同步
	results := make([]map[string]interface{}, 0)
	for _, accountId := range req.AccountIds {
		account, err := h.service.GetCloudAccount(c.Request.Context(), accountId)
		if err != nil || account == nil {
			results = append(results, map[string]interface{}{
				"account_id": accountId,
				"success":    false,
				"message":    "account not found",
			})
			continue
		}

		stats, err := h.service.SyncResourcesWithMode(c.Request.Context(), account, syncMode, model.SyncTriggerManual, nil, resourceTypes)
		if err != nil {
			results = append(results, map[string]interface{}{
				"account_id": accountId,
				"success":    false,
				"message":    err.Error(),
			})
			continue
		}

		results = append(results, map[string]interface{}{
			"account_id": accountId,
			"success":    true,
			"message":    "sync completed",
			"statistics": stats,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "batch sync completed",
		"total":       len(req.AccountIds),
		"success":     len(results),
		"sync_mode":   string(syncMode),
		"results":     results,
	})
}

// Export 导出云账号列表
func (h *CloudAccountHandler) Export(c *gin.Context) {
	// 获取所有云账号
	accounts, _, err := h.service.ListCloudAccounts(c.Request.Context(), 1000, 0)
	if err != nil {
		h.logger.Error("failed to list cloud accounts for export", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建导出数据
	exportData := make([]map[string]interface{}, 0)
	for _, account := range accounts {
		exportData = append(exportData, map[string]interface{}{
			"id":            account.ID,
			"name":          account.Name,
			"provider_type": account.ProviderType,
			"status":        account.Status,
			"enabled":       account.Enabled,
			"health_status": account.HealthStatus,
			"balance":       account.Balance,
			"account_number": account.AccountNumber,
			"domain_id":     account.DomainID,
			"last_sync":     account.LastSync,
			"created_at":    account.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": exportData,
		"total": len(exportData),
	})
}
