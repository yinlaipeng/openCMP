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

// MessageSubscriptionHandler 消息订阅 Handler
type MessageSubscriptionHandler struct {
	subscriptionService *service.MessageSubscriptionService
	logger              *zap.Logger
	db                  *gorm.DB
}

// NewMessageSubscriptionHandler 创建消息订阅 Handler
func NewMessageSubscriptionHandler(db *gorm.DB, logger *zap.Logger) *MessageSubscriptionHandler {
	return &MessageSubscriptionHandler{
		subscriptionService: service.NewMessageSubscriptionService(db),
		logger:              logger,
		db:                  db,
	}
}

// List 获取消息订阅列表
func (h *MessageSubscriptionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userIDStr := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var subs []model.MessageSubscription
	query := h.db.Model(&model.MessageSubscription{}).Order("created_at DESC")

	// 如果提供了用户ID，则过滤特定用户的数据
	if userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		query = query.Where("user_id = ?", uint(userID))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("failed to count subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&subs).Error; err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     subs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get 获取消息订阅详情
func (h *MessageSubscriptionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.subscriptionService.GetSubscription(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// Create 创建消息订阅
func (h *MessageSubscriptionHandler) Create(c *gin.Context) {
	var sub model.MessageSubscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.subscriptionService.CreateSubscription(c.Request.Context(), &sub); err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// Update 更新消息订阅
func (h *MessageSubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.subscriptionService.GetSubscription(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	if err := c.ShouldBindJSON(sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.ID = uint(id)

	if err := h.subscriptionService.UpdateSubscription(c.Request.Context(), sub); err != nil {
		h.logger.Error("failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// Delete 删除消息订阅
func (h *MessageSubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.subscriptionService.DeleteSubscription(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ListMessageTypes 获取消息类型列表
func (h *MessageSubscriptionHandler) ListMessageTypes(c *gin.Context) {
	var messageTypes []model.MessageType
	if err := h.db.WithContext(c.Request.Context()).Find(&messageTypes).Error; err != nil {
		h.logger.Error("failed to list message types", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": messageTypes,
		"total": len(messageTypes),
	})
}
