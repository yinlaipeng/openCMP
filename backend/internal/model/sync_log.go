package model

import (
	"time"

	"gorm.io/gorm"
)

// SyncLog 同步日志
type SyncLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`                 // 云账户ID
	CloudAccountName string        `gorm:"size:100" json:"cloud_account_name"`                     // 云账户名称（冗余存储）
	SyncType        string         `gorm:"size:20;not null" json:"sync_type"`                      // 同步类型: full/incremental
	SyncMode        string         `gorm:"size:20;default:'incremental'" json:"sync_mode"`         // 同步模式: incremental/full
	ResourceType    string         `gorm:"size:50;index" json:"resource_type"`                     // 资源类型: vm/vpc/subnet/security_group/eip/disk/image/rds/redis
	SyncStartTime   time.Time      `json:"sync_start_time"`                                        // 同步开始时间
	SyncEndTime     *time.Time     `json:"sync_end_time,omitempty"`                                // 同步结束时间
	SyncDuration    int            `gorm:"default:0" json:"sync_duration"`                         // 同步耗时（秒）
	Status          string         `gorm:"size:20;default:'running'" json:"status"`                // 状态: running/success/partial_failure/failed
	TotalCount      int            `gorm:"default:0" json:"total_count"`                           // 总资源数
	NewCount        int            `gorm:"default:0" json:"new_count"`                             // 新增资源数
	UpdatedCount    int            `gorm:"default:0" json:"updated_count"`                         // 更新资源数
	DeletedCount    int            `gorm:"default:0" json:"deleted_count"`                         // 删除（标记）资源数
	SkippedCount    int            `gorm:"default:0" json:"skipped_count"`                         // 跳过资源数（增量模式）
	ErrorCount      int            `gorm:"default:0" json:"error_count"`                           // 失败资源数
	ErrorMessage    string         `gorm:"type:text" json:"error_message"`                         // 错误信息汇总
	Details         string         `gorm:"type:text" json:"details"`                               // 详细日志（JSON格式）
	DomainID        uint           `gorm:"index;default:1" json:"domain_id"`                       // 所属域
	TriggeredBy     string         `gorm:"size:50" json:"triggered_by"`                            // 触发方式: manual/scheduled
	ScheduledTaskID *uint          `json:"scheduled_task_id,omitempty"`                            // 关联定时任务ID
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SyncLog) TableName() string {
	return "sync_logs"
}

// SyncLogStatus 同步日志状态
type SyncLogStatus string

const (
	SyncLogStatusRunning      SyncLogStatus = "running"
	SyncLogStatusSuccess      SyncLogStatus = "success"
	SyncLogStatusPartialFail  SyncLogStatus = "partial_failure"
	SyncLogStatusFailed       SyncLogStatus = "failed"
)

// SyncMode 同步模式
type SyncMode string

const (
	SyncModeIncremental SyncMode = "incremental"
	SyncModeFull        SyncMode = "full"
)

// SyncResourceType 同步资源类型
type SyncResourceType string

const (
	SyncResourceTypeVM            SyncResourceType = "vm"
	SyncResourceTypeVPC           SyncResourceType = "vpc"
	SyncResourceTypeSubnet        SyncResourceType = "subnet"
	SyncResourceTypeSecurityGroup SyncResourceType = "security_group"
	SyncResourceTypeEIP           SyncResourceType = "eip"
	SyncResourceTypeDisk          SyncResourceType = "disk"
	SyncResourceTypeSnapshot      SyncResourceType = "snapshot"
	SyncResourceTypeImage         SyncResourceType = "image"
	SyncResourceTypeLoadBalancer  SyncResourceType = "load_balancer"
	SyncResourceTypeRDS           SyncResourceType = "rds"
	SyncResourceTypeRedis         SyncResourceType = "redis"
)

// SyncTriggerType 触发方式
type SyncTriggerType string

const (
	SyncTriggerManual    SyncTriggerType = "manual"
	SyncTriggerScheduled SyncTriggerType = "scheduled"
)

// SyncLogDetail 同步日志详情条目
type SyncLogDetail struct {
	ResourceID   string `json:"resource_id"`
	ResourceName string `json:"resource_name"`
	Action       string `json:"action"` // create/update/delete/skip/error
	OldStatus    string `json:"old_status,omitempty"`
	NewStatus    string `json:"new_status,omitempty"`
	ProjectID    uint   `json:"project_id,omitempty"`
	ProjectName  string `json:"project_name,omitempty"`
	MatchedRule  string `json:"matched_rule,omitempty"`
	Tags         string `json:"tags,omitempty"` // JSON格式
	Error        string `json:"error,omitempty"`
	Timestamp    string `json:"timestamp"`
}