package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduledTask 定时同步任务
type ScheduledTask struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`                           // 任务名称
	Type        string         `gorm:"size:50;not null" json:"type"`                           // 任务类型，如sync_cloud_account
	Frequency   string         `gorm:"size:20;default:daily" json:"frequency"`                // 触发频次：once, daily, weekly, monthly, custom
	TriggerTime string         `gorm:"size:10;default:'02:00'" json:"trigger_time"`           // 触发时间，格式HH:mm
	ValidFrom   *time.Time     `json:"valid_from,omitempty"`                                   // 有效开始时间
	ValidUntil  *time.Time     `json:"valid_until,omitempty"`                                  // 有效结束时间
	Status      string         `gorm:"size:20;default:active" json:"status"`                  // 状态：active, inactive
	CloudAccountID *uint       `gorm:"index" json:"cloud_account_id"`                         // 关联的云账户ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

// ScheduledTaskStatus 定时任务状态类型
type ScheduledTaskStatus string

const (
	ScheduledTaskStatusActive   ScheduledTaskStatus = "active"
	ScheduledTaskStatusInactive ScheduledTaskStatus = "inactive"
)

// ScheduledTaskFrequency 任务频次类型
type ScheduledTaskFrequency string

const (
	ScheduledTaskFreqOnce    ScheduledTaskFrequency = "once"
	ScheduledTaskFreqDaily   ScheduledTaskFrequency = "daily"
	ScheduledTaskFreqWeekly  ScheduledTaskFrequency = "weekly"
	ScheduledTaskFreqMonthly ScheduledTaskFrequency = "monthly"
	ScheduledTaskFreqCustom  ScheduledTaskFrequency = "custom"
)