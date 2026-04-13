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

// ScheduledTaskHandler 定时同步任务 Handler
type ScheduledTaskHandler struct {
	taskService    *service.ScheduledTaskService
	accountService *service.CloudAccountService
	logger         *zap.Logger
}

// NewScheduledTaskHandler 创建定时同步任务 Handler
func NewScheduledTaskHandler(db *gorm.DB, logger *zap.Logger) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{
		taskService:    service.NewScheduledTaskService(db),
		accountService: service.NewCloudAccountService(db),
		logger:         logger,
	}
}

// CreateScheduledTaskRequest 创建定时同步任务请求
type CreateScheduledTaskRequest struct {
	Name           string  `json:"name" binding:"required"`
	Type           string  `json:"type" binding:"required"`
	Frequency      string  `json:"frequency" binding:"required"`
	TriggerTime    string  `json:"trigger_time" binding:"required"`
	ValidFrom      *string `json:"valid_from"`
	ValidUntil     *string `json:"valid_until"`
	Status         string  `json:"status"`
	CloudAccountID *uint   `json:"cloud_account_id"`
}

// UpdateScheduledTaskRequest 更新定时同步任务请求
type UpdateScheduledTaskRequest struct {
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Frequency      string  `json:"frequency"`
	TriggerTime    string  `json:"trigger_time"`
	ValidFrom      *string `json:"valid_from"`
	ValidUntil     *string `json:"valid_until"`
	Status         string  `json:"status"`
	CloudAccountID *uint   `json:"cloud_account_id"`
}

// ConvertStringToTime 将字符串转换为 *time.Time
func (h *ScheduledTaskHandler) convertStringToTime(timeStr *string) (*time.Time, error) {
	if timeStr == nil || *timeStr == "" {
		return nil, nil
	}

	parsedTime, err := time.Parse("2006-01-02", *timeStr)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}

// Create 创建定时同步任务
func (h *ScheduledTaskHandler) Create(c *gin.Context) {
	var req CreateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &model.ScheduledTask{
		Name:           req.Name,
		Type:           req.Type,
		Frequency:      req.Frequency,
		TriggerTime:    req.TriggerTime,
		Status:         req.Status,
		CloudAccountID: req.CloudAccountID,
	}

	var err error
	task.ValidFrom, err = h.convertStringToTime(req.ValidFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_from format, expected YYYY-MM-DD"})
		return
	}

	task.ValidUntil, err = h.convertStringToTime(req.ValidUntil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_until format, expected YYYY-MM-DD"})
		return
	}

	if task.Status == "" {
		task.Status = string(model.ScheduledTaskStatusActive)
	}

	if err := h.taskService.CreateScheduledTask(c.Request.Context(), task); err != nil {
		h.logger.Error("failed to create scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// List 列出定时同步任务
func (h *ScheduledTaskHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	tasks, total, err := h.taskService.ListScheduledTasks(c.Request.Context(), pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list scheduled tasks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get 获取定时同步任务详情
func (h *ScheduledTaskHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.taskService.GetScheduledTask(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update 更新定时同步任务
func (h *ScheduledTaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.taskService.GetScheduledTask(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	var req UpdateScheduledTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Type != "" {
		task.Type = req.Type
	}
	if req.Frequency != "" {
		task.Frequency = req.Frequency
	}
	if req.TriggerTime != "" {
		task.TriggerTime = req.TriggerTime
	}
	if req.ValidFrom != nil {
		task.ValidFrom, err = h.convertStringToTime(req.ValidFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_from format, expected YYYY-MM-DD"})
			return
		}
	}
	if req.ValidUntil != nil {
		task.ValidUntil, err = h.convertStringToTime(req.ValidUntil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid valid_until format, expected YYYY-MM-DD"})
			return
		}
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.CloudAccountID != nil {
		task.CloudAccountID = req.CloudAccountID
	}

	if err := h.taskService.UpdateScheduledTask(c.Request.Context(), task); err != nil {
		h.logger.Error("failed to update scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete 删除定时同步任务
func (h *ScheduledTaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.taskService.DeleteScheduledTask(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// UpdateStatus 更新定时同步任务状态
func (h *ScheduledTaskHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.UpdateScheduledTaskStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		h.logger.Error("failed to update scheduled task status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// Execute 执行定时任务（手动触发）
func (h *ScheduledTaskHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.taskService.GetScheduledTask(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get scheduled task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// 检查任务状态
	if task.Status != string(model.ScheduledTaskStatusActive) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task is not active"})
		return
	}

	// 检查任务类型
	switch task.Type {
	case "sync_cloud_account":
		// 执行云账号同步
		if task.CloudAccountID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no cloud account associated"})
			return
		}

		account, err := h.accountService.GetCloudAccount(c.Request.Context(), *task.CloudAccountID)
		if err != nil {
			h.logger.Error("failed to get cloud account", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if account == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
			return
		}

		stats, err := h.accountService.SyncResources(c.Request.Context(), account)
		if err != nil {
			h.logger.Error("failed to sync cloud account", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "task executed successfully",
			"task_id":    task.ID,
			"task_name":  task.Name,
			"statistics": stats,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown task type: " + task.Type})
	}
}
