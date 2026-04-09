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

// ReceiverHandler 接收人 Handler
type ReceiverHandler struct {
	receiverService *service.ReceiverService
	logger          *zap.Logger
}

// NewReceiverHandler 创建接收人 Handler
func NewReceiverHandler(db *gorm.DB, logger *zap.Logger) *ReceiverHandler {
	return &ReceiverHandler{
		receiverService: service.NewReceiverService(db),
		logger:          logger,
	}
}

// List 获取接收人列表
func (h *ReceiverHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	receivers, total, err := h.receiverService.ListReceivers(c.Request.Context(), pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list receivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": receivers,
		"total": total,
	})
}

// Get 获取接收人详情
func (h *ReceiverHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	receiver, err := h.receiverService.GetReceiver(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if receiver == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
		return
	}

	c.JSON(http.StatusOK, receiver)
}

// Create 创建接收人
func (h *ReceiverHandler) Create(c *gin.Context) {
	var receiver model.Receiver
	if err := c.ShouldBindJSON(&receiver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.receiverService.CreateReceiver(c.Request.Context(), &receiver); err != nil {
		h.logger.Error("failed to create receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, receiver)
}

// Update 更新接收人
func (h *ReceiverHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	receiver, err := h.receiverService.GetReceiver(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if receiver == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
		return
	}

	if err := c.ShouldBindJSON(receiver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	receiver.ID = uint(id)

	if err := h.receiverService.UpdateReceiver(c.Request.Context(), receiver); err != nil {
		h.logger.Error("failed to update receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receiver)
}

// Delete 删除接收人
func (h *ReceiverHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.receiverService.DeleteReceiver(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用接收人
func (h *ReceiverHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.receiverService.EnableReceiver(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用接收人
func (h *ReceiverHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.receiverService.DisableReceiver(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable receiver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// GetChannels 获取接收人的通知渠道
func (h *ReceiverHandler) GetChannels(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	channels, err := h.receiverService.GetNotificationChannelsByReceiver(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get receiver channels", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": channels})
}

// SetChannels 设置接收人的通知渠道
func (h *ReceiverHandler) SetChannels(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		ChannelIDs []uint `json:"channel_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.receiverService.SetNotificationChannelsForReceiver(c.Request.Context(), uint(id), req.ChannelIDs); err != nil {
		h.logger.Error("failed to set receiver channels", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channels updated"})
}

// GetWithChannels 获取接收人详情及通知渠道
func (h *ReceiverHandler) GetWithChannels(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	receiver, err := h.receiverService.GetReceiverWithChannels(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get receiver with channels", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if receiver == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
		return
	}

	c.JSON(http.StatusOK, receiver)
}
