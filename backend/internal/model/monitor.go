package model

import (
	"time"

	"gorm.io/gorm"
)

// AlertPolicy 告警策略
type AlertPolicy struct {
	gorm.Model
	Name          string         `gorm:"size:100;not null" json:"name"`
	Status        string         `gorm:"size:20;default:'正常'" json:"status"`
	Enabled       bool           `gorm:"default:true" json:"enabled"`
	ResourceType  string         `gorm:"size:50" json:"resource_type"` // 虚拟机/数据库/网络等
	Metric        string         `gorm:"size:50" json:"metric"`        // cpu_usage/memory_usage等
	Threshold     float64        `gorm:"type:decimal(10,2)" json:"threshold"`
	Duration      int            `gorm:"default:5" json:"duration"` // 持续时间(分钟)
	Level         string         `gorm:"size:20" json:"level"`      // 信息/警告/严重
	Owner         string         `gorm:"size:50" json:"owner"`      // 系统/自定义
	DomainID      uint           `json:"domain_id"`
	ProjectID     uint           `json:"project_id"`
	ResourceTags  string         `gorm:"size:500" json:"resource_tags"` // JSON格式的标签匹配规则
	Description   string         `gorm:"size:500" json:"description"`
	NotifyChannel string         `gorm:"size:100" json:"notify_channel"` // 通知渠道
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// AlertHistory 告警历史
type AlertHistory struct {
	gorm.Model
	PolicyID     uint      `json:"policy_id"`
	PolicyName   string    `gorm:"size:100" json:"policy_name"`
	ResourceID   string    `gorm:"size:100" json:"resource_id"`
	ResourceName string    `gorm:"size:100" json:"resource_name"`
	ResourceType string    `gorm:"size:50" json:"resource_type"`
	Level        string    `gorm:"size:20" json:"level"`
	Status       string    `gorm:"size:20" json:"status"` // pending/resolved/ignored
	MetricValue  float64   `gorm:"type:decimal(10,2)" json:"metric_value"`
	Message      string    `gorm:"size:500" json:"message"`
	TriggeredAt  time.Time `json:"triggered_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	DomainID     uint      `json:"domain_id"`
	ProjectID    uint      `json:"project_id"`
}

// MonitorResource 监控资源
type MonitorResource struct {
	gorm.Model
	ResourceID     string    `gorm:"size:100;uniqueIndex" json:"resource_id"`
	ResourceName   string    `gorm:"size:100" json:"resource_name"`
	ResourceType   string    `gorm:"size:50" json:"resource_type"` // vm/rds/redis等
	MonitorStatus  string    `gorm:"size:20" json:"monitor_status"` // 正常/告警/异常
	AccountID      uint      `json:"account_id"`
	Platform       string    `gorm:"size:50" json:"platform"`
	Region         string    `gorm:"size:50" json:"region"`
	ProjectID      uint      `json:"project_id"`
	LastSyncAt     time.Time `json:"last_sync_at"`
	Metrics        string    `gorm:"type:text" json:"metrics"` // JSON格式的最新指标数据
}