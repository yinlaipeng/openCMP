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

// NotificationChannelHandler 通知渠道 Handler
type NotificationChannelHandler struct {
	channelService *service.NotificationChannelService
	logger         *zap.Logger
}

// NewNotificationChannelHandler 创建通知渠道 Handler
func NewNotificationChannelHandler(db *gorm.DB, logger *zap.Logger) *NotificationChannelHandler {
	return &NotificationChannelHandler{
		channelService: service.NewNotificationChannelService(db),
		logger:         logger,
	}
}

// List 获取通知渠道列表
func (h *NotificationChannelHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	channelType := c.Query("type")

	channels, total, err := h.channelService.ListChannels(c.Request.Context(), channelType, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list notification channels", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": channels,
		"total": total,
	})
}

// Get 获取通知渠道详情
func (h *NotificationChannelHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	channel, err := h.channelService.GetChannel(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if channel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification channel not found"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// Create 创建通知渠道
func (h *NotificationChannelHandler) Create(c *gin.Context) {
	var channel model.NotificationChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.channelService.CreateChannel(c.Request.Context(), &channel); err != nil {
		h.logger.Error("failed to create notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// Update 更新通知渠道
func (h *NotificationChannelHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	channel, err := h.channelService.GetChannel(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if channel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification channel not found"})
		return
	}

	if err := c.ShouldBindJSON(channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.ID = uint(id)

	if err := h.channelService.UpdateChannel(c.Request.Context(), channel); err != nil {
		h.logger.Error("failed to update notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// Delete 删除通知渠道
func (h *NotificationChannelHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.channelService.DeleteChannel(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用通知渠道
func (h *NotificationChannelHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.channelService.EnableChannel(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用通知渠道
func (h *NotificationChannelHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.channelService.DisableChannel(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// Test 测试通知渠道连通性
func (h *NotificationChannelHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	channel, err := h.channelService.GetChannel(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get notification channel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if channel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification channel not found"})
		return
	}

	// 简单的连通性测试：验证渠道存在且启用
	if !channel.Enabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "channel is disabled"})
		return
	}

	// 根据渠道类型进行具体的连通性测试
	testPassed := true
	testMessage := "channel test passed"

	switch channel.Type {
	case "email":
		// 测试邮件服务器连通性
		configBytes := []byte(channel.Config)
		emailCfg, err := service.UnmarshalEmailConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse email config: " + err.Error()
		} else {
			// 简单验证必要的字段是否存在
			if emailCfg.SMTPHost == "" || emailCfg.SMTPPort == 0 || emailCfg.SMTPUser == "" {
				testPassed = false
				testMessage = "missing required email config fields (host, port, username)"
			} else {
				// 在实际应用中，这里应该尝试连接到SMTP服务器
				testMessage = "email config validation passed"
			}
		}
	case "sms":
		// 测试短信服务配置
		configBytes := []byte(channel.Config)
		smsCfg, err := service.UnmarshalSMSConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse SMS config: " + err.Error()
		} else {
			// 验证短信服务必需的字段
			if smsCfg.Provider == "" || smsCfg.AccessKeyID == "" || smsCfg.AccessKeySecret == "" {
				testPassed = false
				testMessage = "missing required SMS config fields (provider, access key, secret)"
			} else {
				// 在实际应用中，这里应该尝试连接到短信服务提供商
				testMessage = "SMS config validation passed"
			}
		}
	case "dingtalk":
		// 测试钉钉机器人
		configBytes := []byte(channel.Config)
		dingtalkCfg, err := service.UnmarshalDingTalkConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse DingTalk config: " + err.Error()
		} else {
			if dingtalkCfg.WebhookURL == "" {
				testPassed = false
				testMessage = "missing required DingTalk config field (webhook URL)"
			} else {
				// 在实际应用中，这里应该尝试向Webhook发送测试消息
				testMessage = "DingTalk webhook validation passed"
			}
		}
	case "wechat":
		// 测试企业微信机器人
		configBytes := []byte(channel.Config)
		weChatCfg, err := service.UnmarshalWeChatConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse WeChat config: " + err.Error()
		} else {
			if weChatCfg.WebhookURL == "" {
				testPassed = false
				testMessage = "missing required WeChat config field (webhook URL)"
			} else {
				// 在实际应用中，这里应该尝试向Webhook发送测试消息
				testMessage = "WeChat webhook validation passed"
			}
		}
	case "feishu":
		// 测试飞书机器人
		configBytes := []byte(channel.Config)
		feishuCfg, err := service.UnmarshalFeishuConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse Feishu config: " + err.Error()
		} else {
			if feishuCfg.WebhookURL == "" {
				testPassed = false
				testMessage = "missing required Feishu config field (webhook URL)"
			} else {
				// 在实际应用中，这里应该尝试向Webhook发送测试消息
				testMessage = "Feishu webhook validation passed"
			}
		}
	case "lark":
		// 测试Lark机器人
		configBytes := []byte(channel.Config)
		larkCfg, err := service.UnmarshalLarkConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse Lark config: " + err.Error()
		} else {
			if larkCfg.WebhookURL == "" {
				testPassed = false
				testMessage = "missing required Lark config field (webhook URL)"
			} else {
				// 在实际应用中，这里应该尝试向Webhook发送测试消息
				testMessage = "Lark webhook validation passed"
			}
		}
	case "webhook":
		// 测试通用Webhook
		configBytes := []byte(channel.Config)
		webhookCfg, err := service.UnmarshalWebhookConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse Webhook config: " + err.Error()
		} else {
			if webhookCfg.URL == "" {
				testPassed = false
				testMessage = "missing required Webhook config field (URL)"
			} else {
				// 在实际应用中，这里应该尝试向Webhook发送测试请求
				testMessage = "Webhook validation passed"
			}
		}
	default:
		testMessage = "unsupported channel type: " + channel.Type
		testPassed = false
	}

	c.JSON(http.StatusOK, gin.H{"success": testPassed, "message": testMessage})
}
