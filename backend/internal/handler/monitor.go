package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// MonitorHandler 监控 Handler
type MonitorHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMonitorHandler 创建监控 Handler
func NewMonitorHandler(db *gorm.DB, logger *zap.Logger) *MonitorHandler {
	return &MonitorHandler{
		db:     db,
		logger: logger,
	}
}

// AlertPolicyRequest 告警策略请求
type AlertPolicyRequest struct {
	Name          string  `json:"name" binding:"required"`
	ResourceType  string  `json:"resource_type" binding:"required"`
	Metric        string  `json:"metric" binding:"required"`
	Threshold     float64 `json:"threshold" binding:"required"`
	Duration      int     `json:"duration"`
	Level         string  `json:"level"`
	Enabled       bool    `json:"enabled"`
	Owner         string  `json:"owner"`
	DomainID      uint    `json:"domain_id"`
	ProjectID     uint    `json:"project_id"`
	Description   string  `json:"description"`
	NotifyChannel string  `json:"notify_channel"`
}

// ListAlertPolicies 列出告警策略
func (h *MonitorHandler) ListAlertPolicies(c *gin.Context) {
	var policies []model.AlertPolicy

	query := h.db.Model(&model.AlertPolicy{})

	// 过滤条件
	resourceType := c.Query("resource_type")
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	owner := c.Query("owner")
	if owner != "" {
		query = query.Where("owner = ?", owner)
	}

	level := c.Query("level")
	if level != "" {
		query = query.Where("level = ?", level)
	}

	enabled := c.Query("enabled")
	if enabled == "true" {
		query = query.Where("enabled = ?", true)
	} else if enabled == "false" {
		query = query.Where("enabled = ?", false)
	}

	if err := query.Find(&policies).Error; err != nil {
		h.logger.Error("failed to list alert policies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policies)
}

// CreateAlertPolicy 创建告警策略
func (h *MonitorHandler) CreateAlertPolicy(c *gin.Context) {
	var req AlertPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := model.AlertPolicy{
		Name:          req.Name,
		Status:        "正常",
		Enabled:       req.Enabled,
		ResourceType:  req.ResourceType,
		Metric:        req.Metric,
		Threshold:     req.Threshold,
		Duration:      req.Duration,
		Level:         req.Level,
		Owner:         req.Owner,
		DomainID:      req.DomainID,
		ProjectID:     req.ProjectID,
		Description:   req.Description,
		NotifyChannel: req.NotifyChannel,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if policy.Duration == 0 {
		policy.Duration = 5
	}
	if policy.Level == "" {
		policy.Level = "警告"
	}
	if policy.Owner == "" {
		policy.Owner = "自定义"
	}

	if err := h.db.Create(&policy).Error; err != nil {
		h.logger.Error("failed to create alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, policy)
}

// GetAlertPolicy 获取告警策略详情
func (h *MonitorHandler) GetAlertPolicy(c *gin.Context) {
	id := c.Param("id")

	var policy model.AlertPolicy
	if err := h.db.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		h.logger.Error("failed to get alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// UpdateAlertPolicy 更新告警策略
func (h *MonitorHandler) UpdateAlertPolicy(c *gin.Context) {
	id := c.Param("id")

	var req AlertPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var policy model.AlertPolicy
	if err := h.db.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		h.logger.Error("failed to get alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	policy.Name = req.Name
	policy.ResourceType = req.ResourceType
	policy.Metric = req.Metric
	policy.Threshold = req.Threshold
	policy.Duration = req.Duration
	policy.Level = req.Level
	policy.Enabled = req.Enabled
	policy.Description = req.Description
	policy.NotifyChannel = req.NotifyChannel
	policy.UpdatedAt = time.Now()

	if err := h.db.Save(&policy).Error; err != nil {
		h.logger.Error("failed to update alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// DeleteAlertPolicy 删除告警策略
func (h *MonitorHandler) DeleteAlertPolicy(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.AlertPolicy{}, id).Error; err != nil {
		h.logger.Error("failed to delete alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy deleted successfully"})
}

// ToggleAlertPolicy 启用/禁用告警策略
func (h *MonitorHandler) ToggleAlertPolicy(c *gin.Context) {
	id := c.Param("id")

	var policy model.AlertPolicy
	if err := h.db.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		h.logger.Error("failed to get alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	policy.Enabled = !policy.Enabled
	policy.UpdatedAt = time.Now()

	if err := h.db.Save(&policy).Error; err != nil {
		h.logger.Error("failed to toggle alert policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// AlertHistoryFilter 告警历史过滤条件
type AlertHistoryFilter struct {
	Status      string `json:"status"`
	Level       string `json:"level"`
	ResourceType string `json:"resource_type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// ListAlertHistory 列出告警历史
func (h *MonitorHandler) ListAlertHistory(c *gin.Context) {
	var history []model.AlertHistory

	query := h.db.Model(&model.AlertHistory{})

	// 过滤条件
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	level := c.Query("level")
	if level != "" {
		query = query.Where("level = ?", level)
	}

	resourceType := c.Query("resource_type")
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	startTime := c.Query("start_time")
	if startTime != "" {
		query = query.Where("triggered_at >= ?", startTime)
	}

	endTime := c.Query("end_time")
	if endTime != "" {
		query = query.Where("triggered_at <= ?", endTime)
	}

	// 排序和分页
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query = query.Order("triggered_at desc").Offset(offset).Limit(pageSize)

	if err := query.Find(&history).Error; err != nil {
		h.logger.Error("failed to list alert history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取总数
	var total int64
	h.db.Model(&model.AlertHistory{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"data": history,
		"total": total,
		"page": page,
		"page_size": pageSize,
	})
}

// ResolveAlert 解决告警
func (h *MonitorHandler) ResolveAlert(c *gin.Context) {
	id := c.Param("id")

	var history model.AlertHistory
	if err := h.db.First(&history, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		h.logger.Error("failed to get alert history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	history.Status = "resolved"
	history.ResolvedAt = &now

	if err := h.db.Save(&history).Error; err != nil {
		h.logger.Error("failed to resolve alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// IgnoreAlert 忽略告警
func (h *MonitorHandler) IgnoreAlert(c *gin.Context) {
	id := c.Param("id")

	var history model.AlertHistory
	if err := h.db.First(&history, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		h.logger.Error("failed to get alert history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	history.Status = "ignored"

	if err := h.db.Save(&history).Error; err != nil {
		h.logger.Error("failed to ignore alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// MonitorResourcesRequest 监控资源请求
type MonitorResourcesRequest struct {
	AccountID uint `json:"account_id" binding:"required"`
}

// ListMonitorResources 列出监控资源
func (h *MonitorHandler) ListMonitorResources(c *gin.Context) {
	var resources []model.MonitorResource

	query := h.db.Model(&model.MonitorResource{})

	accountID := c.Query("account_id")
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}

	resourceType := c.Query("resource_type")
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	monitorStatus := c.Query("monitor_status")
	if monitorStatus != "" {
		query = query.Where("monitor_status = ?", monitorStatus)
	}

	if err := query.Find(&resources).Error; err != nil {
		h.logger.Error("failed to list monitor resources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resources)
}

// SyncMonitorResources 同步监控资源
func (h *MonitorHandler) SyncMonitorResources(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	// 获取云账号信息和资源
	monitorService := service.NewMonitorService(h.db)
	resources, err := monitorService.SyncMonitorResources(c.Request.Context(), uint(accountID))
	if err != nil {
		h.logger.Error("failed to sync monitor resources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count": len(resources),
	})
}

// GetResourceMetrics 获取资源监控指标
func (h *MonitorHandler) GetResourceMetrics(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	resourceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	monitorService := service.NewMonitorService(h.db)
	metrics, err := monitorService.GetResourceMetrics(c.Request.Context(), uint(accountID), resourceID)
	if err != nil {
		h.logger.Error("failed to get resource metrics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}