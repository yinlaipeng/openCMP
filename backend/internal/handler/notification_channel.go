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

	h.testChannel(c, channel)
}

// TestNew 测试新建通知渠道配置（无需保存）
func (h *NotificationChannelHandler) TestNew(c *gin.Context) {
	var channel model.NotificationChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.testChannel(c, &channel)
}

// testChannel 测试通知渠道配置的通用方法
func (h *NotificationChannelHandler) testChannel(c *gin.Context, channel *model.NotificationChannel) {
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
			// 验证必要的字段
			if emailCfg.SMTPHost == "" || emailCfg.SMTPPort == 0 || emailCfg.SMTPUser == "" || emailCfg.SMTPPassword == "" || emailCfg.FromAddress == "" {
				testPassed = false
				testMessage = "missing required email config fields (host, port, username, password, from_address)"
			} else {
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
			if smsCfg.Provider == "" || smsCfg.AccessKeyID == "" || smsCfg.AccessKeySecret == "" || smsCfg.Signature == "" {
				testPassed = false
				testMessage = "missing required SMS config fields (provider, access_key_id, access_key_secret, signature)"
			} else {
				// 检查是否有验证码模板
				if smsCfg.VerifyCodeTemplate == "" && smsCfg.DomesticTemplates.VerifyCode == "" {
					testPassed = false
					testMessage = "missing SMS verify code template"
				} else {
					testMessage = "SMS config validation passed"
				}
			}
		}
	case "dingtalk":
		// 测试钉钉应用配置
		configBytes := []byte(channel.Config)
		dingtalkCfg, err := service.UnmarshalDingTalkConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse DingTalk config: " + err.Error()
		} else {
			// 支持应用凭证模式和 Webhook 模式
			if dingtalkCfg.AgentId != "" && dingtalkCfg.AppKey != "" && dingtalkCfg.AppSecret != "" {
				// 应用凭证模式
				testMessage = "DingTalk app config validation passed"
			} else if dingtalkCfg.WebhookURL != "" {
				// Webhook 模式（向后兼容）
				testMessage = "DingTalk webhook validation passed"
			} else {
				testPassed = false
				testMessage = "missing DingTalk config: need (agent_id, app_key, app_secret) or webhook_url"
			}
		}
	case "wechat", "workwx":
		// 测试企业微信应用配置
		configBytes := []byte(channel.Config)
		weChatCfg, err := service.UnmarshalWeChatConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse WeChat config: " + err.Error()
		} else {
			// 支持应用凭证模式和 Webhook 模式
			if weChatCfg.CorpId != "" && weChatCfg.AgentId != "" && weChatCfg.Secret != "" {
				// 应用凭证模式
				testMessage = "WeCom app config validation passed"
			} else if weChatCfg.WebhookURL != "" {
				// Webhook 模式（向后兼容）
				testMessage = "WeCom webhook validation passed"
			} else {
				testPassed = false
				testMessage = "missing WeChat config: need (corp_id, agent_id, secret) or webhook_url"
			}
		}
	case "feishu", "lark":
		// 测试飞书/Lark应用配置
		configBytes := []byte(channel.Config)
		feishuCfg, err := service.UnmarshalFeishuConfig(json.RawMessage(configBytes))
		if err != nil {
			testPassed = false
			testMessage = "failed to parse Feishu config: " + err.Error()
		} else {
			// 支持应用凭证模式和 Webhook 模式
			if feishuCfg.AppId != "" && feishuCfg.AppSecret != "" {
				// 应用凭证模式
				testMessage = "Feishu/Lark app config validation passed"
			} else if feishuCfg.WebhookURL != "" {
				// Webhook 模式（向后兼容）
				testMessage = "Feishu/Lark webhook validation passed"
			} else {
				testPassed = false
				testMessage = "missing Feishu config: need (app_id, app_secret) or webhook_url"
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
				testMessage = "missing required Webhook config field (url)"
			} else {
				testMessage = "Webhook validation passed"
			}
		}
	default:
		testMessage = "unsupported channel type: " + channel.Type
		testPassed = false
	}

	c.JSON(http.StatusOK, gin.H{"success": testPassed, "message": testMessage})
}
