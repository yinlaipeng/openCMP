package service

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// NotificationChannelService 通知渠道服务
type NotificationChannelService struct {
	db *gorm.DB
}

// NewNotificationChannelService 创建通知渠道服务
func NewNotificationChannelService(db *gorm.DB) *NotificationChannelService {
	return &NotificationChannelService{db: db}
}

// CreateChannel 创建通知渠道
func (s *NotificationChannelService) CreateChannel(ctx context.Context, channel *model.NotificationChannel) error {
	return s.db.WithContext(ctx).Create(channel).Error
}

// GetChannel 获取通知渠道
func (s *NotificationChannelService) GetChannel(ctx context.Context, id uint) (*model.NotificationChannel, error) {
	var channel model.NotificationChannel
	err := s.db.WithContext(ctx).First(&channel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &channel, nil
}

// GetChannelByName 根据名称获取通知渠道
func (s *NotificationChannelService) GetChannelByName(ctx context.Context, name string) (*model.NotificationChannel, error) {
	var channel model.NotificationChannel
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &channel, nil
}

// ListChannels 列出通知渠道
func (s *NotificationChannelService) ListChannels(ctx context.Context, channelType string, limit, offset int) ([]*model.NotificationChannel, int64, error) {
	var channels []*model.NotificationChannel
	var total int64

	query := s.db.Model(&model.NotificationChannel{})
	if channelType != "" {
		query = query.Where("type = ?", channelType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&channels).Error

	return channels, total, err
}

// UpdateChannel 更新通知渠道
func (s *NotificationChannelService) UpdateChannel(ctx context.Context, channel *model.NotificationChannel) error {
	return s.db.WithContext(ctx).Save(channel).Error
}

// DeleteChannel 删除通知渠道
func (s *NotificationChannelService) DeleteChannel(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.NotificationChannel{}, id).Error
}

// EnableChannel 启用通知渠道
func (s *NotificationChannelService) EnableChannel(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.NotificationChannel{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableChannel 禁用通知渠道
func (s *NotificationChannelService) DisableChannel(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.NotificationChannel{}).Where("id = ?", id).Update("enabled", false).Error
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"` // 启用 TLS
	UseSSL       bool   `json:"use_ssl"` // 启用 SSL（Cloudpods 兼容）
}

// DingTalkConfig 钉钉配置（应用凭证模式）
type DingTalkConfig struct {
	AgentId  string `json:"agent_id"`  // 应用 AgentId
	AppKey   string `json:"app_key"`   // 应用 AppKey
	AppSecret string `json:"app_secret"` // 应用 AppSecret
	// Webhook 模式（向后兼容）
	WebhookURL string `json:"webhook_url,omitempty"` // Webhook 地址（可选）
	Secret     string `json:"secret,omitempty"`      // Webhook 密钥（可选）
}

// WeChatConfig 企业微信配置（应用凭证模式）
type WeChatConfig struct {
	CorpId   string `json:"corp_id"`   // 企业 CorpId
	AgentId  string `json:"agent_id"`  // 应用 AgentId
	Secret   string `json:"secret"`    // 应用 Secret
	// Webhook 模式（向后兼容）
	WebhookURL string `json:"webhook_url,omitempty"` // Webhook 地址（可选）
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	AuthToken string            `json:"auth_token"`
}

// UnmarshalEmailConfig 解析邮件配置
func UnmarshalEmailConfig(configJSON json.RawMessage) (*EmailConfig, error) {
	var cfg EmailConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalDingTalkConfig 解析钉钉配置
func UnmarshalDingTalkConfig(configJSON json.RawMessage) (*DingTalkConfig, error) {
	var cfg DingTalkConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalWeChatConfig 解析企业微信配置
func UnmarshalWeChatConfig(configJSON json.RawMessage) (*WeChatConfig, error) {
	var cfg WeChatConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalWorkwxConfig 解析企业微信配置（workwx 类型）
func UnmarshalWorkwxConfig(configJSON json.RawMessage) (*WeChatConfig, error) {
	var cfg WeChatConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalWebhookConfig 解析 Webhook 配置
func UnmarshalWebhookConfig(configJSON json.RawMessage) (*WebhookConfig, error) {
	var cfg WebhookConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// SMSConfig 短信配置
type SMSConfig struct {
	Provider          string             `json:"provider"` // aliyun, huawei
	AccessKeyID       string             `json:"access_key_id"`
	AccessKeySecret   string             `json:"access_key_secret"`
	Signature         string             `json:"signature"`
	// 简化模板配置（Cloudpods 兼容）
	VerifyCodeTemplate       string `json:"verify_code_template"`       // 验证码模板CODE
	AlertTemplate            string `json:"alert_template"`             // 告警模板CODE
	AbnormalLoginTemplate    string `json:"abnormal_login_template"`    // 异常登录模板CODE
	// 嵌套模板配置（向后兼容）
	DomesticTemplates SMSTemplatesConfig `json:"domestic_templates,omitempty"`
	IntlTemplates     SMSTemplatesConfig `json:"intl_templates,omitempty"`
}

// SMSTemplatesConfig 短信模板配置
type SMSTemplatesConfig struct {
	VerifyCode    string `json:"verify_code"`    // 验证码
	Alert         string `json:"alert"`          // 告警
	AbnormalLogin string `json:"abnormal_login"` // 异常登录
}

// FeishuConfig 飞书配置（应用凭证模式）
type FeishuConfig struct {
	AppId     string `json:"app_id"`     // 应用 AppID
	AppSecret string `json:"app_secret"` // 应用 AppSecret
	// Webhook 模式（向后兼容）
	WebhookURL string `json:"webhook_url,omitempty"` // Webhook 地址（可选）
	Secret     string `json:"secret,omitempty"`      // Webhook 密钥（可选）
}

// LarkConfig Lark配置（与飞书相同结构）
type LarkConfig struct {
	AppId     string `json:"app_id"`     // 应用 AppID
	AppSecret string `json:"app_secret"` // 应用 AppSecret
	// Webhook 模式（向后兼容）
	WebhookURL string `json:"webhook_url,omitempty"` // Webhook 地址（可选）
	Secret     string `json:"secret,omitempty"`      // Webhook 密钥（可选）
}

// UnmarshalSMSConfig 解析短信配置
func UnmarshalSMSConfig(configJSON json.RawMessage) (*SMSConfig, error) {
	var cfg SMSConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalFeishuConfig 解析飞书配置
func UnmarshalFeishuConfig(configJSON json.RawMessage) (*FeishuConfig, error) {
	var cfg FeishuConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}

// UnmarshalLarkConfig 解析Lark配置
func UnmarshalLarkConfig(configJSON json.RawMessage) (*LarkConfig, error) {
	var cfg LarkConfig
	err := json.Unmarshal(configJSON, &cfg)
	return &cfg, err
}
