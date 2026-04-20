package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/service"
)

// MessageHandler 消息 Handler
type MessageHandler struct {
	messageService *service.MessageService
	logger         *zap.Logger
}

// NewMessageHandler 创建消息 Handler
func NewMessageHandler(db *gorm.DB, logger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: service.NewMessageService(db),
		logger:         logger,
	}
}

// List 获取消息列表
func (h *MessageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offsetStr := c.DefaultQuery("offset", "0")

	// 优先使用 limit/offset 参数（新前端），兼容 page/page_size（旧前端）
	if limit > 0 {
		pageSize = limit
		offset, _ := strconv.Atoi(offsetStr)
		if offset >= 0 {
			page = offset/pageSize + 1
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 从 JWT token 中获取用户 ID（更安全）
	// 注意：中间件使用 user_id（下划线），不是 userID（驼峰）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	receiverID := userIDValue.(uint)

	// 也支持前端传递 user_id（兼容模式）
	receiverIDStr := c.Query("user_id")
	if receiverIDStr != "" {
		parsedID, err := strconv.ParseUint(receiverIDStr, 10, 32)
		if err == nil {
			receiverID = uint(parsedID)
		}
	}

	messages, total, err := h.messageService.ListMessages(c.Request.Context(), receiverID, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list messages", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": messages,
		"total": total,
	})
}

// Get 获取单条消息
func (h *MessageHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	msg, err := h.messageService.GetMessage(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get message", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if msg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// MarkRead 标记消息为已读
func (h *MessageHandler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.messageService.MarkAsRead(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to mark message as read", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

// MarkAllRead 全部标记为已读
func (h *MessageHandler) MarkAllRead(c *gin.Context) {
	// 从 JWT token 中获取用户 ID（更安全）
	// 注意：中间件使用 user_id（下划线），不是 userID（驼峰）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	receiverID := userIDValue.(uint)

	// 也支持前端传递 user_id（兼容模式）
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		parsedID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			receiverID = uint(parsedID)
		}
	}

	if err := h.messageService.MarkAllAsRead(c.Request.Context(), receiverID); err != nil {
		h.logger.Error("failed to mark all messages as read", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all messages marked as read"})
}

// Delete 删除消息
func (h *MessageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete message", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetUnreadCount 获取未读消息数量
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	// 从 JWT token 中获取用户 ID（更安全）
	// 注意：中间件使用 user_id（下划线），不是 userID（驼峰）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	receiverID := userIDValue.(uint)

	// 也支持前端传递 user_id（兼容模式）
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		parsedID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			receiverID = uint(parsedID)
		}
	}

	count, err := h.messageService.GetUnreadCount(c.Request.Context(), receiverID)
	if err != nil {
		h.logger.Error("failed to get unread count", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
