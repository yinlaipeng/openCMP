package model

import (
	"time"

	"gorm.io/gorm"
)

// SyncPolicyExecutionLog 同步策略执行日志
type SyncPolicyExecutionLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SyncPolicyID    uint           `gorm:"index;not null" json:"sync_policy_id"`        // 关联同步策略
	ExecutionTime   time.Time      `gorm:"index" json:"execution_time"`                 // 执行时间
	TriggerType     string         `gorm:"size:20;not null" json:"trigger_type"`        // 触发方式: manual, auto, scheduled
	ResourceCount   int            `gorm:"default:0" json:"resource_count"`             // 资源总数
	MatchedCount    int            `gorm:"default:0" json:"matched_count"`              // 匹配数量
	MappedCount     int            `gorm:"default:0" json:"mapped_count"`               // 映射成功数量
	Result          string         `gorm:"size:20;default:'success'" json:"result"`     // 执行结果: success, partial, failed
	Duration        int            `gorm:"default:0" json:"duration"`                   // 执行耗时(毫秒)
	ErrorMessage    string         `gorm:"size:500" json:"error_message"`               // 错误信息
	Operator        string         `gorm:"size:100" json:"operator"`                    // 操作人
	CloudAccountID  *uint          `json:"cloud_account_id,omitempty"`                  // 云账号ID（可选，指定范围执行时）
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SyncPolicyExecutionLog) TableName() string {
	return "sync_policy_execution_logs"
}

// SyncPolicyMappingResult 同步策略映射结果
type SyncPolicyMappingResult struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SyncPolicyID    uint           `gorm:"index;not null" json:"sync_policy_id"`        // 关联同步策略
	ExecutionLogID  uint           `gorm:"index;not null" json:"execution_log_id"`      // 关联执行日志
	ResourceName    string         `gorm:"size:200;not null" json:"resource_name"`      // 资源名称
	ResourceType    string         `gorm:"size:50;not null" json:"resource_type"`       // 资源类型: vm, vpc, subnet, etc.
	CloudAccountID  uint           `gorm:"index" json:"cloud_account_id"`               // 云账号ID
	MatchedRuleID   uint           `json:"matched_rule_id"`                             // 匹配的规则ID
	MatchedTags     string         `gorm:"size:500" json:"matched_tags"`                // 匹配的标签JSON
	TargetProjectID uint           `json:"target_project_id"`                           // 映射目标项目ID
	TargetProjectName string       `gorm:"size:100" json:"target_project_name"`         // 映射目标项目名称
	MappedAt        time.Time      `json:"mapped_at"`                                   // 映射时间
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SyncPolicyMappingResult) TableName() string {
	return "sync_policy_mapping_results"
}

// 执行结果类型
type ExecutionResult string

const (
	ExecutionResultSuccess  ExecutionResult = "success"
	ExecutionResultPartial  ExecutionResult = "partial"
	ExecutionResultFailed   ExecutionResult = "failed"
)

// 触发类型
type TriggerType string

const (
	TriggerTypeManual    TriggerType = "manual"
	TriggerTypeAuto      TriggerType = "auto"
	TriggerTypeScheduled TriggerType = "scheduled"
)