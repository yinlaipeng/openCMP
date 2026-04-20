package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// MessageTypeHandler 消息类型 Handler（参考 Cloudpods Topic 设计）
type MessageTypeHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMessageTypeHandler 创建消息类型 Handler
func NewMessageTypeHandler(db *gorm.DB, logger *zap.Logger) *MessageTypeHandler {
	return &MessageTypeHandler{
		db:     db,
		logger: logger,
	}
}

// MessageTypeResponse 消息类型响应（参考 Cloudpods）
type MessageTypeResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	Type          string `json:"type"`              // security/resource/automated_process
	Enabled       bool   `json:"enabled"`
	TitleCN       string `json:"title_cn"`
	TitleEN       string `json:"title_en"`
	ContentCN     string `json:"content_cn"`
	ContentEN     string `json:"content_en"`
	ResourceTypes string `json:"resource_types"`
	GroupKeys     string `json:"group_keys"`
	AdvanceDays   string `json:"advance_days"`
	IsSystem      bool   `json:"is_system"`
	CanDelete     bool   `json:"can_delete"`
	CanUpdate     bool   `json:"can_update"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// List 获取消息类型列表
func (h *MessageTypeHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	enabledStr := c.Query("enabled")
	typeFilter := c.Query("type")
	keyword := c.Query("keyword")
	searchField := c.Query("search_field")

	// 兼容 limit/offset 和 page/page_size
	if limit > 0 {
		pageSize = limit
	} else {
		offset = (page - 1) * pageSize
	}

	var messageTypes []model.MessageType
	var total int64

	query := h.db.Model(&model.MessageType{}).Order("created_at ASC")

	// 状态筛选
	if enabledStr != "" {
		enabled := enabledStr == "true"
		query = query.Where("enabled = ?", enabled)
	}

	// 类型筛选
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}

	// 关键词搜索
	if keyword != "" {
		switch searchField {
		case "name":
			query = query.Where("name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
		case "description":
			query = query.Where("description LIKE ?", "%"+keyword+"%")
		default:
			// 默认搜索名称和显示名称
			query = query.Where("name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("failed to count message types", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&messageTypes).Error; err != nil {
		h.logger.Error("failed to list message types", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	items := make([]MessageTypeResponse, len(messageTypes))
	for i, mt := range messageTypes {
		items[i] = MessageTypeResponse{
			ID:            mt.ID,
			Name:          mt.Name,
			DisplayName:   mt.DisplayName,
			Description:   mt.Description,
			Type:          mt.Type,
			Enabled:       mt.Enabled,
			TitleCN:       mt.TitleCN,
			TitleEN:       mt.TitleEN,
			ContentCN:     mt.ContentCN,
			ContentEN:     mt.ContentEN,
			ResourceTypes: mt.ResourceTypes,
			GroupKeys:     mt.GroupKeys,
			AdvanceDays:   mt.AdvanceDays,
			IsSystem:      mt.IsSystem,
			CanDelete:     mt.CanDelete,
			CanUpdate:     mt.CanUpdate,
			CreatedAt:     mt.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     mt.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"items": items,
		"total": total,
		"limit": pageSize,
		"offset": offset,
	})
}

// Get 获取消息类型详情
func (h *MessageTypeHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mt model.MessageType
	if err := h.db.First(&mt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "message type not found"})
			return
		}
		h.logger.Error("failed to get message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageTypeResponse{
		ID:            mt.ID,
		Name:          mt.Name,
		DisplayName:   mt.DisplayName,
		Description:   mt.Description,
		Type:          mt.Type,
		Enabled:       mt.Enabled,
		TitleCN:       mt.TitleCN,
		TitleEN:       mt.TitleEN,
		ContentCN:     mt.ContentCN,
		ContentEN:     mt.ContentEN,
		ResourceTypes: mt.ResourceTypes,
		GroupKeys:     mt.GroupKeys,
		AdvanceDays:   mt.AdvanceDays,
		IsSystem:      mt.IsSystem,
		CanDelete:     mt.CanDelete,
		CanUpdate:     mt.CanUpdate,
		CreatedAt:     mt.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     mt.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Enable 启用消息类型
func (h *MessageTypeHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mt model.MessageType
	if err := h.db.First(&mt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "message type not found"})
			return
		}
		h.logger.Error("failed to get message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mt.Enabled = true
	if err := h.db.Save(&mt).Error; err != nil {
		h.logger.Error("failed to enable message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用消息类型
func (h *MessageTypeHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mt model.MessageType
	if err := h.db.First(&mt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "message type not found"})
			return
		}
		h.logger.Error("failed to get message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 系统内置消息类型不可禁用
	if mt.IsSystem {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot disable system message type"})
		return
	}

	mt.Enabled = false
	if err := h.db.Save(&mt).Error; err != nil {
		h.logger.Error("failed to disable message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// Update 更新消息类型
func (h *MessageTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mt model.MessageType
	if err := h.db.First(&mt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "message type not found"})
			return
		}
		h.logger.Error("failed to get message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 系统内置消息类型不可修改核心字段
	if mt.IsSystem {
		var req struct {
			DisplayName string `json:"display_name"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mt.DisplayName = req.DisplayName
		mt.Description = req.Description
	} else {
		var req model.MessageType
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mt.DisplayName = req.DisplayName
		mt.Description = req.Description
		mt.TitleCN = req.TitleCN
		mt.TitleEN = req.TitleEN
		mt.ContentCN = req.ContentCN
		mt.ContentEN = req.ContentEN
		mt.ResourceTypes = req.ResourceTypes
		mt.GroupKeys = req.GroupKeys
		mt.AdvanceDays = req.AdvanceDays
	}

	if err := h.db.Save(&mt).Error; err != nil {
		h.logger.Error("failed to update message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mt)
}

// Delete 删除消息类型
func (h *MessageTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var mt model.MessageType
	if err := h.db.First(&mt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "message type not found"})
			return
		}
		h.logger.Error("failed to get message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 系统内置消息类型不可删除
	if mt.IsSystem || !mt.CanDelete {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete system message type"})
		return
	}

	if err := h.db.Delete(&mt).Error; err != nil {
		h.logger.Error("failed to delete message type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}