package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/opencmp/opencmp/internal/service"
)

// SyncLogHandler 同步日志 Handler
type SyncLogHandler struct {
	service *service.SyncLogService
	logger  *zap.Logger
}

// NewSyncLogHandler 创建同步日志 Handler
func NewSyncLogHandler(service *service.SyncLogService, logger *zap.Logger) *SyncLogHandler {
	return &SyncLogHandler{
		service: service,
		logger:  logger,
	}
}

// ListSyncLogs 获取同步日志列表
func (h *SyncLogHandler) ListSyncLogs(c *gin.Context) {
	// 获取云账户ID
	cloudAccountIDStr := c.Query("cloud_account_id")
	if cloudAccountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cloud_account_id is required"})
		return
	}

	cloudAccountID, err := strconv.ParseUint(cloudAccountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud_account_id"})
		return
	}

	// 获取limit参数
	limit := 20
	limitStr := c.Query("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.service.GetSyncLogs(c.Request.Context(), uint(cloudAccountID), limit)
	if err != nil {
		h.logger.Error("failed to get sync logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": logs,
		"total": len(logs),
	})
}

// GetSyncLog 获取同步日志详情
func (h *SyncLogHandler) GetSyncLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	log, err := h.service.GetSyncLogByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get sync log", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "sync log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetSyncStatistics 获取同步统计信息
func (h *SyncLogHandler) GetSyncStatistics(c *gin.Context) {
	cloudAccountIDStr := c.Query("cloud_account_id")
	if cloudAccountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cloud_account_id is required"})
		return
	}

	cloudAccountID, err := strconv.ParseUint(cloudAccountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud_account_id"})
		return
	}

	// 获取统计天数参数
	days := 30
	daysStr := c.Query("days")
	if daysStr != "" {
		d, err := strconv.Atoi(daysStr)
		if err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.service.GetSyncStatistics(c.Request.Context(), uint(cloudAccountID), days)
	if err != nil {
		h.logger.Error("failed to get sync statistics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetLatestSyncLog 获取最新同步日志
func (h *SyncLogHandler) GetLatestSyncLog(c *gin.Context) {
	cloudAccountIDStr := c.Query("cloud_account_id")
	if cloudAccountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cloud_account_id is required"})
		return
	}

	cloudAccountID, err := strconv.ParseUint(cloudAccountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud_account_id"})
		return
	}

	log, err := h.service.GetLatestSyncLog(c.Request.Context(), uint(cloudAccountID))
	if err != nil {
		h.logger.Error("failed to get latest sync log", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "no sync log found"})
		return
	}

	c.JSON(http.StatusOK, log)
}