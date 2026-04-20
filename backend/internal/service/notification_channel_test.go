package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

// TestNotificationChannelService 基本CRUD测试
func TestNotificationChannelService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewNotificationChannelService(db)
	ctx := context.Background()

	// Create a test notification channel
	channel := &model.NotificationChannel{
		Name:        "Test Channel",
		Type:        "email",
		Description: "Test email channel",
		Config:      datatypes.JSON(`{"smtp_host":"smtp.example.com","smtp_port":587,"smtp_user":"test@example.com","smtp_password":"password","from_address":"test@example.com","use_ssl":true}`),
		Enabled:     true,
	}

	// Test CreateChannel
	err := service.CreateChannel(ctx, channel)
	assert.NoError(t, err)
	assert.NotZero(t, channel.ID)

	// Test GetChannel
	retrievedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Equal(t, channel.Name, retrievedChannel.Name)
	assert.Equal(t, channel.Type, retrievedChannel.Type)

	// Test GetChannelByName
	channelByName, err := service.GetChannelByName(ctx, "Test Channel")
	assert.NoError(t, err)
	assert.NotNil(t, channelByName)
	assert.Equal(t, channel.ID, channelByName.ID)

	// Test ListChannels
	channels, total, err := service.ListChannels(ctx, "email", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, channels, 1)
	assert.Equal(t, channel.ID, channels[0].ID)

	// Test UpdateChannel
	channel.Description = "Updated description"
	err = service.UpdateChannel(ctx, channel)
	assert.NoError(t, err)

	updatedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", updatedChannel.Description)

	// Test DisableChannel
	err = service.DisableChannel(ctx, channel.ID)
	assert.NoError(t, err)

	disabledChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.False(t, disabledChannel.Enabled)

	// Test EnableChannel
	err = service.EnableChannel(ctx, channel.ID)
	assert.NoError(t, err)

	enabledChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.True(t, enabledChannel.Enabled)

	// Test DeleteChannel
	err = service.DeleteChannel(ctx, channel.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedChannel)
}

