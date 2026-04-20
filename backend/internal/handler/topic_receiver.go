package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// TopicReceiverHandler 消息类型接收人 Handler（参考 Cloudpods）
type TopicReceiverHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewTopicReceiverHandler 创建消息类型接收人 Handler
func NewTopicReceiverHandler(db *gorm.DB, logger *zap.Logger) *TopicReceiverHandler {
	return &TopicReceiverHandler{
		db:     db,
		logger: logger,
	}
}

// TopicReceiverResponse 接收人响应
type TopicReceiverResponse struct {
	ID           uint   `json:"id"`
	TopicID      uint   `json:"topic_id"`
	ReceiverType string `json:"receiver_type"` // user/group/role
	ReceiverID   uint   `json:"receiver_id"`
	ReceiverName string `json:"receiver_name"`
	Inbox        bool   `json:"inbox"`
	Email        bool   `json:"email"`
	Wechat       bool   `json:"wechat"`
	Dingtalk     bool   `json:"dingtalk"`
	Webhook      bool   `json:"webhook"`
	Enabled      bool   `json:"enabled"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// List 获取消息类型的接收人列表
func (h *TopicReceiverHandler) List(c *gin.Context) {
	topicIDStr := c.Query("topic_id")
	if topicIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic_id is required"})
		return
	}

	topicID, err := strconv.ParseUint(topicIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic_id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var receivers []model.TopicReceiver
	var total int64

	query := h.db.Model(&model.TopicReceiver{}).Where("topic_id = ?", topicID).Order("created_at ASC")

	if err := query.Count(&total).Error; err != nil {
		h.logger.Error("failed to count topic receivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := query.Offset(offset).Limit(limit).Find(&receivers).Error; err != nil {
		h.logger.Error("failed to list topic receivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	items := make([]TopicReceiverResponse, len(receivers))
	for i, r := range receivers {
		items[i] = TopicReceiverResponse{
			ID:           r.ID,
			TopicID:      r.TopicID,
			ReceiverType: r.ReceiverType,
			ReceiverID:   r.ReceiverID,
			ReceiverName: r.ReceiverName,
			Inbox:        r.Inbox,
			Email:        r.Email,
			Wechat:       r.Wechat,
			Dingtalk:     r.Dingtalk,
			Webhook:      r.Webhook,
			Enabled:      r.Enabled,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"items": items,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// Create 创建消息类型接收人
func (h *TopicReceiverHandler) Create(c *gin.Context) {
	var req struct {
		TopicID      uint   `json:"topic_id" binding:"required"`
		ReceiverType string `json:"receiver_type" binding:"required"`
		ReceiverID   uint   `json:"receiver_id" binding:"required"`
		ReceiverName string `json:"receiver_name"`
		Inbox        bool   `json:"inbox"`
		Email        bool   `json:"email"`
		Wechat       bool   `json:"wechat"`
		Dingtalk     bool   `json:"dingtalk"`
		Webhook      bool   `json:"webhook"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取接收人名称（如果未提供）
	if req.ReceiverName == "" {
		req.ReceiverName = h.getReceiverName(req.ReceiverType, req.ReceiverID)
	}

	receiver := model.TopicReceiver{
		TopicID:      req.TopicID,
		ReceiverType: req.ReceiverType,
		ReceiverID:   req.ReceiverID,
		ReceiverName: req.ReceiverName,
		Inbox:        req.Inbox,
		Email:        req.Email,
		Wechat:       req.Wechat,
		Dingtalk:     req.Dingtalk,
		Webhook:      req.Webhook,
		Enabled:      true,
	}

	if err := h.db.Create(&receiver).Error; err != nil {
		h.logger.Error("failed to create topic receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TopicReceiverResponse{
		ID:           receiver.ID,
		TopicID:      receiver.TopicID,
		ReceiverType: receiver.ReceiverType,
		ReceiverID:   receiver.ReceiverID,
		ReceiverName: receiver.ReceiverName,
		Inbox:        receiver.Inbox,
		Email:        receiver.Email,
		Wechat:       receiver.Wechat,
		Dingtalk:     receiver.Dingtalk,
		Webhook:      receiver.Webhook,
		Enabled:      receiver.Enabled,
		CreatedAt:    receiver.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    receiver.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Update 更新消息类型接收人
func (h *TopicReceiverHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var receiver model.TopicReceiver
	if err := h.db.First(&receiver, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
			return
		}
		h.logger.Error("failed to get topic receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		ReceiverName string `json:"receiver_name"`
		Inbox        bool   `json:"inbox"`
		Email        bool   `json:"email"`
		Wechat       bool   `json:"wechat"`
		Dingtalk     bool   `json:"dingtalk"`
		Webhook      bool   `json:"webhook"`
		Enabled      bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiver.ReceiverName = req.ReceiverName
	receiver.Inbox = req.Inbox
	receiver.Email = req.Email
	receiver.Wechat = req.Wechat
	receiver.Dingtalk = req.Dingtalk
	receiver.Webhook = req.Webhook
	receiver.Enabled = req.Enabled

	if err := h.db.Save(&receiver).Error; err != nil {
		h.logger.Error("failed to update topic receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete 删除消息类型接收人
func (h *TopicReceiverHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.db.Delete(&model.TopicReceiver{}, id).Error; err != nil {
		h.logger.Error("failed to delete topic receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// getReceiverName 获取接收人名称
func (h *TopicReceiverHandler) getReceiverName(receiverType string, receiverID uint) string {
	switch receiverType {
	case "user":
		var user model.User
		if h.db.First(&user, receiverID).Error == nil {
			if user.DisplayName != "" {
				return user.DisplayName
			}
			return user.Name
		}
	case "group":
		var group model.Group
		if h.db.First(&group, receiverID).Error == nil {
			return group.Name
		}
	case "role":
		var role model.Role
		if h.db.First(&role, receiverID).Error == nil {
			if role.DisplayName != "" {
				return role.DisplayName
			}
			return role.Name
		}
	}
	return ""
}