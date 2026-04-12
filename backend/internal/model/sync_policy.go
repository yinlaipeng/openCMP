package model

import (
	"time"

	"gorm.io/gorm"
)

// Rule 同步规则定义
type Rule struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	SyncPolicyID      uint           `gorm:"index" json:"sync_policy_id"`              // 关联同步策略
	ConditionType     string         `gorm:"size:50;not null" json:"condition_type"`   // 条件类型: all_match, any_match, key_match
	ResourceMapping   string         `gorm:"size:50;not null" json:"resource_mapping"` // 资源映射: specify_project, specify_name
	TargetProjectID   *uint          `json:"target_project_id,omitempty"`              // 目标项目ID
	TargetProjectName string         `gorm:"size:100" json:"target_project_name"`      // 目标项目名称
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Tags []RuleTag `gorm:"foreignKey:RuleID" json:"tags,omitempty"` // 规则对应的标签
}

// RuleTag 规则标签关联表
type RuleTag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RuleID    uint           `gorm:"index" json:"rule_id"`               // 关联规则
	TagKey    string         `gorm:"size:100;not null" json:"tag_key"`   // 标签键
	TagValue  string         `gorm:"size:100;not null" json:"tag_value"` // 标签值
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SyncPolicy 同步策略配置
type SyncPolicy struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Remarks   string         `gorm:"size:500" json:"remarks"`
	Status    string         `gorm:"size:20;default:active" json:"status"` // active/inactive
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	Scope     string         `gorm:"size:100" json:"scope"`  // 应用范围
	DomainID  uint           `gorm:"index" json:"domain_id"` // 所属域
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Rules []Rule `gorm:"foreignKey:SyncPolicyID" json:"rules,omitempty"` // 策略对应规则
}

// TableName 指定表名
func (SyncPolicy) TableName() string {
	return "sync_policies"
}

// TableName 指定规则表名
func (Rule) TableName() string {
	return "sync_policy_rules"
}

// TableName 指定规则标签表名
func (RuleTag) TableName() string {
	return "sync_policy_rule_tags"
}

// SyncPolicyStatus 同步策略状态类型
type SyncPolicyStatus string

const (
	SyncPolicyStatusActive   SyncPolicyStatus = "active"
	SyncPolicyStatusInactive SyncPolicyStatus = "inactive"
)

// RuleConditionType 规则条件类型
type RuleConditionType string

const (
	RuleConditionTypeAllMatch RuleConditionType = "all_match"
	RuleConditionTypeAnyMatch RuleConditionType = "any_match"
	RuleConditionTypeKeyMatch RuleConditionType = "key_match"
)

// RuleResourceMapping 资源映射类型
type RuleResourceMapping string

const (
	RuleResourceMappingSpecifyProject RuleResourceMapping = "specify_project"
	RuleResourceMappingSpecifyName    RuleResourceMapping = "specify_name"
)