// TestEmailConfigParsing 邮件配置解析测试
func TestEmailConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *EmailConfig
		wantErr bool
	}{
		{
			name:   "完整配置",
			config: `{"smtp_host":"smtp.gmail.com","smtp_port":465,"smtp_user":"user@gmail.com","smtp_password":"password","from_address":"user@gmail.com","from_name":"Test","use_tls":true,"use_ssl":true}`,
			want: &EmailConfig{
				SMTPHost:     "smtp.gmail.com",
				SMTPPort:     465,
				SMTPUser:     "user@gmail.com",
				SMTPPassword: "password",
				FromAddress:  "user@gmail.com",
				FromName:     "Test",
				UseTLS:       true,
				UseSSL:       true,
			},
			wantErr: false,
		},
		{
			name:   "最小配置",
			config: `{"smtp_host":"smtp.example.com","smtp_port":25}`,
			want: &EmailConfig{
				SMTPHost: "smtp.example.com",
				SMTPPort: 25,
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalEmailConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestSMSConfigParsing 短信配置解析测试
func TestSMSConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *SMSConfig
		wantErr bool
	}{
		{
			name:   "阿里云简化模板",
			config: `{"provider":"aliyun","access_key_id":"LTAIxxxx","access_key_secret":"secret","signature":"测试签名","verify_code_template":"SMS_123456","alert_template":"SMS_234567","abnormal_login_template":"SMS_345678"}`,
			want: &SMSConfig{
				Provider:              "aliyun",
				AccessKeyID:           "LTAIxxxx",
				AccessKeySecret:       "secret",
				Signature:             "测试签名",
				VerifyCodeTemplate:    "SMS_123456",
				AlertTemplate:         "SMS_234567",
				AbnormalLoginTemplate: "SMS_345678",
			},
			wantErr: false,
		},
		{
			name:   "嵌套模板配置",
			config: `{"provider":"aliyun","access_key_id":"key","access_key_secret":"secret","signature":"签","domestic_templates":{"verify_code":"SMS_111","alert":"SMS_222","abnormal_login":"SMS_333"},"intl_templates":{"verify_code":"SMS_444"}}`,
			want: &SMSConfig{
				Provider:        "aliyun",
				AccessKeyID:     "key",
				AccessKeySecret: "secret",
				Signature:       "签",
				DomesticTemplates: SMSTemplatesConfig{
					VerifyCode:    "SMS_111",
					Alert:         "SMS_222",
					AbnormalLogin: "SMS_333",
				},
				IntlTemplates: SMSTemplatesConfig{
					VerifyCode: "SMS_444",
				},
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalSMSConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestDingTalkConfigParsing 钉钉配置解析测试
func TestDingTalkConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *DingTalkConfig
		wantErr bool
	}{
		{
			name:   "应用凭证模式",
			config: `{"agent_id":"217947123","app_key":"dingo9s3gzs5123456","app_secret":"secret123"}`,
			want: &DingTalkConfig{
				AgentId:    "217947123",
				AppKey:     "dingo9s3gzs5123456",
				AppSecret:  "secret123",
			},
			wantErr: false,
		},
		{
			name:   "Webhook模式（向后兼容）",
			config: `{"webhook_url":"https://oapi.dingtalk.com/robot/send?access_token=xxx","secret":"SECxxx"}`,
			want: &DingTalkConfig{
				WebhookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
				Secret:     "SECxxx",
			},
			wantErr: false,
		},
		{
			name:   "混合模式",
			config: `{"agent_id":"123","app_key":"key","app_secret":"secret","webhook_url":"https://example.com"}`,
			want: &DingTalkConfig{
				AgentId:    "123",
				AppKey:     "key",
				AppSecret:  "secret",
				WebhookURL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalDingTalkConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestWeChatConfigParsing 企业微信配置解析测试
func TestWeChatConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *WeChatConfig
		wantErr bool
	}{
		{
			name:   "应用凭证模式",
			config: `{"corp_id":"ww2c41e47d2d3b13cb","agent_id":"1000002","secret":"ZgyVyfr2Mvd0zzy6bE5p..."}`,
			want: &WeChatConfig{
				CorpId:   "ww2c41e47d2d3b13cb",
				AgentId:  "1000002",
				Secret:   "ZgyVyfr2Mvd0zzy6bE5p...",
			},
			wantErr: false,
		},
		{
			name:   "Webhook模式（向后兼容）",
			config: `{"webhook_url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"}`,
			want: &WeChatConfig{
				WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalWeChatConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestFeishuConfigParsing 飞书配置解析测试
func TestFeishuConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *FeishuConfig
		wantErr bool
	}{
		{
			name:   "应用凭证模式",
			config: `{"app_id":"cli_9adbc25c4cb2020d","app_secret":"ccyaskdfjLKjkJN5jngse..."}`,
			want: &FeishuConfig{
				AppId:     "cli_9adbc25c4cb2020d",
				AppSecret: "ccyaskdfjLKjkJN5jngse...",
			},
			wantErr: false,
		},
		{
			name:   "Webhook模式（向后兼容）",
			config: `{"webhook_url":"https://open.feishu.cn/open-apis/bot/v2/hook/xxx","secret":"secret"}`,
			want: &FeishuConfig{
				WebhookURL: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
				Secret:     "secret",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalFeishuConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestWorkwxConfigParsing 企业微信(workwx类型)配置解析测试
func TestWorkwxConfigParsing(t *testing.T) {
	config := `{"corp_id":"ww2c41e47d2d3b13cb","agent_id":"1000002","secret":"ZgyVyfr2Mvd0zzy6bE5p..."}`
	got, err := UnmarshalWorkwxConfig(json.RawMessage(config))
	assert.NoError(t, err)
	assert.Equal(t, "ww2c41e47d2d3b13cb", got.CorpId)
	assert.Equal(t, "1000002", got.AgentId)
	assert.Equal(t, "ZgyVyfr2Mvd0zzy6bE5p...", got.Secret)
}

// TestLarkConfigParsing Lark配置解析测试
func TestLarkConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *LarkConfig
		wantErr bool
	}{
		{
			name:   "应用凭证模式",
			config: `{"app_id":"cli_xxx","app_secret":"secret_xxx"}`,
			want: &LarkConfig{
				AppId:     "cli_xxx",
				AppSecret: "secret_xxx",
			},
			wantErr: false,
		},
		{
			name:   "Webhook模式",
			config: `{"webhook_url":"https://open.larksuite.com/open-apis/bot/v2/hook/xxx","secret":"secret"}`,
			want: &LarkConfig{
				WebhookURL: "https://open.larksuite.com/open-apis/bot/v2/hook/xxx",
				Secret:     "secret",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalLarkConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestWebhookConfigParsing Webhook配置解析测试
func TestWebhookConfigParsing(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    *WebhookConfig
		wantErr bool
	}{
		{
			name:   "完整配置",
			config: `{"url":"https://example.com/webhook","method":"POST","headers":{"Content-Type":"application/json"},"auth_token":"token123"}`,
			want: &WebhookConfig{
				URL:       "https://example.com/webhook",
				Method:    "POST",
				Headers:   map[string]string{"Content-Type": "application/json"},
				AuthToken: "token123",
			},
			wantErr: false,
		},
		{
			name:   "最小配置",
			config: `{"url":"https://example.com/hook"}`,
			want: &WebhookConfig{
				URL: "https://example.com/hook",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			config:  `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalWebhookConfig(json.RawMessage(tt.config))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestCreateChannelWithTypeByType 各类型创建测试
func TestCreateChannelWithTypeByType(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewNotificationChannelService(db)
	ctx := context.Background()

	testCases := []struct {
		name   string
		type_  string
		config string
	}{
		{
			name:   "邮件渠道",
			type_:  "email",
			config: `{"smtp_host":"smtp.gmail.com","smtp_port":465,"smtp_user":"user@gmail.com","smtp_password":"password","from_address":"user@gmail.com","use_ssl":true}`,
		},
		{
			name:   "短信渠道",
			type_:  "sms",
			config: `{"provider":"aliyun","access_key_id":"LTAIxxxx","access_key_secret":"secret","signature":"测试签名","verify_code_template":"SMS_123456"}`,
		},
		{
			name:   "钉钉渠道-应用模式",
			type_:  "dingtalk",
			config: `{"agent_id":"217947123","app_key":"dingo9s3gzs5123456","app_secret":"secret123"}`,
		},
		{
			name:   "钉钉渠道-Webhook模式",
			type_:  "dingtalk",
			config: `{"webhook_url":"https://oapi.dingtalk.com/robot/send?access_token=xxx","secret":"SECxxx"}`,
		},
		{
			name:   "飞书渠道-应用模式",
			type_:  "feishu",
			config: `{"app_id":"cli_9adbc25c4cb2020d","app_secret":"ccyaskdfjLKjkJN5jngse..."}`,
		},
		{
			name:   "企业微信渠道-应用模式",
			type_:  "workwx",
			config: `{"corp_id":"ww2c41e47d2d3b13cb","agent_id":"1000002","secret":"ZgyVyfr2Mvd0zzy6bE5p..."}`,
		},
		{
			name:   "企业微信渠道-Webhook模式",
			type_:  "wechat",
			config: `{"webhook_url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"}`,
		},
		{
			name:   "Webhook渠道",
			type_:  "webhook",
			config: `{"url":"https://example.com/webhook","method":"POST","headers":{"Content-Type":"application/json"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			channel := &model.NotificationChannel{
				Name:    tc.name,
				Type:    tc.type_,
				Config:  datatypes.JSON(tc.config),
				Enabled: true,
			}

			err := service.CreateChannel(ctx, channel)
			assert.NoError(t, err)
			assert.NotZero(t, channel.ID)

			// 验证可以正确获取
			retrieved, err := service.GetChannel(ctx, channel.ID)
			assert.NoError(t, err)
			assert.Equal(t, tc.type_, retrieved.Type)

			// 清理
			service.DeleteChannel(ctx, channel.ID)
		})
	}
}