package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// SecurityAlertHandler 安全告警 Handler
type SecurityAlertHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewSecurityAlertHandler 创建安全告警 Handler
func NewSecurityAlertHandler(db *gorm.DB, logger *zap.Logger) *SecurityAlertHandler {
	return &SecurityAlertHandler{
		db:     db,
		logger: logger,
	}
}

// SecurityAlertResponse 安全告警响应结构（参考 Cloudpods）
type SecurityAlertResponse struct {
	ID            uint                   `json:"id"`
	AlertID       string                 `json:"alert_id"`
	AlertName     string                 `json:"alert_name"`
	AlertState    string                 `json:"alert_state"` // alerting/ok
	Level         string                 `json:"level"`       // important/critical/normal
	Metric        string                 `json:"metric"`
	Data          map[string]interface{} `json:"data"`
	Tags          map[string]string      `json:"tags"`
	IsSetShield   bool                   `json:"is_set_shield"`
	CanDelete     bool                   `json:"can_delete"`
	CanUpdate     bool                   `json:"can_update"`
	Type          string                 `json:"type"`
	Title         string                 `json:"title"`
	Message       string                 `json:"message"`
	UserID        *uint                  `json:"user_id"`
	SourceIP      string                 `json:"source_ip"`
	Status        string                 `json:"status"` // active/resolved/ignored
	HandledAt     *time.Time             `json:"handled_at"`
	HandledBy     *uint                  `json:"handled_by"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// List 列出安全告警
func (h *SecurityAlertHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	alerting := c.Query("alerting") == "true"
	status := c.Query("status")
	level := c.Query("level")

	var alerts []model.SecurityAlert
	var total int64

	query := h.db.Model(&model.SecurityAlert{})

	if alerting {
		query = query.Where("status = ?", "active")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}

	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("failed to count security alerts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&alerts).Error; err != nil {
		h.logger.Error("failed to list security alerts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应
	var responses []SecurityAlertResponse
	for _, a := range alerts {
		resp := h.buildAlertResponse(a)
		responses = append(responses, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": responses,
		"total": total,
	})
}

// buildAlertResponse 构建告警响应
func (h *SecurityAlertHandler) buildAlertResponse(a model.SecurityAlert) SecurityAlertResponse {
	alertState := "ok"
	if a.Status == "active" {
		alertState = "alerting"
	}

	return SecurityAlertResponse{
		ID:          a.ID,
		AlertID:     strconv.FormatUint(uint64(a.ID), 10),
		AlertName:   a.Title,
		AlertState:  alertState,
		Level:       a.Level,
		Metric:      a.Type,
		Data: map[string]interface{}{
			"message":   a.Message,
			"source_ip": a.SourceIP,
		},
		Tags:        map[string]string{},
		IsSetShield: false,
		CanDelete:   a.Status == "resolved" || a.Status == "ignored",
		CanUpdate:   true,
		Type:        a.Type,
		Title:       a.Title,
		Message:     a.Message,
		UserID:      a.UserID,
		SourceIP:    a.SourceIP,
		Status:      a.Status,
		HandledAt:   a.HandledAt,
		HandledBy:   a.HandledBy,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// Get 获取安全告警详情
func (h *SecurityAlertHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var alert model.SecurityAlert
	if err := h.db.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		h.logger.Error("failed to get security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := h.buildAlertResponse(alert)
	c.JSON(http.StatusOK, resp)
}

// Resolve 处理告警
func (h *SecurityAlertHandler) Resolve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var alert model.SecurityAlert
	if err := h.db.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		h.logger.Error("failed to get security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	alert.Status = "resolved"
	alert.HandledAt = &now

	if err := h.db.Save(&alert).Error; err != nil {
		h.logger.Error("failed to resolve security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert resolved"})
}

// Ignore 忽略告警
func (h *SecurityAlertHandler) Ignore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var alert model.SecurityAlert
	if err := h.db.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		h.logger.Error("failed to get security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	alert.Status = "ignored"
	alert.HandledAt = &now

	if err := h.db.Save(&alert).Error; err != nil {
		h.logger.Error("failed to ignore security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert ignored"})
}

// Delete 删除告警
func (h *SecurityAlertHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.db.Delete(&model.SecurityAlert{}, id).Error; err != nil {
		h.logger.Error("failed to delete security alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetStats 获取告警统计
func (h *SecurityAlertHandler) GetStats(c *gin.Context) {
	var total, active, critical, handled int64

	h.db.Model(&model.SecurityAlert{}).Count(&total)
	h.db.Model(&model.SecurityAlert{}).Where("status = ?", "active").Count(&active)
	h.db.Model(&model.SecurityAlert{}).Where("level = ? OR level = ?", "critical", "high").Count(&critical)
	h.db.Model(&model.SecurityAlert{}).Where("status IN ?", []string{"resolved", "ignored"}).Count(&handled)

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"active":   active,
		"critical": critical,
		"handled":  handled,
	})
}